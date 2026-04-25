package cmd

import (
	"fmt"

	"github.com/bwpge/loot/internal/state"
	"github.com/bwpge/loot/internal/ui"

	"github.com/spf13/cobra"
)

var importYes bool

var importCmd = &cobra.Command{
	Use:   "import file [file...]",
	Short: "Import entries and flags from one or more loot files",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		s, f := loadLootFile()

		sources := make([]*state.State, len(args))
		var total int
		for i, path := range args {
			src, err := state.Load(path)
			if err != nil {
				bail(err)
			}
			sources[i] = src
			total += len(src.Data) + len(src.Flags)
		}

		if total == 0 {
			warn("nothing to import")
			return
		}

		if !importYes &&
			!confirm(fmt.Sprintf("import %d items from %d file(s)?", total, len(args))) {
			return
		}

		var added, merged, skipped, flagsAdded, flagsSkipped int
		for i, path := range args {
			src := sources[i]

			fmt.Println("importing from", path)

			for _, e := range src.Data {
				_, found := s.Find(e.Value)
				if found {
					id, changed := s.Merge(e)
					if changed {
						fmt.Println("merged", id, "->", truncate(e.Value))
						merged++
					} else {
						fmt.Println(ui.Comment("skipped duplicate: " + truncate(e.Value)))
						skipped++
					}
				} else {
					id := s.Add(e)
					fmt.Println("added", id, "->", truncate(e.Value))
					added++
				}
			}

			for k, flag := range src.Flags {
				if _, exists := s.Flags[k]; exists {
					fmt.Println(ui.Comment("skipped duplicate flag: " + k))
					flagsSkipped++
					continue
				}
				s.Capture(k, flag)
				fmt.Println("added flag", k)
				flagsAdded++
			}
		}

		s.Save(f)
		fmt.Printf(
			"\nimported: %d added, %d merged, %d skipped (%d flags added, %d flags skipped)\n",
			added, merged, skipped, flagsAdded, flagsSkipped,
		)
	},
}

func init() {
	importCmd.Flags().BoolVarP(&importYes, "yes", "y", false, "Skip confirmation prompts")
	rootCmd.AddCommand(importCmd)
}
