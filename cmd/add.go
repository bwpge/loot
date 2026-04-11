package cmd

import (
	"fmt"
	"os"
	"strings"

	"loot/internal/ui"
	"loot/loot"

	"github.com/spf13/cobra"
)

var (
	addForce      bool
	addTags       []string
	addHosts      []string
	addInputFiles []string
	addComment    string
)

var addCmd = &cobra.Command{
	Use:   "add <value> [value...]",
	Short: "Add a new entry to loot file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && len(addInputFiles) == 0 {
			bail("at least one value or input file must be provided")
		}

		s, f := loadLootFile()
		doAdd := func(e loot.Entry, skipDup bool) {
			if s.ContainsValue(e.Value) {
				if skipDup {
					return
				}
				if addForce {
					warn("adding duplicate value")
				} else {
					bail("entry value already exists (use", ui.Cli("-f"), "to add anyway)")
				}
			}

			id := s.Add(e)
			tags := ""
			if len(e.Tags) > 0 {
				tags = ui.Comment(listString(e.Tags))
			}

			fmt.Println("added "+id, "->", truncate(e.Value), tags)
		}

		for _, arg := range args {
			user, pass, found := strings.Cut(arg, ":")
			if loot.Config().DetectType && found && !strings.HasPrefix(pass, "//") {
				fmt.Println("detected username:password format")
				doAdd(
					loot.Entry{
						Value:   user,
						Comment: addComment,
						Tags:    append(addTags, "username"),
						Hosts:   addHosts,
					},
					true,
				)
				doAdd(
					loot.Entry{
						Value:   pass,
						Comment: addComment,
						Tags:    append(addTags, "password"),
						Hosts:   addHosts,
					},
					true,
				)
				doAdd(
					loot.Entry{
						Value:   arg,
						Comment: addComment,
						Tags:    append(addTags, "credential"),
						Hosts:   addHosts,
					},
					false,
				)
			} else {
				doAdd(loot.Entry{Value: arg, Comment: addComment, Tags: addTags, Hosts: addHosts}, false)
			}
		}
		for _, f := range addInputFiles {
			bytes, err := os.ReadFile(f)
			if err != nil {
				bail(err)
			}
			doAdd(
				loot.Entry{
					Value:   string(bytes),
					Comment: addComment,
					Tags:    addTags,
					Hosts:   addHosts,
				},
				false,
			)
		}
		s.Save(f)
	},
	ValidArgsFunction: emptyNoFileCompletion,
}

func init() {
	addCmd.Flags().BoolVarP(&addForce, "force", "f", false, "Allow adding duplicate entry values")
	addCmd.Flags().
		StringSliceVarP(&addInputFiles, "input", "i", []string{}, "Add an entry value by file (useful for e.g., ssh keys)")
	addCmd.Flags().
		StringSliceVarP(&addTags, "type", "t", []string{}, "Type of entry (used for filtering)")
	addCmd.Flags().
		StringVarP(&addComment, "comment", "c", "", "Additional note to store with the entry")
	addCmd.Flags().
		StringSliceVarP(&addHosts, "hosts", "H", []string{}, "Host attribution for new entries")
	rootCmd.AddCommand(addCmd)
}
