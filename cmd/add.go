package cmd

import (
	"fmt"

	"github.com/selcux/embarkdiff/diff"
	"github.com/spf13/cobra"
)

var source, target string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds the defined resource folder.",
	Long:  `Adds the defined resource folder which can be either source or target.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if source == "" && target == "" {
			fmt.Println(cmd.Flags().FlagUsages())
			return nil
		}

		err := addSource()
		if err != nil {
			return err
		}

		err = addTarget()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().StringVarP(&source, "source", "", "", "Source directory path")
	addCmd.Flags().StringVarP(&target, "target", "", "", "Target directory path")
}

//addSource to the resource file
func addSource() error {
	if source == "" {
		return nil
	}

	res := diff.NewResource()
	err := res.Load()
	if err != nil {
		return err
	}

	err = res.SetSource(source)
	if err != nil {
		return err
	}

	err = res.Write()
	if err != nil {
		return err
	}

	return nil
}

//addTarget to the resource file
func addTarget() error {
	if target == "" {
		return nil
	}

	res := diff.NewResource()
	err := res.Load()
	if err != nil {
		return err
	}

	err = res.SetTarget(target)
	if err != nil {
		return err
	}

	err = res.Write()
	if err != nil {
		return err
	}

	return nil
}
