package audio

import (
	"bytes"
	"fmt"

	"github.com/rs/zerolog"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type AudioService struct {
	Log zerolog.Logger
}

// ConvertToMonoPCM converts audio bytes to raw PCM format
func (a AudioService) ConvertToMonoPCM(audioBytes []byte) ([]byte, error) {
	input := bytes.NewReader(audioBytes)
	output := &bytes.Buffer{}

	a.Log.Info().Msg("Converting audio to PCM")
	err := ffmpeg_go.Input("pipe:0").
		Output("pipe:1", ffmpeg_go.KwArgs{"ac": "1", "ar": "16000", "f": "s16le", "t": "30"}).
		WithInput(input).
		WithOutput(output).
		Run()
	if err != nil {
		a.Log.Error().Err(err).Msg("Failed to convert audio to PCM")
		return nil, fmt.Errorf("audio conversion to PCM failed: %w", err)
	}

	a.Log.Info().Msg("Successfully converted audio to PCM")
	return output.Bytes(), nil
}
