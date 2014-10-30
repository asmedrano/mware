package main

import (
	"flag"
	"github.com/asmedrano/mware/mware"
	"log"
	"strings"
	"fmt"
)

func main() {
	task := flag.String("t", "import", "What task to run. Options are <import>")
	// TODO it would be nice to get per task from always being declared
	importType := flag.String("b", "simple", "Document Source Bank i.e Simple | CapOne")
	docPath := flag.String("p", "example.csv", "Path to document")
	dbPath := flag.String("d", "transactions.db", "Path to db file")
    startDate := flag.String("start", "", "Start Date, used when displaying transactions")
    endDate := flag.String("end", "", "End Date, used when displaying transactions")
	flag.Parse()
	// TODO: validate task input
	if *task == "import" {
		db, err := mware.GetDb(*dbPath)
		if err != nil {
			log.Fatal("Could not open db")
		}
		defer db.Close()

		iT := strings.ToLower(*importType)

		if iT == "simple" {
			log.Println("Importing Simple Bank CSV...")
			i := mware.SimpleImporter{}
			i.Import(*docPath, db)
		} else if iT == "capone" {
			log.Println("Importing CapOne OFX...")

			i := mware.CapOneImporter{}
			i.Import(*docPath, db)
		}

		log.Print("Done!")
	} else if *task == "show" {
        log.Printf("Listing Transactions -- Staring from: %v", *startDate)
		db, err := mware.GetDb(*dbPath)
		if err != nil {
			log.Fatal("Could not open db")
		}
		defer db.Close()
        results, err := mware.GetResultsFilterDate(db, strings.Trim(*startDate, " "), strings.Trim(*endDate, " "))
        if err == nil {
            for i := range(results) { 
                fmt.Print(results[i])
            }
        }else{
            log.Println(err)
        }

	}

}
