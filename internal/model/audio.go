package model

type AudioFormatType string

const (
	AudioFormatTypeWav AudioFormatType = "wav"
	AudioFormatTypeM4a AudioFormatType = "m4a"
)

var audioProps = map[AudioFormatType]struct {
	FileExtension string
	MimeType      string
}{
	AudioFormatTypeWav: {FileExtension: ".wav", MimeType: "audio/wav"},
	// https://stackoverflow.com/questions/39885749/is-a-m4a-file-considered-as-of-mime-type-audio-m4a-or-audio-mp4
	AudioFormatTypeM4a: {FileExtension: ".m4a", MimeType: "audio/mp4"},
}

func (a AudioFormatType) String() string {
	return string(a)
}

func (a AudioFormatType) GetFileExtension() string {
	return audioProps[a].FileExtension
}

func (a AudioFormatType) GetMimeType() string {
	return audioProps[a].MimeType
}
