package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/CGamesPlay/chess2/pkg/chess2"
)

type requestStruct struct {
	Armies string `json:"armies"`
	Epd    string `json:"epd"`
	Move   string `json:"move"`
}

func formatGame(game chess2.Game) map[string]interface{} {
	response := make(map[string]interface{})
	response["epd"] = chess2.EncodeEpd(game)
	legalMoves := game.GenerateLegalMoves()
	legalMovesStr := make([]string, len(legalMoves))
	for i := 0; i < len(legalMoves); i++ {
		legalMovesStr[i] = legalMoves[i].String()
	}
	sort.Strings(legalMovesStr)
	response["legal_moves"] = legalMovesStr
	response["game_over"] = game.GameState() != chess2.GameInProgress
	switch game.GameState() {
	case chess2.GameInProgress:
		response["winner"] = nil
	case chess2.GameOverWhite:
		response["winner"] = "white"
	case chess2.GameOverBlack:
		response["winner"] = "black"
	case chess2.GameOverDraw:
		response["winner"] = "draw"
	}
	return response
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		data := scanner.Text()
		var request requestStruct
		var game chess2.Game
		var err error
		if err = json.Unmarshal([]byte(data), &request); err != nil {
			err = fmt.Errorf("Invalid JSON input")
		} else if request.Armies != "" && request.Epd != "" {
			err = fmt.Errorf("Either `epd` or `armies` must be provided; but not both")
		} else if request.Armies != "" {
			if len(request.Armies) == 2 {
				white, foundWhite := chess2.FindArmySymbol(rune(request.Armies[0]))
				black, foundBlack := chess2.FindArmySymbol(rune(request.Armies[1]))
				if !foundWhite || !foundBlack {
					err = fmt.Errorf("Invalid armies")
				} else {
					game = chess2.GameFromArmies(white, black)
				}
			} else {
				err = fmt.Errorf("Invalid armies")
			}
		} else {
			game, err = chess2.ParseEpd(request.Epd)
		}
		var response map[string]interface{}
		if err == nil {
			if request.Move != "" {
				var move chess2.Move
				move, err = chess2.ParseUci(request.Move)
				if err == nil {
					if err = game.ValidateLegalMove(move); err != nil {
						err = fmt.Errorf("illegal move: %s", err.Error())
					} else {
						nextGame := game.ApplyMove(move)
						duels := game.GenerateDuels(move)
						duelsStr := make([]string, len(duels))
						for i := 0; i < len(duels); i++ {
							duelsStr[i] = duels[i].String()
						}
						response = formatGame(nextGame)
						response["available_duels"] = duelsStr
					}
				}
			} else {
				response = formatGame(game)
			}
		}
		if err != nil {
			response = map[string]interface{}{"error": err.Error()}
		}
		json, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(json))
	}
}
