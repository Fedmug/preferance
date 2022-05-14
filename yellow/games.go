package yellow

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type FileName string

const (
	YellowGames   FileName = "yellow_games.csv"
	YellowInfo    FileName = "yellow_info.csv"
	YellowDeals   FileName = "yellow_deals.csv"
	YellowBidPlay FileName = "yellow_bid_play.txt"
)

func ReadGameIds(filePath string, gameType string) ([]int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("can't open file %s: %v", filePath, err)
	}
	defer f.Close()
	gameIds := make([]int, 0, 1000)
	csvReader := csv.NewReader(f)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return gameIds, fmt.Errorf("unexpected error while reading record: %v", err)
		}
		if record[1] == gameType {
			gameId, err := strconv.Atoi(record[0])
			if err != nil {
				return gameIds, fmt.Errorf("cannot convert %s to int: %v", record[1], err)
			}
			gameIds = append(gameIds, gameId)
		}
	}
	return gameIds, nil
}
