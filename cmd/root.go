package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	userEmail   string
	userPass    string
	assistantID string
	machineType string
	userContext string
	debug       bool
	quietMode   bool
)

var rootCmd = &cobra.Command{
	Use:   "configgen",
	Short: "Configuration Generator using AI.YOU",
	Long: `A configuration generator tool that uses AI.YOU to create 
           YAML configurations for various types of IT infrastructure.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.configgen.yaml)")
	rootCmd.PersistentFlags().StringVar(&userEmail, "email", "", "AI.YOU email")
	rootCmd.PersistentFlags().StringVar(&userPass, "password", "", "AI.YOU password")
	rootCmd.PersistentFlags().StringVar(&assistantID, "assistant", "", "AI.YOU assistant ID")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug mode")
	rootCmd.PersistentFlags().BoolVar(&quietMode, "quiet", false, "Disable status messages")
}