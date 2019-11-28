# goengine

Chess2 implemented in Go.

To test:

```bash
make test
```

#### Interpretation of Chess 2 rules

- There is a duel each time a move other than a king's captures an opponent's piece. The duel can be skipped, meaning that the defender does not issue a challenge. There can be multiple duels for a single move in the case of an Elephant's rampage.
- A duel is "legal" if the defender has enough stones to initiate a challenge and pay for their bid; and the attacker has enough stones to pay for their bid.
- A move is "pseudo-legal" if it:
  - moves one of the player to move's pieces;
  - follows the distance and direction rules for the piece moved;
  - all intermediate squares are empty (except for jump moves) or capturable (for the Elephant's rampage);
  - the target square is empty or contains a capturable piece; and
  - all duels are legal.
  - Additionally, a pass move is pseudo-legal during a king turn.
- A move is "into check" if it leaves the board in a state where any of the player's kings are threatened.
- A move is "legal" if it is pseudo-legal and not into check.

## TODO

The engine presently incorrectly considers a color in check if an enemy elephant could capture the king, even if that capture is illegal because the rampage would capture the enemy's king.
