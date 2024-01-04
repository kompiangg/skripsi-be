package main

import (
	"fmt"
	"time"
)

func main() {
	db := []int{0, 6, 29, 89}

	paramStart, err := time.Parse("2006-01-02", "2024-01-01")
	if err != nil {
		panic(err)
	}

	paramEnd, err := time.Parse("2006-01-02", "2024-01-31")
	if err != nil {
		panic(err)
	}

	now, err := time.Parse("2006-01-02", "2024-02-01")
	if err != nil {
		panic(err)
	}

	diffStart := now.Sub(paramStart)
	diffStartInDay := int(diffStart.Hours() / 24)

	diffEnd := now.Sub(paramEnd)
	diffEndInDay := int(diffEnd.Hours() / 24)
	fmt.Println(diffStartInDay, diffEndInDay)

	var startIdx int
	var endIdx int
	isBreak := make([]bool, 2)

	for i, v := range db {
		if diffStartInDay > db[len(db)-1] {
			diffStartInDay = db[len(db)-1]
			startIdx = len(db) - 1
			isBreak[0] = true
		}

		if diffEndInDay < db[0] {
			diffEndInDay = db[0]
			endIdx = 0
			isBreak[1] = true
		}

		if diffStartInDay-v <= 0 && !isBreak[0] {
			startIdx = i
			isBreak[0] = true
		}

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
			rangeDate = append(rangeDate, now.Add(-time.Duration(int(time.Hour)*24*(diffStartInDay+1))))
		} else {
			if i == 0 {
				rangeDate = append(rangeDate, now.Add(-time.Duration(int(time.Hour)*24*diffEndInDay)))
				rangeDate = append(rangeDate, now.Add(-time.Duration(int(time.Hour)*24*(v+1))))
			} else if i == len(choosedDB)-1 {
				rangeDate = append(rangeDate, now.Add(-time.Duration(int(time.Hour)*24*(choosedDB[i-1]))))
				rangeDate = append(rangeDate, now.Add(-time.Duration(int(time.Hour)*24*(diffStartInDay+1))))
			} else {
				rangeDate = append(rangeDate, now.Add(-time.Duration(int(time.Hour)*24*(choosedDB[i-1]))))
				rangeDate = append(rangeDate, now.Add(-time.Duration(int(time.Hour)*24*(choosedDB[i]+1))))
			}
		}

		allRangeDate = append(allRangeDate, rangeDate...)
	}

	fmt.Println(allRangeDate)
}
