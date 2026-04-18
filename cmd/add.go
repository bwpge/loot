package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/bwpge/loot/internal/entry"
	"github.com/bwpge/loot/internal/ui"
	"github.com/spf13/cobra"
)

var (
	addTags         []string
	addHosts        []string
	addInputFiles   []string
	addInputLines   []string
	addComment      string
	addNoDetectType bool
)

var addCmd = &cobra.Command{
	Use:   "add value [value...]",
	Short: "Add a new entry to loot file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && len(addInputFiles) == 0 && len(addInputLines) == 0 {
			bail("at least one value or", ui.Cli("--input"), "or", ui.Cli("--lines"), "is required")
		}

		s, f := loadLootFile()
		doAdd := func(e entry.Entry) {
			_, found := s.Find(e.Value)
			var id, op string
			if found {
				var changed bool
				id, changed = s.Merge(e)
				if !changed {
					fmt.Println(ui.Comment("skipped duplicate: " + truncate(e.Value)))
					return
				}
				op = "merged"
			} else {
				id = s.Add(e)
				op = "added"
			}

			tags := ""
			if len(e.Tags) > 0 {
				tags = ui.Comment(listString(e.Tags))
			}

			fmt.Println(op, id, "->", truncate(e.Value), tags)
		}

		values := args
		for _, path := range addInputLines {
			lines, err := os.ReadFile(path)
			if err != nil {
				bail(err)
			}

			for line := range strings.SplitSeq(string(lines), "\n") {
				line := strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "#") {
					continue
				}
				values = append(values, line)
			}
		}

		for _, arg := range values {
			e := entry.Entry{Value: arg, Comment: addComment, Tags: addTags, Hosts: addHosts}
			if addNoDetectType {
				doAdd(e)
				continue
			}

			entries, detected := entry.DetectValues(e)
			if detected != "" {
				fmt.Println("detected format:", ui.Cli(detected))
			}
			for _, e := range entries {
				doAdd(e)
			}
		}

		for _, path := range addInputFiles {
			bytes, err := os.ReadFile(path)
			if err != nil {
				bail(err)
			}
			doAdd(entry.Entry{
				Value:   string(bytes),
				Comment: addComment,
				Tags:    addTags,
				Hosts:   addHosts,
			})
		}
		s.Save(f)
	},
	ValidArgsFunction: cobra.NoFileCompletions,
}

func init() {
	addCmd.Flags().
		BoolVarP(&addNoDetectType, "no-detect", "n", false, "Do not create additional entries by detecting common formats like user@domain")
	addCmd.Flags().
		StringSliceVarP(&addInputFiles, "input", "i", []string{}, "Add file contents as entry value (useful for e.g., ssh keys)")
	addCmd.Flags().
		StringSliceVarP(&addInputLines, "lines", "l", []string{}, "Add an entry per line in the given file (empty and # ignored)")
	addCmd.Flags().
		StringSliceVarP(&addTags, "tag", "t", []string{}, "Additional data to store with the entry (used for filtering)")
	addCmd.Flags().
		StringVarP(&addComment, "comment", "c", "", "Additional note to store with the entry")
	addCmd.Flags().
		StringSliceVarP(&addHosts, "host", "H", []string{}, "Host attribution for the entry")

	addCmd.RegisterFlagCompletionFunc("tag", completeTag)
	addCmd.RegisterFlagCompletionFunc("host", completeHost)

	rootCmd.AddCommand(addCmd)
}
