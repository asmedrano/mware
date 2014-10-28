package mware

import (
	"github.com/asmedrano/mWare/ofx"
	"log"
	"path"
	"time"
	"fmt"
)

type Importer interface {
	Import() // Extract and Transform Data here
	Save()
}

type SimpleImporter struct {
	ImportTime time.Time // when this got imported
}

// In theory maybe something could happen here but I think all importers will look like this
func (s *SimpleImporter) Import(path string) {
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
		vals = append(vals, RowVal{
			Date:        date,
			Amount:      amount,
			Description: description,
			Category:    category,
		})
	}

	s.Save(vals)
}

func (s *SimpleImporter) Save(data []RowVal) {
	db, err := getDb("/tmp/transactions.db") // TODO: This should end up in some sort of config var
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	// save the rows
	insertRows(db, data)

}

type CapOneImporter struct {
	ImportTime time.Time // when this got imported
}

// Cap One importores should convert .ofx dumps to csv first
func (s *CapOneImporter) Import(ofxPath string) {
    tempPath := fmt.Sprintf("/tmp/%v.csv", path.Base(ofxPath))
	ofx.ConvertToCSV(ofxPath, tempPath)
	data, err := Read(tempPath)
	if err != nil {
		log.Print("Error importing "+ ofxPath, "\n", err)
		return
	}

	vals := []RowVal{}

	// Transform data from simple to our nomalized version
    //FIELDS: TRNTYPE, DTPOSTED, TRNAMT, FITID, NAME, MEMO
	for i := range data.Results {
		date, _ := data.GetVal("DTPOSTED", data.Results[i])
		description, _ := data.GetVal("NAME", data.Results[i])
		amount, _ := data.GetVal("TRNAMT", data.Results[i])
		vals = append(vals, RowVal{
			Date:        date,
			Amount:      amount,
			Description: description,
			Category:    "", // TODO: need to parse description to get this
		})
	}

	s.Save(vals)
}

func (s *CapOneImporter) Save(data []RowVal) {
	db, err := getDb("/tmp/transactions.db") // TODO: This should end up in some sort of config var
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	// save the rows
	insertRows(db, data)

}
