package utils

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var Logger zerolog.Logger

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	Logger = zerolog.New(os.Stderr).With().Caller().Timestamp().Stack().Logger()
}
