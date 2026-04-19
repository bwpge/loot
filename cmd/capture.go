package cmd

import (
	"fmt"
	"os"

	"github.com/bwpge/loot/internal/entry"
	"github.com/bwpge/loot/internal/ui"
	"github.com/spf13/cobra"
)

var (
	captureUser  string
	captureRoot  string
	captureAdmin bool
	captureHost  string
)

var captureCmd = &cobra.Command{
	Use:     "capture flag",
	Short:   "Store a captured flag",
	Aliases: []string{"cap"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		userSet := cmd.Flags().Changed("user")
		rootSet := cmd.Flags().Changed("root")

		n := 0
		for _, set := range []bool{userSet, rootSet, captureAdmin} {
			if set {
				n++
			}
		}
		if n != 1 {
			bail(
				"flag must be marked with exactly one of",
				ui.Cli("--user")+",",
				ui.Cli("--root")+",",
				"or",
				ui.Cli("--admin"),
			)
		}

		s, f := loadLootFile()
		var ty, owner string
		switch {
		case userSet:
			ty = "user"
			owner = captureUser
		case rootSet:
			ty = "root"
			owner = captureRoot
			if owner == "a" {
				owner = "Administrator"
			}
		case captureAdmin:
			ty = "root"
			owner = "Administrator"
		}

		target := os.Getenv("TARGET")
		if target != "" && captureHost == "" {
			fmt.Println("using TARGET environment variable as HOST")
			captureHost = target
		}

		s.Capture(args[0], entry.Flag{
			Type:  ty,
			Owner: owner,
			Host:  captureHost,
		})
		s.Save(f)
		fmt.Println("captured flag", args[0])
	},
	ValidArgsFunction: cobra.NoFileCompletions,
}

func init() {
	captureCmd.Flags().StringVarP(&captureUser, "user", "u", "", "Mark as a user/local flag")
	captureCmd.Flags().StringVarP(&captureRoot, "root", "r", "", "Mark as root/proof flag")
	captureCmd.Flags().Lookup("root").NoOptDefVal = "root"
	captureCmd.Flags().BoolVarP(&captureAdmin, "admin", "a", false, "Same as --root=Administrator")
	captureCmd.Flags().StringVarP(&captureHost, "host", "H", "", "Host this flag belongs to")
	rootCmd.AddCommand(captureCmd)
}
