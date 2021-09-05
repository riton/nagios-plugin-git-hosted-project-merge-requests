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
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/riton/nagios-plugin-git-hosted-project-merge-requests/nagios"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

type rootCmdFlags struct {
	Host                    string        `mapstructure:"host"`
	Debug                   bool          `mapstructure:"debug"`
	Timeout                 time.Duration `mapstructure:"timeout"`
	ConfigFile              string
	GitProvider             string        `mapstructure:"git-provider"`
	APIToken                string        `mapstructure:"api-token"`
	Project                 string        `mapstructure:"project"`
	TargetBranch            string        `mapstructure:"target-branch"`
	WarningLastUpdateDelay  time.Duration `mapstructure:"delay-warning-last-update"`
	CriticalLastUpdateDelay time.Duration `mapstructure:"delay-critical-last-update"`
}

var (
	// Default command flags
	cmdFlags = rootCmdFlags{
		Timeout: 30 * time.Second,
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nagios-plugin-git-hosted-project-merge-requests",
	Short: "Checks that a github / gitlab / gitea project has opened merge requests",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		nagios.ProbeCobraAdapter(cmd, args, nagiosConfigViperAdapter())
	},
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

	rootCmd.PersistentFlags().StringVarP(&cmdFlags.ConfigFile, "config", "c", "", "config file (default is /etc/nagios-plugin-git-hosted-project-merge-requests/config.yaml)")

	rootCmd.Flags().StringVarP(&cmdFlags.Host, "host", "H", "", "host to check (API endpoint)")
	rootCmd.MarkFlagRequired("host")

	rootCmd.Flags().StringVarP(&cmdFlags.Project, "project", "P", "", "project to check for opened MergeRequests")
	rootCmd.MarkFlagRequired("project")

	rootCmd.Flags().StringVarP(&cmdFlags.GitProvider, "git-provider", "p", "", fmt.Sprintf("git provider can be one of %s", strings.Join([]string{nagios.GitlabGitProvider}, ",")))

	rootCmd.PersistentFlags().DurationVarP(&cmdFlags.Timeout, "timeout", "t", 30*time.Second, "Global timeout")
	rootCmd.PersistentFlags().BoolVarP(&cmdFlags.Debug, "debug", "d", false, "Enable debug")

	rootCmd.Flags().StringVar(&cmdFlags.APIToken, "api-token", "", "API Token used for authentication")
	rootCmd.Flags().StringVar(&cmdFlags.TargetBranch, "target-branch", "master", "Only consider merge requests with this target-branch")

	rootCmd.Flags().DurationVar(&cmdFlags.WarningLastUpdateDelay, "warning-last-update", 6*time.Hour, "warning if last-update was that delay ago")
	rootCmd.Flags().DurationVar(&cmdFlags.CriticalLastUpdateDelay, "critical-last-update", 24*time.Hour, "critical if last-update was that delay ago")

	viper.BindPFlag("host", rootCmd.Flags().Lookup("host"))
	viper.BindPFlag("project", rootCmd.Flags().Lookup("project"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("api-token", rootCmd.Flags().Lookup("api-token"))
	viper.BindPFlag("git-provider", rootCmd.Flags().Lookup("git-provider"))
	viper.BindPFlag("target-branch", rootCmd.Flags().Lookup("target-branch"))
	viper.BindPFlag("warning-last-update", rootCmd.Flags().Lookup("warning-last-update"))
	viper.BindPFlag("critical-last-update", rootCmd.Flags().Lookup("critical-last-update"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log.SetOutput(os.Stderr)

	if cmdFlags.ConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cmdFlags.ConfigFile)
	} else {
		// Search config in home directory with name ".nagios-plugin-check_ssh_interactive_connect" (without extension).
		viper.AddConfigPath("/etc/nagios-plugin-git-hosted-project-merge-requests")
		viper.SetConfigName("config")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	if cmdFlags.Debug {
		log.SetLevel(log.DebugLevel)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Debug("Using config file:", viper.ConfigFileUsed())
	}
}

func nagiosConfigViperAdapter() nagios.ProbeConfig {
	return nagios.ProbeConfig{
		Timeout:                 viper.GetDuration("timeout"),
		APIEndpoint:             viper.GetString("host"),
		Project:                 viper.GetString("project"),
		Debug:                   viper.GetBool("debug"),
		APIToken:                viper.GetString("api-token"),
		GitProvider:             viper.GetString("git-provider"),
		TargetBranch:            viper.GetString("target-branch"),
		WarningLastUpdateDelay:  viper.GetDuration("warning-last-update"),
		CriticalLastUpdateDelay: viper.GetDuration("critical-last-update"),
	}
}
