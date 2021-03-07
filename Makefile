export GO111MODULE=on

test:
	go test --coverprofile cover.out ./cmd/... ./internal/...
	go tool cover -html=cover.out -o coverage/go-coverage.html

financesd: $(wildcard internal/*/*.go) cmd/financesd/financesd.go
	go build ./cmd/financesd

clean:
	rm -f financesd cover.out coverage/go-coverage.html

# install: financesd
# 	go install ./cmd/financesd/
