/*
Copyright © 2022 Teh|Striker
*/
package cmd

import (
	"github.com/spf13/cobra"
	authClient "github.com/strikersk/user-auth/src/client"
	userJWT "github.com/strikersk/user-auth/src/jwt"
	userRepository "github.com/strikersk/user-auth/src/repository"
	userGrpcServer "github.com/strikersk/user-auth/src/server"
	userServices "github.com/strikersk/user-auth/src/service"
)

// grpcCmd represents the grpc command
var grpcCmd = &cobra.Command{
	Use:   "grpc",
	Short: "Run gRPC base server for user authorization",
	Run: func(cmd *cobra.Command, args []string) {
		mode, err := cmd.Flags().GetString("mode")
		if err != nil {
			panic(err)
		}

		switch mode {
		case "server":
			userRepo := userRepository.NewLocalUserRepository()
			userCache := userRepository.NewCacheConfig()
			jwtConfig := userJWT.NewConfigStruct()

			userService := userServices.NewLocalUserRepository(&userRepo, userCache, jwtConfig)
			userGrpcServer.NewAuthorizationServer(&userService).RunServer()
		case "client":
			client := authClient.NewAuthorizationSample()
			client.RegisterUser()
		default:
			panic("Unknown mode")
		}

	},
}

func init() {
	rootCmd.AddCommand(grpcCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// grpcCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	grpcCmd.Flags().StringP("mode", "m", "server", "Create a server or a client")
}
