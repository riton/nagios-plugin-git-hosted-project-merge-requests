/*
Copyright © 2021 Remi Ferrand

Contributor(s): Remi Ferrand <riton.github_at_gmail(dot)com>, 2021

This software is a computer program whose purpose is to [describe
functionalities and technical features of your software].

This software is governed by the CeCILL-B license under French law and
abiding by the rules of distribution of free software.  You can  use,
modify and/ or redistribute the software under the terms of the CeCILL-B
license as circulated by CEA, CNRS and INRIA at the following URL
"http://www.cecill.info".

As a counterpart to the access to the source code and  rights to copy,
modify and redistribute granted by the license, users are provided only
with a limited warranty  and the software's author,  the holder of the
economic rights,  and the successive licensors  have only  limited
liability.

In this respect, the user's attention is drawn to the risks associated
with loading,  using,  modifying and/or developing or reproducing the
software by the user in light of its specific status of free software,
that may mean  that it is complicated to manipulate,  and  that  also
therefore means  that it is reserved for developers  and  experienced
professionals having in-depth computer knowledge. Users are therefore
encouraged to load and test the software's suitability as regards their
requirements in conditions enabling the security of their systems and/or
data to be ensured and,  more generally, to use and operate it in the
same conditions as regards security.

The fact that you are presently reading this means that you have had
knowledge of the CeCILL-B license and that you accept its terms.

*/
package cmd

import (
	"context"
	"os"
	"time"

	"github.com/riton/nagios-plugin-git-hosted-project-merge-requests/nagios"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "nagios-plugin-git-hosted-project-merge-requests",
	Short:         "Checks that a github / gitlab / gitea project has opened merge requests",
	Long:          ``,
	Run:           nagios.ProbeCobraAdapter,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ctx context.Context) {
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		panic(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/nagios-plugin-git-hosted-project-merge-requests/config.yaml)")

	rootCmd.Flags().StringP("host", "H", "", "host to check")
	rootCmd.MarkFlagRequired("host")

	rootCmd.PersistentFlags().DurationP("timeout", "t", 30*time.Second, "Global timeout")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("debug", "d", false, "Enable debug")

	viper.BindPFlag("host", rootCmd.Flags().Lookup("host"))
	//viper.BindPFlag("ssh.connect_timeout", rootCmd.PersistentFlags().Lookup("connect-timeout"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log.SetOutput(os.Stderr)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".nagios-plugin-check_ssh_interactive_connect" (without extension).
		viper.AddConfigPath("/etc/nagios-plugin-git-hosted-project-merge-requests")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	debug, _ := rootCmd.Flags().GetBool("debug")
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Using config file:", viper.ConfigFileUsed())
	}
}
