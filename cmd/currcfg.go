package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

// currcfgCmd represents the currcfg command
var currcfgCmd = &cobra.Command{
	Use:   "currcfg",
	Short: "Get the absolute path of the current configuration file.",
	Long:  "Print the absolute path of the current configuration file to the standard output.",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		defer func() {
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}()
		// open the configuration file
		user, err := user.Current()
		if err != nil {
			return
		}
		homeDir := user.HomeDir
		p := filepath.Join(homeDir, "vtool.json")
		byt, err := os.ReadFile(p)
		if err != nil && os.IsNotExist(err) {
			fmt.Printf("Configuration file is empty.\nPlease run 'vtool setcfg <path>' to set the configuration file path.\n")
			err = nil
			return
		} else if err != nil {
			return
		}

		var base struct {
			Base string `json:"base"`
		}
		err = json.Unmarshal(byt, &base)
		if err != nil {
			return
		}

		fmt.Println("Path:", base.Base)
	},
}

func init() {
	rootCmd.AddCommand(currcfgCmd)
}
