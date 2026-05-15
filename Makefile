.PHONY: server client dev build clean count test test-e2e

# Run Go server (serves on :8080)
server:
	go run ./cmd/server -addr :8080 -client ./client/dist

# Run Vite dev server (proxies /ws to Go server)
client:
	cd client && npm run dev

# Run both (start server in background, then client)
dev:
	@echo "Start server:  make server"
	@echo "Start client:  make client"
	@echo "Run both in separate terminals for best experience."

# Build client for production
build:
	cd client && npm run build

# Clean build artifacts
clean:
	rm -rf client/dist client/node_modules/.vite

count:
	git ls-files | grep -P ".*(go|ts)" | xargs wc -l

# Run Go unit tests
test:
	go test ./...

# Run Playwright E2E tests (starts both servers automatically)
test-e2e:
	npx playwright test
