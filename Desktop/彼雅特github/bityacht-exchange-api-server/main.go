package main

import (
	"bityacht-exchange-api-server/cmd/make"
	"bityacht-exchange-api-server/cmd/migrate"
	"bityacht-exchange-api-server/cmd/seed"
	"bityacht-exchange-api-server/cmd/server"
	"bityacht-exchange-api-server/configs"
	"bityacht-exchange-api-server/internal/pkg/logger"

	_ "bityacht-exchange-api-server/cmd/migrate/migrations"
	_ "bityacht-exchange-api-server/cmd/seed/seeders"
	_ "bityacht-exchange-api-server/docs"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Long: "This project is a layout for common api server.",
	Use:  "go-api-layout",
}

func init() {
	cobra.OnInitialize(configs.Init, logger.Init)

	rootCmd.PersistentFlags().StringVarP(&configs.ConfigFile, "config", "c", "", "config file (default is configs/config.yaml)")

	rootCmd.AddCommand(server.Cmd)
	rootCmd.AddCommand(migrate.Cmd)
	rootCmd.AddCommand(make.Cmd)
	rootCmd.AddCommand(seed.Cmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

// Usage: https://github.com/swaggo/swag/blob/master/README.md
// @title						BitYacht Exchange API Doc
// @version						0.1.0
// @description					This is a RESTful API documentation of BitYacht Exchange Backend.
// @host						localhost:9000
// @BasePath					/api/v1
// @schemes						http
// @externalDocs.description	BitYacht Exchange API External Doc
// @externalDocs.url			https://hackmd.io/@9iHsj5PNTbKIqyrL54-ffQ/B1-ktNCqh
// @query.collection.format 	multi

// @securityDefinitions.apikey  BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer abcde12345".
func main() {
	rootCmd.Execute()
}
