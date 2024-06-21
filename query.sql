-- name: GetSteamUser :one
select * from steam_users where id = $1 limit 1;

-- name: GetSteamUserBySteamId :one
select * from steam_users where steam_id = $1;

-- name: ListSteamUsers :many
select * from steam_users order by steam_id;

-- name: CreateSteamUser :one
insert into steam_users (steam_id, steam_community_id) values ($1, $2) returning *;

-- name: UpdateSteamUser :exec
update steam_users set steam_id = $2, steam_community_id = $3 where id = $1;

-- name: DeleteSteamUser :exec
delete from steam_users where id = $1;

-- name: GetPlayer :one
select * from players where id = $1 limit 1;

-- name: GetPlayerByName :one
select * from players where name = $1 limit 1;

-- name: ListPlayers :many
select * from players order by name;

-- name: CreatePlayer :one
insert into players (steam_user_id, name, bot) values ($1, $2, $3) returning *;

-- name: UpdatePlayer :exec
update players set steam_user_id = $2, name = $3, bot = $4 where id = $1;

-- name: DeletePlayer :exec
delete from players where id = $1;
