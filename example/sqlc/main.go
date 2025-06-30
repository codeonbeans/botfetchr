package main

import (
	"botmediasaver/generated/sqlc"
	"botmediasaver/internal/client/pgxpool"
	"context"
	"fmt"
)

func main() {
	pgpool, err := pgxpool.NewPgxpool(pgxpool.PgxpoolOptions{
		Url:            "postgres://khoakomlem:peakkhoakomlempassword@localhost:5432/botmediasaver",
		MaxConnections: 10,
	})
	if err != nil {
		panic(err)
	}
	defer pgpool.Close()

	if err = pgpool.Ping(context.TODO()); err != nil {
		panic(err)
	}

	fmt.Println("Connected to the database successfully!")

	a := sqlc.New(pgpool)
	rows, err := a.ListAccountTelegrams(context.TODO(), sqlc.ListAccountTelegramsParams{
		// Orderby: "telegram_id",
		OrderBy: "created_at_desc",
		Offset:  0,
		Limit:   10,
	})
	if err != nil {
		panic(err)
	}

	for _, row := range rows {
		fmt.Println(row.ID, row.TelegramID)
	}
}
