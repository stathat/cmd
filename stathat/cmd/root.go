// Copyright © 2016 Numerotron Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/stathat/cmd/stathat/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var debug string
var posthost string
var host string

// this example not behaving in Long:
// stathat post value loadavg `uptime | cut -d " " -f 12`

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "stathat",
	Short: "stathat is a command line interface to StatHat (www.stathat.com)",
	Long: `stathat is a command line interface to StatHat (www.stathat.com).
This application is a tool that allows efficient, powerful use
of the StatHat service.

For example, you can post stat data:

    stathat post value "load average" 1.3
    stathat post count "script run count" 1

Download a dataset:

    stathat dataset --tf 1w3h JfFv

`,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.stathat/config.yaml)")
	RootCmd.PersistentFlags().StringVar(&debug, "debug", "", "debug flags (comma separated)")
	RootCmd.PersistentFlags().StringVar(&posthost, "posthost", "https://api.stathat.com", "specify an api post host")
	RootCmd.PersistentFlags().StringVar(&host, "host", "https://www.stathat.com", "specify a host")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("config")         // name of config file (without extension)
	viper.AddConfigPath(".")              // adding current directory as first search path
	viper.AddConfigPath("$HOME/.stathat") // adding home directory as first search path
	viper.SetEnvPrefix("stathat")
	viper.AutomaticEnv() // read in environment variables that match

	viper.BindPFlag("host", RootCmd.Flags().Lookup("host"))
	viper.BindPFlag("posthost", RootCmd.Flags().Lookup("posthost"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error in config file %q: %s\n", viper.ConfigFileUsed(), err)
	}

	names := strings.Split(debug, ",")
	for _, name := range names {
		config.SetDebug(name, true)
	}
}
