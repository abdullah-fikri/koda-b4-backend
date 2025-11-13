package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Db *pgxpool.Pool

func ConnectDb() {
	if Db != nil {
		return
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		fmt.Println("DATABASE_URL not found")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Println("failed to create pool:", err)
		return
	}

	if err = conn.Ping(ctx); err != nil {
		fmt.Println("cannot ping database:", err)
		return
	}

	Db = conn
	fmt.Println("database connected")
}
