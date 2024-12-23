# How Audio Processing Works

Given the product requirements that users can only retrieve the audio file using the same format as the original file, 
the audio processing is done in a background process, and it can reduce the latency for API to submit audio. 

### Submit Audio Recording

When a user submits an audio file:

1. The original M4A audio file is stored immediately
2. A background process starts to:
   - Convert the M4A file to WAV format using FFmpeg
   - Store the converted WAV file

This asynchronous processing ensures:
- Fast response times for audio submissions
- No blocking while conversion happens


### Conversion Process

The conversion process is handled by a background process using FFmpeg. The process is as follows:

1. The original M4A file is downloaded from the storage (can be local file system or cloud storage), and stored in a temporary file
2. The temporary M4A file is converted to WAV format using FFmpeg, and stored in a temporary file
3. The temporary WAV file is stored to the storage (can be local file system or cloud storage)
4. Delete the temporary files

The reason for using FFmpeg is that it is a powerful and flexible tool for audio processing, and it is widely used and supported. FFmpeg  is also comes with FFprobe, which can be used to get the audio duration and other metadata.


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

The benchmark results for files in ./sample directory:
```
BenchmarkAudioService_ConvertAudio_AllSamples/eleven.m4a-12         	       6	 168978898 ns/op	 2076373 B/op	      96 allocs/op
BenchmarkAudioService_ConvertAudio_AllSamples/ninety.m4a-12         	       3	 489590436 ns/op	16044032 B/op	      98 allocs/op
BenchmarkAudioService_ConvertAudio_AllSamples/three.m4a-12          	       8	 136805618 ns/op	  715934 B/op	      95 allocs/op
```

From the benchmark results, we can see that the conversion process is quite fast, but it may still may affect the user experience if the audio conversion is done in the main goroutine. This is the reason why the audio conversion is done in a background process.

### Limitations

#### File Descriptor Limit

Given the expected traffic is less than 10rps, and 250K requests per day, the audio conversion process using temporary files is not an issue. We may need to note that the conversion using temporary files is also limited by the file descriptor limit. 

#### Storage

In this implementation, the audio data is stored in local file system, so there is no additional data transfer to storage because the bytes is already in memory and just need to persist to storage.
But local storage is not production ready, and we may need to use dedicated storage that can be accessed by multiple instances, for example AWS S3.

If we use AWS S3, we also get benefits to upload file directly to S3 bucket from user's device, and we can use AWS S3 pre-signed URL to upload file to S3 bucket. This way, we can reduce the data transfer to/from storage.

Please note that both M4A and WAV files are stored in the storage, so the storage size will be doubled. It is to make sure we don't need to convert the audio file again if the user wants to retrieve the audio file in the same format. 