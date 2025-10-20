package cli

import (
    "errors"
    "os"
    "anima/internal/storage"
    "github.com/spf13/cobra"
)


const animaArt = `
                                                      
                     @@@@@@@@@@@@@@                   
                 @@@@@@@@@@@@@@@@@@@@@                
               @@@@@@@@@@@@@@@@@@@@@@@@@@             
             @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@           
           @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@          
          @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@         
         @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@        
         @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@       
        @@@@@@@@@@@@@@@@@@       @@@@@@@@@@@@@@       
        @@@@@@@@@@@@@@   @@                
        @@@@@@@@@@@@       @@@    @@@@@@@@@@@@@       
                             @@ @@@@@@@@@@@@@@@       
         @@@@@@@@@@@@@@@    @@@@@@@@@@@@@@@@@@@       
         @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@        
          @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@         
           @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@          
             @@@@@@@@@@@@@@@@@@@@@@@@@@@@@@           
               @@@@@@@@@@@@@@@@@@@@@@@@@@             
                 @@@@@@@@@@@@@@@@@@@@@                
                     @@@@@@@@@@@@@@                   
`


func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "anima [command] [flags]",
		Long: animaArt + `
This is a command-line tool that serves you as a simple personal journal. 
You can write your diary entries, and they will be saved securely in a JSON file on your local device.
The more you write, the better you and Anima will get to know yourself.`,
		SilenceUsage: true,

		// PersistentPreRunE runs before any subcommand.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
            // Check for setup on every commmand
			ctx, err := initServices(cmd.Context())
			if err != nil {
				return err
			}
			// Set the command's context to the new one
			cmd.SetContext(ctx)
			
            services, err := GetServices(ctx)
            if err != nil {
                return err
            }

            isSetup := services.Config.IsSetup()
            cmdName := cmd.Name()

            if !isSetup && cmdName != "setup" && cmdName != "recover" && cmdName != "config" {
                return errors.New("Anima has not been set up. Please run 'anima setup' first")
            }
            return nil
		},

		// PersistentPostRunE runs *after* any subcommand.
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			// Retrieve the store from the context and close it.
			if store, ok := cmd.Context().Value(dbStoreKey).(*storage.Storage); ok && store != nil {
				store.Close()
			}
            if services, err := GetServices(cmd.Context()); err == nil && services.Auth != nil {
                services.Auth.Clear()
            }
			return nil
		},
	}

	cmd.CompletionOptions.DisableDefaultCmd = true

    cmd.AddCommand(SetupCmd())
    cmd.AddCommand(LoginCmd())
    cmd.AddCommand(LogoutCmd())
    cmd.AddCommand(RecoverCmd())

	cmd.AddCommand(ConfigCmd())
	cmd.AddCommand(TodayCmd())
	cmd.AddCommand(YesterdayCmd())
	cmd.AddCommand(DateCmd())

	return cmd
}

func Execute() {
	if err := RootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
