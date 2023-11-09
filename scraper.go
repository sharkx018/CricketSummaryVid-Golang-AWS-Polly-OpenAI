package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type CommentaryList struct {
	CommentaryList []Commentary    `json:"commentaryList"`
	MatchHeader    MatchHeaderInfo `json:"matchHeader"`
}

type MatchHeaderInfo struct {
	Status      string          `json:"status"`
	TossResults TossResultsInfo `json:"tossResults"`
	SeriesName  string          `json:"seriesName"`
	Result      ResultInfo      `json:"result"`
	TeamOne     TeamInfo        `json:"team1"`
	TeamTwo     TeamInfo        `json:"team2"`
}

type ResultInfo struct {
	ResultType    string `json:"resultType"`
	WinningTeam   string `json:"winningTeam"`
	WinningteamId int    `json:"winningteamId"`
	WinningMargin int    `json:"winningMargin"`
	WinByRuns     bool   `json:"winByRuns"`
	WinByInnings  bool   `json:"winByInnings"`
}

type TossResultsInfo struct {
	TossWinnerId   int    `json:"tossWinnerId"`
	TossWinnerName string `json:"tossWinnerName"`
	Decision       string `json:"decision"`
}

type Commentary struct {
	CommText          string                 `json:"commText"`
	Timestamp         int64                  `json:"timestamp"`
	BallNumber        int                    `json:"ballNbr"`
	OverNumber        float64                `json:"overNumber"`
	InningsID         int                    `json:"inningsId"`
	Event             string                 `json:"event"`
	BatTeamName       string                 `json:"batTeamName"`
	CommentaryFormats map[string]interface{} `json:"commentaryFormats"`
	BatsmanStriker    Batsman                `json:"batsmanStriker"`
	BowlerStriker     Bowler                 `json:"bowlerStriker"`
	BatTeamScore      int                    `json:"batTeamScore"`
	OverSeperator     OverSeparatorInfo      `json:"overSeparator"`
}

type TeamInfo struct {
	Id            int           `json:"id"`
	Name          string        `json:"name"`
	PlayerDetails []interface{} `json:"playerDetails"`
	ShortName     string        `json:"shortName"`
}

type OverSeparatorInfo struct {
	BatStrikerNames []string `json:"batStrikerNames"`
	BatStrikerRuns  int      `json:"batStrikerRuns"`

	BatNonStrikerNames []string `json:"batNonStrikerNames"`
	BatNonStrikerRuns  int      `json:"batNonStrikerRuns"`
}

type Batsman struct {
	BatBalls      int     `json:"batBalls"`
	BatDots       int     `json:"batDots"`
	BatFours      int     `json:"batFours"`
	BatID         int     `json:"batId"`
	BatName       string  `json:"batName"`
	BatMins       int     `json:"batMins"`
	BatRuns       int     `json:"batRuns"`
	BatSixes      int     `json:"batSixes"`
	BatStrikeRate float64 `json:"batStrikeRate"`
}

type Bowler struct {
	BowlID      int     `json:"bowlId"`
	BowlName    string  `json:"bowlName"`
	BowlMaidens int     `json:"bowlMaidens"`
	BowlNoballs int     `json:"bowlNoballs"`
	BowlOvs     float64 `json:"bowlOvs"`
	BowlRuns    int     `json:"bowlRuns"`
	BowlWides   int     `json:"bowlWides"`
	BowlWkts    int     `json:"bowlWkts"`
	BowlEcon    float64 `json:"bowlEcon"`
}

const (
	matchUrl      = "https://www.cricbuzz.com/api/cricket-match/commentary/"
	PaginationUrl = "https://www.cricbuzz.com/api/cricket-match/commentary-pagination/"
)

func scrapeData(link string) ([]string, error) {

	var timestampPagination int64 = -1
	var currentInningID int = 2

	matchID, err := extractMatchID(link)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err

	}
	visitedBall := map[int]bool{}

	commentaryData, err := callMatchAPI(matchUrl + matchID)

	finalResponse := []string{
		" " + commentaryData.MatchHeader.Status, // define who match
		" Innings 2 ended",
	}

	for {

		done := false
		// loop for response
		for _, comment := range commentaryData.CommentaryList {

			// pagination update
			timestampPagination = comment.Timestamp

			// take care of inning changes, reset the visited balls
			if comment.InningsID != currentInningID {
				finalResponse = append(finalResponse, " Innings 1 ended")
				currentInningID = 1
				visitedBall = map[int]bool{}
			}

			ballNum := getBallNoFromOverNumber(comment.OverNumber)

			// Only store the ball commentary
			if ballNum == 0 || visitedBall[ballNum] == true {
				continue
			}
			// visitBall
			visitedBall[ballNum] = true

			if strings.ToLower(comment.Event) == strings.ToLower("WICKET") {
				finalResponse = append(finalResponse, comment.CommText+" "+comment.BatsmanStriker.BatName+" is OUT "+comment.Event)
			} else if strings.ToLower(comment.Event) == strings.ToLower("FOUR") {
				finalResponse = append(finalResponse, comment.CommText+" "+comment.BatsmanStriker.BatName+" hits "+comment.Event)
			} else if strings.ToLower(comment.Event) == strings.ToLower("SIX") {
				finalResponse = append(finalResponse, comment.CommText+" "+comment.BatsmanStriker.BatName+" hits "+comment.Event)
			} else if strings.Contains(strings.ToLower(comment.Event), "over-break") {

				str := strings.ToLower(comment.Event)

				if strings.Contains(str, strings.ToLower("WICKET")) {
					finalResponse = append(finalResponse, comment.CommText+" "+comment.BatsmanStriker.BatName+" is OUT and now its OVER-BREAK")
				} else if strings.Contains(str, strings.ToLower("FOUR")) {
					finalResponse = append(finalResponse, comment.CommText+" "+comment.BatsmanStriker.BatName+" hits FOUR and now its OVER-BREAK")
				} else if strings.Contains(str, strings.ToLower("SIX")) {
					finalResponse = append(finalResponse, comment.CommText+" "+comment.BatsmanStriker.BatName+" hits SIX and now its OVER-BREAK")
				} else if str == strings.ToLower("over-break") {
					finalResponse = append(finalResponse, comment.CommText+" and now its "+strings.ToUpper(comment.Event))
				}

				// add the match summary
				// team batting, total runs made ,total overs, inningID
				finalResponse = append(finalResponse, fmt.Sprintf(" batting team: %s, Score: %d, Overs: %.1f, Innings: %d, batStrikerName: %s, batStrikerRuns: %d, non-batStrikerName: %s, nonbatStrikerRuns: %d,",
					comment.BatTeamName,
					comment.BatTeamScore,
					comment.OverNumber,
					comment.InningsID,
					comment.OverSeperator.BatStrikerNames[0],
					comment.OverSeperator.BatStrikerRuns,
					comment.OverSeperator.BatNonStrikerNames[0],
					comment.OverSeperator.BatNonStrikerRuns))
			}

			//DEBUG STATMENT
			// finalResponse[len(finalResponse)-1] = fmt.Sprintf("Innings: %d, Ball: %d ", currentInningID, ballNum) + finalResponse[len(finalResponse)-1]

			// match ended from reverse side
			if currentInningID == 1 && ballNum == 1 {
				done = true
				break
			}

		}

		// check for termination
		if done {

			finalResponse = append(finalResponse, " "+commentaryData.MatchHeader.TossResults.TossWinnerName+" wins the toss and decided to "+commentaryData.MatchHeader.TossResults.Decision)
			finalResponse = append(finalResponse, " "+commentaryData.MatchHeader.SeriesName+" ("+commentaryData.MatchHeader.TeamOne.Name+" vs "+commentaryData.MatchHeader.TeamTwo.Name+")")

			break
		}

		// call to paginationAPI
		//commentaryData, status, err := callMatchAPI(PaginationUrl + matchID + "/" + to_string(currentInningID) + "/" + timestampPagination)
		commentaryData, err = callMatchAPI(fmt.Sprintf("%s%s/%d/%d", PaginationUrl, matchID, currentInningID, timestampPagination))
		if err != nil {
			fmt.Println("Error while calling pagination api: ", fmt.Sprintf("%s%s/%d/%d", PaginationUrl, matchID, currentInningID, timestampPagination))
		}

	}

	// Specify substrings to remove
	substringsToRemove := []string{"B0$", "B1$"}

	// Remove specified substrings from the string slice
	finalResponse = removeSubstrings(finalResponse, substringsToRemove...)

	// reverse the timeline
	reverseStringSlice(finalResponse)
	//fmt.Println(commentaryData)

	err = writeStringSliceToTextFile("commentaryLogs.txt", finalResponse)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return finalResponse, nil

}

func reverseStringSlice(slice []string) {
	length := len(slice)
	for i := 0; i < length/2; i++ {
		// Swap elements from the beginning with corresponding elements from the end
		slice[i], slice[length-i-1] = slice[length-i-1], slice[i]
	}
}

func writeStringSliceToTextFile(fileName string, data []string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, str := range data {
		_, err := fmt.Fprintln(file, str)
		if err != nil {
			return err
		}
	}

	fmt.Println("Data written to", fileName)
	return nil
}

func getBallNoFromOverNumber(overNumber float64) int {

	if overNumber == 0 {
		return 0
	}

	overs := int(overNumber)
	balls := int((overNumber - float64(overs)) * 10)
	return overs*6 + balls

}

func callMatchAPI(url string) (CommentaryList, error) {

	// URL for the cricket match commentary (replace 82374 with the actual match ID)
	//url := "https://www.cricbuzz.com/api/cricket-match/commentary/"+matchID

	// Make GET request
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making the request:", err)
		return CommentaryList{}, err
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error reading the response body:", err)
		return CommentaryList{}, err
	}

	// Unmarshal the JSON response into CommentaryList
	var commentaryData CommentaryList
	err = json.Unmarshal(body, &commentaryData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return CommentaryList{}, err
	}

	// Use the commentaryList as needed
	fmt.Printf("%+v\n", commentaryData)

	return commentaryData, nil

}

func extractMatchID(urlString string) (string, error) {
	// Parse the URL
	u, err := url.Parse(urlString)
	if err != nil {
		return "", err
	}

	// Split the path and get the last part
	pathParts := strings.Split(u.Path, "/")
	if len(pathParts) < 2 {
		return "", fmt.Errorf("invalid URL path")
	}

	// Get the match ID
	matchID := pathParts[len(pathParts)-2]

	return matchID, nil
}

func removeSubstrings(slice []string, substrings ...string) []string {
	var result []string

	for _, str := range slice {
		// Remove specified substrings from each string
		for _, sub := range substrings {
			str = strings.Replace(str, sub, "", -1)
		}
		result = append(result, str)
	}

	return result
}
