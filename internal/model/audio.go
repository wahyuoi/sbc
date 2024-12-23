package model

import "fmt"

type AudioFormatType string

const (
	AudioFormatTypeWav AudioFormatType = "wav"
	AudioFormatTypeM4a AudioFormatType = "m4a"
)

func GetAudioFormatType(fileExtOrFormat string) (AudioFormatType, error) {
	for k, v := range audioProps {
		if v.fileExtension == fileExtOrFormat || string(k) == fileExtOrFormat {
			return k, nil
		}
	}
	return "", fmt.Errorf("invalid audio format type")
}

var audioProps = map[AudioFormatType]struct {
	fileExtension string
	mimeType      string
	isForUpload   bool // true if the format is allowed for user to upload
	isForDownload bool // true if the format is allowed for user to download
	codec         string
	sampleRate    int
	channel       int
}{
	AudioFormatTypeWav: {
		fileExtension: ".wav",
		mimeType:      "audio/wav",
		isForUpload:   false,
		isForDownload: true,
		codec:         "pcm_s16le",
		sampleRate:    44100,
		channel:       2,
	},
	AudioFormatTypeM4a: {
		fileExtension: ".m4a",
		mimeType:      "audio/mp4", // https://stackoverflow.com/questions/39885749/is-a-m4a-file-considered-as-of-mime-type-audio-m4a-or-audio-mp4
		isForUpload:   true,
		isForDownload: true,
		codec:         "aac",
		sampleRate:    44100,
		channel:       2,
	},
}

func (a AudioFormatType) String() string {
	return string(a)
}

func (a AudioFormatType) GetFileExtension() string {
	return audioProps[a].fileExtension
}

func (a AudioFormatType) GetMimeType() string {
	return audioProps[a].mimeType
}

func (a AudioFormatType) IsForUpload() bool {
	return audioProps[a].isForUpload
}

func (a AudioFormatType) IsForDownload() bool {
	return audioProps[a].isForDownload
}

func (a AudioFormatType) GetCodec() string {
	return audioProps[a].codec
}

func (a AudioFormatType) GetSampleRate() int {
	return audioProps[a].sampleRate
}

func (a AudioFormatType) GetChannel() int {
	return audioProps[a].channel
}
