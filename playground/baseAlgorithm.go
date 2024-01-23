package main

import (
	"fmt"
	"time"
)

func main() {
	db := []int{0, 6, 29, 89}

	paramStart, err := time.Parse("2006-01-02", "2023-12-01")
	if err != nil {
		panic(err)
	}

	paramStart = time.Date(paramStart.Year(), paramStart.Month(), paramStart.Day(), 0, 0, 0, 0, paramStart.Location())

	paramEnd, err := time.Parse("2006-01-02", "2023-12-31")
	if err != nil {
		panic(err)
	}

	paramEnd = time.Date(paramEnd.Year(), paramEnd.Month(), paramEnd.Day(), 0, 0, 0, 0, paramEnd.Location())

	now, err := time.Parse("2006-01-02", "2023-12-31")
	if err != nil {
		panic(err)
	}

	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	diffStart := now.Sub(paramStart)
	diffStartInDay := int(diffStart.Hours() / 24)

	diffEnd := now.Sub(paramEnd)
	diffEndInDay := int(diffEnd.Hours() / 24)
	fmt.Println(diffStartInDay, diffEndInDay)

	var startIdx int
	var endIdx int
	isBreak := make([]bool, 2)

	// If the latest date is out of range, then set the latest date to the latest range
	if diffStartInDay > db[len(db)-1] {
		diffStartInDay = db[len(db)-1]
		startIdx = len(db) - 1
		isBreak[0] = true
	}

	// If the earliest date is out of range, then set the earliest date to the earliest range
	if diffEndInDay < db[0] {
		diffEndInDay = db[0]
		endIdx = 0
		isBreak[1] = true
	}

	for i, v := range db {
		// If the latest date is in the range, then set the latest date to the latest range
		if diffStartInDay-v <= 0 && !isBreak[0] {
			startIdx = i
			isBreak[0] = true
		}

		// If the earliest date is in the range, then set the earliest date to the earliest range
		if diffEndInDay-v <= 0 && !isBreak[1] {
			endIdx = i
			isBreak[1] = true
		}

		if isBreak[0] && isBreak[1] {
			break
		}
	}

	choosedDB := db[endIdx : startIdx+1]
	fmt.Println(choosedDB)

	allRangeDate := []time.Time{}
	for i, v := range choosedDB {
		rangeDate := []time.Time{}

		if len(choosedDB) == 1 {
			rangeDate = append(rangeDate, now.Add(-time.Duration(int(time.Hour)*24*diffEndInDay)))
			rangeDate = append(rangeDate, now.Add(-time.Duration(int(time.Hour)*24*diffStartInDay)))
		} else {
			if i == 0 {
				rangeDate = append(rangeDate, now.Add(-time.Duration(int(time.Hour)*24*diffEndInDay)))
				rangeDate = append(rangeDate, now.Add(-time.Duration(int(time.Hour)*24*v)))
			} else if i == len(choosedDB)-1 {
				rangeDate = append(rangeDate, now.Add(-time.Duration(int(time.Hour)*24*(choosedDB[i-1]+1))))
				rangeDate = append(rangeDate, now.Add(-time.Duration(int(time.Hour)*24*(diffStartInDay))))
			} else {
				rangeDate = append(rangeDate, now.Add(-time.Duration(int(time.Hour)*24*(choosedDB[i-1]+1))))
				rangeDate = append(rangeDate, now.Add(-time.Duration(int(time.Hour)*24*v)))
			}
		}

		rangeDate[0] = time.Date(rangeDate[0].Year(), rangeDate[0].Month(), rangeDate[0].Day(), 23, 59, 59, 0, rangeDate[0].Location())
		rangeDate[1] = time.Date(rangeDate[1].Year(), rangeDate[1].Month(), rangeDate[1].Day(), 0, 0, 0, 0, rangeDate[1].Location())

		allRangeDate = append(allRangeDate, rangeDate...)
	}

	fmt.Println(allRangeDate)
}
