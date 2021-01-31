package cmd

import (
	"fmt"

	"github.com/selcux/embarkdiff/diff"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the given folders.",
	Long:  `Lists the given source and target folders.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		res := diff.NewResource()
		err := res.Load()
		if err != nil {
			return err
		}

		fmt.Printf("Source:\t%s\n", res.Source())
		fmt.Printf("Target:\t%s\n", res.Target())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
