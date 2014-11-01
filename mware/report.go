// Some prepackaged reports
package mware

import (
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// TODO: Aggregate functions
/*
   - Total Income, as filter
   - Total Expenses, as filters
   - Unique Vendors
   - Expense By Vendor
   - Biggest Expense
   - Filter all results
*/

// Filter transactions by a start and end date with optional modifiers
// start and end date are formatted like this mm-dd-yyyy.
// start date cannot be ""
// end date can also be "". Which just treats it as no upper limit
// filters is a slice of field<op>? ex: date > 100099939
// filterArgs is a slice of interfaces that gets passed the the query
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

// Get All credits from start date to end date
func GetCreditsFilterDate(db *sql.DB, start string, end string, filters []string, filterArgs []interface{}) ([]RowVal, error) {
	f := filters
	f = append(f, "CAST(amount as float) > ?")
	fA := filterArgs
	fA = append(fA, 0)
	return GetResultsFilterDate(db, start, end, f, fA)
}

// Get All debits  from start date to end date
func GetDebitsFilterDate(db *sql.DB, start string, end string, filters []string, filterArgs []interface{}) ([]RowVal, error) {
	f := filters
	f = append(f, "CAST(amount as float) < ?")
	fA := filterArgs
	fA = append(fA, 0)
	return GetResultsFilterDate(db, start, end, f, fA)
}

// Get the Total Amount of a slice of RowVal
func Total(r []RowVal) float64 {
	total := 0.00
	var v float64
	var err error
	for i := range r {
		v, err = strconv.ParseFloat(r[i].Amount, 64)
		if err == nil {
			total += v
		}
	}
	return total
}

// Get the Biggest transaction
func Max(r []RowVal) RowVal {
	m := 0
	lastMax := 0.00
	var v float64
	var err error
	var abs float64

	for i := range r {
		v, err = strconv.ParseFloat(r[i].Amount, 64)
		if err == nil {
			abs = math.Abs(v)
			if abs > lastMax {
				m = i
				lastMax = abs
			}
		}

	}

	return r[m]
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
