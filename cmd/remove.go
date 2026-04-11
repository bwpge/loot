package cmd

import (
	"fmt"

	"github.com/bwpge/loot/internal/ui"
	"github.com/spf13/cobra"
)

var (
	removeAll  bool
	removeFlag []string
)

var removeCmd = &cobra.Command{
	Use:     "remove id [id...]",
	Aliases: []string{"rm"},
	Short:   "Remove an entry from the loot file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(removeFlag)+len(args) == 0 && !removeAll {
			bail("at least one id or", ui.Cli("--all"), "is required")
		}
		if len(removeFlag)+len(args) > 0 && removeAll {
			bail("cannot use any arguments with", ui.Cli("--all"))
		}

		s, f := loadLootFile()
		if len(s.Data) == 0 {
			warn("no entries to remove")
			return
		}

		if removeAll {
			count := len(s.Data) + len(s.Flags)
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
		for _, arg := range removeFlag {
			flag, err := s.FindFlag(arg)
			if err != nil {
				bail(err)
			}
			delete(s.Flags, flag)
			fmt.Println("removed", flag)
		}

		s.Save(f)
	},
	ValidArgsFunction: idCompletion,
}

func init() {
	removeCmd.Flags().
		BoolVarP(&removeAll, "all", "a", false, "Remove all entries in the loot file (exclusive with -f)")
	removeCmd.Flags().
		StringSliceVarP(&removeFlag, "flag", "f", []string{}, "One or more flags to remove (exclusive with -a)")
	removeCmd.RegisterFlagCompletionFunc("flag", flagCompletion)
	rootCmd.AddCommand(removeCmd)
}
