// Copyright 2021 l1b0k
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:   "cs",
	Short: "cs",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {

	},
}

var (
	cfgPath     string
	cfgFilePath string
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	_ = flag.Set("v", "0")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgPath, "configPath", "", "config file (default is $HOME/.cs)")
	RootCmd.PersistentFlags().StringVar(&cfgFilePath, "config", "", "config file (default is $HOME/.cs/config)")
	RootCmd.PersistentFlags().String("region", "cn-hangzhou", "specific the region for api process")
	RootCmd.PersistentFlags().Bool("private", false, "specific the api to use VPC endpoint or public endpoint.Default is public access.This also affect kubeconfig access type.")

	home, err := homedir.Dir()
	cobra.CheckErr(err)

	_ = viper.BindPFlag("kubeconfig", pflag.CommandLine.Lookup("kubeconfig"))
	viper.SetDefault("kubeconfig", filepath.Join(home, ".kube/config"))

	_ = viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))

	_ = viper.BindEnv("region", "ALIBABA_CLOUD_REGION_ID")
	_ = viper.BindPFlag("region", RootCmd.PersistentFlags().Lookup("region"))
	_ = viper.BindPFlag("private", RootCmd.PersistentFlags().Lookup("private"))
}

func initConfig() {
	if cfgFilePath != "" {
		viper.SetConfigFile(cfgFilePath)
	} else {
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		viper.AddConfigPath(filepath.Join(home, ".cs"))
		viper.SetConfigType("yaml")
		viper.SetConfigName(".cs")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
