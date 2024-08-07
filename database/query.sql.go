// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createAttack = `-- name: CreateAttack :one
insert into attacks (attacker_id, attacked_id, round_id, attack_time, attack_date, attacker_team_id, attacked_team_id, attacker_position, attacked_position, attacker_weapon_id, damage, damage_armor, health, armor, hit_group_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) returning id, attacker_id, attacked_id, round_id, attack_time, attack_date, attacker_team_id, attacked_team_id, attacker_position, attacked_position, attacker_weapon_id, damage, damage_armor, health, armor, hit_group_id
`

type CreateAttackParams struct {
	AttackerID       pgtype.Int8
	AttackedID       pgtype.Int8
	RoundID          pgtype.Int8
	AttackTime       pgtype.Time
	AttackDate       pgtype.Date
	AttackerTeamID   pgtype.Int8
	AttackedTeamID   pgtype.Int8
	AttackerPosition []byte
	AttackedPosition []byte
	AttackerWeaponID pgtype.Int8
	Damage           int32
	DamageArmor      int32
	Health           int32
	Armor            int32
	HitGroupID       pgtype.Int8
}

func (q *Queries) CreateAttack(ctx context.Context, arg CreateAttackParams) (Attack, error) {
	row := q.db.QueryRow(ctx, createAttack,
		arg.AttackerID,
		arg.AttackedID,
		arg.RoundID,
		arg.AttackTime,
		arg.AttackDate,
		arg.AttackerTeamID,
		arg.AttackedTeamID,
		arg.AttackerPosition,
		arg.AttackedPosition,
		arg.AttackerWeaponID,
		arg.Damage,
		arg.DamageArmor,
		arg.Health,
		arg.Armor,
		arg.HitGroupID,
	)
	var i Attack
	err := row.Scan(
		&i.ID,
		&i.AttackerID,
		&i.AttackedID,
		&i.RoundID,
		&i.AttackTime,
		&i.AttackDate,
		&i.AttackerTeamID,
		&i.AttackedTeamID,
		&i.AttackerPosition,
		&i.AttackedPosition,
		&i.AttackerWeaponID,
		&i.Damage,
		&i.DamageArmor,
		&i.Health,
		&i.Armor,
		&i.HitGroupID,
	)
	return i, err
}

const createGameType = `-- name: CreateGameType :one
insert into game_types (name) values ($1) returning id, name
`

func (q *Queries) CreateGameType(ctx context.Context, name string) (GameType, error) {
	row := q.db.QueryRow(ctx, createGameType, name)
	var i GameType
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const createHitGroup = `-- name: CreateHitGroup :one
insert into hit_groups (name) values ($1) returning id, name
`

func (q *Queries) CreateHitGroup(ctx context.Context, name string) (HitGroup, error) {
	row := q.db.QueryRow(ctx, createHitGroup, name)
	var i HitGroup
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const createKill = `-- name: CreateKill :one
insert into kills (killer_id, killed_id, round_id, kill_time, kill_date, killer_team_id, killed_team_id, killer_position, killed_position, killer_weapon_id, special_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) returning id, killer_id, killed_id, round_id, kill_time, kill_date, killer_team_id, killed_team_id, killer_position, killed_position, killer_weapon_id, special_id
`

type CreateKillParams struct {
	KillerID       pgtype.Int8
	KilledID       pgtype.Int8
	RoundID        pgtype.Int8
	KillTime       pgtype.Time
	KillDate       pgtype.Date
	KillerTeamID   pgtype.Int8
	KilledTeamID   pgtype.Int8
	KillerPosition []byte
	KilledPosition []byte
	KillerWeaponID pgtype.Int8
	SpecialID      pgtype.Int8
}

func (q *Queries) CreateKill(ctx context.Context, arg CreateKillParams) (Kill, error) {
	row := q.db.QueryRow(ctx, createKill,
		arg.KillerID,
		arg.KilledID,
		arg.RoundID,
		arg.KillTime,
		arg.KillDate,
		arg.KillerTeamID,
		arg.KilledTeamID,
		arg.KillerPosition,
		arg.KilledPosition,
		arg.KillerWeaponID,
		arg.SpecialID,
	)
	var i Kill
	err := row.Scan(
		&i.ID,
		&i.KillerID,
		&i.KilledID,
		&i.RoundID,
		&i.KillTime,
		&i.KillDate,
		&i.KillerTeamID,
		&i.KilledTeamID,
		&i.KillerPosition,
		&i.KilledPosition,
		&i.KillerWeaponID,
		&i.SpecialID,
	)
	return i, err
}

const createMap = `-- name: CreateMap :one
insert into maps (name) values ($1) returning id, name
`

func (q *Queries) CreateMap(ctx context.Context, name string) (Map, error) {
	row := q.db.QueryRow(ctx, createMap, name)
	var i Map
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const createMatch = `-- name: CreateMatch :one
insert into matches (start_date, start_time, map_id) values ($1, $2, $3) returning id, start_date, start_time, end_date, end_time, map_id, score_ct, score_t, game_type_id
`

type CreateMatchParams struct {
	StartDate pgtype.Date
	StartTime pgtype.Time
	MapID     pgtype.Int8
}

func (q *Queries) CreateMatch(ctx context.Context, arg CreateMatchParams) (Match, error) {
	row := q.db.QueryRow(ctx, createMatch, arg.StartDate, arg.StartTime, arg.MapID)
	var i Match
	err := row.Scan(
		&i.ID,
		&i.StartDate,
		&i.StartTime,
		&i.EndDate,
		&i.EndTime,
		&i.MapID,
		&i.ScoreCt,
		&i.ScoreT,
		&i.GameTypeID,
	)
	return i, err
}

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

const createRound = `-- name: CreateRound :one
insert into rounds (start_date, start_time, match_id) values ($1, $2, $3) returning id, start_date, start_time, end_date, end_time, match_id, winner_team_id
`

type CreateRoundParams struct {
	StartDate pgtype.Date
	StartTime pgtype.Time
	MatchID   pgtype.Int8
}

func (q *Queries) CreateRound(ctx context.Context, arg CreateRoundParams) (Round, error) {
	row := q.db.QueryRow(ctx, createRound, arg.StartDate, arg.StartTime, arg.MatchID)
	var i Round
	err := row.Scan(
		&i.ID,
		&i.StartDate,
		&i.StartTime,
		&i.EndDate,
		&i.EndTime,
		&i.MatchID,
		&i.WinnerTeamID,
	)
	return i, err
}

const createRoundTeamEntry = `-- name: CreateRoundTeamEntry :one
insert into round_teams (player_id, team_id, round_id) values ($1, $2, $3) returning id, player_id, team_id, round_id
`

type CreateRoundTeamEntryParams struct {
	PlayerID pgtype.Int8
	TeamID   pgtype.Int8
	RoundID  pgtype.Int8
}

func (q *Queries) CreateRoundTeamEntry(ctx context.Context, arg CreateRoundTeamEntryParams) (RoundTeam, error) {
	row := q.db.QueryRow(ctx, createRoundTeamEntry, arg.PlayerID, arg.TeamID, arg.RoundID)
	var i RoundTeam
	err := row.Scan(
		&i.ID,
		&i.PlayerID,
		&i.TeamID,
		&i.RoundID,
	)
	return i, err
}

const createSpecialKill = `-- name: CreateSpecialKill :one
insert into special_kills (name) values ($1) returning id, name
`

func (q *Queries) CreateSpecialKill(ctx context.Context, name string) (SpecialKill, error) {
	row := q.db.QueryRow(ctx, createSpecialKill, name)
	var i SpecialKill
	err := row.Scan(&i.ID, &i.Name)
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

const createTeam = `-- name: CreateTeam :one
insert into teams (name) values ($1) returning id, name
`

func (q *Queries) CreateTeam(ctx context.Context, name string) (Team, error) {
	row := q.db.QueryRow(ctx, createTeam, name)
	var i Team
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const createTeamSwitchEvent = `-- name: CreateTeamSwitchEvent :exec
insert into team_switch (player_id, from_team_id, to_team_id, switch_date, switch_time, round_id) values ($1, $2, $3, $4, $5, $6) returning id, player_id, from_team_id, to_team_id, switch_date, switch_time, round_id
`

type CreateTeamSwitchEventParams struct {
	PlayerID   pgtype.Int8
	FromTeamID pgtype.Int8
	ToTeamID   pgtype.Int8
	SwitchDate pgtype.Date
	SwitchTime pgtype.Time
	RoundID    pgtype.Int8
}

func (q *Queries) CreateTeamSwitchEvent(ctx context.Context, arg CreateTeamSwitchEventParams) error {
	_, err := q.db.Exec(ctx, createTeamSwitchEvent,
		arg.PlayerID,
		arg.FromTeamID,
		arg.ToTeamID,
		arg.SwitchDate,
		arg.SwitchTime,
		arg.RoundID,
	)
	return err
}

const createWeapon = `-- name: CreateWeapon :one
insert into weapons (name) values ($1) returning id, name
`

func (q *Queries) CreateWeapon(ctx context.Context, name string) (Weapon, error) {
	row := q.db.QueryRow(ctx, createWeapon, name)
	var i Weapon
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getAttack = `-- name: GetAttack :one

select id, attacker_id, attacked_id, round_id, attack_time, attack_date, attacker_team_id, attacked_team_id, attacker_position, attacked_position, attacker_weapon_id, damage, damage_armor, health, armor, hit_group_id from attacks where id = $1 limit 1
`

// Attacks
func (q *Queries) GetAttack(ctx context.Context, id int64) (Attack, error) {
	row := q.db.QueryRow(ctx, getAttack, id)
	var i Attack
	err := row.Scan(
		&i.ID,
		&i.AttackerID,
		&i.AttackedID,
		&i.RoundID,
		&i.AttackTime,
		&i.AttackDate,
		&i.AttackerTeamID,
		&i.AttackedTeamID,
		&i.AttackerPosition,
		&i.AttackedPosition,
		&i.AttackerWeaponID,
		&i.Damage,
		&i.DamageArmor,
		&i.Health,
		&i.Armor,
		&i.HitGroupID,
	)
	return i, err
}

const getAttackByAttackerAttackedRoundDateTime = `-- name: GetAttackByAttackerAttackedRoundDateTime :one
select id, attacker_id, attacked_id, round_id, attack_time, attack_date, attacker_team_id, attacked_team_id, attacker_position, attacked_position, attacker_weapon_id, damage, damage_armor, health, armor, hit_group_id from attacks where attacker_id = $1 and attacked_id = $2 and round_id = $3 and attack_date = $4 and attack_time = $5 limit 1
`

type GetAttackByAttackerAttackedRoundDateTimeParams struct {
	AttackerID pgtype.Int8
	AttackedID pgtype.Int8
	RoundID    pgtype.Int8
	AttackDate pgtype.Date
	AttackTime pgtype.Time
}

func (q *Queries) GetAttackByAttackerAttackedRoundDateTime(ctx context.Context, arg GetAttackByAttackerAttackedRoundDateTimeParams) (Attack, error) {
	row := q.db.QueryRow(ctx, getAttackByAttackerAttackedRoundDateTime,
		arg.AttackerID,
		arg.AttackedID,
		arg.RoundID,
		arg.AttackDate,
		arg.AttackTime,
	)
	var i Attack
	err := row.Scan(
		&i.ID,
		&i.AttackerID,
		&i.AttackedID,
		&i.RoundID,
		&i.AttackTime,
		&i.AttackDate,
		&i.AttackerTeamID,
		&i.AttackedTeamID,
		&i.AttackerPosition,
		&i.AttackedPosition,
		&i.AttackerWeaponID,
		&i.Damage,
		&i.DamageArmor,
		&i.Health,
		&i.Armor,
		&i.HitGroupID,
	)
	return i, err
}

const getGameType = `-- name: GetGameType :one

select id, name from game_types where id = $1 limit 1
`

// Game Type
func (q *Queries) GetGameType(ctx context.Context, id int64) (GameType, error) {
	row := q.db.QueryRow(ctx, getGameType, id)
	var i GameType
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getGameTypeByName = `-- name: GetGameTypeByName :one
select id, name from game_types where name = $1 limit 1
`

func (q *Queries) GetGameTypeByName(ctx context.Context, name string) (GameType, error) {
	row := q.db.QueryRow(ctx, getGameTypeByName, name)
	var i GameType
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getHitGroup = `-- name: GetHitGroup :one

select id, name from hit_groups where id = $1 limit 1
`

// Hit Groups
func (q *Queries) GetHitGroup(ctx context.Context, id int64) (HitGroup, error) {
	row := q.db.QueryRow(ctx, getHitGroup, id)
	var i HitGroup
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getHitGroupByName = `-- name: GetHitGroupByName :one
select id, name from hit_groups where name = $1 limit 1
`

func (q *Queries) GetHitGroupByName(ctx context.Context, name string) (HitGroup, error) {
	row := q.db.QueryRow(ctx, getHitGroupByName, name)
	var i HitGroup
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getKill = `-- name: GetKill :one

select id, killer_id, killed_id, round_id, kill_time, kill_date, killer_team_id, killed_team_id, killer_position, killed_position, killer_weapon_id, special_id from kills where id = $1 limit 1
`

// Kills
func (q *Queries) GetKill(ctx context.Context, id int64) (Kill, error) {
	row := q.db.QueryRow(ctx, getKill, id)
	var i Kill
	err := row.Scan(
		&i.ID,
		&i.KillerID,
		&i.KilledID,
		&i.RoundID,
		&i.KillTime,
		&i.KillDate,
		&i.KillerTeamID,
		&i.KilledTeamID,
		&i.KillerPosition,
		&i.KilledPosition,
		&i.KillerWeaponID,
		&i.SpecialID,
	)
	return i, err
}

const getKillByKillerKilledRoundDateTime = `-- name: GetKillByKillerKilledRoundDateTime :one
select id, killer_id, killed_id, round_id, kill_time, kill_date, killer_team_id, killed_team_id, killer_position, killed_position, killer_weapon_id, special_id from kills where killer_id = $1 and killed_id = $2 and round_id = $3 and kill_date = $4 and kill_time = $5 limit 1
`

type GetKillByKillerKilledRoundDateTimeParams struct {
	KillerID pgtype.Int8
	KilledID pgtype.Int8
	RoundID  pgtype.Int8
	KillDate pgtype.Date
	KillTime pgtype.Time
}

func (q *Queries) GetKillByKillerKilledRoundDateTime(ctx context.Context, arg GetKillByKillerKilledRoundDateTimeParams) (Kill, error) {
	row := q.db.QueryRow(ctx, getKillByKillerKilledRoundDateTime,
		arg.KillerID,
		arg.KilledID,
		arg.RoundID,
		arg.KillDate,
		arg.KillTime,
	)
	var i Kill
	err := row.Scan(
		&i.ID,
		&i.KillerID,
		&i.KilledID,
		&i.RoundID,
		&i.KillTime,
		&i.KillDate,
		&i.KillerTeamID,
		&i.KilledTeamID,
		&i.KillerPosition,
		&i.KilledPosition,
		&i.KillerWeaponID,
		&i.SpecialID,
	)
	return i, err
}

const getMap = `-- name: GetMap :one

select id, name from maps where id = $1 limit 1
`

// Maps
func (q *Queries) GetMap(ctx context.Context, id int64) (Map, error) {
	row := q.db.QueryRow(ctx, getMap, id)
	var i Map
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getMapByName = `-- name: GetMapByName :one
select id, name from maps where name = $1 limit 1
`

func (q *Queries) GetMapByName(ctx context.Context, name string) (Map, error) {
	row := q.db.QueryRow(ctx, getMapByName, name)
	var i Map
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getMatch = `-- name: GetMatch :one

select id, start_date, start_time, end_date, end_time, map_id, score_ct, score_t, game_type_id from matches where id = $1 limit 1
`

// Matches
func (q *Queries) GetMatch(ctx context.Context, id int64) (Match, error) {
	row := q.db.QueryRow(ctx, getMatch, id)
	var i Match
	err := row.Scan(
		&i.ID,
		&i.StartDate,
		&i.StartTime,
		&i.EndDate,
		&i.EndTime,
		&i.MapID,
		&i.ScoreCt,
		&i.ScoreT,
		&i.GameTypeID,
	)
	return i, err
}

const getMatchByMapAndStartDateTime = `-- name: GetMatchByMapAndStartDateTime :one
select id, start_date, start_time, end_date, end_time, map_id, score_ct, score_t, game_type_id from matches where map_id = $1 and start_date = $2 and start_time = $3 limit 1
`

type GetMatchByMapAndStartDateTimeParams struct {
	MapID     pgtype.Int8
	StartDate pgtype.Date
	StartTime pgtype.Time
}

func (q *Queries) GetMatchByMapAndStartDateTime(ctx context.Context, arg GetMatchByMapAndStartDateTimeParams) (Match, error) {
	row := q.db.QueryRow(ctx, getMatchByMapAndStartDateTime, arg.MapID, arg.StartDate, arg.StartTime)
	var i Match
	err := row.Scan(
		&i.ID,
		&i.StartDate,
		&i.StartTime,
		&i.EndDate,
		&i.EndTime,
		&i.MapID,
		&i.ScoreCt,
		&i.ScoreT,
		&i.GameTypeID,
	)
	return i, err
}

const getPlayer = `-- name: GetPlayer :one

select id, steam_user_id, name, bot from players where id = $1 limit 1
`

// Players
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

const getRound = `-- name: GetRound :one

select id, start_date, start_time, end_date, end_time, match_id, winner_team_id from rounds where id = $1 limit 1
`

// Rounds
func (q *Queries) GetRound(ctx context.Context, id int64) (Round, error) {
	row := q.db.QueryRow(ctx, getRound, id)
	var i Round
	err := row.Scan(
		&i.ID,
		&i.StartDate,
		&i.StartTime,
		&i.EndDate,
		&i.EndTime,
		&i.MatchID,
		&i.WinnerTeamID,
	)
	return i, err
}

const getRoundByStartDateTime = `-- name: GetRoundByStartDateTime :one
select id, start_date, start_time, end_date, end_time, match_id, winner_team_id from rounds where start_date = $1 and start_time = $2 limit 1
`

type GetRoundByStartDateTimeParams struct {
	StartDate pgtype.Date
	StartTime pgtype.Time
}

func (q *Queries) GetRoundByStartDateTime(ctx context.Context, arg GetRoundByStartDateTimeParams) (Round, error) {
	row := q.db.QueryRow(ctx, getRoundByStartDateTime, arg.StartDate, arg.StartTime)
	var i Round
	err := row.Scan(
		&i.ID,
		&i.StartDate,
		&i.StartTime,
		&i.EndDate,
		&i.EndTime,
		&i.MatchID,
		&i.WinnerTeamID,
	)
	return i, err
}

const getRoundTeamEntry = `-- name: GetRoundTeamEntry :one

select id, player_id, team_id, round_id from round_teams where id = $1 limit 1
`

// Round Teams
func (q *Queries) GetRoundTeamEntry(ctx context.Context, id int64) (RoundTeam, error) {
	row := q.db.QueryRow(ctx, getRoundTeamEntry, id)
	var i RoundTeam
	err := row.Scan(
		&i.ID,
		&i.PlayerID,
		&i.TeamID,
		&i.RoundID,
	)
	return i, err
}

const getRoundTeamEntryByPlayerTeamRound = `-- name: GetRoundTeamEntryByPlayerTeamRound :one
select id, player_id, team_id, round_id from round_teams where player_id = $1 and team_id = $2 and round_id = $3 limit 1
`

type GetRoundTeamEntryByPlayerTeamRoundParams struct {
	PlayerID pgtype.Int8
	TeamID   pgtype.Int8
	RoundID  pgtype.Int8
}

func (q *Queries) GetRoundTeamEntryByPlayerTeamRound(ctx context.Context, arg GetRoundTeamEntryByPlayerTeamRoundParams) (RoundTeam, error) {
	row := q.db.QueryRow(ctx, getRoundTeamEntryByPlayerTeamRound, arg.PlayerID, arg.TeamID, arg.RoundID)
	var i RoundTeam
	err := row.Scan(
		&i.ID,
		&i.PlayerID,
		&i.TeamID,
		&i.RoundID,
	)
	return i, err
}

const getSpecialKill = `-- name: GetSpecialKill :one

select id, name from special_kills where id = $1 limit 1
`

// Special Kills
func (q *Queries) GetSpecialKill(ctx context.Context, id int64) (SpecialKill, error) {
	row := q.db.QueryRow(ctx, getSpecialKill, id)
	var i SpecialKill
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getSpecialKillByName = `-- name: GetSpecialKillByName :one
select id, name from special_kills where name = $1 limit 1
`

func (q *Queries) GetSpecialKillByName(ctx context.Context, name string) (SpecialKill, error) {
	row := q.db.QueryRow(ctx, getSpecialKillByName, name)
	var i SpecialKill
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getSteamUser = `-- name: GetSteamUser :one

select id, steam_id, steam_community_id from steam_users where id = $1 limit 1
`

// Steam User
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

const getTeam = `-- name: GetTeam :one

select id, name from teams where id = $1 limit 1
`

// Teams
func (q *Queries) GetTeam(ctx context.Context, id int64) (Team, error) {
	row := q.db.QueryRow(ctx, getTeam, id)
	var i Team
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getTeamByName = `-- name: GetTeamByName :one
select id, name from teams where name = $1 limit 1
`

func (q *Queries) GetTeamByName(ctx context.Context, name string) (Team, error) {
	row := q.db.QueryRow(ctx, getTeamByName, name)
	var i Team
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getTeamSwitchByPlayerAndDateTime = `-- name: GetTeamSwitchByPlayerAndDateTime :one

select id, player_id, from_team_id, to_team_id, switch_date, switch_time, round_id from team_switch where player_id = $1 and switch_date = $2 and switch_time = $3 limit 1
`

type GetTeamSwitchByPlayerAndDateTimeParams struct {
	PlayerID   pgtype.Int8
	SwitchDate pgtype.Date
	SwitchTime pgtype.Time
}

// Team Switch
func (q *Queries) GetTeamSwitchByPlayerAndDateTime(ctx context.Context, arg GetTeamSwitchByPlayerAndDateTimeParams) (TeamSwitch, error) {
	row := q.db.QueryRow(ctx, getTeamSwitchByPlayerAndDateTime, arg.PlayerID, arg.SwitchDate, arg.SwitchTime)
	var i TeamSwitch
	err := row.Scan(
		&i.ID,
		&i.PlayerID,
		&i.FromTeamID,
		&i.ToTeamID,
		&i.SwitchDate,
		&i.SwitchTime,
		&i.RoundID,
	)
	return i, err
}

const getWeapon = `-- name: GetWeapon :one

select id, name from weapons where id = $1 limit 1
`

// Weapons
func (q *Queries) GetWeapon(ctx context.Context, id int64) (Weapon, error) {
	row := q.db.QueryRow(ctx, getWeapon, id)
	var i Weapon
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getWeaponByName = `-- name: GetWeaponByName :one
select id, name from weapons where name = $1 limit 1
`

func (q *Queries) GetWeaponByName(ctx context.Context, name string) (Weapon, error) {
	row := q.db.QueryRow(ctx, getWeaponByName, name)
	var i Weapon
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const updateMatch = `-- name: UpdateMatch :exec
update matches set end_date=$2, end_time=$3, score_ct=$4, score_t=$5, game_type_id=$6 where id=$1
`

type UpdateMatchParams struct {
	ID         int64
	EndDate    pgtype.Date
	EndTime    pgtype.Time
	ScoreCt    int32
	ScoreT     int32
	GameTypeID pgtype.Int8
}

func (q *Queries) UpdateMatch(ctx context.Context, arg UpdateMatchParams) error {
	_, err := q.db.Exec(ctx, updateMatch,
		arg.ID,
		arg.EndDate,
		arg.EndTime,
		arg.ScoreCt,
		arg.ScoreT,
		arg.GameTypeID,
	)
	return err
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

const updateRound = `-- name: UpdateRound :exec
update rounds set end_date=$2, end_time=$3, winner_team_id=$4 where id = $1
`

type UpdateRoundParams struct {
	ID           int64
	EndDate      pgtype.Date
	EndTime      pgtype.Time
	WinnerTeamID pgtype.Int8
}

func (q *Queries) UpdateRound(ctx context.Context, arg UpdateRoundParams) error {
	_, err := q.db.Exec(ctx, updateRound,
		arg.ID,
		arg.EndDate,
		arg.EndTime,
		arg.WinnerTeamID,
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
