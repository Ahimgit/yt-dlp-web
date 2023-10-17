# A simple yt-dlp Web UI

Web UI for yt-dlp for music download with metadata and cover embed.

###Features:
- Simple Web UI with yt-dlp stdout output
- Automatically updates tag metadata and cover through Deezer API integration

![ui](https://github.com/Ahimgit/yt-dlp-web/assets/6353365/1cf9cf7b-290f-4c18-952e-b7932d2e1064)

## Building

Checkout and build with `go build ./cmd/yt-dlp-web`

## Running locally

- Download or install package of [ffmpeg](https://ffmpeg.org/download.html)
- Download [yt-dlp](https://github.com/yt-dlp/yt-dlp/releases/)
- Make ffmpeg and yt-dlp available in PATH or provide their location via params
- Start yt-dlp-web optionally providing params

```
yt-dlp-web --help
Usage of yt-dlp-web.exe:
  -ffmpegPath string
        full path to ffmpeg (default "ffmpeg")
  -outputFolder string
        output folder for downloads (default "./downloads/")
  -port string
        port to listen (default "8801")
  -ytdpPath string
        full path to yt-dlp (default "yt-dlp")
```

