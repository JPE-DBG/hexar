# Feedback Intake Checklist

Use this checklist for each tester session to ensure all feedback is processed and nothing is missed.

---

## Tester A (Czech) — 2026-05-17

- [x] Read session file completely (`.claude/feedback/sessions/tester_czech_2026-05-17.md`)
- [x] Extract all issues (9 identified)
- [x] Categorize each (HIGH/MEDIUM/LOW)
- [x] Add to CLAUDE.md Feedback Log
- [x] Make decision: IMPLEMENT / REJECT / DEFER for each
- [x] Code changes + test (all implemented)
- [x] Update tester file status block (PROCESSED)
- [x] Close checklist

**Status:** 8/8 complete ✅

**Summary:**
- 9 issues submitted
- 9 implemented ✅
- 0 rejected ❌
- 0 deferred ⏸️

---

## Template for Beta Cycle 2+

```markdown
## Tester [Name] — [YYYY-MM-DD]

- [ ] Read session file completely (`.claude/feedback/sessions/tester_*.md`)
- [ ] Extract all issues (count: ___)
- [ ] Categorize each (HIGH/MEDIUM/LOW)
- [ ] Add to CLAUDE.md Feedback Log
- [ ] Make decision: IMPLEMENT / REJECT / DEFER for each
- [ ] Code changes + test (if implemented)
- [ ] Update tester file status block
- [ ] Close checklist

**Status:** _/8 complete

**Summary:**
- ___ issues submitted
- ___ implemented ✅
- ___ rejected ❌
- ___ deferred ⏸️
```

---

## Audit Anytime

**Verify no feedback session was skipped:**

```bash
# Count total tester session files
ls .claude/feedback/sessions/tester_*.md | wc -l

# Count sessions marked as PROCESSED
grep -l "status: PROCESSED" .claude/feedback/sessions/tester_*.md | wc -l

# If these don't match → unprocessed sessions exist
# Find which ones:
ls .claude/feedback/sessions/tester_*.md | while read f; do
  grep -q "status: PROCESSED" "$f" || echo "UNPROCESSED: $f"
done
```

**Verify all feedback was added to CLAUDE.md:**

```bash
# Count total issues across all tester files
grep "^## Issue" .claude/feedback/sessions/tester_*.md | wc -l

# Count issues in CLAUDE.md Feedback Log
grep "^| ✅\|^| ❌\|^| ⏸️" CLAUDE.md | wc -l

# If these don't match → some feedback wasn't logged yet
```

**Quick status snapshot:**

```bash
echo "=== Feedback Processing Status ==="
echo "Tester sessions:"
ls .claude/feedback/sessions/tester_*.md
echo ""
echo "Processed:"
grep "status: PROCESSED" .claude/feedback/sessions/tester_*.md | cut -d: -f1
echo ""
echo "Issues logged in CLAUDE.md:"
grep "^| ✅\|^| ❌\|^| ⏸️" CLAUDE.md | wc -l
```

---
