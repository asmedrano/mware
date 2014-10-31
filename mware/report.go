// Some prepackaged reports
package mware

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"
	"fmt"
)

// TODO: Aggregate functions
/*
   - Total Income
   - Total Expenses, (filters)
   - Unique Vendors
   - Expense By Vendor
   - Biggest Expense
*/

// Filter transactions by a start and end date
// start and end date are formatted like this mm-dd-yyyy.
// start date cannot be ""
// end date can also be "". Which just treats it as no upper limit
func GetResultsFilterDate(db *sql.DB, start string, end string, filters []string, filterArgs []interface{}) ([]RowVal, error) {
	results := []RowVal{}
    filter := filters
    args := filterArgs
	if start == "" {
		return results, errors.New("Start Date is required")
	}

	startDateTime, err := ParseDateString(start)
	if err != nil {
		return results, err
	}
	filter = append(filter, "date > ?")
	args = append(args, startDateTime.Unix())

	if end != "" {
		endDateTime, err := ParseDateString(end)
        if err != nil {
            return results, err
        }
		filter = append(filter, "date < ?")
		args = append(args, endDateTime.Unix())
	}

	res, err := getRowsWhere(db, filter, args)
	if err == nil {
		return res, nil
	}

	return results, err
}

// Parse a date that looks like this  mm-dd-yyyy
func ParseDateString(date string) (time.Time, error) {
	dateParts := strings.Split(date, "-") // start date parts
	if len(dateParts) < 3 {
		return time.Time{}, errors.New(fmt.Sprintf("Invalid date format %v. Use mm-dd-yyyy", date))
	}
	y, err := strconv.Atoi(dateParts[2])
	m, err := strconv.Atoi(dateParts[0])
	d, err := strconv.Atoi(dateParts[1])

	if err != nil {
		return time.Time{}, errors.New("Invalid date format. Use mm-dd-yyyy")
	}

	return time.Date(y, getMonth(m), d, 0, 0, 0, 0, time.Local), nil

}
