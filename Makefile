.PHONY: server client dev build clean count

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
