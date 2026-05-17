.PHONY: server client dev build clean count test test-e2e audit-feedback audit-feedback-sessions audit-feedback-logged

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

# Feedback audit targets

# Audit which tester sessions have been processed
audit-feedback-sessions:
	@echo "=== Tester Sessions ===" && \
	echo "Total sessions:" && \
	ls .claude/feedback/sessions/tester_*.md 2>/dev/null | wc -l && \
	echo "" && \
	echo "Processed sessions:" && \
	grep -l "status: PROCESSED" .claude/feedback/sessions/tester_*.md 2>/dev/null | wc -l && \
	echo "" && \
	echo "Unprocessed sessions (if any):" && \
	ls .claude/feedback/sessions/tester_*.md 2>/dev/null | while read f; do \
		grep -q "status: PROCESSED" "$$f" || echo "  - $$f"; \
	done

# Audit if all feedback issues were logged to CLAUDE.md
audit-feedback-logged:
	@echo "=== Feedback Issues ===" && \
	echo "Total issues across all sessions:" && \
	grep "^## Issue" .claude/feedback/sessions/tester_*.md 2>/dev/null | wc -l && \
	echo "" && \
	echo "Issues logged in CLAUDE.md:" && \
	grep "^| ✅\|^| ❌\|^| ⏸️" CLAUDE.md 2>/dev/null | wc -l

# Quick feedback status snapshot
audit-feedback:
	@echo "=== Feedback Processing Status ===" && \
	echo "" && \
	echo "Session files:" && \
	ls -1 .claude/feedback/sessions/tester_*.md 2>/dev/null | sed 's|.*/||' || echo "  (none)" && \
	echo "" && \
	echo "Processed:" && \
	grep -h "^tester:" .claude/feedback/sessions/tester_*.md 2>/dev/null | sort -u | while read line; do \
		file=$$(grep -l "$$line" .claude/feedback/sessions/tester_*.md 2>/dev/null | head -1); \
		[ -n "$$file" ] && (grep -q "status: PROCESSED" "$$file" && echo "  ✅ $$line" || echo "  ⏳ $$line"); \
	done && \
	echo "" && \
	echo "Issues logged in CLAUDE.md: $$(grep -c '^| ✅\|^| ❌\|^| ⏸️' CLAUDE.md 2>/dev/null || echo 0)"
