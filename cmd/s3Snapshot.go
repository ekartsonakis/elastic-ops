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
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var elasticUser string
var elasticPass string
var elasticServer string
var indices string
var s3repo string
var snap string

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
		username, _ := cmd.Flags().GetString("username")
		fmt.Println(username)
		password, _ := cmd.Flags().GetString("password")
		fmt.Println(password)
		server, _ := cmd.Flags().GetString("server")
		fmt.Println(server)
		s3Repo, _ := cmd.Flags().GetString("s3repo")
		fmt.Println(s3Repo)
		onlyList, _ := cmd.Flags().GetBool("list")
		fmt.Println(onlyList)
		deleteSnap, _ := cmd.Flags().GetBool("delete")
		fmt.Println(deleteSnap)
		indices, _ := cmd.Flags().GetString("indices")

		if !onlyList {
			fmt.Println(indices)
			if indices == "" {
				log.Fatal("Indices to snapshot should be given")
			}
		}
		globalState, _ := cmd.Flags().GetBool("include_global_state")
		fmt.Println(globalState)

		snapName, _ := cmd.Flags().GetString("snap_name")
		fmt.Println(snapName)

		s3snapshot.Snap(username, password, server, s3Repo, indices, snapName, onlyList, globalState, deleteSnap)
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
	s3SnapshotCmd.Flags().StringVarP(&s3repo, "s3repo", "r", "", "S3 repository")
	s3SnapshotCmd.MarkFlagRequired("s3repo")

	s3SnapshotCmd.Flags().BoolP("list", "l", false, "Only list available snapshots. If given no new snapshot will be taken.")
	s3SnapshotCmd.Flags().BoolP("delete", "d", false, "Delete a specific snapshot. If given no new snapshot will be taken.")

	s3SnapshotCmd.Flags().StringVarP(&indices, "indices", "i", "", "Elastic indices. If special character is included put the string in double quotes.")
	s3SnapshotCmd.Flags().StringVarP(&snap, "snap_name", "n", "%3Csnapshot-%7Bnow%2Fd%7D%3E", "New snapshot name. If not given then snapshot-Y.M.D will be given.")
	s3SnapshotCmd.Flags().BoolP("include_global_state", "g", false, "Add Floating Numbers")
}
