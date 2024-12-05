package d05

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func Solution(input io.Reader) error {
	sc := bufio.NewScanner(input)

	var rules []OrderRule

	for sc.Scan() {
		line := sc.Text()
		if line == "" {
			break
		}

		var rule OrderRule

		_, err := fmt.Sscanf(line, "%d|%d", &rule[0], &rule[1])
		if err != nil {
			return fmt.Errorf("read rule %d: %w",
				len(rules)+1, err)
		}

		rules = append(rules, rule)
	}

	var (
		updateCount    int
		middleSum      int
		invalidUpdates [][]int
	)

	for sc.Scan() {
		updateCount++

		update, err := checkUpdate(sc.Text(), rules)
		if errors.Is(err, errInvalidUpdate) {
			invalidUpdates = append(invalidUpdates, update)

			continue
		} else if err != nil {
			return fmt.Errorf("verify update %d: %w",
				updateCount, err)
		}

		// Integer division of length should put us in the middle.
		middleSum += update[len(update)/2]
	}

	fmt.Printf("Sum of middles of valid updates: %d\n", middleSum)

	var invalidMiddleSum int

	for i, update := range invalidUpdates {
		err := fixUpdate(update, rules)
		if err != nil {
			return fmt.Errorf("fix invalid update %d: %w", i+1, err)
		}

		invalidMiddleSum += update[len(update)/2]
	}

	fmt.Printf("Sum of middles of corrected updates: %d\n", invalidMiddleSum)

	return nil
}

func fixUpdate(update []int, rules []OrderRule) error {
	return _fixUpdate(0, update, rules)
}

func _fixUpdate(depth int, update []int, rules []OrderRule) error {
	var edits int

	for i := range update {
		idx, violation := getViolation(update, rules, i)
		if !violation {
			continue
		}

		n := update[idx]
		copy(update[idx:i], update[idx+1:i])
		update[i-1] = update[i]
		update[i] = n

		edits++
	}

	switch {
	case edits == 0:
		return nil
	case depth < 50:
		return _fixUpdate(depth+1, update, rules)
	default:
		return fmt.Errorf("gave up at depth %d", depth)
	}
}

func getViolation(update []int, rules []OrderRule, idx int) (int, bool) {
	for i := range rules {
		if rules[i][0] != update[idx] {
			continue
		}

		for j := range update[:idx] {
			if update[j] == rules[i][1] {
				return j, true
			}
		}
	}

	return 0, false
}

type OrderRule [2]int

var errInvalidUpdate = errors.New("invalid update")

func checkUpdate(update string, rules []OrderRule) ([]int, error) {
	pages := strings.Split(update, ",")
	pageNums := make([]int, len(pages))

	var violation bool

	for i := range pages {
		n, err := strconv.Atoi(pages[i])
		if err != nil {
			return nil, fmt.Errorf(
				"parse page number %d: %w",
				i+i, err)
		}

		pageNums[i] = n

		if containsAnyRuleTarget(pageNums[:i], n, rules) {
			violation = true
		}
	}

	if violation {
		return pageNums, errInvalidUpdate
	}

	return pageNums, nil
}

func containsAnyRuleTarget(haystack []int, current int, rules []OrderRule) bool {
	for i := range rules {
		if rules[i][0] != current {
			continue
		}

		for _, needle := range haystack {
			if needle == rules[i][1] {
				return true
			}
		}
	}

	return false
}
