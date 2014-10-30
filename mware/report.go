// Some prepackaged reports
package mware

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"
)

// TODO:
// Filter by dates

// Filter transactions by a start and end date
// start and end date are formatted like this mm-dd-yyyy.
// start date cannot be ""
// end date can also be "". Which just treats it as no upper limit
func GetResultsFilterDate(db *sql.DB, start string, end string) ([]RowVal, error) {
	results := []RowVal{}
	var endDateTime time.Time
	if start == "" {
		return results, errors.New("Start Date is required")
	}
	sdp := strings.Split(start, "-") // start date parts
	if len(sdp) < 3 {
		return results, errors.New("Invalid start date format. Use mm-dd-yyyy")
	}
	sdY, _ := strconv.Atoi(sdp[2])
	sdM, _ := strconv.Atoi(sdp[0])
	sdD, _ := strconv.Atoi(sdp[1])
	edp := strings.Split(end, "-") // start date parts
	startDateTime := time.Date(sdY, getMonth(sdM), sdD, 0, 0, 0, 0, time.UTC)

	if len(edp) > 0 {
		if len(edp) < 3 {
			return results, errors.New("Invalid end date format. Use mm-dd-yyyy")
		}
		edY, _ := strconv.Atoi(edp[2])
		edM, _ := strconv.Atoi(edp[0])
		edD, _ := strconv.Atoi(edp[1])
		endDateTime = time.Date(edY, getMonth(edM), edD, 0, 0, 0, 0, time.UTC)
	}

    res, err := getRowsWhere(db, []string{"date > ?", "date < ?"}, []interface{}{startDateTime.Unix(), endDateTime.Unix()})
    if err == nil {
        return res, nil
    }

	return results, err
}
