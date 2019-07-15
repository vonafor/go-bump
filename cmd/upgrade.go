package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"

	"github.com/vonafor/go-bump/tools"
)

var dependency string
var version string

var upgradeCmd = &cobra.Command{
	Use: "upgrade",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !funk.ContainsString(config.Libraries, dependency) {
			return fmt.Errorf("dependency shoud be from libraries list")
		}

		var libraries []*tools.Library
		for _, name := range config.Libraries {
			libraries = append(libraries, tools.NewLibrary(libsDir, name))
		}

		if err := os.MkdirAll(libsDir, 0700); err != nil {
			return err
		}

		for _, library := range libraries {
			if err := library.Prepare(); err != nil {
				return err
			}
		}

		var mrs []string
		for _, library := range libraries {
			mr, err := library.UpdateDependency(dependency, version)
			if err != nil {
				return err
			}
			if mr != "" {
				mrs = append(mrs, mr)
			}
		}

		for _, mr := range mrs {
			fmt.Println("MR:", mr)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)

	upgradeCmd.Flags().StringVarP(&dependency, "dependency", "d", "", "upgraded dependency")
	upgradeCmd.Flags().StringVarP(&version, "version", "v", "", "new version")

	err := upgradeCmd.MarkFlagRequired("dependency")
	if err != nil {
		panic(err)
	}
	err = upgradeCmd.MarkFlagRequired("version")
	if err != nil {
		panic(err)
	}
}
