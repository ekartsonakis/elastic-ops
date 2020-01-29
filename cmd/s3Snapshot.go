/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"elastic-ops/cmd/s3Snapshot"
	"github.com/spf13/cobra"
)

var elasticUser string
var elasticPass string
var elasticServer string
var indices string
var s3repo string
// var includeGlobalState bool

// s3SnapshotCmd represents the s3Snapshot command
var s3SnapshotCmd = &cobra.Command{
	Use:   "s3Snapshot",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("s3Snapshot called")
		username, _:= cmd.Flags().GetString("username")
		fmt.Println(username)
		password, _:= cmd.Flags().GetString("password")
		fmt.Println(password)
		server, _:= cmd.Flags().GetString("server")
		fmt.Println(server)
		indices, _:= cmd.Flags().GetString("indices")
		fmt.Println(indices)
		globalState, _:= cmd.Flags().GetBool("include_global_state")
		fmt.Println(globalState)
		s3Repo, _:= cmd.Flags().GetString("s3repo")
		fmt.Println(s3Repo)
		s3snapshot.Snap(username, password, server, indices, s3Repo, globalState)
	},
}

func init() {
	rootCmd.AddCommand(s3SnapshotCmd)
	s3SnapshotCmd.Flags().StringVarP(&elasticUser, "username", "u", "", "Elasticsearch username")
	s3SnapshotCmd.MarkFlagRequired("username")
	s3SnapshotCmd.Flags().StringVarP(&elasticPass, "password", "p", "", "Elasticsearch password")
	s3SnapshotCmd.MarkFlagRequired("password")
	s3SnapshotCmd.Flags().StringVarP(&elasticServer, "server", "s", "", "IP of elasticsearch instance")
	s3SnapshotCmd.MarkFlagRequired("server")
	s3SnapshotCmd.Flags().StringVarP(&indices, "indices", "i", "", "Elastic indices")
	s3SnapshotCmd.MarkFlagRequired("indices")
	s3SnapshotCmd.Flags().StringVarP(&s3repo, "s3repo", "r", "", "S3 repository")
	s3SnapshotCmd.MarkFlagRequired("s3repo")
	s3SnapshotCmd.Flags().BoolP("include_global_state", "g", false, "Add Floating Numbers")
}
