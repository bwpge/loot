package cmd

import (
	"fmt"
	"os"

	"github.com/bwpge/loot/internal/entry"
	"github.com/bwpge/loot/internal/ui"
	"github.com/spf13/cobra"
)

var (
	captureUser  bool
	captureRoot  bool
	captureOwner string
	captureHost  string
)

var captureCmd = &cobra.Command{
	Use:     "capture flag",
	Short:   "Store a captured flag",
	Aliases: []string{"cap"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if (!captureUser && !captureRoot) || (captureUser && captureRoot) {
			bail(
				"flag must be marked with",
				ui.Cli("--user"),
				"or",
				ui.Cli("--root"),
				"(not both)",
			)
		}

		s, f := loadLootFile()
		ty := "user"
		if captureRoot {
			ty = "root"
		}

		target := os.Getenv("TARGET")
		if target != "" && captureHost == "" {
			fmt.Println("using TARGET environment variable as HOST")
			captureHost = target
		}

		s.Capture(args[0], entry.Flag{
			Type:  ty,
			Owner: captureOwner,
			Host:  captureHost,
		})
		s.Save(f)
		fmt.Println("captured flag", args[0])
	},
	ValidArgsFunction: cobra.NoFileCompletions,
}

func init() {
	captureCmd.Flags().
		BoolVarP(&captureUser, "user", "u", false, "Mark this as a user/local flag")
	captureCmd.Flags().
		BoolVarP(&captureRoot, "root", "r", false, "Mark this as a root/proof flag")
	captureCmd.Flags().
		StringVarP(&captureOwner, "owner", "o", "", "Owner of the flag (e.g., the user's name)")
	captureCmd.Flags().StringVarP(&captureHost, "host", "H", "", "Host this flag belongs to")
	rootCmd.AddCommand(captureCmd)
}
