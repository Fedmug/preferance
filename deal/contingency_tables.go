package deal

import (
	"context"
	"log"
	"sync"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

const nTables = 564020

type ContingencyTable [NumberOfSuits][NumberOfHands]int8

var (
	contingencyTables []ContingencyTable
	muTables          sync.Mutex
)

var (
	contingencyTableMap map[ContingencyTable]int
	muTableMap          sync.Mutex
)

var wg sync.WaitGroup

func initContingencyTablesBatch(conn *pgxpool.Pool, begin, end int) {
	defer wg.Done()
	rows, err := conn.Query(context.Background(),
		"SELECT index, matrix::int[][] FROM contingency_tables WHERE index >= $1 AND index < $2", begin, end)
	if err != nil {
		log.Fatalln("query row failed:", err)
	}
	for rows.Next() {
		var index int
		var matrix ContingencyTable
		if err = rows.Scan(&index, &matrix); err != nil {
			log.Fatal(err)
		}
		muTables.Lock()
		contingencyTables[index-1] = matrix
		muTables.Unlock()

		muTableMap.Lock()
		contingencyTableMap[matrix] = index - 1
		muTableMap.Unlock()
	}
	if rows.Err() != nil {
		log.Fatalln("rows error:", err)
	}
}

func initContingencyTablesPool() {
	pgConn, err := pgxpool.Connect(context.Background(), "postgres://rafael:VfPLiCASXsMd7Y@localhost/preferance")
	if err != nil {
		log.Fatalln("pgconn failed to connect:", err)
	}
	defer pgConn.Close()
	contingencyTables = make([]ContingencyTable, nTables)
	contingencyTableMap = make(map[ContingencyTable]int, nTables)

	less4Index := 13293
	nThreads := 18
	batchSize := (nTables - less4Index) / nThreads
	initContingencyTables(less4Index)
	for i := less4Index + 1; i <= nTables; i += batchSize {
		wg.Add(1)
		go initContingencyTablesBatch(pgConn, i, i+batchSize)
	}
	wg.Wait()
}

func initContingencyTables(maxIndex int) {
	pgConn, err := pgx.Connect(context.Background(), "postgres://rafael:VfPLiCASXsMd7Y@localhost/preferance")
	if err != nil {
		log.Fatalln("pgconn failed to connect:", err)
	}
	defer pgConn.Close(context.Background())
	contingencyTables = make([]ContingencyTable, nTables)
	contingencyTableMap = make(map[ContingencyTable]int, nTables)

	for nSuits := 1; nSuits <= NumberOfSuits; nSuits++ {
		rows, err := pgConn.Query(context.Background(),
			"SELECT index, matrix::int[][] FROM contingency_tables WHERE ARRAY_LENGTH(matrix, 1) = $1"+
				"AND index <= $2", nSuits, maxIndex)
		if err != nil {
			log.Fatalln("query row failed:", err)
		}
		for rows.Next() {
			var index int
			matrix := make([][NumberOfHands]int8, nSuits)
			if err = rows.Scan(&index, &matrix); err != nil {
				log.Fatal(err)
			}
			copy(contingencyTables[index-1][:nSuits], matrix)
			contingencyTableMap[contingencyTables[index-1]] = index - 1
		}
		if rows.Err() != nil {
			log.Fatalln("rows error:", err)
		}
	}
}

func countTableLines(tableName string) int {
	pgConn, err := pgx.Connect(context.Background(), "postgres://rafael:VfPLiCASXsMd7Y@localhost/preferance")
	if err != nil {
		log.Fatalln("pgconn failed to connect:", err)
	}
	defer pgConn.Close(context.Background())
	var count int
	err = pgConn.QueryRow(context.Background(), "SELECT count(*) FROM "+tableName).Scan(&count)
	if err != nil {
		log.Fatalln("query row failed:", err)
	}
	return count
}

func getMatrix() [][NumberOfSuits][NumberOfHands]int8 {
	pgConn, err := pgx.Connect(context.Background(), "postgres://rafael:VfPLiCASXsMd7Y@localhost/preferance")
	if err != nil {
		log.Fatalln("pgconn failed to connect:", err)
	}
	defer pgConn.Close(context.Background())
	matrices := make([][NumberOfSuits][NumberOfHands]int8, 0, 6e5)
	rows, err := pgConn.Query(context.Background(),
		"SELECT matrix::int[][] FROM contingency_tables WHERE ARRAY_LENGTH(matrix, 1) = $1", NumberOfSuits)
	if err != nil {
		log.Fatalln("query row failed:", err)
	}
	i := 0
	for ; rows.Next(); i++ {
		var matrix [NumberOfSuits][NumberOfHands]int8
		if err = rows.Scan(&matrix); err != nil {
			log.Fatal(err)
		}
		matrices = append(matrices, matrix)
	}
	return matrices
}
