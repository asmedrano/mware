package mware

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"github.com/asmedrano/mWare/ofx"
	"io"
	"log"
	"path"
	"strconv"
	"strings"
	"time"
)

type Importer interface {
	Import() // Extract and Transform Data here
	Save()
	Date() int64 // should return UNIX time
	Bank() string
}

type SimpleImporter struct {
	ImportTime time.Time // when this got imported
}

// In theory maybe something could happen here but I think all importers will look like this
func (s *SimpleImporter) Import(path string, db *sql.DB) {
	data, err := Read(path)
	if err != nil {
		log.Print("Error importing "+path, "\n", err)
		return
	}

	vals := []RowVal{}

	// Transform data from simple to our nomalized version
	for i := range data.Results {
		date, _ := data.GetVal("Date", data.Results[i])
		amount, _ := data.GetVal("Amount", data.Results[i])
		description, _ := data.GetVal("Description", data.Results[i])
		category, _ := data.GetVal("Category", data.Results[i])
		recorded_at, _ := data.GetVal("Recorded at", data.Results[i])
		vals = append(vals, RowVal{
			Date:        s.Date(date),
			Amount:      amount,
			Description: description,
			Category:    category,
			Key:         s.MakeKey(date + recorded_at + amount + description),
			Bank:        s.Bank(),
		})
	}

	s.Save(db, vals)
}

func (s *SimpleImporter) Save(db *sql.DB, data []RowVal) {
	insertRows(db, data)
}

func (s *SimpleImporter) MakeKey(raw string) string {
	h := md5.New()
	io.WriteString(h, raw)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (s *SimpleImporter) Date(raw string) int64 {
	// raw looks like: 2014/08/04
	dp := strings.Split(raw, "/")
	year, _ := strconv.Atoi(dp[0])  // lets assume all these transactions happend in the year 2000+
	month, _ := strconv.Atoi(dp[1]) // lets assume all these transactions happend in the year 2000+
	day, _ := strconv.Atoi(dp[2])   // lets assume all these transactions happend in the year 2000+
	t := time.Date(year, getMonth(month), day, 0, 0, 0, 0, time.Local)
	return t.Unix()

}

func (s *SimpleImporter) Bank() string {
	return "Simple Bank"
}

type CapOneImporter struct {
	ImportTime time.Time // when this got imported
}

// Cap One importores should convert .ofx dumps to csv first
func (s *CapOneImporter) Import(ofxPath string, db *sql.DB) {
	tempPath := fmt.Sprintf("/tmp/%v.csv", path.Base(ofxPath))
	ofx.ConvertToCSV(ofxPath, tempPath)
	data, err := Read(tempPath)
	if err != nil {
		log.Print("Error importing "+ofxPath, "\n", err)
		return
	}

	vals := []RowVal{}

	// Transform data from simple to our nomalized version
	//FIELDS: TRNTYPE, DTPOSTED, TRNAMT, FITID, NAME, MEMO
	for i := range data.Results {
		date, _ := data.GetVal("DTPOSTED", data.Results[i])
		description, _ := data.GetVal("NAME", data.Results[i])
		amount, _ := data.GetVal("TRNAMT", data.Results[i])
		key, _ := data.GetVal("FITID", data.Results[i])
		vals = append(vals, RowVal{
			Date:        s.Date(date),
			Amount:      amount,
			Description: description,
			Category:    "", // TODO: need to parse description to get this
			Key:         key,
			Bank:        s.Bank(),
		})
	}

	s.Save(db, vals)
}

func (s *CapOneImporter) Save(db *sql.DB, data []RowVal) {
	insertRows(db, data)
}

func (s *CapOneImporter) Date(raw string) int64 {
	// raw looks like: 20140930170000.000
	year, _ := strconv.Atoi(raw[:4])   // lets assume all these transactions happend in the year 2000+
	month, _ := strconv.Atoi(raw[4:6]) // lets assume all these transactions happend in the year 2000+
	day, _ := strconv.Atoi(raw[6:8])   // lets assume all these transactions happend in the year 2000+
	t := time.Date(year, getMonth(month), day, 0, 0, 0, 0, time.Local)
	return t.Unix()

}

func (s *CapOneImporter) Bank() string {
	return "CapitalOne"
}

func getMonth(num int) time.Month {
	var m time.Month
	switch num {
	case 1:
		m = time.January
	case 2:
		m = time.February
	case 3:
		m = time.March
	case 4:
		m = time.April
	case 5:
		m = time.May
	case 6:
		m = time.June
	case 7:
		m = time.July
	case 8:
		m = time.August
	case 9:
		m = time.September
	case 10:
		m = time.October
	case 11:
		m = time.November
	case 12:
		m = time.December
	}
	return m
}
