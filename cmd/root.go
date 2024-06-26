package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/taylormonacelli/goldbug"
)

var (
	cfgFile      string
	verbose      bool
	logFormat    string
	sentinelPath string
)

var rootCmd = &cobra.Command{
	Use:   "mylime",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	var err error

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mylime.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose mode")
	err = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	if err != nil {
		slog.Error("error binding verbose flag", "error", err)
		os.Exit(1)
	}

	rootCmd.PersistentFlags().StringVar(&logFormat, "log-format", "", "json or text (default is text)")
	err = viper.BindPFlag("log-format", rootCmd.PersistentFlags().Lookup("log-format"))
	if err != nil {
		slog.Error("error binding log-format flag", "error", err)
		os.Exit(1)
	}

	rootCmd.PersistentFlags().StringVar(&sentinelPath, "sentinel", "", "Path to the sentinel file")
	err = viper.BindPFlag("sentinel", rootCmd.PersistentFlags().Lookup("sentinel"))
	if err != nil {
		slog.Error("error binding sentinel flag", "error", err)
		os.Exit(1)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mylime")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	logFormat = viper.GetString("log-format")
	verbose = viper.GetBool("verbose")
	sentinelPath = viper.GetString("sentinel")

	slog.Debug("using config file", "path", viper.ConfigFileUsed())
	slog.Debug("log-format", "value", logFormat)
	slog.Debug("log-format", "value", viper.GetString("log-format"))
	slog.Debug("sentinel", "value", sentinelPath)
	slog.Debug("sentinel", "value", viper.GetString("sentinel"))

	setupLogging()
}

func setupLogging() {
	if verbose || logFormat != "" {
		if logFormat == "json" {
			goldbug.SetDefaultLoggerJson(slog.LevelDebug)
		} else {
			goldbug.SetDefaultLoggerText(slog.LevelDebug)
		}

		slog.Debug("setup", "verbose", verbose)
	}
}
