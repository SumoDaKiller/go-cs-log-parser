package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"go-cs-log-parser/database"
	"html/template"
	"log"
	"os"
)

func generateAllPages() error {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, k.String("database_url"))
	if err != nil {
		log.Fatalf("error connecting to database: %v\n", err)
	}
	defer conn.Close(ctx)

	err = generatePlayersPage(ctx, conn)
	if err != nil {
		return err
	}

	err = generatePlayerPages(ctx, conn)
	if err != nil {
		return err
	}

	return nil
}

func generatePlayersPage(ctx context.Context, conn *pgx.Conn) error {
	queries := database.New(conn)
	players, err := queries.ListPlayers(ctx)
	if err != nil {
		return err
	}

	outputPath := k.String("output_path")
	outputFile, err := os.Create(outputPath + "players.html")
	if err != nil {
		return err
	}
	defer outputFile.Close()

	tmpl, err := template.ParseFiles("templates/players.html")
	if err != nil {
		return err
	}

	err = tmpl.Execute(outputFile, players)
	if err != nil {
		return err
	}

	return nil
}

type PlayerPageData struct {
	Player         database.GetPlayerWithStatsRow
	PlayerAccolade []database.GetAccoladeForPlayerRow
}

func generatePlayerPages(ctx context.Context, conn *pgx.Conn) error {
	queries := database.New(conn)
	outputPath := k.String("output_path")

	players, err := queries.ListPlayers(ctx)
	if err != nil {
		return err
	}

	for _, player := range players {
		playerStats, err := queries.GetPlayerWithStats(ctx, player.ID)
		if err != nil {
			return err
		}
		playerAccolade, err := queries.GetAccoladeForPlayer(ctx, player.ID)
		if err != nil {
			return err
		}
		outputFile, err := os.Create(fmt.Sprint(outputPath, "player_", player.ID, ".html"))
		if err != nil {
			return err
		}
		tmpl, err := template.ParseFiles("templates/player.html")
		if err != nil {
			return err
		}
		err = tmpl.Execute(outputFile, PlayerPageData{
			Player:         playerStats,
			PlayerAccolade: playerAccolade,
		})
		if err != nil {
			return err
		}
		outputFile.Close()
	}

	return nil
}
