# Chess 2

[![Documentation](https://godoc.org/github.com/CGamesPlay/chess2/pkg/chess2?status.svg)](https://godoc.org/github.com/CGamesPlay/chess2/pkg/chess2)

Chess2 implemented in Go.

To start a Chess 2 API server and hit it with a request:

```bash
make serve &
http -v :8080/new white==c black==k
http -v :8080/move epd="rnbkkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 ck 33" move=d2d4
```

To test the engine:

```bash
make test perft
```

## Interpretation of Chess 2 rules

- There is a duel each time a move other than a king's captures an opponent's piece. The duel can be skipped, meaning that the defender does not issue a challenge. There can be multiple duels for a single move in the case of an Elephant's rampage.
- A duel is "legal" if the defender has enough stones to initiate a challenge and pay for their bid; and the attacker has enough stones to pay for their bid.
- A move is "pseudo-legal" if it:
  - moves one of the player to move's pieces;
  - follows the distance and direction rules for the piece moved;
  - all intermediate squares are empty (except for jump moves) or capturable (for the Elephant's rampage);
  - the target square is empty or contains a capturable piece; and
  - all duels are legal.
  - Additionally, a pass move is pseudo-legal during a king turn.
- A piece is "threatened" if there is a pseudo-legal move which results in the capture of the piece.
- A move is "into check" if it leaves the board in a state where any of the player's kings are threatened.
- A move is "legal" if it is pseudo-legal and not into check.

#### Bugs

The engine presently incorrectly considers a color in check if an enemy elephant could capture the king, even if that capture is illegal because the rampage would capture the enemy's king.

## Performance tests

Here are the timings for `chess2_perft` at depth 3.

| Configuration          | Speed  | Relative |
| ---------------------- | ------ | -------- |
| Python engine          | 636.73 | 1.00000  |
| Go engine, brute force | 21.20  | 0.03330  |
| Go engine, fast        | 3.92   | 0.00616  |
