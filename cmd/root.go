package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var homeDir string
var workDir string
var libsDir string

var config struct {
	Libraries []string
}

var rootCmd = &cobra.Command{
	Use: "go-bump",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initWorkDir)

	rootCmd.PersistentFlags().StringVar(&workDir, "workdir", "", "working directory (default is $HOME/.go-bump)")
}

func initWorkDir() {
	var err error
	homeDir, err = homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if workDir == "" {
		workDir = filepath.Join(homeDir, ".go-bump")
	}
	cfgFile := filepath.Join(workDir, "config.yaml")
	libsDir = filepath.Join(workDir, "libraries")

	fmt.Println("Config file:", cfgFile)
	fmt.Println("Libraries directory:", libsDir)
	fmt.Println()

	viper.SetConfigFile(cfgFile)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
