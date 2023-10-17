package main

import (
	"flag"
	"log"
	"net/http"
)

type Config struct {
	outputFolder string
	ytdlpPath    string
	ffMpegPath   string
}

func main() {
	paramOutputFolder := flag.String("outputFolder", "./downloads/", "output folder for downloads")
	paramYtdlpPath := flag.String("ytdpPath", "yt-dlp", "full path to yt-dlp")
	paramFfmpegPath := flag.String("ffmpegPath", "ffmpeg", "full path to ffmpeg")
	paramPort := flag.String("port", "8801", "port to listen")
	flag.Parse()

	log.Println("Starting up...")
	log.Println("Using output folder: " + *paramOutputFolder)
	log.Println("Using yt-dlp path: " + *paramYtdlpPath)
	log.Println("Using ffmpeg path: " + *paramFfmpegPath)
	log.Println("Listening on port: " + *paramPort)

	handler := &Config{
		outputFolder: *paramOutputFolder,
		ytdlpPath:    *paramYtdlpPath,
		ffMpegPath:   *paramFfmpegPath,
	}
	http.HandleFunc("/", ServeHTML)
	http.HandleFunc("/ws", handler.HandleWSConnection)

	log.Fatal(http.ListenAndServe(":"+*paramPort, nil))
}
