package main

import (
	"context"
	"fmt"
	"github.com/codeonbeans/botfetchr/generated/sqlc"
	"github.com/codeonbeans/botfetchr/internal/client/pgxpool"
)

func main() {
	pgpool, err := pgxpool.NewPgxpool(pgxpool.PgxpoolOptions{
		//Url:            "",
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
