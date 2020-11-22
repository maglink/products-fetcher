package server

import (
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

func TestFetchCsv(t *testing.T) {
	go func() {
		fs := http.FileServer(http.Dir("../../tests"))
		http.Handle("/", fs)

		err := http.ListenAndServe(":3000", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	productsEntries, err := FetchCsv(context.Background(), "http://localhost:3000/products_example.csv")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 1000, len(productsEntries))
	assert.Greater(t, len(productsEntries[0].Name), 0)
	assert.Greater(t, productsEntries[0].Price, 0.0)
}
