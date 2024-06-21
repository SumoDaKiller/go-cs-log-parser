-- grant all privileges on database gocslogparser to cslogparser;
-- grant all on schema public to cslogparser;

CREATE TABLE steam_users (
    id BIGSERIAL PRIMARY KEY,
    steam_id varchar(50) not null,
    steam_community_id BIGINT not null
);

CREATE TABLE players (
    id BIGSERIAL PRIMARY KEY,
    steam_user_id BIGINT references steam_users(id),
    name varchar(255) not null,
    bot boolean default false
);

alter table steam_users add unique (steam_id);
alter table players add unique (name);
