package main

import (
	"encoding/csv"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"log"
	"math/rand"
	"os"
)

func main() {
	file, err := os.Create("../products_example.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	writer.Comma = ';'
	defer writer.Flush()

	for i := 0; i < 1000; i++ {
		name := randomdata.SillyName()
		price := rand.Float64() * 50000
		writer.Write([]string{name, fmt.Sprintf("%.2f", price)})
	}
}
