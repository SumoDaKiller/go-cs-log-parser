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
	"time"
)

func initializeRegexpPatterns() map[string]*regexp.Regexp {
	re := make(map[string]*regexp.Regexp)
	re["matchStart"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): World triggered "Match_Start" on "(?P<map>\w+)"$`)
	re["matchStatus"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): MatchStatus: Score: (?P<ctScore>\d+):(?P<tScore>\d+) on map "(?P<map>\w+)" RoundsPlayed: (?P<roundsPlayed>\d+)$`)
	re["roundStart"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): World triggered "Round_Start"$`)
	re["roundEnd"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): World triggered "Round_End"$`)
	re["sfuiNotice"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): Team "(?P<winnerTeam>CT|TERRORIST)?" triggered "SFUI_Notice_(CTs|Terrorists)?_Win" \(CT "(?P<ctScore>\d+)"\) \(T "(?P<tScore>\d+)"\)$`)
	re["gameOver"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): Game Over: (?P<gameType>\w+) (?P<mapPool>\w+) (?P<map>\w+) score (?P<ctScore>\d+):(?P<tScore>\d+) after (?P<mapTime>\d+) min$`)
	re["switchTeam"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): "(?P<playerName>[^<]+)<\d+><(?P<steamId>BOT|STEAM[^>]+)>(<(?P<currentTeam>CT|TERRORIST|Unassigned)>)?" switched from team <(?P<fromTeam>CT|TERRORIST|Unassigned)> to <(?P<toTeam>CT|TERRORIST|Unassigned)>$`)
	re["killed"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): "(?P<killerName>[^<]+)<\d+><(?P<killerSteamId>BOT|STEAM[^>]+)><(?P<killerTeam>CT|TERRORIST)>" \[(?P<killerXYZ>-?\d+ -?\d+ -?\d+)] killed "(?P<killedName>[^<]+)<\d+><(?P<killedSteamId>BOT|STEAM[^>]+)><(?P<killedTeam>CT|TERRORIST)>" \[(?P<killedXYZ>-?\d+ -?\d+ -?\d+)] with "(?P<killerWeapon>\w+)"\s?\(?(?P<special>\w*)\)?$`)
	re["assistedKill"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): "(?P<killerName>[^<]+)<\d+><(?P<killerSteamId>BOT|STEAM[^>]+)><(?P<killerTeam>CT|TERRORIST)>" assisted killing "(?P<killedName>[^<]+)<\d+><(?P<killedSteamId>BOT|STEAM[^>]+)><(?P<killedTeam>CT|TERRORIST)>"$`)
	re["killedOther"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): "(?P<killerName>[^<]+)<\d+><(?P<killerSteamId>BOT|STEAM[^>]+)><(?P<killerTeam>CT|TERRORIST)>" \[(?P<killerXYZ>-?\d+ -?\d+ -?\d+)] killed other "(?P<killedName>[^<]+)<\d+>" \[(?P<killedXYZ>-?\d+ -?\d+ -?\d+)] with "(?P<killerWeapon>\w+)"$`)
	re["threw"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): "(?P<playerName>[^<]+)<\d+><(?P<playerSteamId>BOT|STEAM[^>]+)><(?P<playerTeam>CT|TERRORIST)>" threw (?P<object>\S+) \[(?P<playerXYZ>-?\d+ -?\d+ -?\d+)]( flashbang entindex \d+)?$`)
	re["blinded"] = regexp.MustCompile(`^L (?P<date>\d{2}/\d{2}/\d{4}) - (?P<time>\d{2}:\d{2}:\d{2}): "(?P<blindedName>[^<]+)<\d+><(?P<blindedSteamId>BOT|STEAM[^>]+)><(?P<blindedTeam>CT|TERRORIST)>" blinded for (?P<blindedTime>\S+) by "(?P<byName>[^<]+)<\d+><(?P<bySteamId>BOT|STEAM[^>]+)><(?P<byTeam>CT|TERRORIST)>" from flashbang entindex \d+$`)
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

func parseDate(dateStr string) (pgtype.Date, error) {
	dbDate, err := time.Parse("01/02/2006", dateStr)
	if err != nil {
		return pgtype.Date{Valid: false}, err
	} else {
		return pgtype.Date{
			Time:             dbDate,
			InfinityModifier: 0,
			Valid:            true,
		}, nil
	}
}

func parseTime(dateStr string, timeStr string) (pgtype.Time, error) {
	dbTime, err := time.Parse("01/02/2006 15:04:05", dateStr+" "+timeStr)
	if err != nil {
		return pgtype.Time{Valid: false}, err
	}
	dbTimeMidnight, err := time.Parse("01/02/2006 15:04:05", dateStr+" 00:00:00")
	if err != nil {
		return pgtype.Time{Valid: false}, err
	}

	diff := dbTime.Sub(dbTimeMidnight)

	return pgtype.Time{Microseconds: diff.Microseconds(), Valid: true}, nil
}

func parseFile(r io.Reader, re map[string]*regexp.Regexp) error {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, k.String("database_url"))
	if err != nil {
		log.Fatalf("error connecting to database: %v\n", err)
	}
	defer conn.Close(ctx)

	var currentMatch database.Match
	var currentRound database.Round
	var currentWinnerID int64

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		for key, rgx := range re {
			if rgx.MatchString(line) {
				switch key {
				case "matchStart":
					matches := rgx.FindStringSubmatch(line)
					dateIdx := rgx.SubexpIndex("date")
					timeIdx := rgx.SubexpIndex("time")
					mapIdx := rgx.SubexpIndex("map")
					currentMatch, err = handleMatchStart(ctx, conn, matches[dateIdx], matches[timeIdx], matches[mapIdx])
					if err != nil {
						log.Println("error handling matchStart: ", err)
					}
				case "matchStatus":
					fmt.Println("matchStatus")
				case "roundStart":
					matches := rgx.FindStringSubmatch(line)
					dateIdx := rgx.SubexpIndex("date")
					timeIdx := rgx.SubexpIndex("time")
					currentRound, err = handleRoundStart(ctx, conn, currentMatch.ID, matches[dateIdx], matches[timeIdx])
					if err != nil {
						log.Println("error handling roundStart: ", err)
					}
				case "roundEnd":
					matches := rgx.FindStringSubmatch(line)
					dateIdx := rgx.SubexpIndex("date")
					timeIdx := rgx.SubexpIndex("time")
					err = handleRoundEnd(ctx, conn, currentRound.ID, matches[dateIdx], matches[timeIdx], currentWinnerID)
					if err != nil {
						log.Println("error handling roundEnd: ", err)
					}
				case "sfuiNotice":
					matches := rgx.FindStringSubmatch(line)
					winnerIdx := rgx.SubexpIndex("winnerTeam")
					currentWinnerID, err = getTeamID(ctx, conn, matches[winnerIdx])
					if err != nil {
						log.Println("error getting team id: ", err)
					}
				case "gameOver":
					matches := rgx.FindStringSubmatch(line)
					dateIdx := rgx.SubexpIndex("date")
					timeIdx := rgx.SubexpIndex("time")
					gameTypeIdx := rgx.SubexpIndex("gameType")
					ctScoreIdx := rgx.SubexpIndex("ctScore")
					tScoreIdx := rgx.SubexpIndex("tScore")
					var ctScore int
					var tScore int
					ctScore, err = strconv.Atoi(matches[ctScoreIdx])
					if err != nil {
						log.Println("error parsing ctScore: ", err)
					}
					tScore, err = strconv.Atoi(matches[tScoreIdx])
					if err != nil {
						log.Println("error parsing tScore: ", err)
					}
					err = handleGameOver(ctx, conn, currentMatch.ID, matches[dateIdx], matches[timeIdx], int32(ctScore), int32(tScore), matches[gameTypeIdx])
				case "switchTeam":
					matches := rgx.FindStringSubmatch(line)
					dateIdx := rgx.SubexpIndex("date")
					timeIdx := rgx.SubexpIndex("time")
					playerNameIdx := rgx.SubexpIndex("playerName")
					steamIdIdx := rgx.SubexpIndex("steamId")
					fromTeamIdx := rgx.SubexpIndex("fromTeam")
					toTeamIdx := rgx.SubexpIndex("toTeam")
					err = handleSwitchTeam(ctx, conn, matches[dateIdx], matches[timeIdx], matches[playerNameIdx], matches[steamIdIdx], matches[fromTeamIdx], matches[toTeamIdx], currentRound.ID)
					if err != nil {
						log.Println("error handling switchTeam: ", err)
					}
				case "attacking":
					fmt.Println("attacking")
				case "killed":
					fmt.Println("killed")
				case "assistedKill":
					fmt.Println("assistedKill")
				case "shopping":
					fmt.Println("shopping")
				case "moneyChange":
					fmt.Println("moneyChange")
				case "leftBuyZone":
					fmt.Println("leftBuyZone")
				case "suicide":
					fmt.Println("suicide")
				case "disconnected":
					fmt.Println("disconnected")
				case "accolade":
					fmt.Println("accolade")
				case "triggered":
					fmt.Println("triggered")
				default:
					fmt.Println("No match found for line: ", line)
				}
			}
		}
	}
	return nil
}

func getTeamID(ctx context.Context, conn *pgx.Conn, teamName string) (int64, error) {
	queries := database.New(conn)
	team, err := queries.GetTeamByName(ctx, teamName)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			team, err = queries.CreateTeam(ctx, teamName)
			if err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}
	return team.ID, nil
}

func handleRoundStart(ctx context.Context, conn *pgx.Conn, matchId int64, dateStr string, timeStr string) (database.Round, error) {
	dbDate, err := parseDate(dateStr)
	if err != nil {
		return database.Round{}, err
	}
	dbTime, err := parseTime(dateStr, timeStr)
	if err != nil {
		return database.Round{}, err
	}
	queries := database.New(conn)
	round, err := queries.GetRoundByStartDateTime(ctx, database.GetRoundByStartDateTimeParams{
		StartDate: dbDate,
		StartTime: dbTime,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			round, err = queries.CreateRound(ctx, database.CreateRoundParams{
				StartDate: dbDate,
				StartTime: dbTime,
				MatchID:   pgtype.Int8{Int64: matchId, Valid: true},
			})
			if err != nil {
				return database.Round{}, err
			}
		} else {
			return database.Round{}, err
		}
	}
	return round, nil
}

func handleRoundEnd(ctx context.Context, conn *pgx.Conn, roundId int64, dateStr string, timeStr string, winnerId int64) error {
	dbDate, err := parseDate(dateStr)
	if err != nil {
		return err
	}
	dbTime, err := parseTime(dateStr, timeStr)
	if err != nil {
		return err
	}
	queries := database.New(conn)
	err = queries.UpdateRound(ctx, database.UpdateRoundParams{
		ID:           roundId,
		EndDate:      dbDate,
		EndTime:      dbTime,
		WinnerTeamID: pgtype.Int8{Int64: winnerId, Valid: true},
	})
	if err != nil {
		return err
	}
	return nil
}

func handleMatchStart(ctx context.Context, conn *pgx.Conn, dateStr string, timeStr string, csMap string) (database.Match, error) {
	dbDate, err := parseDate(dateStr)
	if err != nil {
		return database.Match{}, err
	}
	dbTime, err := parseTime(dateStr, timeStr)
	if err != nil {
		return database.Match{}, err
	}
	queries := database.New(conn)
	dbMap, err := queries.GetMapByName(ctx, csMap)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// First time we see this map, so create it in the database
			dbMap, err = queries.CreateMap(ctx, csMap)
			if err != nil {
				return database.Match{}, err
			}
		} else {
			return database.Match{}, err
		}
	}
	match, err := queries.GetMatchByMapAndStartDateTime(ctx, database.GetMatchByMapAndStartDateTimeParams{
		MapID:     pgtype.Int8{Int64: dbMap.ID, Valid: true},
		StartDate: dbDate,
		StartTime: dbTime,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			match, err = queries.CreateMatch(ctx, database.CreateMatchParams{
				StartDate: dbDate,
				StartTime: dbTime,
				MapID:     pgtype.Int8{Int64: dbMap.ID, Valid: true},
			})
			if err != nil {
				return database.Match{}, err
			}
		} else {
			return database.Match{}, err
		}
	}
	return match, nil
}

func handleGameOver(ctx context.Context, conn *pgx.Conn, matchId int64, dateStr string, timeStr string, ctScore int32, tScore int32, gameTypeStr string) error {
	dbDate, err := parseDate(dateStr)
	if err != nil {
		return err
	}
	dbTime, err := parseTime(dateStr, timeStr)
	if err != nil {
		return err
	}
	queries := database.New(conn)
	gameType, err := queries.GetGameTypeByName(ctx, gameTypeStr)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			gameType, err = queries.CreateGameType(ctx, gameTypeStr)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	err = queries.UpdateMatch(ctx, database.UpdateMatchParams{
		ID:         matchId,
		EndDate:    dbDate,
		EndTime:    dbTime,
		ScoreCt:    ctScore,
		ScoreT:     tScore,
		GameTypeID: pgtype.Int8{Int64: gameType.ID, Valid: true},
	})
	if err != nil {
		return err
	}
	return nil
}

func handleSwitchTeam(ctx context.Context, conn *pgx.Conn, dateStr string, timeStr string, playerName string, steamId string, fromTeam string, toTeam string, roundID int64) error {
	dbDate, err := parseDate(dateStr)
	if err != nil {
		return err
	}
	dbTime, err := parseTime(dateStr, timeStr)
	if err != nil {
		return err
	}
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
			steamUser, err = queries.CreateSteamUser(ctx, database.CreateSteamUserParams{
				SteamID:          steamId,
				SteamCommunityID: steamCommunityID,
			})
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
			player, err = queries.CreatePlayer(ctx, database.CreatePlayerParams{
				Name:        playerName,
				SteamUserID: pgtype.Int8{Int64: steamUser.ID, Valid: true},
				Bot:         pgtype.Bool{Bool: bot, Valid: true},
			})
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	_, err = queries.GetTeamSwitchByPlayerAndDateTime(ctx, database.GetTeamSwitchByPlayerAndDateTimeParams{
		PlayerID:   pgtype.Int8{Int64: player.ID, Valid: true},
		SwitchDate: dbDate,
		SwitchTime: dbTime,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			fromTeamID, err2 := getTeamID(ctx, conn, fromTeam)
			if err2 != nil {
				return err2
			}
			toTeamID, err3 := getTeamID(ctx, conn, toTeam)
			if err3 != nil {
				return err3
			}
			err = queries.CreateTeamSwitchEvent(ctx, database.CreateTeamSwitchEventParams{
				PlayerID:   pgtype.Int8{Int64: player.ID, Valid: true},
				FromTeamID: pgtype.Int8{Int64: fromTeamID, Valid: true},
				ToTeamID:   pgtype.Int8{Int64: toTeamID, Valid: true},
				SwitchDate: dbDate,
				SwitchTime: dbTime,
				RoundID:    pgtype.Int8{Int64: roundID, Valid: true},
			})
		} else {
			return err
		}
	}
	return nil
}
