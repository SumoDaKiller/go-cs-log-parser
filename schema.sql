-- grant all privileges on database gocslogparser to cslogparser;
-- grant all on schema public to cslogparser;

CREATE TABLE IF NOT EXISTS steam_users (
    id BIGSERIAL PRIMARY KEY,
    steam_id varchar(50) not null unique,
    steam_community_id BIGINT not null
);

CREATE TABLE IF NOT EXISTS players (
    id BIGSERIAL PRIMARY KEY,
    steam_user_id BIGINT references steam_users(id),
    name varchar(255) not null unique,
    bot boolean default false
);

CREATE TABLE IF NOT EXISTS maps (
    id BIGSERIAL PRIMARY KEY,
    name varchar(255) not null unique
);

CREATE TABLE IF NOT EXISTS teams (
    id BIGSERIAL PRIMARY KEY,
    name varchar(20) not null unique
);

CREATE TABLE IF NOT EXISTS weapons (
    id BIGSERIAL PRIMARY KEY,
    name varchar(255) not null unique
);

CREATE TABLE IF NOT EXISTS special_kills (
    id BIGSERIAL PRIMARY KEY,
    name varchar(255) not null unique
);

CREATE TABLE IF NOT EXISTS game_types (
    id BIGSERIAL PRIMARY KEY,
    name varchar(50) not null unique
);

CREATE TABLE IF NOT EXISTS hit_groups (
    id BIGSERIAL PRIMARY KEY,
    name varchar(255) not null unique
);

CREATE TABLE IF NOT EXISTS matches (
    id BIGSERIAL PRIMARY KEY,
    start_date date not null,
    start_time time not null,
    end_date date,
    end_time time,
    map_id BIGINT references maps(id),
    score_ct int not null default 0,
    score_t int not null default 0,
    game_type_id BIGINT references game_types(id)
);

CREATE TABLE IF NOT EXISTS rounds (
    id BIGSERIAL PRIMARY KEY,
    start_date date not null,
    start_time time not null,
    end_date date,
    end_time time,
    match_id BIGINT references matches(id),
    winner_team_id BIGINT references teams(id)
);

CREATE TABLE IF NOT EXISTS team_switch (
    id BIGSERIAL PRIMARY KEY,
    player_id BIGINT references players(id),
    from_team_id BIGINT references teams(id),
    to_team_id BIGINT references teams(id),
    switch_date date not null,
    switch_time time not null,
    round_id BIGINT references rounds(id)
);

CREATE TABLE IF NOT EXISTS round_teams (
    id BIGSERIAL PRIMARY KEY,
    player_id BIGINT references players(id),
    team_id BIGINT references teams(id),
    round_id BIGINT references rounds(id)
);

CREATE TABLE IF NOT EXISTS kills (
    id BIGSERIAL PRIMARY KEY,
    killer_id BIGINT references players(id),
    killed_id BIGINT references players(id),
    round_id BIGINT references rounds(id),
    kill_time time not null,
    kill_date date not null,
    killer_team_id BIGINT references teams(id),
    killed_team_id BIGINT references teams(id),
    killer_position json not null,
    killed_position json not null,
    killer_weapon_id BIGINT references weapons(id),
    special_id BIGINT references special_kills(id)
);

CREATE TABLE IF NOT EXISTS attacks (
    id BIGSERIAL PRIMARY KEY,
    attacker_id BIGINT references players(id),
    attacked_id BIGINT references players(id),
    round_id BIGINT references players(id),
    attack_time time not null,
    attack_date date not null,
    attacker_team_id BIGINT references teams(id),
    attacked_team_id BIGINT references teams(id),
    attacker_position json not null,
    attacked_position json not null,
    attacker_weapon_id BIGINT references weapons(id),
    damage INT not null,
    damage_armor INT not null,
    health INT not null,
    armor INT not null,
    hit_group_id BIGINT references hit_groups(id)
);
