package exporter

import (
	"encoding/csv"
	"os"

	"github.com/yeo/bima/dto"
)

func Export(password string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	tokens, err := dto.LoadTokens()

	for _, t := range tokens {
		err := writer.Write([]string{t.URL, t.Name, t.DecryptToken(password)})
		if err != nil {
			return err
		}
	}

	return nil
}
