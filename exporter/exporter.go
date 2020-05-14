package exporter

import (
	"encoding/csv"
	"os"

	"log"

	"github.com/yeo/bima/dto"
)

func Import(password []byte, path string) error {
	f, e := os.Open(path)
	if e != nil {
		log.Fatal(e)
	}

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range records {
		token := dto.Token{
			Name:     row[1],
			URL:      row[0],
			RawToken: row[2],
		}
		dto.AddSecret(&token, password)
	}

	log.Println(records)

	return nil
}

func Export(password []byte, path string) error {

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
