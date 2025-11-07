package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Db *pgxpool.Pool

func ConnectDb() {
	godotenv.Load()
	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("failed Connect Database", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal("cannot ping database:", err)
	}

	Db = conn
	fmt.Println("sukses ")
}
