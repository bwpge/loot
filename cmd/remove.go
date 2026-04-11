package cmd

import (
	"fmt"

	"loot/internal/ui"

	"github.com/spf13/cobra"
)

var removeAll bool

var removeCmd = &cobra.Command{
	Use:     "remove id [id...]",
	Aliases: []string{"rm"},
	Short:   "Remove an entry from the loot file",
	Run: func(cmd *cobra.Command, args []string) {
		s, f := loadLootFile()

		if len(args) == 0 && !removeAll {
			bail("at least one entry ID or", ui.Cli("--all"), "is required")
		}
		if len(args) > 0 && removeAll {
			bail("entry ID cannot be used with", ui.Cli("--all"))
		}
		if len(s.Data) == 0 {
			warn("no entries to remove")
			return
		}

		if removeAll {
			count := len(s.Data)
			s.Clear()
			s.Save(f)
			fmt.Printf("removed all entries (%d total)\n", count)

			return
		}

		for _, arg := range args {
			id, err := s.FindID(arg)
			if err != nil {
				bail(err)
			}
			err = s.Remove(id)
			if err != nil {
				bail(err)
			}
			fmt.Println("removed", id)
		}

		s.Save(f)
	},
	ValidArgsFunction: idCompletion,
}

func init() {
	removeCmd.Flags().BoolVarP(&removeAll, "all", "a", false, "Remove all entries in the loot file")
	rootCmd.AddCommand(removeCmd)
}
