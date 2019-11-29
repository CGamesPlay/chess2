package main

import (
	"bufio"
	"chess2/internal/chess2"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

var (
	maxDepth   int
	bruteforce bool
)

func main() {
	pflag.IntVarP(&maxDepth, "depth", "d", 2, "depth for perft")
	pflag.BoolVarP(&bruteforce, "brute-force", "b", false, "use brute force search")
	pflag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	failed := false
	for scanner.Scan() {
		epd := scanner.Text()
		result, err := runPerft(epd)
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

// runPerft takes in a formatted input string and runs a perft test. The input
// string is an EPD string, optionally followed by a semicolon and
// slash-delimited list of numbers, corresponding to the perft at each depth.
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

	game, err := chess2.ParseEpd(epd)
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
