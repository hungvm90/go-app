/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/hungvm90/go-app/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

var cfgFile string
var appConfig internal.AppConfig

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-app",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.Getwd()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Stack().Err(err).Msgf("fail to read config file")
	}
	log.Info().Msgf("Using config file: %v", viper.ConfigFileUsed())
	if err := viper.Unmarshal(&appConfig); err != nil {
		log.Fatal().Stack().Err(err).Msgf("fail to parse config file")
	}
	log.Info().Msgf("Running app with config: %+v", appConfig.Version)
	setupLogger()
	log.Debug().Msgf("Running app with config: %+v", appConfig)
}

func setupLogger() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	var writers []io.Writer
	writers = append(writers, os.Stdout)
	if appConfig.LogToFile {
		file := &lumberjack.Logger{
			Filename:   "runtime.log",
			MaxSize:    50, // megabytes
			MaxBackups: 10,
			MaxAge:     30,   //days
			Compress:   true, // disabled by default
		}
		writers = append(writers, file)
	}
	mw := io.MultiWriter(writers...)
	logger := zerolog.New(mw).With().Timestamp().Caller().Logger().Level(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = time.RFC3339Nano
	if appConfig.Debug {
		logger = logger.Level(zerolog.DebugLevel)
	}
	log.Logger = logger
}
