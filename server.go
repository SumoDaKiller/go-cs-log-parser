package main

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-cs-log-parser/database"
	"log"
	"log/slog"
	"net/http"
	"strconv"
)

func runServer() error {
	ctx = context.Background()
	conn, err := pgx.Connect(ctx, k.String("database_url"))
	if err != nil {
		log.Fatalf("error connecting to database: %v\n", err)
	}
	defer conn.Close(ctx)

	queries = database.New(conn)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/players", getAllPlayers)
	e.GET("/players/:id", getPlayer)
	e.GET("/players/:id/accolade", getPlayerAccolade)

	if err := e.Start(":" + k.String("server_port")); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}

	return nil
}

func getAllPlayers(c echo.Context) error {
	players, err := queries.ListPlayers(ctx)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, players)
}

func getPlayer(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	playerStats, err := queries.GetPlayerWithStats(ctx, intID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, playerStats)
}

func getPlayerAccolade(c echo.Context) error {
	id := c.Param("id")
	intID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	playerAccolade, err := queries.GetAccoladeForPlayer(ctx, intID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, playerAccolade)
}
