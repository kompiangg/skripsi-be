package config

import "time"

type Date struct {
	UsingCustomNow bool      `yaml:"UsingCustomDate"`
	Current        time.Time `yaml:"Current"`
}

func (d Date) Now() time.Time {
	if d.UsingCustomNow {
		return d.Current
	}

	return time.Now()
}
