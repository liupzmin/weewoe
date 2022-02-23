package main

import (
	"fmt"

	"os"

	"github.com/liupzmin/weewoe/internal/color"
	"github.com/liupzmin/weewoe/internal/config"
	"github.com/liupzmin/weewoe/internal/ui"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func infoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Print configuration info",
		Long:  "Print configuration information",
		Run: func(cmd *cobra.Command, args []string) {
			printInfo()
		},
	}
}

func printInfo() {
	const fmat = "%-25s %s\n"

	printLogo(color.Cyan)
	printTuple(fmat, "Configuration", config.W2ConfigFile, color.Cyan)
	printTuple(fmat, "Logs", config.DefaultLogFile, color.Cyan)
	printTuple(fmat, "Screen Dumps", getScreenDumpDirForInfo(), color.Cyan)
}

func printLogo(c color.Paint) {
	for _, l := range ui.LogoSmall {
		fmt.Fprintln(out, color.Colorize(l, c))
	}
	fmt.Fprintln(out)
}

// getScreenDumpDirForInfo get default screen dump config dir or from config.W2ConfigFile configuration.
func getScreenDumpDirForInfo() string {
	if config.W2ConfigFile == "" {
		return config.W2DefaultScreenDumpDir
	}

	f, err := os.ReadFile(config.W2ConfigFile)
	if err != nil {
		log.Error().Err(err).Msgf("Reads k9s config file %v", err)
		return config.W2DefaultScreenDumpDir
	}

	var cfg config.Config
	if err := yaml.Unmarshal(f, &cfg); err != nil {
		log.Error().Err(err).Msgf("Unmarshal k9s config %v", err)
		return config.W2DefaultScreenDumpDir
	}
	return cfg.W2.GetScreenDumpDir()
}
