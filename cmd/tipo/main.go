// Copyright (C) 2023 CGI France
//
// This file is part of TIPO.
//
// TIPO is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// TIPO is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with TIPO.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"os"

	"github.com/cgi-fr/tipo/pkg/swap"
	"github.com/cgi-fr/tipo/pkg/swap/jsonline"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

// Provisioned by ldflags.
var (
	name      string //nolint: gochecknoglobals
	version   string //nolint: gochecknoglobals
	commit    string //nolint: gochecknoglobals
	buildDate string //nolint: gochecknoglobals
	builtBy   string //nolint: gochecknoglobals
)

type pdef struct {
	configuration string // path to the configuration file
	logLevel      string // number of the loglevel
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}) //nolint:exhaustruct

	var definition pdef

	//nolint:exhaustruct
	rootCmd := &cobra.Command{
		Use:   "tipo",
		Short: "Command line to swap elements of a dataset",
		Version: fmt.Sprintf(`%v (commit=%v date=%v by=%v)
	Copyright (C) 2023 CGI France \n License GPLv3: GNU GPL version 3 <https://gnu.org/licenses/gpl.html>.
	This is free software: you are free to change and redistribute it.
	There is NO WARRANTY, to the extent permitted by law.`, version, commit, buildDate, builtBy),
		Run: func(cmd *cobra.Command, args []string) {
			initLog(definition)

			log.Info().Msgf("Version %v %v (commit=%v date=%v by=%v)", name, version, commit, buildDate, builtBy)
			run(definition)
		},
	}

	rootCmd.PersistentFlags().
		StringVarP(&definition.configuration, "configuration", "c",
			"swap.yml", "location of the configuration file")

	rootCmd.PersistentFlags().
		StringVarP(&definition.logLevel, "logLevel",
			"v", "warn", "value of the loglevel")

	if err := rootCmd.Execute(); err != nil {
		log.Err(err).AnErr("error", err).Msg("Error when executing command")
		os.Exit(1)
	}
}

func run(definition pdef) {
	log.Info().
		Str("config", definition.configuration).
		Msg("Start Swap")

	log.Info().Msg("Reading the configuration file")

	// read configuration file
	configuration, err := swap.LoadConfigurationFromYAML(definition.configuration)
	if err != nil {
		panic(err)
	}

	driver := configuration.BuildDriver()
	reader := jsonline.NewReader(os.Stdin)
	collector := jsonline.NewCollector(os.Stdout)

	log.Info().Msg("Starting process")

	driver.Run(reader, collector)
}

func initLog(definition pdef) {
	switch definition.logLevel {
	case "trace", "5":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
		log.Info().Msg("Logger level set to trace")
	case "debug", "4":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Info().Msg("Logger level set to debug")
	case "info", "3":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Info().Msg("Logger level set to info")
	case "warn", "2":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error", "1":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}
}
