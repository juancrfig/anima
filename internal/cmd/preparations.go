package cmd

import (
	"path/filepath"
	"os"
	"context"

	"github.com/spf13/cobra"
)

func ensureAnimaDir(cmd *cobra.Command) error {
 	homePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	p := filepath.Join(homePath, ".anima", "entries")
	
	if err := os.MkdirAll(p, 0700); err != nil {
		return err
	}

	ctx := context.WithValue(cmd.Context(), pathKey, p)
	cmd.SetContext(ctx)

	return nil
}
