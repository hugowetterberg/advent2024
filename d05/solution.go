package d05

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func SolutionOne(input io.Reader) error {
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
		updateCount int
		middleSum   int
	)

	for sc.Scan() {
		updateCount++

		mid, err := checkUpdate(sc.Text(), rules)
		if errors.Is(err, errInvalidUpdate) {
			continue
		} else if err != nil {
			return fmt.Errorf("verify update %d: %w",
				updateCount, err)
		}

		middleSum += mid
	}

	fmt.Printf("Sum of middles: %d\n", middleSum)

	return nil
}

type OrderRule [2]int

var errInvalidUpdate = errors.New("invalid update")

func checkUpdate(update string, rules []OrderRule) (int, error) {
	pages := strings.Split(update, ",")
	pageNums := make([]int, len(pages))

	for i := range pages {
		n, err := strconv.Atoi(pages[i])
		if err != nil {
			return 0, fmt.Errorf(
				"parse page number %d: %w",
				i+i, err)
		}

		pageNums[i] = n

		if containsAnyRuleTarget(pageNums[:i], n, rules) {
			return 0, errInvalidUpdate
		}
	}

	// Integer division of length should put us in the middle.
	return pageNums[len(pageNums)/2], nil
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
