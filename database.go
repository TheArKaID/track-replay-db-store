package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
)

func connect() (clickhouse.Conn, error) {
	CH_HOST := os.Getenv("CLICKHOUSE_HOST")
	if CH_HOST == "" {
		CH_HOST = "localhost"
	}
	CH_DB := os.Getenv("CLICKHOUSE_DB")
	if CH_DB == "" {
		CH_DB = "default"
	}
	CH_USER := os.Getenv("CLICKHOUSE_USER")
	if CH_USER == "" {
		CH_USER = "default"
	}
	CH_PASS := os.Getenv("CLICKHOUSE_PASS")
	if CH_PASS == "" {
		CH_PASS = ""
	}

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", CH_HOST, 9000)},
		Auth: clickhouse.Auth{
			Database: CH_DB,
			Username: CH_USER,
			Password: CH_PASS,
		},
		DialContext: func(ctx context.Context, addr string) (net.Conn, error) {
			// dialCount++
			var d net.Dialer
			return d.DialContext(ctx, "tcp", addr)
		},
		Debug: false,
		Debugf: func(format string, v ...interface{}) {
			fmt.Printf(format, v)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:      time.Duration(10) * time.Second,
		MaxOpenConns:     50,
		MaxIdleConns:     25,
		ConnMaxLifetime:  time.Duration(10) * time.Minute,
		ConnOpenStrategy: clickhouse.ConnOpenInOrder,
		BlockBufferSize:  10,
	})
	return conn, err
}

func createTable() {
	conn, err := connect()
	if err != nil {
		panic(err)
	}

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

// func xConn() {
// 	fmt.Println(sql.Drivers())
// 	db, _ := sql.Open("clickhouse", "clickhouse://default:@localhost:9000/default?dial_timeout=200ms&max_execution_time=60")
// 	defer db.Close()

// 	row, err := db.Query("SELECT * FROM devices")
// 	if err != nil {
// 		panic(err)
// 	}

// 	columns, _ := row.Columns()
// 	fmt.Println(columns)
// }
