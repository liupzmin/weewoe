package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/debug"

	"github.com/liupzmin/weewoe/internal/client"
	"github.com/liupzmin/weewoe/internal/color"
	"github.com/liupzmin/weewoe/internal/config"
	"github.com/liupzmin/weewoe/internal/view"
	"github.com/mattn/go-colorable"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

const (
	appName      = "w2ctl"
	shortAppDesc = "A graphical CLI for your distributed processes management."
	longAppDesc  = "w2ctl is a CLI to view and manage your distributed processes."
)

var (
	version, commit, date = "dev", "dev", client.NA
	w2Flags               *config.Flags

	rootCmd = &cobra.Command{
		Use:   appName,
		Short: shortAppDesc,
		Long:  longAppDesc,
		Run:   run,
	}

	out = colorable.NewColorableStdout()
)

func init() {
	rootCmd.AddCommand(versionCmd(), infoCmd())
	initW2Flags()
}

// Execute root command.
func Execute() {
	go func() {
		_ = http.ListenAndServe("0.0.0.0:6060", nil)
	}()
	if err := rootCmd.Execute(); err != nil {
		log.Panic().Err(err)
	}
}

func run(cmd *cobra.Command, args []string) {
	config.EnsurePath(*w2Flags.LogFile, config.DefaultDirMod)
	mod := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	file, err := os.OpenFile(*w2Flags.LogFile, mod, config.DefaultFileMod)
	if err != nil {
		panic(err)
	}
	defer func() {
		if file != nil {
			_ = file.Close()
		}
	}()
	defer func() {
		if err := recover(); err != nil {
			log.Error().Msgf("Boom! %v", err)
			log.Error().Msg(string(debug.Stack()))
			printLogo(color.Red)
			fmt.Printf("%s", color.Colorize("Boom!! ", color.Red))
			fmt.Println(color.Colorize(fmt.Sprintf("%v.", err), color.LightGray))
		}
	}()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: file})

	zerolog.SetGlobalLevel(parseLevel(*w2Flags.LogLevel))
	app := view.NewApp(loadConfiguration())
	if err := app.Init(version, *w2Flags.RefreshRate); err != nil {
		panic(fmt.Sprintf("app init failed -- %v", err))
	}
	if err := app.Run(); err != nil {
		panic(fmt.Sprintf("app run failed %v", err))
	}
	if view.ExitStatus != "" {
		panic(fmt.Sprintf("view exit status %s", view.ExitStatus))
	}
}

func loadConfiguration() *config.Config {
	log.Info().Msg("🐶 W2 starting up...")

	// Load W2 config file...
	w2Cfg := config.NewConfig()

	if err := w2Cfg.Load(config.W2ConfigFile); err != nil {
		log.Warn().Msg("Unable to locate W2 config. Generating new configuration...")
	}

	if *w2Flags.RefreshRate != config.DefaultRefreshRate {
		w2Cfg.W2.OverrideRefreshRate(*w2Flags.RefreshRate)
	}

	w2Cfg.W2.OverrideHost(*w2Flags.Host)
	w2Cfg.W2.OverridePort(*w2Flags.Port)
	w2Cfg.W2.OverrideHeadless(*w2Flags.Headless)
	w2Cfg.W2.OverrideLogoless(*w2Flags.Logoless)
	w2Cfg.W2.OverrideCrumbsless(*w2Flags.Crumbsless)
	w2Cfg.W2.OverrideReadOnly(*w2Flags.ReadOnly)
	w2Cfg.W2.OverrideWrite(*w2Flags.Write)
	w2Cfg.W2.OverrideCommand(*w2Flags.Command)
	w2Cfg.W2.OverrideScreenDumpDir(*w2Flags.ScreenDumpDir)

	if err := w2Cfg.Refine(); err != nil {
		log.Error().Err(err).Msgf("refine failed")
	}

	if err := w2Cfg.Save(); err != nil {
		log.Error().Err(err).Msg("Config save")
	}

	return w2Cfg
}

func parseLevel(level string) zerolog.Level {
	switch level {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}

func initW2Flags() {
	w2Flags = config.NewFlags()
	rootCmd.Flags().StringVarP(
		w2Flags.Host,
		"server", "s",
		config.DefaultHost,
		"Specify the server to connect",
	)
	rootCmd.Flags().StringVarP(
		w2Flags.Port,
		"port", "p",
		config.DefaultPort,
		"Specify the server port",
	)
	/*rootCmd.Flags().IntVarP(
		w2Flags.RefreshRate,
		"refresh", "r",
		config.DefaultRefreshRate,
		"Specify the default refresh rate as an integer (sec)",
	)*/
	rootCmd.Flags().StringVarP(
		w2Flags.LogLevel,
		"logLevel", "l",
		config.DefaultLogLevel,
		"Specify a log level (info, warn, debug, trace, error)",
	)
	rootCmd.Flags().StringVarP(
		w2Flags.LogFile,
		"logFile", "",
		config.DefaultLogFile,
		"Specify the log file",
	)
	rootCmd.Flags().BoolVar(
		w2Flags.Headless,
		"headless",
		false,
		"Turn W2 header off",
	)
	rootCmd.Flags().BoolVar(
		w2Flags.Logoless,
		"logoless",
		false,
		"Turn W2 logo off",
	)
	rootCmd.Flags().BoolVar(
		w2Flags.Crumbsless,
		"crumbsless",
		false,
		"Turn W2 crumbs off",
	)
	rootCmd.Flags().BoolVarP(
		w2Flags.AllNamespaces,
		"all-namespaces", "A",
		false,
		"Launch W2 in all namespaces",
	)
	/*rootCmd.Flags().StringVarP(
		w2Flags.Command,
		"command", "c",
		config.DefaultCommand,
		"Overrides the default resource to load when the application launches",
	)*/
	rootCmd.Flags().BoolVar(
		w2Flags.ReadOnly,
		"readonly",
		false,
		"Sets readOnly mode by overriding readOnly configuration setting",
	)
	rootCmd.Flags().BoolVar(
		w2Flags.Write,
		"write",
		false,
		"Sets write mode by overriding the readOnly configuration setting",
	)
	rootCmd.Flags().StringVar(
		w2Flags.ScreenDumpDir,
		"screen-dump-dir",
		"",
		"Sets a path to a dir for a screen dumps",
	)
	rootCmd.Flags()
}
