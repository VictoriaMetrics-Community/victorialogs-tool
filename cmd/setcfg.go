/*
Copyright © 2024 here-Leslie-Lau <i.leslie.lau@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// setcfgCmd represents the setcfg command
var setcfgCmd = &cobra.Command{
	Use:   "setcfg [toml file]",
	Short: "Set up the configuration file for query logs",
	Long: `Set up the configuration file for query logs.
	Accept a TOML file, format reference cfgs/example.toml.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 || args[0] == "" || !strings.Contains(args[0], ".toml") {
			_ = cmd.Usage()
			return
		}

		// save toml to base.json
		if err := saveTomlToBase(args[0]); err != nil {
			fmt.Println("Error: ", err)
			return
		}

		fmt.Println("Set up the configuration file successfully.")
	},
}

func saveTomlToBase(filePath string) error {
	// generated configuration file path
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	homeDir := user.HomeDir
	p := filepath.Join(homeDir, "vtool.json")

	var base struct {
		Base string `json:"base"`
	}
	base.Base = filePath

	file, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	byt, _ := json.Marshal(base)
	_, err = file.WriteString(string(byt))
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(setcfgCmd)
}
