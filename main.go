package main

import (
	"context"
	"database/sql"
	"io"
	"log"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

// db1 -> db2 へ ストリーム を使って移動

func main() {
	pr, pw := io.Pipe()

	// dsn := "unix://user:pass@dbname/var/run/postgresql/.s.PGSQL.5432"
	dsn1 := "postgres://postgres:@localhost:15432/postgres?sslmode=disable"
	sqldb1 := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn1)))

	db1 := bun.NewDB(sqldb1, pgdialect.New())

	conn1, err := db1.Conn(context.Background())
	if err != nil {
		log.Fatalln(err.Error())
	}

	ddl := `
		CREATE TABLE IF NOT EXISTS sample (
			a text,
			b text,
			c text,
			d text,
			e text,
			f text,
			g text
		)
		`

	_, err = conn1.QueryContext(context.Background(), ddl)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for i := 0; i < 5; i++ {
		_, err = conn1.QueryContext(context.Background(), "INSERT INTO sample VALUES ('a', 'b', 'c', 'd', 'e', 'f', 'g')")
		if err != nil {
			log.Fatalln(err.Error())
		}
	}

	go func() {
		_, err = pgdriver.CopyTo(context.Background(), conn1, pw, "COPY sample TO STDOUT")
		if err != nil {
			log.Fatalln(err.Error())
		}
		pw.Close()
		conn1.Close()
	}()

	dsn2 := "postgres://postgres:@localhost:25432/postgres?sslmode=disable"
	sqldb2 := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn2)))

	db2 := bun.NewDB(sqldb2, pgdialect.New())

	conn2, err := db2.Conn(context.Background())
	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = conn2.QueryContext(context.Background(), ddl)
	if err != nil {
		log.Fatalln(err.Error())
	}

	_, err = pgdriver.CopyFrom(context.Background(), conn2, pr, "COPY sample FROM STDIN")
	if err != nil {
		panic(err)
	}
}
