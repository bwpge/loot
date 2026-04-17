package cmd

import (
	"fmt"
	"slices"

	"github.com/bwpge/loot/internal/ui"
	"github.com/spf13/cobra"
)

var (
	updateValue   string
	updateComment string
	updateTags    []string
	updateHosts   []string
)

var updateCmd = &cobra.Command{
	Use:   "update id",
	Short: "Update an entry in the loot file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		changed := func(v string) bool {
			return cmd.Flags().Lookup(v).Changed
		}

		if !changed("value") && !changed("comment") && !changed("tag") && !changed("host") {
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

		changes := [][]string{}

		if changed("value") && e.Value != updateValue {
			changes = append(changes, []string{"value", truncate(e.Value), truncate(updateValue)})
			e.Value = updateValue
		}
		if changed("comment") && e.Comment != updateComment {
			changes = append(
				changes,
				[]string{"comment", truncate(e.Comment), truncate(updateComment)},
			)
			e.Comment = updateComment
		}
		if changed("tag") && !sameElems(e.Tags, updateTags) {
			changes = append(
				changes,
				[]string{"tags", truncate(listString(e.Tags)), truncate(listString(updateTags))},
			)
			e.Tags = updateTags
		}
		if changed("host") && !sameElems(e.Hosts, updateHosts) {
			changes = append(
				changes,
				[]string{"hosts", truncate(listString(e.Hosts)), truncate(listString(updateHosts))},
			)
			e.Hosts = updateHosts
		}

		if len(changes) == 0 {
			fmt.Println("nothing to update")
			return
		}

		s.Data[id] = *e

		fmt.Println("updating " + id)
		for _, change := range changes {
			printChange(change[0], change[1], change[2])
		}

		s.Save(f)
	},
	ValidArgsFunction: completeID,
}

func init() {
	updateCmd.Flags().StringVarP(&updateValue, "value", "v", "", "The new value for the entry")
	updateCmd.Flags().
		StringVarP(&updateComment, "comment", "c", "", "The new comment for the entry")
	updateCmd.Flags().
		StringSliceVarP(&updateTags, "tag", "t", []string{}, "The new tags for the entry (replaces all)")
	updateCmd.Flags().
		StringSliceVarP(&updateHosts, "host", "H", []string{}, "The new hosts for the entry (replaces all)")

	updateCmd.RegisterFlagCompletionFunc("tag", completeTag)
	updateCmd.RegisterFlagCompletionFunc("host", completeHost)

	rootCmd.AddCommand(updateCmd)
}

func sameElems(s1 []string, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for _, v := range s1 {
		if !slices.Contains(s2, v) {
			return false
		}
	}
	for _, v := range s2 {
		if !slices.Contains(s1, v) {
			return false
		}
	}

	return true
}

func printChange(property string, old string, new string) {
	if old == "" {
		old = ui.Comment("<empty>")
	} else {
		old = ui.Old(old)
	}
	if new == "" {
		new = ui.Comment("<empty>")
	} else {
		new = ui.New(new)
	}

	fmt.Printf("  %s: %s -> %s\n", property, old, new)
}
