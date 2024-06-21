// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createPlayer = `-- name: CreatePlayer :one
insert into players (steam_user_id, name, bot) values ($1, $2, $3) returning id, steam_user_id, name, bot
`

type CreatePlayerParams struct {
	SteamUserID pgtype.Int8
	Name        string
	Bot         pgtype.Bool
}

func (q *Queries) CreatePlayer(ctx context.Context, arg CreatePlayerParams) (Player, error) {
	row := q.db.QueryRow(ctx, createPlayer, arg.SteamUserID, arg.Name, arg.Bot)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.SteamUserID,
		&i.Name,
		&i.Bot,
	)
	return i, err
}

const createSteamUser = `-- name: CreateSteamUser :one
insert into steam_users (steam_id, steam_community_id) values ($1, $2) returning id, steam_id, steam_community_id
`

type CreateSteamUserParams struct {
	SteamID          string
	SteamCommunityID int64
}

func (q *Queries) CreateSteamUser(ctx context.Context, arg CreateSteamUserParams) (SteamUser, error) {
	row := q.db.QueryRow(ctx, createSteamUser, arg.SteamID, arg.SteamCommunityID)
	var i SteamUser
	err := row.Scan(&i.ID, &i.SteamID, &i.SteamCommunityID)
	return i, err
}

const deletePlayer = `-- name: DeletePlayer :exec
delete from players where id = $1
`

func (q *Queries) DeletePlayer(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deletePlayer, id)
	return err
}

const deleteSteamUser = `-- name: DeleteSteamUser :exec
delete from steam_users where id = $1
`

func (q *Queries) DeleteSteamUser(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteSteamUser, id)
	return err
}

const getPlayer = `-- name: GetPlayer :one
select id, steam_user_id, name, bot from players where id = $1 limit 1
`

func (q *Queries) GetPlayer(ctx context.Context, id int64) (Player, error) {
	row := q.db.QueryRow(ctx, getPlayer, id)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.SteamUserID,
		&i.Name,
		&i.Bot,
	)
	return i, err
}

const getPlayerByName = `-- name: GetPlayerByName :one
select id, steam_user_id, name, bot from players where name = $1 limit 1
`

func (q *Queries) GetPlayerByName(ctx context.Context, name string) (Player, error) {
	row := q.db.QueryRow(ctx, getPlayerByName, name)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.SteamUserID,
		&i.Name,
		&i.Bot,
	)
	return i, err
}

const getSteamUser = `-- name: GetSteamUser :one
select id, steam_id, steam_community_id from steam_users where id = $1 limit 1
`

func (q *Queries) GetSteamUser(ctx context.Context, id int64) (SteamUser, error) {
	row := q.db.QueryRow(ctx, getSteamUser, id)
	var i SteamUser
	err := row.Scan(&i.ID, &i.SteamID, &i.SteamCommunityID)
	return i, err
}

const getSteamUserBySteamId = `-- name: GetSteamUserBySteamId :one
select id, steam_id, steam_community_id from steam_users where steam_id = $1
`

func (q *Queries) GetSteamUserBySteamId(ctx context.Context, steamID string) (SteamUser, error) {
	row := q.db.QueryRow(ctx, getSteamUserBySteamId, steamID)
	var i SteamUser
	err := row.Scan(&i.ID, &i.SteamID, &i.SteamCommunityID)
	return i, err
}

const listPlayers = `-- name: ListPlayers :many
select id, steam_user_id, name, bot from players order by name
`

func (q *Queries) ListPlayers(ctx context.Context) ([]Player, error) {
	rows, err := q.db.Query(ctx, listPlayers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Player
	for rows.Next() {
		var i Player
		if err := rows.Scan(
			&i.ID,
			&i.SteamUserID,
			&i.Name,
			&i.Bot,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSteamUsers = `-- name: ListSteamUsers :many
select id, steam_id, steam_community_id from steam_users order by steam_id
`

func (q *Queries) ListSteamUsers(ctx context.Context) ([]SteamUser, error) {
	rows, err := q.db.Query(ctx, listSteamUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SteamUser
	for rows.Next() {
		var i SteamUser
		if err := rows.Scan(&i.ID, &i.SteamID, &i.SteamCommunityID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePlayer = `-- name: UpdatePlayer :exec
update players set steam_user_id = $2, name = $3, bot = $4 where id = $1
`

type UpdatePlayerParams struct {
	ID          int64
	SteamUserID pgtype.Int8
	Name        string
	Bot         pgtype.Bool
}

func (q *Queries) UpdatePlayer(ctx context.Context, arg UpdatePlayerParams) error {
	_, err := q.db.Exec(ctx, updatePlayer,
		arg.ID,
		arg.SteamUserID,
		arg.Name,
		arg.Bot,
	)
	return err
}

const updateSteamUser = `-- name: UpdateSteamUser :exec
update steam_users set steam_id = $2, steam_community_id = $3 where id = $1
`

type UpdateSteamUserParams struct {
	ID               int64
	SteamID          string
	SteamCommunityID int64
}

func (q *Queries) UpdateSteamUser(ctx context.Context, arg UpdateSteamUserParams) error {
	_, err := q.db.Exec(ctx, updateSteamUser, arg.ID, arg.SteamID, arg.SteamCommunityID)
	return err
}