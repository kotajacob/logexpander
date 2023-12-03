package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

const journalDir = "/home/kota/docs/memex/journal"

// dates maps a date in the form 2006-01-02 to a short optional message.
type dates map[string]string

// list of lines in the journal to be printed.
type journal []string

func (x journal) Len() int           { return len(x) }
func (x journal) Less(i, j int) bool { return x[i] < x[j] }
func (x journal) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func main() {
	dates := getDates()
	journal := sortDates(dates)
	journal = padJournal(journal)
	for _, line := range journal {
		fmt.Println(line)
	}
}

func getDates() map[string]string {
	// Use journal entries to fill out dates.
	dates := make(map[string]string)
	entries, err := os.ReadDir(journalDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed journal entries:", err)
		os.Exit(1)
	}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		date := strings.TrimSuffix(e.Name(), ".md")
		dates[date] = ""
	}

	// Update messages with given journal input.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		date, ok := hasDate(line)
		if !ok {
			continue
		}
		dateStr := date.Format("2006-01-02")
		msg := strings.TrimPrefix(line, dateStr)
		msg = strings.TrimPrefix(msg, " - ")
		dates[dateStr] = msg
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "failed reading STDIN:", err)
		os.Exit(1)
	}
	return dates
}

// hasDate checks if a line begins with a date in the form 2006-01-02 and
// returns the date string and true if it does.
func hasDate(line string) (time.Time, bool) {
	before, _, _ := strings.Cut(line, " ")
	date, err := time.ParseInLocation(
		"2006-01-02",
		before,
		time.Now().Location(),
	)
	if err != nil {
		return date, false
	}
	return date, true
}

func sortDates(dates dates) journal {
	var journal journal
	for d, m := range dates {
		if m == "" {
			journal = append(journal, d)
		} else {
			journal = append(journal, d+" - "+m)
		}
	}
	sort.Sort(journal)
	return journal
}

func padJournal(unpadded journal) journal {
	var padded journal

	var last, current int
	for _, l := range unpadded {
		current = weekNumber(l)
		if current == 0 {
			// Ignore lines that do not begin with a date.
			continue
		}
		if last != current {
			padded = append(padded, "")
		}
		padded = append(padded, l)
		last = current
	}

	return padded
}

// Returns the week number for a line in the journal: 1-53.
// If the line does not begin with a date we return 0.
//
// ISOWeek returns the ISO 8601 year and week number in which t occurs. Week
// ranges from 1 to 53. Jan 01 to Jan 03 of year n might belong to week 52 or
// 53 of year n-1, and Dec 29 to Dec 31 might belong to week 1 of year n+1.
func weekNumber(line string) int {
	date, ok := hasDate(line)
	if !ok {
		return 0
	}
	_, week := date.ISOWeek()
	return week
}
