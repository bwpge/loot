package cmd

import "github.com/spf13/cobra"

func emptyNoFileCompletion(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
	return []string{}, cobra.ShellCompDirectiveNoFileComp
}

func idCompletion(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	values := []string{}
	s := loadLootFileNoErr()
	if s == nil {
		return values, cobra.ShellCompDirectiveNoFileComp
	}

	for k, v := range s.Data {
		desc := truncate(v.Value)
		if v.Comment != "" {
			desc += " (" + v.Comment + ")"
		}
		values = append(values, k+"\t"+desc)
	}

	return values, cobra.ShellCompDirectiveNoFileComp
}

func flagCompletion(
	cmd *cobra.Command,
	args []string,
	toComplete string,
) ([]string, cobra.ShellCompDirective) {
	values := []string{}
	s := loadLootFileNoErr()
	if s == nil {
		return values, cobra.ShellCompDirectiveNoFileComp
	}

	for k, v := range s.Flags {
		desc := truncate(v.Owner)
		if v.Host != "" {
			desc += " (" + v.Host + ")"
		}
		values = append(values, k+"\t"+desc)
	}

	return values, cobra.ShellCompDirectiveNoFileComp
}
