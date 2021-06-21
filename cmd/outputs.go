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
	"aspect.build/cli/bazel"
	"aspect.build/cli/proto/analysis"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

// outputsCmd represents the outputs command
var outputsCmd = &cobra.Command{
	Use:   "outputs [target]",
	Short: "Show the output locations of files produced by an action",
	Long:  ``,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// FIXME: check exit code, like if target does not exist
		reader := bazel.CaptureStdout("aquery", "--output=proto", args[0])
		bytes, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatalln("can't read proto file")
		}
		a := &analysis.ActionGraphContainer{}

		if err := proto.Unmarshal(bytes, a); err != nil {
			log.Fatalln("Failed to parse response:", err)
		}
		log.Printf("Parsed ActionGraph %v", a)
	},
}

func init() {
	rootCmd.AddCommand(outputsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// outputsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// outputsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
