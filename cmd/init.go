package cmd

import (
	"errors"
	"fmt"
	"os"

	"loot/internal"
	"loot/internal/ui"

	"github.com/spf13/cobra"
)

var (
	initForce        bool
	initDefaultHosts []string
	initDetectType   bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new loot file",
	Run: func(cmd *cobra.Command, args []string) {
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

		s := internal.NewState()
		err := s.Save(lootFile)
		if err != nil {
			bail(err)
		}

		forcedStr := ""
		if forced {
			forcedStr = "(forced)"
		}
		fmt.Println("initialized loot file:", lootFile, forcedStr)
	},
}

func init() {
	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "Overwrite existing loot file")
	initCmd.Flags().BoolVarP(&initDetectType, "detect-type", "d", true, "")
	initCmd.Flags().
		StringSliceVarP(&initDefaultHosts, "default-hosts", "H", []string{}, "Default host attribution for new entries")
	rootCmd.AddCommand(initCmd)
}
