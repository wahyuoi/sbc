# How Audio Processing Works

Given the product requirements that users can only retrieve the audio file using the same format as the original file, 
the audio processing is done in a background process, and it can reduce the latency by >180ms. More details about this number can be found in the [Benchmark](#benchmark) section.

## Submit Audio Recording

When a user submits an audio file:

1. The original M4A audio file is stored immediately
2. A background process starts to:
   - Convert the M4A file to WAV format using FFmpeg
   - Store the converted WAV file

This asynchronous processing ensures:
- Fast response times for audio submissions
- No blocking while conversion happens


## Conversion Process

The conversion process is handled by a background process using FFmpeg. The process is as follows:

1. The original M4A file is downloaded from the storage (can be local file system or cloud storage), and stored in a temporary file
2. The temporary M4A file is converted to WAV format using FFmpeg, and stored in a temporary file
3. The temporary WAV file is stored to the storage (can be local file system or cloud storage)
4. Delete the temporary files

The reason for using FFmpeg is that it is a powerful and flexible tool for audio processing, and it is widely used and supported. 


### Benchmark

The conversion process was benchmarked using sample m4a audio files. The benchmark was done on:
```
goos: linux
goarch: amd64
pkg: github.com/wahyuoi/sbc/internal/service
cpu: AMD Ryzen 5 PRO 4650G with Radeon Graphics
audio sample rate: 44100
audio sample channel: 2 (stereo)
```

The benchmark results for 11 seconds audio file, ./sample/eleven.m4a:
```
BenchmarkAudioService_ConvertAudio-12    	       6	 189494051 ns/op	 2079780 B/op	      94 allocs/op
```

The benchmark results for 90 seconds audio file, ./sample/ninety.m4a:
```
BenchmarkAudioService_ConvertAudio-12    	       2	 510177491 ns/op	16047800 B/op	      97 allocs/op
```



