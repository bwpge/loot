package cmd

import (
	"errors"
	"fmt"
	"os"

	"loot/internal/state"
	"loot/internal/ui"

	"github.com/spf13/cobra"
)

var initForce bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new loot file",
	Run: func(cmd *cobra.Command, args []string) {
		err := cobra.NoArgs(cmd, args)
		if err != nil {
			bail(err)
		}

		forced := false
		if _, err := os.Stat(lootFile); err == nil {
			if !initForce {
				bail(
					"loot file already exists:",
					lootFile,
					"(use",
					ui.Cli("-f"),
					"to overwrite)",
				)
			} else {
				forced = true
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			bail(err)
		}

		s := state.New()
		err = s.Save(lootFile)
		if err != nil {
			bail(err)
		}

		forcedStr := ""
		if forced {
			forcedStr = "(forced)"
		}
		fmt.Println("initialized loot file:", lootFile, forcedStr)
	},
	ValidArgsFunction: emptyNoFileCompletion,
}

func init() {
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Overwrite existing loot file")
	rootCmd.AddCommand(initCmd)
}
