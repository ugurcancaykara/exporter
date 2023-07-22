package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ugurcancaykara/exporter/utils"
)

var (
	envName string
	envs    string
)

var (
	rootCmd = &cobra.Command{
		Use:   "exporter",
		Args:  cobra.ArbitraryArgs,
		Short: "Export informations to CLI",
		Long: `If you have different export credentials like me and don't want to copy and paste it everytime 
and don't want to waste your time by searching them you need them, you can configure them to exporter your terminal.
		
Exporter is a CLI library for saving different kind of export values with different names.For example:
Adding a new environment:
'exporter set --name my-a-company-infra-creds --envs "key1=value1,key2=value2,key3=value3"'

Exporting an existing environment:
'exporter environmentname'`,
		Run: func(cmd *cobra.Command, args []string) {
			var envCfgs utils.EnvConfigs

			// listing existing environments
			if len(args) == 0 {
				if err := utils.ListEnvironments(); err != nil {
					fmt.Printf("Failed to list environments: %s\n", err)
					return
				}
				return
			}

			// if a environment name provided, it will be exported
			envName = args[0]

			// fetch environment configuration file path
			path, err := utils.GetConfigPath()
			if err != nil {
				fmt.Printf("Failed to get environment config: %s\n", err)
				return
			}

			// read configuration file
			bytes, err := ioutil.ReadFile(path)
			if err != nil && !os.IsNotExist(err) {
				fmt.Printf("there is no configuration file set an environment, it will automatically be created")
				return
			}

			if len(bytes) > 0 {
				if err := json.Unmarshal(bytes, &envCfgs); err != nil {
					return
				}
			} else {
				fmt.Printf("there is no configuration for %s\n", envName)
			}

			values, ok := envCfgs[envName]
			if !ok {
				fmt.Printf("there is no configuration for environment: %s\n", envName)
				return
			}

			for key, value := range values {
				fmt.Printf("export %s=%s\n", key, value)
			}
		},
	}
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set environment variables for a given name",
	Long:  `Set environment variables for a given name`,
	Run: func(cmd *cobra.Command, args []string) {
		envName, _ := cmd.Flags().GetString("name")
		kvPairs := strings.Split(envs, ",")
		kvMap := make(map[string]string, len(kvPairs))
		for _, pair := range kvPairs {
			kv := strings.Split(pair, "=")
			if len(kv) != 2 {
				fmt.Printf("Invalid key=value pair: %s\n", pair)
				return
			}
			kvMap[kv[0]] = kv[1]
		}

		if err := utils.SaveConfig(envName, kvMap); err != nil {
			fmt.Printf("failed to save environment config: %s\n", err)
			return
		}
		fmt.Printf("Successfully saved environment: %s\n", envName)
	},
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete environment variables for a given name",
	Long:  `Delete environment variables for a given name`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := utils.DeleteConfig(envName); err != nil {
			fmt.Printf("Failed to delete environment config: %s\n", err)
			return
		}
		fmt.Printf("Successfully deleted environment config: %s\n", envName)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update environment variables for a given name",
	Long:  `Update environment variables for a given name`,
	Run: func(cmd *cobra.Command, args []string) {
		kvPairs := strings.Split(envs, ",")
		kvMap := make(map[string]string, len(kvPairs))
		for _, pair := range kvPairs {
			kv := strings.Split(pair, "=")
			if len(kv) != 2 {
				fmt.Printf("Invalid key=value pair: %s\n", pair)
				return
			}
			kvMap[kv[0]] = kv[1]
		}

		if err := utils.UpdateConfig(envName, kvMap); err != nil {
			fmt.Printf("Failed to update environment config: %s\n", err)
			return
		}
		fmt.Printf("Successfully updated environment config: %s\n", envName)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(updateCmd)

	setCmd.Flags().StringVar(&envName, "name", "", "Environment name (required)")
	setCmd.Flags().StringVar(&envs, "envs", "", "Comma-seperated key=value pairs (required)")
	setCmd.MarkFlagRequired("name")
	setCmd.MarkFlagRequired("envs")

	deleteCmd.Flags().StringVar(&envName, "name", "", "Environment name (required)")
	deleteCmd.MarkFlagRequired("name")

	updateCmd.Flags().StringVar(&envName, "name", "", "Environment name (required)")
	updateCmd.Flags().StringVar(&envs, "envs", "", "Comma-separated key=value pairs (required)")
	updateCmd.MarkFlagRequired("name")
	updateCmd.MarkFlagRequired("envs")

	// list functionality provided by rootCmd

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
