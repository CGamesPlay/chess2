PKG=github.com/CGamesPlay/chess2

.PHONY: install
install:
	go install $(PKG)/...

.PHONY: test
test: install
	go test $(PKG)/...
	./test/json_server.sh

.PHONY: perft
perft:
	cat test/chess2_perft.epd | go run $(PKG)/cmd/chess2_perft -d 3 >/dev/null
	cat test/perft.epd | go run $(PKG)/cmd/chess2_perft --classic -d 3 >/dev/null

.PHONY: serve
serve: install
	chess2_api
