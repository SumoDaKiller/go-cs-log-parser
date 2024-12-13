package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"testing"
)

var re map[string]*regexp.Regexp

func TestMain(m *testing.M) {
	fmt.Println("Setup test environment")
	re = initializeRegexpPatterns()
	exitVal := m.Run()
	os.Exit(exitVal)
}

func Test_calculateSteamCommunityId(t *testing.T) {
	result, err := calculateSteamCommunityId("STEAM_1:0:402610")
	if err != nil {
		t.Fatal(err)
	}
	if result != 76561197961070948 {
		t.Error("Expected 76561197961070948, got ", result)
	}
}

func TestRegExp(t *testing.T) {
	f, err := os.Open("testdata/L000_000_000_000_27015_202211011929_000.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var matchStart, roundStart, roundEnd, gameOver, sfuiNotice, switchTeam, attacking int
	var killed, killedOther, assistedKill, shopping, moneyChange, suicide, triggered int
	var threw, blinded, accolade, leftBuyZone, disconnected, matchStatus int

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		for key, rgx := range re {
			if rgx.MatchString(line) {
				switch key {
				case "matchStart":
					matchStart++
				case "roundStart":
					roundStart++
				case "roundEnd":
					roundEnd++
				case "sfuiNotice":
					sfuiNotice++
				case "gameOver":
					gameOver++
				case "switchTeam":
					switchTeam++
				case "attacking":
					attacking++
				case "killed":
					killed++
				case "killedOther":
					killedOther++
				case "assistedKill":
					assistedKill++
				case "shopping":
					shopping++
				case "moneyChange":
					moneyChange++
				case "suicide":
					suicide++
				case "triggered":
					triggered++
				case "threw":
					threw++
				case "blinded":
					blinded++
				case "accolade":
					accolade++
				case "leftBuyZone":
					leftBuyZone++
				case "disconnected":
					disconnected++
				case "matchStatus":
					matchStatus++
				default:
					fmt.Println("No match found for line: ", line)
				}
			}
		}
	}
	if matchStart != 5 {
		t.Error("Expected 5 matches for matchStart, got ", matchStart)
	}
	if matchStatus != 61 {
		t.Error("Expected 183 matches for matchStatus, got ", matchStatus)
	}
	if roundStart != 32 {
		t.Error("Expected 32 matches for roundStart, got ", roundStart)
	}
	if roundEnd != 28 {
		t.Error("Expected 28 matches for roundEnd, got ", roundEnd)
	}
	if sfuiNotice != 28 {
		t.Error("Expected 28 matches for sfuiNotice, got ", sfuiNotice)
	}
	if gameOver != 2 {
		t.Error("Expected 2 matches for gameOver, got ", gameOver)
	}
	if switchTeam != 63 {
		t.Error("Expected 63 matches for switchTeam, got ", switchTeam)
	}
	if killed != 231 {
		t.Error("Expected 634 matches for killed, got ", killed)
	}
	if killedOther != 388 {
		t.Error("Expected 388 matches for killedOther, got ", killedOther)
	}
	if assistedKill != 57 {
		t.Error("Expected 57 matches for assistedKill, got ", assistedKill)
	}
	if threw != 184 {
		t.Error("Expected 184 matches for threw, got ", threw)
	}
	if blinded != 209 {
		t.Error("Expected 209 matches for blinded, got ", blinded)
	}
	if shopping != 4464 {
		t.Error("Expected 4464 matches for shopping, got ", shopping)
	}
	if moneyChange != 1316 {
		t.Error("Expected 1316 matches for moneyChange, got ", moneyChange)
	}
	if leftBuyZone != 386 {
		t.Error("Expected 386 matches for leftBuyZone, got ", leftBuyZone)
	}
	if attacking != 993 {
		t.Error("Expected 993 matches for attacking, got ", attacking)
	}
	if suicide != 6 {
		t.Error("Expected 6 matches for suicide, got ", suicide)
	}
	if disconnected != 18 {
		t.Error("Expected 18 matches for disconnected, got ", disconnected)
	}
	if accolade != 20 {
		t.Error("Expected 20 matches for accolade, got ", accolade)
	}
	if triggered != 99 {
		t.Error("Expected 99 matches for triggered, got ", triggered)
	}
}
