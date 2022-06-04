package sound

import (
	"os"

	"github.com/skycoin/cx-game/lib/openal"
	"github.com/mjibson/go-dsp/wav"
)

//get wav information such as number of channels, bitrate and bit depth
func getWavInfo(filename string) (*wav.Header, error) {
	r, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	wavv, err := wav.New(r)
	if err != nil {
		return nil, err
	}

	return &wavv.Header, nil
}

//only stereo and mono format, bit depth of 8 or 16 depth supported
func determineFormat(wavInfo *wav.Header) openal.Format {
	printDebug(wavInfo.NumChannels, wavInfo.BitsPerSample)
	if wavInfo.NumChannels == 1 {
		if wavInfo.BitsPerSample == 8 {
			return openal.FormatMono8
		} else {
			return openal.FormatMono16
		}
	} else if wavInfo.NumChannels == 2 {
		if wavInfo.BitsPerSample == 8 {
			return openal.FormatStereo8
		} else {
			return openal.FormatStereo16
		}
	}
	printDebug("Bad WAV format!")
	return openal.FormatMono16
}
