# goengine

Chess2 implemented in Go.

To test:

```bash
go test chess2/...
```

#### Interpretation of Chess 2 rulues

- There is a duel each time a move other than a king's captures an opponent's piece. The duel can be skipped, meaning that the defender does not issue a challenge. There can be multiple duels for a single move in the case of an Elephant's rampage.
- A duel is "valid" if the defender has enough stones to initiate a challenge and pay for their bid; and the attacker has enough stones to pay for their bid.
- A duel is "legal" if it does not result in the capture any of the defender's kings or if any of the valid duels following it are legal.
- A move is "pseudo-legal" if it is one of the player to move's pieces; follows the distance and direction rules for the piece moved; all intermediate squares are empty (except for jump moves) or capturable (for the Elephant's rampage); and the target square is empty or contains a capturable piece. Additionally, a pass move is pseudo-legal during a king turn.
- A move is "king-threatening" if it is pseudo-legal; captures at least one of the opponent's kings; and there are no legal duels.
- A move is "into check" if it would leave the board in a state where there the opponent has a king-threatening move (regardless of if the move is for a standard or king turn).
- A move is "winning" if it is king-threatening; or if it is pseudo-legal, not into check, and results in all of the player's kings being past the midline.
- A move is "losing" if it would leave the board in a state where the opponent has a winning move, or if it captures one of the player's own kings without also capturing one of the other player's kings.
- A move is "legal" if it is pseudo-legal; is not into check; and is not a losing move.
- If the player to move has any winning move, that player wins the game. If the player to move has no legal moves, that player loses the game.

The following points are non-normative:

- A player is not in check if there are legal duels that could prevent the opponent's capture of all of their kings.
- A defender must always select a legal duel. If there are no legal duels, the move is a winning move and thus the attacker either already won the game (or the move leading up to it was a losing move).

## TODO

- Add support for legal moves wihtout duels.