.PHONY: install
install:
	go install chess2/...

.PHONY: test
test:
	go test chess2/...

.PHONY: perft
perft: install
	cat test/chess2_perft.epd | chess2_perft -d 2 -b
