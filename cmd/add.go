package cmd

import (
	"fmt"
	"os"
	"strings"

	"loot/internal"
	"loot/internal/ui"

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
		doAdd := func(e internal.Entry, src string, skipDup bool) {
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

			fmt.Println(ui.ID(id), "->", truncate(src), tags)
		}

		for _, arg := range args {
			user, pass, found := strings.Cut(arg, ":")
			if s.Config.DetectType && found && !strings.HasPrefix(pass, "//") {
				fmt.Println("detected username:password format")
				doAdd(
					internal.Entry{
						Value:   user,
						Comment: addComment,
						Tags:    append(addTags, "username"),
						Hosts:   addHosts,
					},
					user,
					true,
				)
				doAdd(
					internal.Entry{
						Value:   pass,
						Comment: addComment,
						Tags:    append(addTags, "password"),
						Hosts:   addHosts,
					},
					pass,
					true,
				)
				doAdd(
					internal.Entry{
						Value:   arg,
						Comment: addComment,
						Tags:    append(addTags, "credential"),
						Hosts:   addHosts,
					},
					arg,
					false,
				)
			} else {
				doAdd(internal.Entry{Value: arg, Comment: addComment, Tags: addTags, Hosts: addHosts}, arg, false)
			}
		}
		for _, f := range addInputFiles {
			bytes, err := os.ReadFile(f)
			if err != nil {
				bail(err)
			}
			doAdd(
				internal.Entry{
					Value:   string(bytes),
					Comment: addComment,
					Tags:    addTags,
					Hosts:   addHosts,
				},
				f,
				false,
			)
		}
		s.Save(f)
	},
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
		StringSliceVarP(&addHosts, "default-hosts", "H", []string{}, "Default host attribution for new entries")
	rootCmd.AddCommand(addCmd)
}
