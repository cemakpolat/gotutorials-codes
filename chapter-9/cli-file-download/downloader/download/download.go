package download

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error making get request %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid response code %d", resp.StatusCode)
	}

	file, err := os.Create(filepath)
	if err != nil {
		return fmt.Errorf("error creating file %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("error copying data %w", err)
	}
	return nil
}
