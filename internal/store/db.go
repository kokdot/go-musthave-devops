package store

import (
	"context"
	"fmt"
	"time"

	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)


func GetPing(dataBaseDSN string) (bool, error) {
	// urlExample := "postgres://postgres:postgres@localhost:5432/postgres"
    db, err := sql.Open("pgx", dataBaseDSN)
    if err != nil {
        return false, err
    }
    defer db.Close()
    // работаем с базой
    // ...
    // можем продиагностировать соединение
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    if err = db.PingContext(ctx); err != nil {
        return false, err
    }
	fmt.Println("Ping Ok")
    // в процессе работы
	return true, nil
} 
