package domain

import "time"

type HistoryElem struct {
	UserID    string
	Slug      string
	Operation string
	DateTime  time.Time
}
