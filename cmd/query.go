/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/VictoriaMetrics-Community/victorialogs-tool/internal"
	"github.com/spf13/cobra"
)

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query logs from victoriametrics",
	Long: `query logs from victoriametrics.
	The query source comes from the configuration file set by the 'vtools setcfg' command.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		tail, _ := cmd.Flags().GetBool("tail")
		if tail {
			// TODO leslie: tail the logs
			fmt.Println("tail the logs")
			return
		}

		// just query the logs from the victoriametrics
		list, err := internal.QueryLogs()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Result
		for _, res := range list {
			fmt.Println(res)
		}
	},
}

func init() {
	rootCmd.AddCommand(queryCmd)

	queryCmd.Flags().BoolP("tail", "t", false, "tail the logs")
}
