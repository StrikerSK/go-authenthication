/*
Copyright © 2022 Teh|Striker
*/
package cmd

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	userHandlers "github.com/strikersk/user-auth/src/handlers"
	userJWT "github.com/strikersk/user-auth/src/jwt"
	userRepository "github.com/strikersk/user-auth/src/repository"
	userServices "github.com/strikersk/user-auth/src/service"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "user-auth",
	Short: "A brief description of your application",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cookies, err := cmd.Flags().GetBool("cookies")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		myRouter := mux.NewRouter()

		jwtConfig := userJWT.NewConfigStruct()
		userRepo := userRepository.NewLocalUserRepository()
		userCache := userRepository.NewCacheConfig()
		userService := userServices.NewLocalUserRepository(&userRepo, userCache)

		userHandling := userHandlers.NewUserHandler(&userService)
		jwtHandling := userHandlers.NewJwtHandler(&userService, jwtConfig)

		if cookies {
			cookiesHandling := userHandlers.NewCookiesHandler("session_token", &userService)
			cookiesHandling.EnrichRouter(myRouter)
		}

		userHandling.EnrichRouter(myRouter)
		jwtHandling.EnrichRouter(myRouter)

		fmt.Println(http.ListenAndServe(":"+port, cors.AllowAll().Handler(myRouter)))
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.user-auth.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("port", "p", "5000", "Assign port value")
	rootCmd.Flags().BoolP("cookies", "c", false, "Enable cookies endpoint")
}
