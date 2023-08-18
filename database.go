package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func connect() (clickhouse.Conn, error) {
	dbhost := os.Getenv("CLICKHOUSE_HOST")
	if dbhost == "" {
		// dbhost = "clickhouse-db"
		dbhost = "localhost"
	}
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", dbhost, 9000)},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
	})
	return conn, err
}

func createTable() {
	conn, err := connect()
	if err != nil {
		panic(err)
	}

	// if err := conn.Exec(context.Background(), `DROP TABLE IF EXISTS devices`); err != nil {
	// 	panic(err)
	// }
	// if err := conn.Exec(context.Background(), `DROP TABLE IF EXISTS positions`); err != nil {
	// 	panic(err)
	// }

	// Create table based on struct
	const ddl = `
		CREATE TABLE IF NOT EXISTS devices (
			id UUID DEFAULT generateUUIDv4() PRIMARY KEY,
			name String NULL DEFAULT NULL,
			model String NULL DEFAULT NULL,
			phone String NULL DEFAULT NULL,
			status String NULL DEFAULT NULL,
			contact String NULL DEFAULT NULL,
			category String NULL DEFAULT NULL,
			disabled Boolean NULL DEFAULT NULL,
			unique_id String NULL DEFAULT NULL,
			attributes String NULL DEFAULT NULL,
			last_update DateTime DEFAULT now(),
			expiration_time String NULL DEFAULT NULL
		) ENGINE = MergeTree()
		PARTITION BY toYYYYMMDD(toDate(last_update))
		ORDER BY (id, last_update);
	`
	const ddl2 = `
		CREATE TABLE IF NOT EXISTS positions (
			id UUID DEFAULT generateUUIDv4() PRIMARY KEY,
			speed Integer,
			valid Boolean,
			course Integer,
			address String NULL DEFAULT NULL,
			fix_time String,
			network String,
			accuracy Integer,
			altitude Integer,
			device_id String NULL DEFAULT NULL,
			latitude Double,
			outdated Boolean,
			protocol String NULL DEFAULT NULL,
			longitude Double,
			attributes String NULL DEFAULT NULL,
			device_time String,
			server_time String
		) ENGINE = MergeTree()
		ORDER BY (id, server_time);
	`

	if err := conn.Exec(context.Background(), ddl); err != nil {
		panic(err)
	}
	if err := conn.Exec(context.Background(), ddl2); err != nil {
		panic(err)
	}
}
