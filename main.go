package main

import (
	"flag"
	"fmt"
	"github.com/asmedrano/mware/mware"
	"log"
	"strings"
)

func main() {

	task := flag.String("t", "import", "What task to run. Options are <import|show>") // task can be used in conjuntion with task modifiers
	tm_TransTypeFilter := flag.String("tt", "", "Transaction type filter, credit|debit")
	// TODO it would be nice to get per task from always being declared
	importType := flag.String("b", "simple", "Document Source Bank i.e Simple | CapOne")
	docPath := flag.String("p", "example.csv", "Path to document")
	dbPath := flag.String("d", "transactions.db", "Path to db file")
	startDate := flag.String("start", "", "Start Date, when using <show> task")
	endDate := flag.String("end", "", "End Date, when using <show>task")

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
		var results []mware.RowVal
		log.Printf("Listing Transactions -- Staring from: %v", *startDate)
		var filters = []string{}
		var filterArgs = []interface{}{}
		db, err := mware.GetDb(*dbPath)
		if err != nil {
			log.Fatal("Could not open db")
		}
		defer db.Close()
		// TODO: Refactor this
		switch *tm_TransTypeFilter {
		case "credits":
			results, err = mware.GetCreditsFilterDate(db, strings.Trim(*startDate, " "), strings.Trim(*endDate, " "))
		case "debits":
			results, err = mware.GetDebitsFilterDate(db, strings.Trim(*startDate, " "), strings.Trim(*endDate, " "))
		default:
			results, err = mware.GetResultsFilterDate(db, strings.Trim(*startDate, " "), strings.Trim(*endDate, " "), filters, filterArgs)

		}
		if err == nil {
		    // Run some aggregation methods
			for i := range results {
				fmt.Print(results[i])
			}
            total := mware.Total(results)
            fmt.Print("\n--------------------------------------\n")
            fmt.Printf("Total: %.2f", total)
            max := mware.Max(results)
            fmt.Printf("\nLargest Transaction:%v", max)
		} else {
			log.Println(err)
		}

	}

}
