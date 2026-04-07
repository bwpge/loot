package cmd

import (
	"fmt"

	"loot/internal/ui"

	"github.com/spf13/cobra"
)

var (
	updateValue   string
	updateComment string
	updateTags    []string
	updateHosts   []string
)

var cmdUpdate = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an entry in the loot file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		changed := func(v string) bool {
			return cmd.Flags().Lookup(v).Changed
		}

		if !changed("value") && !changed("comment") && !changed("tags") && !changed("hosts") {
			bail("at least one update flag is required")
		}
		s, f := loadLootFile()

		id, err := s.FindID(args[0])
		if err != nil {
			bail(err)
		}
		e, err := s.Get(id)
		if err != nil {
			bail(err)
		}

		fmt.Println("updating:", ui.ID(id))

		// special case, need to update hashes
		if changed("value") {
			fmt.Printf(
				"  %s %s -> %s\n",
				ui.Header("Value:"),
				truncate(e.Value),
				truncate(updateValue),
			)
			s.UpdateValue(id, updateValue)
		}
		if changed("comment") {
			fmt.Printf(
				"  %s %s -> %s\n",
				ui.Header("Comment:"),
				truncate(e.Comment),
				truncate(updateComment),
			)
			e.Comment = updateComment
		}
		if changed("tags") {
			fmt.Printf(
				"  %s %s -> %s\n",
				ui.Header("Tags:"),
				truncate(listString(e.Tags)),
				truncate(listString(updateTags)),
			)
			e.Tags = updateTags
		}
		if changed("hosts") {
			fmt.Printf(
				"  %s %s -> %s\n",
				ui.Header("Hosts:"),
				truncate(listString(e.Hosts)),
				truncate(listString(updateHosts)),
			)
			e.Hosts = updateHosts
		}

		s.Data[id] = *e

		s.Save(f)
	},
	ValidArgsFunction: idCompletion,
}

func init() {
	cmdUpdate.Flags().StringVarP(&updateValue, "value", "v", "", "The new value for the entry")
	cmdUpdate.Flags().
		StringVarP(&updateComment, "comment", "c", "", "The new comment for the entry")
	cmdUpdate.Flags().
		StringSliceVarP(&updateTags, "tags", "t", []string{}, "The new tags for the entry (replaces all)")
	cmdUpdate.Flags().
		StringSliceVarP(&updateHosts, "hosts", "H", []string{}, "The new hosts for the entry (replaces all)")
	rootCmd.AddCommand(cmdUpdate)
}
