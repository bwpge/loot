package cmd

import (
	"fmt"
	"strings"

	"loot/internal/entry"

	"github.com/spf13/cobra"
)

var (
	showNoNewline bool
	showTags      []string
	showHosts     []string
	showSep       string
)

var showCmd = &cobra.Command{
	Use:   "show [filter]",
	Short: "Display one or more entry values",
	Run: func(cmd *cobra.Command, args []string) {
		s, _ := loadLootFile()
		filtered := s.Filter(entry.Filter{
			ID:    args,
			Tags:  showTags,
			Hosts: showHosts,
		})

		values := []string{}
		for _, e := range filtered {
			values = append(values, e.Value)
		}

		if showSep == "" {
			showSep = "\n"
		}
		fmt.Print(strings.Join(values, showSep))
		if !showNoNewline {
			fmt.Println("")
		}
	},
	ValidArgsFunction: idCompletion,
}

func init() {
	showCmd.Flags().
		BoolVarP(&showNoNewline, "no-newline", "n", false, "Do not display a trailing newline")
	showCmd.Flags().
		StringSliceVarP(&showTags, "tag", "t", []string{}, "Only display entries matching given tags")
	showCmd.Flags().
		StringSliceVarP(&showHosts, "host", "H", []string{}, "Only display entries matching given hosts")
	showCmd.Flags().
		StringVarP(&showSep, "separator", "s", "", "Separator used when displaying multiple values")
	rootCmd.AddCommand(showCmd)
}
