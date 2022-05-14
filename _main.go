package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// @see https://future-architect.github.io/articles/20210727a/

// var _ pgx.CopyFromSource = &InputRows{}

// const uri = "postgres://postgres:password@localhost:15432/postgres?sslmode=disable"

func main() {
	f, err := os.Create("./output.csv")
	if err != nil {
		log.Fatalln(err.Error())
	}

	// w := csv.NewWriter(f)

	// r, w := io.Pipe()

	// client, err := pgxpool.Connect(context.Background(), uri)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }

	// ddl := `
	// CREATE TABLE IF NOT EXISTS sample (
	// 	a text,
	// 	b text,
	// 	c text,
	// 	d text,
	// 	e text,
	// 	f text,
	// 	g text
	// )
	// `
	// _, err = client.Query(context.Background(), ddl)
	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }

	// for i := 0; i < 5; i++ {
	// 	_, err = client.Query(context.Background(), "INSERT INTO sample VALUES ('a', 'b', 'c', 'd', 'e', 'f', 'g')")
	// 	if err != nil {
	// 		log.Fatalln(err.Error())
	// 	}
	// }

	// _, err = client.CopyFrom(context.Background(), pgx.Identifier{"sample"}, []string{
	// 	"a",
	// 	"b",
	// 	"c",
	// 	"d",
	// 	"e",
	// 	"f",
	// 	"g",
	// }, &InputRows{r: r})

	// if err != nil {
	// 	log.Fatalln(err.Error())
	// }

	// @see https://www.asobou.co.jp/blog/web/go-csv

	dsn := "postgres://postgres:@localhost:15432/postgres?sslmode=disable"
	// dsn := "unix://user:pass@dbname/var/run/postgresql/.s.PGSQL.5432"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	conn, err := db.Conn(context.Background())
	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = pgdriver.CopyTo(context.Background(), conn, f, "COPY sample TO STDOUT DELIMITER ',' CSV HEADER")
	if err != nil {
		log.Fatalln(err.Error())
	}
}

// type InputRows struct {
// 	r       *csv.Reader
// 	nextRow []interface{}
// 	err     error
// }

// func (i *InputRows) Next() bool {
// 	i.nextRow = nil
// 	i.err = nil
// 	record, err := i.r.Read()
// 	log.Println(record)
// 	if err == io.EOF {
// 		return false
// 	} else if err != nil {
// 		i.err = err
// 		return false
// 	}

// 	i.nextRow = []interface{}{
// 		record[0],
// 		record[1],
// 		record[2],
// 		record[3],
// 		record[4],
// 		record[5],
// 		record[6],
// 	}

// 	return true
// }

// func (i InputRows) Values() ([]interface{}, error) {
// 	if i.err != nil {
// 		return nil, i.err
// 	}
// 	return i.nextRow, nil
// }

// func (i InputRows) Err() error {
// 	return i.err
// }
