.PHONY: install
install:
	go install chess2/...

.PHONY: test
test:
	go test chess2/...

.PHONY: perft
perft:
	cat test/chess2_perft.epd | go run chess2/cmd/chess2_perft -d 3 >/dev/null
	cat test/perft.epd | go run chess2/cmd/chess2_perft --classic -d 3 >/dev/null

.PHONY: serve
serve: install
	chess2_api
