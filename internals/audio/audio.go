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

// ConvertToMonoWAV converts audio bytes to mono WAV format
func (a AudioService) ConvertToMonoWAV(audioBytes []byte) ([]byte, error) {
	input := bytes.NewReader(audioBytes)
	output := &bytes.Buffer{}

	a.Log.Info().Msg("converting audio to mono WAV")

	err := ffmpeg_go.Input("pipe:0").
		Output("pipe:1", ffmpeg_go.KwArgs{"ac": "1", "ar": "16000", "f": "wav"}).
		WithInput(input).
		WithOutput(output).
		Run()

	if err != nil {
		a.Log.Error().Err(err).Msg("Failed to convert audio to mono WAV")
		return nil, fmt.Errorf("audio conversion failed: %w", err)
	}

	a.Log.Info().Msg("Successfully converted audio to mono WAV")
	return output.Bytes(), nil
}
