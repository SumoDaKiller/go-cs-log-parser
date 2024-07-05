package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go-cs-log-parser/database"
	"io"
	"log"
	"regexp"
	"strconv"
)

func initializeRegexpPatterns() map[string]*regexp.Regexp {
	re := make(map[string]*regexp.Regexp)
	re["matchStart"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): World triggered "Match_Start" on "(?P<map>\w+)"$`)
	re["matchStatus"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): MatchStatus: Score: (?P<ctScore>\d+):(?P<tScore>\d+) on map "(?P<map>\w+)" RoundsPlayed: (?P<roundsPlayed>\d+)$`)
	// L 10/25/2022 - 19:18:09: MatchStatus: Team "CT" is unset.
	// L 10/25/2022 - 19:18:09: MatchStatus: Team "TERRORIST" is unset.
	// L 10/25/2022 - 19:18:09: MatchStatus: Score: 0:0 on map "de_anubis" RoundsPlayed: 0
	re["roundStart"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): World triggered "Round_Start"$`)
	re["roundEnd"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): World triggered "Round_End"$`)
	// L 10/25/2022 - 19:19:13: Team "TERRORIST" triggered "SFUI_Notice_Terrorists_Win" (CT "0") (T "1")
	// L 10/25/2022 - 19:19:13: Team "CT" scored "0" with "5" players
	// L 10/25/2022 - 19:19:13: Team "TERRORIST" scored "1" with "6" players
	// L 10/25/2022 - 19:19:13: MatchStatus: Team "CT" is unset.
	// L 10/25/2022 - 19:19:13: MatchStatus: Team "TERRORIST" is unset.
	// L 10/25/2022 - 19:19:13: MatchStatus: Score: 0:1 on map "de_anubis" RoundsPlayed: 1
	// L 10/25/2022 - 19:19:13: World triggered "Round_End"
	re["gameOver"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): Game Over: (?P<gameType>\w+) (?P<mapPool>\w+) (?P<map>\w+) score (?P<ctScore>\d+):(?P<tScore>\d+) after (?P<mapTime>\d+) min$`)
	re["switchTeam"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): "(?P<playerName>[^<]+)<\d+><(?P<steamId>BOT|STEAM[^>]+)>(<(?P<currentTeam>CT|TERRORIST|Unassigned)>)?" switched from team <(?P<fromTeam>CT|TERRORIST|Unassigned)> to <(?P<toTeam>CT|TERRORIST|Unassigned)>$`)
	re["killed"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): "(?P<killerName>[^<]+)<\d+><(?P<killerSteamId>BOT|STEAM[^>]+)><(?P<killerTeam>CT|TERRORIST)>" \[(?P<killerXYZ>-?\d+ -?\d+ -?\d+)] killed "(?P<killedName>[^<]+)<\d+><(?P<killedSteamId>BOT|STEAM[^>]+)><(?P<killedTeam>CT|TERRORIST)>" \[(?P<killedXYZ>-?\d+ -?\d+ -?\d+)] with "(?P<killerWeapon>\w+)"\s?\(?(?P<special>\w*)\)?$`)
	// L 10/25/2022 - 19:19:13: "Sumo da Killer<15><STEAM_1:0:402610><TERRORIST>" assisted killing "MR.BROKEN<27><STEAM_1:0:513509837><CT>"
	// L 10/25/2022 - 19:19:39: "MR.BROKEN<27><STEAM_1:0:513509837><CT>" threw flashbang [-1838 47 916] flashbang entindex 221)
	// L 10/25/2022 - 19:19:39: "Norm<68><BOT><TERRORIST>" blinded for 0.19 by "MR.BROKEN<27><STEAM_1:0:513509837><CT>" from flashbang entindex 221
	// L 10/25/2022 - 19:20:02: "Ziduac<2><STEAM_1:0:202544><CT>" [262 919 -31] killed other "chicken<157>" [234 935 -32] with "knife"
	re["shopping"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): "(?P<playerName>[^<]+)<\d+><(?P<steamId>BOT|STEAM[^>]+)><(?P<currentTeam>CT|TERRORIST)>" (?P<type>dropped|purchased|picked up) "(?P<item>[^"])"$`)
	re["moneyChange"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): "(?P<playerName>[^<]+)<\d+><(?P<steamId>BOT|STEAM[^>]+)><(?P<currentTeam>CT|TERRORIST)>" money change [^=]= \$(?P<newTotal>\d+) \(tracked\)\s?\(?[^:]*:?\s?(?P<item>[^)]*)\)?$`)
	re["leftBuyZone"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): "(?P<playerName>[^<]+)<\d+><(?P<steamId>BOT|STEAM[^>]+)><(?P<currentTeam>CT|TERRORIST)>" left buyzone with \[(?P<items>[^]])]$`)
	re["attacking"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): "(?P<attackerName>[^<]+)<\d+><(?P<attackerSteamId>BOT|STEAM[^>]+)><(?P<attackerTeam>CT|TERRORIST)>" \[(?P<attackerXYZ>-?\d+ -?\d+ -?\d+)] attacked "(?P<attackedName>[^<]+)<\d+><(?P<attackedSteamId>BOT|STEAM[^>]+)><(?P<attackedTeam>CT|TERRORIST)>" \[(?P<attackedXYZ>-?\d+ -?\d+ -?\d+)] with "(?P<weapon>[^"]+)" \(damage "(?P<damage>\d+)"\) \(damage_armor "(?P<damageArmor>\d+)"\) \(health "(?P<health>\d+)"\) \(armor "(?P<armor>\d+)"\) \(hitgroup "(?P<hitgroup>[^"]+)"\)$`)
	re["suicide"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): "(?P<playerName>[^<]+)<\d+><(?P<steamId>BOT|STEAM[^>]+)><(?P<currentTeam>CT|TERRORIST)>" \[(?P<suicideXYZ>-?\d+ -?\d+ -?\d+)] committed suicide with "(?P<item>[^"]+)"$`)
	re["disconnected"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): "(?P<playerName>[^<]+)<\d+><(?P<steamId>BOT|STEAM[^>]+)><(?P<currentTeam>CT|TERRORIST)>" disconnected \(reason "(?P<reason>[^"]+)"\)$`)
	re["accolade"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): ACCOLADE, FINAL: \{(?P<name>\w+)},\s+(?P<playerName>[^<]+)<\d+>,\s+VALUE: (?P<value>\d+\.\d+),\s+POS:\s+(?P<position>[^,]+),\s+SCORE: (?P<score>\d+\.\d+)$`)
	re["triggered"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): "(?P<playerName>[^<]+)<\d+><(?P<steamId>BOT|STEAM[^>]+)><(?P<currentTeam>CT|TERRORIST)>" triggered "(?P<event>[^"]+)"$`)
	return re
}

func calculateSteamCommunityId(steamId string) (int64, error) {
	re := regexp.MustCompile(`^STEAM_(?P<X>\d{1}):(?P<Y>\d{1}):(?P<Z>\d+)`)
	matches := re.FindStringSubmatch(steamId)
	if matches == nil || len(matches) < 1 {
		return 0, errors.New("steam id does not match regex")
	}
	yIdx := re.SubexpIndex("Y")
	zIdx := re.SubexpIndex("Z")
	// https://developer.valvesoftware.com/wiki/SteamID
	// STEAM_X:Y:Z - V = 0x0110000100000000
	// Steam Community ID =Z*2+V+Y
	v := int64(76561197960265728)
	y, _ := strconv.ParseInt(matches[yIdx], 10, 64)
	z, _ := strconv.ParseInt(matches[zIdx], 10, 64)
	fmt.Println(steamId, y, z, v)
	return z*int64(2) + v + y, nil
}

func parseFile(r io.Reader, re map[string]*regexp.Regexp) error {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, k.String("database_url"))
	if err != nil {
		log.Fatalf("error connecting to database: %v\n", err)
	}
	defer conn.Close(ctx)

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		for key, rgx := range re {
			if rgx.MatchString(line) {
				switch key {
				case "matchStart":
					fmt.Println("matchStart")
				case "switchTeam":
					matches := rgx.FindStringSubmatch(line)
					dateIdx := rgx.SubexpIndex("date")
					timeIdx := rgx.SubexpIndex("time")
					playerNameIdx := rgx.SubexpIndex("playerName")
					steamIdIdx := rgx.SubexpIndex("steamId")

					err = handleSwitchTeam(ctx, conn, matches[dateIdx], matches[timeIdx], matches[playerNameIdx], matches[steamIdIdx], "", "", "")
					if err != nil {
						log.Println("error handling switchTeam: ", err)
					}
				}
			}
		}
	}
	return nil
}

func handleSwitchTeam(ctx context.Context, conn *pgx.Conn, date string, time string, playerName string, steamId string, currentTeam string, fromTeam string, toTeam string) error {
	fmt.Println(date, time, playerName, steamId, currentTeam, fromTeam, toTeam)
	queries := database.New(conn)
	steamUser, err := queries.GetSteamUserBySteamId(ctx, steamId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// First time we see this steam user so create it in the database
			steamCommunityID := int64(0)
			if steamId != "BOT" {
				steamCommunityID, err = calculateSteamCommunityId(steamId)
				if err != nil {
					return err
				}
			}
			steamUser, err = queries.CreateSteamUser(ctx, database.CreateSteamUserParams{SteamID: steamId, SteamCommunityID: steamCommunityID})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	player, err := queries.GetPlayerByName(ctx, playerName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// First time we see this player name so create it in the database
			bot := false
			if steamId == "BOT" {
				bot = true
			}
			player, err = queries.CreatePlayer(ctx, database.CreatePlayerParams{Name: playerName, SteamUserID: pgtype.Int8{Int64: steamUser.ID, Valid: true}, Bot: pgtype.Bool{Bool: bot, Valid: true}})
		} else {
			return err
		}
	}
	fmt.Println(steamUser)
	fmt.Println(player)
	return nil
}
