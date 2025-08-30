package utils

import (
	"errors"
	"time"
)

func ConvertBSToAD(bsYear, bsMonth, bsDay int) (time.Time, error) {
	data, ok := BS_TO_AD_MAP[bsYear]
	if !ok {
		return time.Time{}, errors.New("BS year out of supported range")
	}
	if bsMonth < 1 || bsMonth > 12 {
		return time.Time{}, errors.New("invalid BS month")
	}
	if bsDay < 1 || bsDay > data.DaysOnMonth[bsMonth-1] {
		return time.Time{}, errors.New("invalid BS day")
	}

	baseAD, err := time.Parse("2006-01-02", data.FirstBaisakh)
	if err != nil {
		return time.Time{}, err
	}

	days := 0
	for m := 0; m < bsMonth-1; m++ {
		days += data.DaysOnMonth[m]
	}
	days += bsDay - 1

	return baseAD.AddDate(0, 0, days), nil
}
