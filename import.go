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
		log.Print("Error importing " + path)
		return
	}

	s.Save(data)
}

func (s *SimpleImporter) Save(data CSVData) {
	db, err := getDb("/tmp/transactions.db") // TODO: This should end up in some sort of config var
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

}
