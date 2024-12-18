-- --------------------
-- Steam User
-- --------------------

-- name: GetSteamUser :one
select * from steam_users where id = $1 limit 1;

-- name: GetSteamUserBySteamId :one
select * from steam_users where steam_id = $1;

-- name: CreateSteamUser :one
insert into steam_users (steam_id, steam_community_id) values ($1, $2) returning *;

-- name: UpdateSteamUser :exec
update steam_users set steam_id = $2, steam_community_id = $3 where id = $1;

-- name: ListSteamUsers :many
select * from steam_users order by steam_community_id;

-- --------------------
-- Players
-- --------------------

-- name: GetPlayer :one
select * from players where id = $1 limit 1;

-- name: GetPlayerByName :one
select * from players where name = $1 limit 1;

-- name: CreatePlayer :one
insert into players (steam_user_id, name, bot) values ($1, $2, $3) returning *;

-- name: UpdatePlayer :exec
update players set steam_user_id = $2, name = $3, bot = $4 where id = $1;

-- name: ListPlayers :many
select p.id as id, p.steam_user_id as steam_user_id, p.name as name, su.steam_id as steam_id, su.steam_community_id as steam_community_id from players p, steam_users su where p.bot=false and p.steam_user_id=su.id order by su.steam_community_id;

-- name: ListBots :many
select p.id as id, p.steam_user_id as steam_user_id, p.name as name, su.steam_id as steam_id, su.steam_community_id as steam_community_id from players p, steam_users su where p.bot=true and p.steam_user_id=su.id order by su.steam_community_id;

-- name: ListAllPlayers :many
select p.id as id, p.steam_user_id as steam_user_id, p.name as name, su.steam_id as steam_id, su.steam_community_id as steam_community_id from players p, steam_users su where p.steam_user_id=su.id order by su.steam_community_id;

-- name: GetPlayerWithStats :one
select p.id, p.steam_user_id, p.name, su.steam_id, su.steam_community_id,
       (select count(*) from round_teams rt, teams t where rt.player_id=p.id and rt.team_id=t.id and t.name in ('CT', 'TERRORIST')) as rounds,
       (select count(*) from round_teams rt, teams t where rt.player_id=p.id and rt.team_id=t.id and t.name='TERRORIST') as rounds_t,
       (select count(*) from round_teams rt, teams t where rt.player_id=p.id and rt.team_id=t.id and t.name='CT') as rounds_ct,
       (select count(*) from kills where killer_id=p.id) as kills,
       (select count(*) from kills where killed_id=p.id) as killed,
       (select count(*) from kills k, special_kills sk where k.killer_id=p.id and k.special_id=sk.id and sk.name='headshot') as headshot_kills,
       (select count(*) from kills k, special_kills sk where k.killer_id=p.id and k.special_id=sk.id and sk.name='noscope') as noscope_kills,
       (select count(*) from kills k, special_kills sk where k.killer_id=p.id and k.special_id=sk.id and sk.name='throughsmoke') as throughsmoke_kills,
       (select count(*) from player_suicide where player_id=p.id) as suicides,
       (select count(*) from triggered_events te, events e where te.player_id=p.id and te.event_id=e.id and e.name='Planted_The_Bomb') as bomb_plants,
       (select count(*) from triggered_events te, events e where te.player_id=p.id and te.event_id=e.id and e.name='Defused_The_Bomb') as bomb_defuses
from players p, steam_users su where p.id=$1 and p.steam_user_id=su.id limit 1;

-- --------------------
-- Maps
-- --------------------

-- name: GetMap :one
select * from maps where id = $1 limit 1;

-- name: GetMapByName :one
select * from maps where name = $1 limit 1;

-- name: CreateMap :one
insert into maps (name) values ($1) returning *;

-- name: ListMaps :many
select * from maps order by name;

-- --------------------
-- Teams
-- --------------------

-- name: GetTeam :one
select * from teams where id = $1 limit 1;

-- name: GetTeamByName :one
select * from teams where name = $1 limit 1;

-- name: CreateTeam :one
insert into teams (name) values ($1) returning *;

-- name: ListTeams :many
select * from teams order by name;

-- --------------------
-- Weapons
-- --------------------

-- name: GetWeapon :one
select * from weapons where id = $1 limit 1;

-- name: GetWeaponByName :one
select * from weapons where name = $1 limit 1;

-- name: CreateWeapon :one
insert into weapons (name) values ($1) returning *;

-- --------------------
-- Special Kills
-- --------------------

-- name: GetSpecialKill :one
select * from special_kills where id = $1 limit 1;

-- name: GetSpecialKillByName :one
select * from special_kills where name = $1 limit 1;

-- name: CreateSpecialKill :one
insert into special_kills (name) values ($1) returning *;

-- --------------------
-- Game Type
-- --------------------

-- name: GetGameType :one
select * from game_types where id = $1 limit 1;

-- name: GetGameTypeByName :one
select * from game_types where name = $1 limit 1;

-- name: CreateGameType :one
insert into game_types (name) values ($1) returning *;

-- --------------------
-- Hit Groups
-- --------------------

-- name: GetHitGroup :one
select * from hit_groups where id = $1 limit 1;

-- name: GetHitGroupByName :one
select * from hit_groups where name = $1 limit 1;

-- name: CreateHitGroup :one
insert into hit_groups (name) values ($1) returning *;

-- --------------------
-- Other Kills
-- --------------------

-- name: GetOtherKill :one
select * from other_kills where id = $1 limit 1;

-- name: GetOtherKillByName :one
select * from other_kills where name = $1 limit 1;

-- name: CreateOtherKill :one
insert into other_kills (name) values ($1) returning *;

-- --------------------
-- Items
-- --------------------

-- name: GetItem :one
select * from items where id = $1 limit 1;

-- name: GetItemByName :one
select * from items where name = $1 limit 1;

-- name: CreateItem :one
insert into items (name) values ($1) returning *;

-- --------------------
-- Item Actions
-- --------------------

-- name: GetItemAction :one
select * from item_actions where id = $1 limit 1;

-- name: GetItemActionByName :one
select * from item_actions where name = $1 limit 1;

-- name: CreateItemAction :one
insert into item_actions (name) values ($1) returning *;

-- --------------------
-- Events
-- --------------------

-- name: GetEvent :one
select * from events where id = $1 limit 1;

-- name: GetEventByName :one
select * from events where name = $1 limit 1;

-- name: CreateEvent :one
insert into events (name) values ($1) returning *;

-- --------------------
-- Matches
-- --------------------

-- name: GetMatch :one
select * from matches where id = $1 limit 1;

-- name: GetMatchByMapAndStartDateTime :one
select * from matches where map_id = $1 and start_date = $2 and start_time = $3 limit 1;

-- name: CreateMatch :one
insert into matches (start_date, start_time, map_id) values ($1, $2, $3) returning *;

-- name: UpdateMatch :exec
update matches set end_date=$2, end_time=$3, score_ct=$4, score_t=$5, game_type_id=$6 where id=$1;

-- --------------------
-- Rounds
-- --------------------

-- name: GetRound :one
select * from rounds where id = $1 limit 1;

-- name: GetRoundByStartDateTime :one
select * from rounds where start_date = $1 and start_time = $2 limit 1;

-- name: CreateRound :one
insert into rounds (start_date, start_time, match_id) values ($1, $2, $3) returning *;

-- name: UpdateRound :exec
update rounds set end_date=$2, end_time=$3, winner_team_id=$4 where id = $1;

-- --------------------
-- Team Switch
-- --------------------

-- name: GetTeamSwitch :one
select * from team_switch where id = $1 limit 1;

-- name: GetTeamSwitchByPlayerAndDateTime :one
select * from team_switch where player_id = $1 and switch_date = $2 and switch_time = $3 limit 1;

-- name: CreateTeamSwitchEvent :exec
insert into team_switch (player_id, from_team_id, to_team_id, switch_date, switch_time, round_id) values ($1, $2, $3, $4, $5, $6) returning *;

-- --------------------
-- Round Teams
-- --------------------

-- name: GetRoundTeamEntry :one
select * from round_teams where id = $1 limit 1;

-- name: GetRoundTeamEntryByPlayerTeamRound :one
select * from round_teams where player_id = $1 and team_id = $2 and round_id = $3 limit 1;

-- name: CreateRoundTeamEntry :one
insert into round_teams (player_id, team_id, round_id) values ($1, $2, $3) returning *;

-- --------------------
-- Kills
-- --------------------

-- name: GetKill :one
select * from kills where id = $1 limit 1;

-- name: GetKillByKillerKilledRoundDateTime :one
select * from kills where killer_id = $1 and killed_id = $2 and round_id = $3 and kill_date = $4 and kill_time = $5 limit 1;

-- name: CreateKill :one
insert into kills (killer_id, killed_id, round_id, kill_time, kill_date, killer_team_id, killed_team_id, killer_position_x, killer_position_y, killer_position_z, killed_position_x, killed_position_y, killed_position_z, killer_weapon_id, special_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15) returning *;

-- --------------------
-- Attacks
-- --------------------

-- name: GetAttack :one
select * from attacks where id = $1 limit 1;

-- name: GetAttackByAttackerAttackedRoundDateTime :one
select * from attacks where attacker_id = $1 and attacked_id = $2 and round_id = $3 and attack_date = $4 and attack_time = $5 limit 1;

-- name: CreateAttack :one
insert into attacks (attacker_id, attacked_id, round_id, attack_time, attack_date, attacker_team_id, attacked_team_id, attacker_position_x, attacker_position_y, attacker_position_z, attacked_position_x, attacked_position_y, attacked_position_z, attacker_weapon_id, damage, damage_armor, health, armor, hit_group_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19) returning *;

-- --------------------
-- Killed Other
-- --------------------

-- name: GetKillOther :one
select * from kills_other where id = $1 limit 1;

-- name: GetKillOtherByKillerOtherRoundDateTime :one
select * from kills_other where killer_id = $1 and killed_other_id = $2 and round_id = $3 and kill_date = $4 and kill_time = $5 limit 1;

-- name: CreateKillOther :one
insert into kills_other (killer_id, killed_other_id, round_id, kill_time, kill_date, killer_team_id, killer_position_x, killer_position_y, killer_position_z, killed_position_x, killed_position_y, killed_position_z, killer_weapon_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) returning *;

-- --------------------
-- Kills Assisted
-- --------------------

-- name: GetKillAssisted :one
select * from kills_assisted where id = $1 limit 1;

-- name: GetKillAssistedByKillerKilledRoundDateTime :one
select * from kills_assisted where killer_id = $1 and killed_id = $2 and round_id = $3 and kill_date = $4 and kill_time = $5 limit 1;

-- name: CreateKillAssisted :one
insert into kills_assisted (killer_id, killed_id, round_id, kill_time, kill_date, killer_team_id, killed_team_id) values ($1, $2, $3, $4, $5, $6, $7) returning *;

-- --------------------
-- Item Interactions
-- --------------------

-- name: GetItemInteraction :one
select * from item_interactions where id = $1 limit 1;

-- name: GetItemInteractionByPlayerItemInteractionRoundDateTime :one
select * from item_interactions where player_id = $1 and item_id = $2 and item_action = $3 and round_id = $4 and interaction_date = $5 and interaction_time = $6 limit 1;

-- name: CreateItemInteraction :one
insert into item_interactions (player_id, team_id, round_id, item_id, item_action, interaction_time, interaction_date) values ($1, $2, $3, $4, $5, $6, $7) returning *;

-- --------------------
-- Money Change
-- --------------------

-- name: GetMoneyChange :one
select * from money_change where id = $1 limit 1;

-- name: GetMoneyChangeByPlayerNewTotalRoundDateTime :one
select * from money_change where player_id = $1 and new_total = $2 and round_id = $3 and change_date = $4 and change_time = $5 limit 1;

-- name: CreateMoneyChange :one
insert into money_change (player_id, team_id, round_id, item_id, new_total, change_time, change_date) values ($1, $2, $3, $4, $5, $6, $7) returning *;

-- --------------------
-- Player Suicide
-- --------------------

-- name: GetPlayerSuicide :one
select * from player_suicide where id = $1 limit 1;

-- name: GetPlayerSuicideByPlayerItemRoundDateTime :one
select * from player_suicide where player_id = $1 and with_item_id = $2 and round_id = $3 and suicide_date = $4 and suicide_time = $5 limit 1;

-- name: CreatePlayerSuicide :one
insert into player_suicide (player_id, round_id, suicide_time, suicide_date, team_id, player_position_x, player_position_y, player_position_z, with_item_id) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning *;

-- --------------------
-- Triggered Events
-- --------------------

-- name: GetTriggeredEvent :one
select * from triggered_events where id = $1 limit 1;

-- name: GetTriggeredEventByPlayerEventRoundDateTime :one
select * from triggered_events where player_id = $1 and event_id = $2 and round_id = $3 and event_date = $4 and event_time = $5 limit 1;

-- name: CreateTriggeredEvent :one
insert into triggered_events (player_id, team_id, round_id, event_id, event_time, event_date, bombsite) values ($1, $2, $3, $4, $5, $6, $7) returning *;

-- --------------------
-- Player Threw
-- --------------------

-- name: GetThrew :one
select * from threw where id = $1 limit 1;

-- name: GetThrewByPlayerWeaponRoundDateTime :one
select * from threw where player_id = $1 and weapon_id = $2 and round_id = $3 and threw_date = $4 and threw_time = $5 limit 1;

-- name: CreateThrew :one
insert into threw (player_id, team_id, round_id, threw_time, threw_date, position_x, position_y, position_z, weapon_id, entindex) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning *;

-- --------------------
-- Player Blinded
-- --------------------

-- name: GetBlinded :one
select * from blinded where id = $1 limit 1;

-- name: GetBlindedByPlayerRoundDateTime :one
select * from blinded where blinded_id = $1 and round_id = $2 and blinded_date = $3 and blinded_time= $4 limit 1;

-- name: CreateBlinded :one
insert into blinded (blinded_id, blinded_team_id, blinded_by_id, blinded_by_team_id, round_id, blinded_date, blinded_time, blinded_for, entindex) values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning *;

-- --------------------
-- Accolade
-- --------------------

-- name: GetAccolade :one
select * from accolade where id = $1 limit 1;

-- name: GetAccoladeByNamePlayerMatchDateTime :one
select * from accolade where accolade_name = $1 and player_id = $2 and match_id = $3 and accolade_date = $4 and accolade_time = $5 limit 1;

-- name: CreateAccolade :one
insert into accolade (player_id, match_id, accolade_date, accolade_time, accolade_name, accolade_value, accolade_pos, accolade_score) values ($1, $2, $3, $4, $5, $6, $7, $8) returning *;

-- name: GetAccoladeForPlayer :many
select p.id, p.steam_user_id, p.name, su.steam_id, su.steam_community_id, a.accolade_time, a.accolade_date, a.accolade_name, a.accolade_value, a.accolade_score, a.accolade_pos, apn.pretty_name, apn.description from players p, steam_users su, accolade a, accolade_pretty_names apn where p.id=$1 and p.steam_user_id=su.id and a.player_id=p.id and a.accolade_name=apn.accolade_name  order by a.accolade_date;
