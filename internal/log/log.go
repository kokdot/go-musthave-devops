package log

import (
    "flag"
	"strconv"
	"os"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func llog() {
	log.Logger = log.With().Caller().Logger()
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	log.Logger = log.With().Caller().Logger()
	log.Info().Msg("hello world")
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    debug := flag.Bool("debug", false, "sets log level to debug")

    flag.Parse()

    // Default level for this example is info, unless debug flag is present
    zerolog.SetGlobalLevel(zerolog.InfoLevel)
    if *debug {
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    }

    log.Debug().Msg("--------------------------------------------------------------------")
    log.Debug().Msg("This message appears only when log level set to Debug")
    log.Info().Msg("This message appears when log level set to Debug or Info")
	log.Print("-----------------------------565656565---------------------------------------")
    if e := log.Debug(); e.Enabled() {
        // Compute log output only if enabled.
        value := "bar"
        e.Str("foo", value).Msg("some debug message")
    }
}