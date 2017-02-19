package cmd

import (
	"fmt"
	"os"

	"log"

	"github.com/hirondelle-app/listener/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// HirondelleCmd represents the base command when called without any subcommands
var HirondelleCmd = &cobra.Command{
	Use:   "hirondelle",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	AddCommands()

	if err := HirondelleCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

// AddCommands adds child commands to the root command HirondelleCmd.
func AddCommands() {
	HirondelleCmd.AddCommand(cleanCmd)
	HirondelleCmd.AddCommand(runCmd)
}

func init() {
	config.InitLog()
	cobra.OnInitialize(initConfig)

	HirondelleCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to the config file in yaml")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	err := viper.Unmarshal(&config.Cfg)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}
}
