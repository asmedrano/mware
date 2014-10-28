package mware

import (
	"log"
	"time"
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
		log.Print("Error importing " + path, "\n", err)
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
