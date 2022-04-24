/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"os"

	"github.com/reliefeffortslk/protest-tracker-api/pkg/api"
	"github.com/reliefeffortslk/protest-tracker-api/pkg/bot"
	"github.com/reliefeffortslk/protest-tracker-api/pkg/cron"
	"github.com/spf13/cobra"
	"github.com/thejerf/suture/v4"
)

// rootCmd represents the base command when called without any subcommands

var (
	disableBot  = os.Getenv("DISABLE_BOTS") != ""
	disableAPI  = os.Getenv("DISABLE_API") != ""
	disableCron = os.Getenv("DISABLE_CRON") != ""
)

var rootCmd = &cobra.Command{
	Use:   "protest-tracker",
	Short: "telegram bot for protest-tracker",
	Long:  "telegram bot for protest-tracker",
	Run: func(cmd *cobra.Command, args []string) {
		supervisor := suture.NewSimple("Supervisor")

		if !disableBot {
			botService := new(bot.BotService)
			supervisor.Add(botService)
		}

		if !disableAPI {
			apiService := new(api.APIService)
			supervisor.Add(apiService)
		}

		if !disableCron {
			cronService := new(cron.CronService)
			supervisor.Add(cronService)
		}

		ctx, cancel := context.WithCancel(context.Background())
		if err := supervisor.Serve(ctx); err != nil {
			panic(err)
		}

		cancel()
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
