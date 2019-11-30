package main

import (
	"chess2/internal/chess2"
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
)

func parseArmySymbol(value string, param string) (chess2.Army, error) {
	var army chess2.Army
	var found bool
	if value == "" {
		return chess2.ArmyNone, fmt.Errorf("%s is required", param)
	} else if len(value) == 1 {
		army, found = chess2.FindArmySymbol(rune(value[0]))
	}
	if !found {
		return chess2.ArmyNone, fmt.Errorf("%s must be a valid army symbol", param)
	}
	return army, nil
}

func formatGame(game chess2.Game) gin.H {
	response := make(gin.H)
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

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})
	r.GET("/new", func(c *gin.Context) {
		white, whiteErr := parseArmySymbol(c.Query("white"), "white")
		black, blackErr := parseArmySymbol(c.Query("black"), "black")
		if whiteErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": whiteErr.Error()})
		}
		if blackErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": blackErr.Error()})
		}
		game := chess2.GameFromArmies(white, black)
		c.JSON(http.StatusOK, formatGame(game))
	})
	r.POST("/move", func(c *gin.Context) {
		var request gin.H
		if err := c.BindJSON(&request); err != nil {
			return
		}
		game, err := chess2.ParseEpd(request["epd"].(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		move, err := chess2.ParseUci(request["move"].(string))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		nextGame := game.ApplyMove(move)
		response := formatGame(nextGame)
		duels := game.GenerateDuels(move)
		// Strip off the first item, which is always "don't initiate a duel"
		duelsStr := make([]string, len(duels)-1)
		for i := 1; i < len(duels); i++ {
			duelsStr[i-1] = duels[i].String()[4:]
		}
		response["available_duels"] = duelsStr
		c.JSON(http.StatusOK, response)
	})
	return r
}

func main() {
	r := setupRouter()
	r.Run()
}
