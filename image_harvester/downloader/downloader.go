package downloader

import (
	"bytes"
	"image/jpeg"
	"net/http"
	"os"

	"github.com/nfnt/resize"
)

func DownloadAndResizeImage(url string, size uint) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	img, err := jpeg.Decode(resp.Body)
	if err != nil {
		return nil, err
	}

	resizedImg := resize.Resize(size, 0, img, resize.Lanczos3)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, resizedImg, nil); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func SaveImageToFile(filename string, directory string, data []byte) error {

	err := EnsureDir(directory)

	if err != nil {
		return err
	}

	return os.WriteFile(directory+"/"+filename, data, 0644)
}

// EnsureDir checks if a directory exists, and if not, creates it
func EnsureDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}
