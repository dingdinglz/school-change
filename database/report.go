package database

import (
	"time"
)

func ReportCreateNew(change int, message string, userid int) {
	_, _ = DatabaseEngine.Table(new(ReportModel)).Insert(&ReportModel{
		Time:    time.Now(),
		Change:  change,
		Message: message,
		User:    userid,
	})
}
