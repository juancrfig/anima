package cli

import (
	"bytes"
	"errors"
	"fmt"

    "anima/internal/crypto"
	"github.com/spf13/cobra"
)


// readPasswordWithConfirmation is a helper to read and confirm a new password.
func readPasswordWithConfirmation(prompt string) ([]byte, error) {
	password, err := ReadPassword(prompt)
	if err != nil {
		return nil, err
	}
	if len(password) == 0 {
		return nil, errors.New("password cannot be empty")
	}

	confirm, err := ReadPassword("Confirm password: ")
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(password, confirm) {
		return nil, errors.New("passwords do not match")
	}
	return password, nil
}

// clearBuffer wipes a byte slice to prevent leaking secrets.
func clearBuffer(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

// SetupCmd creates the one-time `anima setup` command.
func SetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "Initialize and encrypt your Anima journal.",
		Long:  `This command performs the one-time setup for Anima. It creates your master password and a unique Recovery Phrase.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			services, err := GetServices(cmd.Context())
			if err != nil {
				return err
			}

			// 1. Check if setup is already done
			if services.Config.IsSetup() {
				return errors.New("setup has already been completed")
			}

			fmt.Println("--- Welcome to Anima ---")
            fmt.Print(animaArt)
			fmt.Println("Let's set up your secure journal.")

			// 2. Read and confirm password
			password, err := readPasswordWithConfirmation("Create a new password: ")
			if err != nil {
				return err
			}
			defer clearBuffer(password)

			// 3. Use the KeyManager to generate keys
			fmt.Println("Generating your encryption keys...")
			setupResult, err := services.KeyManager.Setup(password)
			if err != nil {
				return fmt.Errorf("could not generate keys: %w", err)
			}
			defer clearBuffer(setupResult.MasterKey)

			// 4. Save the encrypted keys to config
			if err := services.Config.SetEncryptedMasterKey(setupResult.EncryptedMasterKey); err != nil {
				return fmt.Errorf("could not save master key: %w", err)
			}
			if err := services.Config.SetEncryptedRecoveryKey(setupResult.EncryptedRecoveryKey); err != nil {
				return fmt.Errorf("could not save recovery key: %w", err)
			}

			// 5. CRITICAL: Display recovery phrase and get confirmation
			fmt.Println("\n--- YOUR RECOVERY PHRASE ---")
			fmt.Println("This is the ONLY way to recover your account if you forget your password.")
			fmt.Println("Write it down and store it somewhere safe.")
			fmt.Println("\n" + setupResult.RecoveryPhrase + "\n")

			if _, err := ReadPassword("Press Enter to confirm you have saved this phrase..."); err != nil {
				return err
			}

			fmt.Println("\nSetup complete. You can now use 'anima login' to start a session.")
			return nil
		},
	}
	return cmd
}

// LoginCmd creates the `anima login` command.
func LoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Start a new secure session.",
		Long:  `Prompts for your password to securely decrypt your master key and start an in-memory session.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			services, err := GetServices(cmd.Context())
			if err != nil {
				return err
			}

			// 1. Check if setup is done
			if !services.Config.IsSetup() {
				return errors.New("anima has not been set up. Please run 'anima setup' first")
			}

			// 2. Check if already logged in
			if services.Auth.IsAuthenticated() {
				fmt.Println("You are already logged in.")
				return nil
			}

			// 3. Get password
			password, err := ReadPassword("Enter password: ")
			if err != nil {
				return fmt.Errorf("could not read password: %w", err)
			}
			defer clearBuffer(password)

			// 4. Load the encrypted key from config
			encryptedMasterKey, err := services.Config.GetEncryptedMasterKey()
			if err != nil {
				return fmt.Errorf("could not load master key: %w", err)
			}

			// 5. Use KeyManager to decrypt it
			masterKey, err := services.KeyManager.RecoverMasterKey(encryptedMasterKey, password)
			if err != nil {
				// This catches ErrInvalidCredentials
				return fmt.Errorf("login failed: %w", err)
			}
			defer clearBuffer(masterKey)

			// 6. Get session duration and start the session
			duration, err := services.Config.SessionDuration()
			if err != nil {
				return fmt.Errorf("could not get session duration: %w", err)
			}
			services.Auth.SetPassword(masterKey, duration) // We set the MASTER KEY as the session password

			fmt.Println("Login successful.")
			if duration > 0 {
				fmt.Printf("Session will expire in %v.\n", duration)
			}

			return nil
		},
	}
	return cmd
}

// LogoutCmd creates the `anima logout` command.
func LogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "End the current secure session.",
		Long:  `Logs out of the current session and clears the encryption key from memory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			services, err := GetServices(cmd.Context())
			if err != nil {
				return err
			}

			if !services.Auth.IsAuthenticated() {
				fmt.Println("Not logged in.")
				return nil
			}

			services.Auth.Clear()
			fmt.Println("Session ended. Your key has been cleared from memory.")
			return nil
		},
	}
	return cmd
}

// RecoverCmd creates the `anima recover` command.
func RecoverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "recover",
		Short: "Recover access using your Recovery Phrase.",
		Long:  `Recover access to your journal using your 24-word Recovery Phrase. This will allow you to set a new password.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			services, err := GetServices(cmd.Context())
			if err != nil {
				return err
			}

			if !services.Config.IsSetup() {
				return errors.New("anima has not been set up")
			}

			// 1. Get Recovery Phrase
			fmt.Println("Enter your 24-word Recovery Phrase, separated by spaces:")
			phraseStr, err := ReadPassword("> ") // Use ReadPassword to hide input
			if err != nil {
				return err
			}
			phrase := []byte(phraseStr)
			defer clearBuffer(phrase)

			// 2. Load the encrypted recovery key
			encryptedRecoveryKey, err := services.Config.GetEncryptedRecoveryKey()
			if err != nil {
				return fmt.Errorf("could not load recovery key: %w", err)
			}

			// 3. Attempt to decrypt
			masterKey, err := services.KeyManager.RecoverMasterKey(encryptedRecoveryKey, phrase)
			if err != nil {
				return fmt.Errorf("recovery failed: %w", err)
			}
			defer clearBuffer(masterKey)

			fmt.Println("\nRecovery successful. Please set a new password.")

			// 4. Set a new password
			newPassword, err := readPasswordWithConfirmation("Enter new password: ")
			if err != nil {
				return err
			}
			defer clearBuffer(newPassword)

			// 5. Re-encrypt the master key with the new password
			cryptoParams, err := services.Config.CryptoParams()
			if err != nil {
				return err
			}
			newEncryptedMasterKey, err := crypto.Encrypt(masterKey, newPassword, cryptoParams)
			if err != nil {
				return fmt.Errorf("could not encrypt with new password: %w", err)
			}

			// 6. Save the *new* encrypted master key
			if err := services.Config.SetEncryptedMasterKey(newEncryptedMasterKey); err != nil {
				return fmt.Errorf("could not save new master key: %w", err)
			}

			fmt.Println("Password updated. You can now log in with your new password.")
			return nil
		},
	}
	return cmd
}
