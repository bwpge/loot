package cmd

import (
	"bufio"
	"fmt"
	"maps"
	"os"
	"slices"
	"strings"

	"github.com/bwpge/loot/internal/entry"
	"github.com/bwpge/loot/internal/ui"
	"github.com/spf13/cobra"
)

var (
	removeAll   bool
	removeFlag  []string
	removeTags  []string
	removeHosts []string
	removeYes   bool
)

var removeCmd = &cobra.Command{
	Use:     "remove [id...]",
	Aliases: []string{"rm"},
	Short:   "Remove entries from the loot file",
	Run: func(cmd *cobra.Command, args []string) {
		hasFilter := len(args)+len(removeTags)+len(removeHosts) > 0
		if !hasFilter && len(removeFlag) == 0 && !removeAll {
			bail("at least one id, filter, or", ui.Cli("--all"), "is required")
		}
		if (hasFilter || len(removeFlag) > 0) && removeAll {
			bail("cannot use any arguments with", ui.Cli("--all"))
		}

		s, f := loadLootFile()

		if removeAll {
			count := len(s.Data) + len(s.Flags)
			if count == 0 {
				warn("no entries to remove")
				return
			}
			if !removeYes && !confirm(fmt.Sprintf("remove all entries (%d total)?", count)) {
				return
			}
			s.Clear()
			s.Save(f)
			fmt.Printf("removed all entries (%d total)\n", count)
			return
		}

		if hasFilter {
			matched := s.Filter(entry.Filter{
				ID:    args,
				Tags:  removeTags,
				Hosts: removeHosts,
			})
			if len(matched) == 0 {
				warn("no entries matched filter")
				return
			}

			ids := slices.Sorted(maps.Keys(matched))
			if !removeYes {
				fmt.Println("to be removed:")
				fmt.Println()
				printEntries(matched)
				if !confirm("\nare you sure?") {
					return
				}
			}

			for _, id := range ids {
				if err := s.Remove(id); err != nil {
					bail(err)
				}
				fmt.Println("removed", id)
			}
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
	ValidArgsFunction: completeID,
}

func confirm(prompt string) bool {
	fmt.Fprintf(os.Stderr, "%s [y/N]: ", prompt)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	line = strings.TrimSpace(strings.ToLower(line))
	return line == "y" || line == "yes"
}

func init() {
	removeCmd.Flags().
		BoolVarP(&removeAll, "all", "a", false, "Remove all entries in the loot file (exclusive with other filters)")
	removeCmd.Flags().
		StringSliceVarP(&removeFlag, "flag", "f", []string{}, "One or more flags to remove (exclusive with -a)")
	removeCmd.Flags().
		StringSliceVarP(&removeTags, "tag", "t", []string{}, "Remove entries matching given tags")
	removeCmd.Flags().
		StringSliceVarP(&removeHosts, "host", "H", []string{}, "Remove entries matching given hosts")
	removeCmd.Flags().
		BoolVarP(&removeYes, "yes", "y", false, "Skip confirmation prompts")

	removeCmd.RegisterFlagCompletionFunc("flag", completeFlag)
	removeCmd.RegisterFlagCompletionFunc("tag", completeTag)
	removeCmd.RegisterFlagCompletionFunc("host", completeHost)

	rootCmd.AddCommand(removeCmd)
}
