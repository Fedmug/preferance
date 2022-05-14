package yellow

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
)

type Rate map[string]float64

var PlayerRate map[string]Rate

const tableSize = 24267
const nGames = 4

func GetPlayerRates(rateTablePath string) {
	f, err := os.Open(rateTablePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	PlayerRate = make(map[string]Rate, tableSize)
	csvReader := csv.NewReader(f)
	headers, err := csvReader.Read()
	if err != nil {
		log.Fatal(err)
	}
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		rates := make(Rate, nGames)
		for i := 1; i <= nGames; i++ {
			if record[i] != "" {
				avgRate, err := strconv.ParseFloat(record[i], 32)
				if err != nil {
					log.Fatal(err)
				}
				rates[headers[i]] = avgRate
			}
		}
		PlayerRate[record[0]] = rates
	}
}
