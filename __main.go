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

// csv へ ストリームを使って書き込み

func main() {
	f, err := os.Create("./output.csv")
	if err != nil {
		log.Fatalln(err.Error())
	}

	defer f.Close()

	dsn := "postgres://postgres:@localhost:15432/postgres?sslmode=disable"
	// dsn := "unix://user:pass@dbname/var/run/postgresql/.s.PGSQL.5432"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())

	conn, err := db.Conn(context.Background())
	if err != nil {
		log.Fatalln(err.Error())
	}

	defer conn.Close()

	_, err = pgdriver.CopyTo(context.Background(), conn, f, "COPY sample TO STDOUT DELIMITER ',' CSV HEADER")
	if err != nil {
		log.Fatalln(err.Error())
	}
}
