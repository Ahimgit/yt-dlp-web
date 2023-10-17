package main

import (
	"bufio"
	"log"
	"os/exec"
	"sync"
)

func ExecCmd(command string, params *Config, stdoutCallback func(string)) error {
	cmd := buildYTDLPCommand(command, params)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	stdoutScanner := bufio.NewScanner(stdoutPipe)
	stderrScanner := bufio.NewScanner(stderrPipe)

	group := waitGroup(2)
	go scanLines(stdoutScanner, stdoutCallback, group)
	go scanLines(stderrScanner, stdoutCallback, group)
	defer group.Wait()

	if err := cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

func buildYTDLPCommand(cmdString string, params *Config) *exec.Cmd {
	args := []string{
		"--output", params.outputFolder + "%(artist,creator,uploader,uploader_id)s - %(title)s",
		"--newline",
		"--embed-metadata",
		"--embed-thumbnail",
		"--extract-audio",
		"--no-quiet",
		"--no-simulate",
		"--print", "after_move:DoneFile$%(filepath)s",
		"--audio-format", "mp3",
		cmdString,
	}
	if params.ffMpegPath != "ffmpeg" {
		args = append(args, "--ffmpeg-location", params.ffMpegPath)
	}
	log.Println("Running", params.ytdlpPath, args)
	return exec.Command(params.ytdlpPath, args...)
}

func waitGroup(count int) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(count)
	return &wg
}

func scanLines(scanner *bufio.Scanner, stdoutCallback func(string), wg *sync.WaitGroup) {
	for scanner.Scan() {
		stdoutCallback(scanner.Text())
	}
	wg.Done()
}
