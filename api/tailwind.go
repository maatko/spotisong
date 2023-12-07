package api

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

type TailWind struct {
	Version     string
	Binary      string
	ProgressBar *pb.ProgressBar
}

func (tailwind *TailWind) Watch(input string, output string) error {
	tailwind.Setup()

	attrs := os.ProcAttr{
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	}

	process, err := os.StartProcess(tailwind.Binary, []string{
		tailwind.Binary,
		"-i", input,
		"-o", output,
		"--watch",
	}, &attrs)

	if err != nil {
		return err
	}

	state, err := process.Wait()
	if err != nil {
		return err
	}

	fmt.Printf("Process exited with status %v\n", state.ExitCode())
	return nil
}

func (tailwind *TailWind) Setup() {
	// makes sure that the directory
	// that will hold the tailwind binary exists
	_, err := os.Stat(tailwind.Binary)
	if os.IsNotExist(err) {
		err = os.Mkdir(tailwind.Binary, 0755)
		if err != nil {
			panic(err)
		}
	}

	if strings.ToLower(runtime.GOOS) == "windows" {
		tailwind.Binary += "tailwind.exe"
	} else {
		tailwind.Binary += "tailwind"
	}

	// if the tailwind binary file exists
	// there is no need to re-download it
	_, err = os.Stat(tailwind.Binary)
	if !os.IsNotExist(err) {
		return
	}

	var arch string
	if runtime.GOARCH == "amd64" {
		arch = "x64"
	} else {
		arch = "arm64"
	}

	downloadURL := fmt.Sprintf(
		TAILWIND_DOWNLOAD_URL,
		tailwind.Version,
		strings.ToLower(runtime.GOOS),
		arch,
	)

	if strings.ToLower(runtime.GOOS) == "windows" {
		downloadURL += ".exe"
	}

	err = tailwind.Download(downloadURL)
	if err != nil {
		panic(err)
	}
}

func (tailwind *TailWind) Download(url string) error {
	fmt.Println("> Downloading TailwindCSS")

	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	file, err := os.Create(tailwind.Binary)
	if err != nil {
		return err
	}
	defer file.Close()

	tailwind.ProgressBar = pb.Full.Start64(response.ContentLength)
	defer tailwind.ProgressBar.Finish()

	writer := io.MultiWriter(file, tailwind)
	_, err = io.Copy(writer, response.Body)
	if err != nil {
		return err
	}

	os.Chmod(tailwind.Binary, 0755)
	return nil
}

func (tailwind *TailWind) Write(p []byte) (int, error) {
	tailwind.ProgressBar.Add(len(p))
	return len(p), nil
}

const TAILWIND_DOWNLOAD_URL = "https://github.com/tailwindlabs/tailwindcss/releases/download/%s/tailwindcss-%s-%s"
