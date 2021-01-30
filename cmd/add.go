/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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

	addCmd.Flags().StringVarP(&source, "source", "", "", "Source directory to read from.")
	addCmd.Flags().StringVarP(&target, "target", "", "", "Target directory to read from.")
}

func addSource() error {
	if source == "" {
		return nil
	}

	res, err := diff.Read()
	if err != nil {
		return err
	}

	res.Source = source
	err = res.Write()
	if err != nil {
		return err
	}

	return nil
}

func addTarget() error {
	if target == "" {
		return nil
	}

	res, err := diff.Read()
	if err != nil {
		return err
	}

	res.Target = target
	err = res.Write()
	if err != nil {
		return err
	}

	return nil
}
