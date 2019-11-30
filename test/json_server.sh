#!/bin/bash

function test {
  result=$(echo "$1" | $(go env GOBIN)/chess2_json)
  if [ "$result" != "$2" ]; then
    echo 'Test failed:' $1 >&2
    diff -u <(echo $2) <(echo $result) >&2
    exit 1
  fi
}

test '{' '{"error":"Invalid JSON input"}'
test '{ "armies": "kk" }' '{"epd":"rnbkkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBKKBNR w KQkq - 0 1 kk 33","game_over":false,"legal_moves":["a2a3","a2a4","b1a3","b1c3","b2b3","b2b4","c2c3","c2c4","d2d3","d2d4","e2e3","e2e4","f2f3","f2f4","g1f3","g1h3","g2g3","g2g4","h2h3","h2h4"],"winner":null}'
test '{ "epd": "rnbkkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBKKBNR w KQkq - 0 1 kk 33", "move": "d2d4" }' '{"available_duels":["d2d4"],"epd":"rnbkkbnr/pppppppp/8/8/3P4/8/PPP1PPPP/RNBKKBNR K KQkq d3 0 1 kk 33","game_over":false,"legal_moves":["0000","d1d2","e1d2"],"winner":null}'
test '{ "epd": "rnbqkbnr/pppp1ppp/8/4p3/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 1 cc 11", "move": "d4f5" }' '{"error":"illegal move: unreachable square"}'
