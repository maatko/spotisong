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
	DownloadURL string
	Binary      string
	ProgressBar *pb.ProgressBar
}

func NewTailWind(version string, cache string) (TailWind, error) {
	osName := runtime.GOOS
	osArch := runtime.GOARCH

	if osArch == "amd64" {
		osArch = "x64"
	} else {
		osArch = "arm64"
	}

	return TailWind{
		Version:     version,
		Binary:      cache,
		DownloadURL: "https://github.com/tailwindlabs/tailwindcss/releases/download/%s/tailwindcss-%s-%s",
	}.Setup(osName, osArch)
}

func (tailWind TailWind) Setup(system string, arch string) (TailWind, error) {
	// makes sure that the directory
	// that will hold the tailwind binary exists
	_, err := os.Stat(tailWind.Binary)
	if os.IsNotExist(err) {
		err = os.Mkdir(tailWind.Binary, 0755)
		if err != nil {
			return tailWind, err
		}
	}

	tailWind.Binary = tailWind.Executable("%s/%s", tailWind.Binary, "tailwind")

	// if the binary exists there is no need
	// to proceed with the setup
	_, err = os.Stat(tailWind.Binary)
	if !os.IsNotExist(err) {
		return tailWind, nil
	}

	fmt.Println("[*] Downloading TailWindCSS")
	return tailWind.Download(system, arch)
}

func (tailWind *TailWind) Watch(input string, output string) (*os.ProcessState, error) {
	process, err := os.StartProcess(tailWind.Binary, []string{
		tailWind.Binary,
		"-i", input,
		"-o", output,
		"--watch",
	}, &os.ProcAttr{
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	})

	if err != nil {
		return nil, err
	}

	return process.Wait()
}

func (tailWind TailWind) Download(system string, arch string) (TailWind, error) {
	file, err := os.Create(tailWind.Binary)
	if err != nil {
		return tailWind, err
	}

	response, err := http.Get(tailWind.Executable(
		tailWind.DownloadURL,
		tailWind.Version,
		system,
		arch,
	))

	if err != nil {
		file.Close()
		os.Remove(tailWind.Binary)
		return tailWind, err
	}

	tailWind.ProgressBar = pb.Full.Start64(response.ContentLength)

	defer response.Body.Close()
	defer tailWind.ProgressBar.Finish()

	_, err = io.Copy(io.MultiWriter(file, &tailWind), response.Body)
	if err != nil {
		return tailWind, err
	}

	file.Close()
	return tailWind, os.Chmod(tailWind.Binary, 0755)
}

func (tailWind TailWind) Executable(binary string, args ...any) string {
	path := fmt.Sprintf(binary, args...)
	if strings.ToLower(runtime.GOOS) == "windows" {
		return path + ".exe"
	}
	return path
}

func (tailwind *TailWind) Write(bytes []byte) (int, error) {
	tailwind.ProgressBar.Add(len(bytes))
	return len(bytes), nil
}
