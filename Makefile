export GO111MODULE=on

test:
	go test --coverprofile cover.out ./cmd/... ./internal/...
	go tool cover -html=cover.out -o coverage/go-coverage.html

go_sources := $(shell find internal -name "*.go" ! -name "*_test.go") cmd/financesd/financesd.go
ts_sources := $(shell find web/src/lib/ \( -name "*.ts" -o -name "*.tsx" \) ! -name "*.test.*")
sass_sources := $(wildcard web/src/styles/*)

financesd: $(go_sources)
	go build ./cmd/financesd

web/dist/bundle.js web/dist/finances.css: $(sass_sources) $(ts_sources)
	npm run build

ui: web/dist/bundle.js web/dist/finances.css

clean:
	rm -f financesd cover.out coverage/go-coverage.html web/dist/*

# install: financesd
# 	go install ./cmd/financesd/

install: financesd ui
	cp financesd /opt/financesd/
	mkdir -p /opt/financesd/web/dist
	cp web/dist/* /opt/financesd/web/dist/
	cp -r web/resources/ /opt/financesd/web
