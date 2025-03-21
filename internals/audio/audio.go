package audio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"math/cmplx"

	"github.com/rs/zerolog"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"gonum.org/v1/gonum/dsp/fourier"
)

type AudioService struct {
	Log zerolog.Logger
}

// ConvertToMonoPCM converts audio bytes to raw PCM format
func (a AudioService) ConvertToMonoPCM(audioBytes []byte) ([]float64, error) {
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
	return a.ToFloat64(output.Bytes()), nil
}

// ToFloat64 converts 16-bit PCM bytes to a slice of float64 values
func (a AudioService) ToFloat64(pcmData []byte) []float64 {
	numSamples := len(pcmData) / 2
	Samples := make([]float64, numSamples)

	for i := 0; i < numSamples; i++ {
		sample := int16(binary.LittleEndian.Uint16(pcmData[i*2 : i*2+2])) // Read 16-bit PCM
		Samples[i] = float64(sample) / math.MaxInt16
	}

	return Samples
}

// ApplyFFT converts PCM samples to frequency magnitudes using FFT(Fast Fourier Transform)
func (a AudioService) ApplyFFT(samples []float64) []float64 {
	fft := fourier.NewFFT(len(samples))
	fftData := fft.Coefficients(nil, samples)

	magnitudes := make([]float64, len(fftData))
	for i := range fftData {
		magnitudes[i] = cmplx.Abs(fftData[i])
	}

	return magnitudes
}
