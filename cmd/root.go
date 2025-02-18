/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	srv "github.com/cybrarymin/polkadot_exporter/internals/server"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "polkadot_exporter",
	Short: "polkadot blockchain binary prometheus exporter",
	Long: `This exporter is used to get connected to the polkadot binary api for collection
	of data and converting it prometheus metrics
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if srv.ShowVersion {
			fmt.Printf("Version: %s\n BuildTime:%s\n", srv.Version, srv.BuildTime)
			return
		}
		srv.Start()
	},

	PreRunE: func(cmd *cobra.Command, args []string) error {
		scheme, _, err := srv.ListenAddrParser(srv.ListenAddr)
		if err != nil {
			return err
		}
		cmdValidator.check(scheme == "https" || scheme == "http", "unsupported protocol scheme for listen address")
		if scheme == "https" {
			// check to see if certificate options are provided
			cmdValidator.check(srv.CertPath != "" || srv.CertKeyPath != "", "certiface and privateKey files are mandatory using HTTPS")
		}

		ok := cmdValidator.valid()
		if !ok {
			return errors.Join(cmdValidator.cmdErrors...)
		}
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.polkadot_exporter.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVar(&srv.LogLevel, "log-level", "info", "loglevel of the exporter. possible values are debug, info, warn, error, fatal, panic, and trace")
	rootCmd.Flags().BoolVar(&srv.ShowVersion, "version", false, "show the version and build time of the exporter")
	rootCmd.Flags().StringVar(&srv.ListenAddr, "listen-addr", "http://0.0.0.0:9100", "listen address for the exporter")
	rootCmd.Flags().StringVar(&srv.CertPath, "crt", "", "HTTPs certificate .pem file path")
	rootCmd.Flags().StringVar(&srv.CertKeyPath, "crt-key", "", "HTTPs key .pem file path")
}
