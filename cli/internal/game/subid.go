package game

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// SubIDRange is a single parsed line from /etc/subuid or /etc/subgid,
// in the format "username:startID:count".
type SubIDRange struct {
	Username string
	Start    int64
	Count    int64
}

// MinRootlessSubIDCount is the smallest count Podman's rootless mode
// needs to comfortably map users/groups inside containers (SPEC §7).
const MinRootlessSubIDCount = 65536

// ParseSubIDFile parses the contents of /etc/subuid or /etc/subgid.
func ParseSubIDFile(content string) ([]SubIDRange, error) {
	var ranges []SubIDRange
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 3 {
			return nil, fmt.Errorf("line %d: expected username:start:count, got %q", lineNo, line)
		}

		start, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid start %q: %w", lineNo, parts[1], err)
		}

		count, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid count %q: %w", lineNo, parts[2], err)
		}

		ranges = append(ranges, SubIDRange{Username: parts[0], Start: start, Count: count})
	}

	return ranges, scanner.Err()
}

// HasSufficientSubIDRange reports whether the given user has a subuid/subgid
// range large enough for rootless Podman (>= MinRootlessSubIDCount).
func HasSufficientSubIDRange(ranges []SubIDRange, username string) bool {
	for _, r := range ranges {
		if r.Username == username && r.Count >= MinRootlessSubIDCount {
			return true
		}
	}
	return false
}
