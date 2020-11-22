package server

import (
	"context"
	"encoding/csv"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"strconv"
)

type CsvProductEntry struct {
	Name  string
	Price float64
}

func FetchCsv(ctx context.Context, url string) ([]CsvProductEntry, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	r := csv.NewReader(resp.Body)
	r.Comma = ';'
	var entries []CsvProductEntry
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "can't read the line")
		}

		if len(record) < 2 {
			continue
		}

		price, err := strconv.ParseFloat(record[1], 32)
		if err != nil {
			return nil, errors.Wrap(err, "can't parse price field")
		}

		entries = append(entries, CsvProductEntry{
			Name:  record[0],
			Price: price,
		})
	}

	return entries, nil
}
