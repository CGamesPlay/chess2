package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/CGamesPlay/chess2/pkg/chess2"

	"github.com/spf13/pflag"
)

var (
	maxDepth   int
	bruteforce bool
	classic    bool
	divide     bool
)

func main() {
	pflag.IntVarP(&maxDepth, "depth", "d", 2, "depth for perft")
	pflag.BoolVarP(&bruteforce, "brute-force", "b", false, "use brute force search")
	pflag.BoolVar(&classic, "classic", false, "use classic chess rules")
	pflag.BoolVar(&divide, "divide", false, "split results for first move")
	pflag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	failed := false
	for scanner.Scan() {
		epd := scanner.Text()
		var result string
		var err error
		if divide {
			result, err = dividePerft(epd)
		} else {
			result, err = runPerft(epd)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v (epd: %s)\n", err, epd)
			failed = true
		} else {
			fmt.Println(result)
		}
	}
	if failed {
		os.Exit(1)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "error reading: %v\n", err)
		os.Exit(2)
	}
}

// Parse the epd according to the configured parameters.
func parseEpd(epd string) (chess2.Game, error) {
	flags := chess2.VariantChess2
	usedEpd := epd
	if classic {
		flags = chess2.VariantClassic
		usedEpd += " cc 33"
	}
	game, err := chess2.ParseEpdFlags(usedEpd, flags)
	return game, err
}

// Output a summary of the perfts for maxDepth-1 for each valid move from the
// given epd.
func dividePerft(epd string) (string, error) {
	game, err := parseEpd(epd)
	if err != nil {
		return "", err
	}
	var sb strings.Builder
	sb.WriteString(strings.TrimSpace(epd))
	sb.WriteRune('\n')
	if maxDepth < 1 {
		return sb.String(), nil
	}

	var moves []chess2.Move
	if bruteforce {
		chess2.BruteforceMoveList(func(m chess2.Move) {
			if err := game.ValidateLegalMove(m); err == nil {
				moves = append(moves, m)
			}
		})
	} else {
		moves = game.GenerateLegalMoves()
	}

	for _, m := range moves {
		if maxDepth > 1 {
			child := game.ApplyMove(m)
			var results []uint64
			if bruteforce {
				results = chess2.PerftBruteforce(child, maxDepth-1)
			} else {
				results = chess2.Perft(child, maxDepth-1)
			}
			sb.WriteString(fmt.Sprintf("%v: %d\n", m, results[maxDepth-2]))
		} else {
			sb.WriteString(fmt.Sprintf("%v: 1\n", m))
		}
	}
	return sb.String(), nil
}

// Take in a formatted input string and run a perft test. The input string is
// an EPD string, optionally followed by a semicolon and slash-delimited list
// of numbers, corresponding to the perft at each depth.
func runPerft(input string) (string, error) {
	parts := strings.SplitN(input, ";", 2)
	epd := parts[0]
	var checkValues []uint64
	if len(parts) > 1 {
		perftValues := strings.Split(strings.TrimSpace(parts[1]), "/")
		checkValues = make([]uint64, len(perftValues))
		for i, str := range perftValues {
			val, err := strconv.ParseUint(str, 10, 64)
			if err != nil {
				return "", err
			}
			checkValues[i] = val
		}
	}

	game, err := parseEpd(epd)
	if err != nil {
		return "", err
	}
	var result []uint64
	if bruteforce {
		result = chess2.PerftBruteforce(game, maxDepth)
	} else {
		result = chess2.Perft(game, maxDepth)
	}
	for i := 0; i < len(checkValues) && i < maxDepth; i++ {
		if checkValues[i] != result[i] {
			return "", fmt.Errorf("expected %d, found %d at depth %d", checkValues[i], result[i], i+1)
		}
	}
	if len(result) < len(checkValues) {
		// Preserve deeper but unchecked perfts
		result = checkValues
	}
	var sb strings.Builder
	sb.WriteString(strings.TrimSpace(epd))
	sb.WriteString(" ; ")
	for i, value := range result {
		if i != 0 {
			sb.WriteString("/")
		}
		sb.WriteString(strconv.FormatUint(value, 10))
	}
	return sb.String(), nil
}
