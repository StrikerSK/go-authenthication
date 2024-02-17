/*
Copyright © 2024 Kristian Hanus <kristianhanus@gmail.com>
*/
package cmd

import (
	"github.com/strikersk/user-auth/config"
	"github.com/strikersk/user-auth/constants"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "User authorization",
	Short: "Simple authorization application",
	Long: `Just a simple authorization application using:
- Cobra CLI & Viper
- Databases: InMemory, Postgres and SQLite
- Caches: InMemory, MemCache and Redis`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		configuration := config.ApplicationConfiguration{
			Application: config.Application{
				Port:        "8080",
				ContextPath: "/api",
			},
			Authorization: config.Authorization{
				AuthorizationType:   constants.Cookies,
				AuthorizationHeader: "Authorization",
				TokenEncodingType:   constants.JWT,
				TokenExpiration:     3600,
				JWT: config.JWTConfiguration{
					TokenEncoding: "Wow, much safe",
				},
				Encryption: config.EncryptionConfiguration{
					Cost: bcrypt.DefaultCost,
				},
			},
			Cache: config.CacheConfiguration{
				Name:       constants.InMemory,
				Expiration: 3600,
			},
			Database: config.DatabaseConfiguration{
				Type: constants.InMemory,
			},
		}

		err := viper.GetViper().Unmarshal(&configuration)
		if err != nil {
			log.Fatal("Configuration unmarshal error: ", err)
		}
		runServer(configuration)
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
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.user-auth.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	}
	//} else {
	//	// Find application directory.
	//	directory, err := os.Getwd() //get the current directory using the built-in function
	//	cobra.CheckErr(err)
	//
	//	// Search config in home directory with name ".user-auth" (without extension).
	//	viper.AddConfigPath(directory)
	//	viper.SetConfigName("/examples/config.yaml")
	//	viper.SetConfigType("yaml")
	//}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}
