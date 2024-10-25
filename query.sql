-- Steam User

-- name: GetSteamUser :one
select * from steam_users where id = $1 limit 1;

-- name: GetSteamUserBySteamId :one
select * from steam_users where steam_id = $1;

-- name: CreateSteamUser :one
insert into steam_users (steam_id, steam_community_id) values ($1, $2) returning *;

-- name: UpdateSteamUser :exec
update steam_users set steam_id = $2, steam_community_id = $3 where id = $1;

-- Players

-- name: GetPlayer :one
select * from players where id = $1 limit 1;

-- name: GetPlayerByName :one
select * from players where name = $1 limit 1;

-- name: CreatePlayer :one
insert into players (steam_user_id, name, bot) values ($1, $2, $3) returning *;

-- name: UpdatePlayer :exec
update players set steam_user_id = $2, name = $3, bot = $4 where id = $1;

-- Maps

-- name: GetMap :one
select * from maps where id = $1 limit 1;

-- name: GetMapByName :one
select * from maps where name = $1 limit 1;

-- name: CreateMap :one
insert into maps (name) values ($1) returning *;

-- Teams

-- name: GetTeam :one
select * from teams where id = $1 limit 1;

-- name: GetTeamByName :one
select * from teams where name = $1 limit 1;

-- name: CreateTeam :one
insert into teams (name) values ($1) returning *;

-- Weapons

-- name: GetWeapon :one
select * from weapons where id = $1 limit 1;

-- name: GetWeaponByName :one
select * from weapons where name = $1 limit 1;

-- name: CreateWeapon :one
insert into weapons (name) values ($1) returning *;

-- Special Kills

-- name: GetSpecialKill :one
select * from special_kills where id = $1 limit 1;

-- name: GetSpecialKillByName :one
select * from special_kills where name = $1 limit 1;

-- name: CreateSpecialKill :one
insert into special_kills (name) values ($1) returning *;

-- Game Type

-- name: GetGameType :one
select * from game_types where id = $1 limit 1;

-- name: GetGameTypeByName :one
select * from game_types where name = $1 limit 1;

-- name: CreateGameType :one
insert into game_types (name) values ($1) returning *;

-- Hit Groups

-- name: GetHitGroup :one
select * from hit_groups where id = $1 limit 1;

-- name: GetHitGroupByName :one
select * from hit_groups where name = $1 limit 1;

-- name: CreateHitGroup :one
insert into hit_groups (name) values ($1) returning *;

-- Matches

-- name: GetMatch :one
select * from matches where id = $1 limit 1;

-- name: GetMatchByMapAndStartDateTime :one
select * from matches where map_id = $1 and start_date = $2 and start_time = $3 limit 1;

-- name: CreateMatch :one
insert into matches (start_date, start_time, map_id) values ($1, $2, $3) returning *;

-- name: UpdateMatch :exec
update matches set end_date=$2, end_time=$3, score_ct=$4, score_t=$5, game_type_id=$6 where id=$1;

-- Rounds

-- name: GetRound :one
select * from rounds where id = $1 limit 1;

-- name: GetRoundByStartDateTime :one
select * from rounds where start_date = $1 and start_time = $2 limit 1;

-- name: CreateRound :one
insert into rounds (start_date, start_time, match_id) values ($1, $2, $3) returning *;

-- name: UpdateRound :exec
update rounds set end_date=$2, end_time=$3, winner_team_id=$4 where id = $1;

-- Team Switch

-- name: GetTeamSwitchByPlayerAndDateTime :one
select * from team_switch where player_id = $1 and switch_date = $2 and switch_time = $3 limit 1;

-- name: CreateTeamSwitchEvent :exec
insert into team_switch (player_id, from_team_id, to_team_id, switch_date, switch_time, round_id) values ($1, $2, $3, $4, $5, $6) returning *;

-- Round Teams

-- name: GetRoundTeamEntry :one
select * from round_teams where id = $1 limit 1;

-- name: GetRoundTeamEntryByPlayerTeamRound :one
select * from round_teams where player_id = $1 and team_id = $2 and round_id = $3 limit 1;

-- name: CreateRoundTeamEntry :one
insert into round_teams (player_id, team_id, round_id) values ($1, $2, $3) returning *;

-- Kills

-- name: GetKill :one
select * from kills where id = $1 limit 1;

-- name: GetKillByKillerKilledRoundDateTime :one
select * from kills where killer_id = $1 and killed_id = $2 and round_id = $3 and kill_date = $4 and kill_time = $5 limit 1;

-- name: CreateKill :one
insert into kills (killer_id, killed_id, round_id, kill_time, kill_date, killer_team_id, killed_team_id, killer_position_x, killer_position_y, killer_position_z, killed_position_x, killed_position_y, killed_position_z, killer_weapon_id, special_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) returning *;

-- Attacks

-- name: GetAttack :one
select * from attacks where id = $1 limit 1;

-- name: GetAttackByAttackerAttackedRoundDateTime :one
select * from attacks where attacker_id = $1 and attacked_id = $2 and round_id = $3 and attack_date = $4 and attack_time = $5 limit 1;

-- name: CreateAttack :one
insert into attacks (attacker_id, attacked_id, round_id, attack_time, attack_date, attacker_team_id, attacked_team_id, attacker_position_x, attacker_position_y, attacker_position_z, attacked_position_x, attacked_position_y, attacked_position_z, attacker_weapon_id, damage, damage_armor, health, armor, hit_group_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19) returning *;
