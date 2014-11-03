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
	bank := flag.String("b", "simple", "Document Source Bank i.e Simple | CapOne")
	docPath := flag.String("p", "example.csv", "Path to document")
	dbPath := flag.String("d", "transactions.db", "Path to db file")
	startDate := flag.String("start", "", "Start Date, when using <show> task")
	endDate := flag.String("end", "", "End Date, when using <show>task")
	groupTransactions := flag.Bool("gt", false, "Group Transactions by Descriptions")

	flag.Parse()

	// TODO: validate task input
	if *task == "import" {
		db, err := mware.GetDb(*dbPath)
		if err != nil {
			log.Fatal("Could not open db")
		}
		defer db.Close()

		iT := strings.ToLower(*bank)

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
		var groupRes = map[string][]mware.RowVal{}
		db, err := mware.GetDb(*dbPath)
		if err != nil {
			log.Fatal("Could not open db")
		}
		defer db.Close()

		switch strings.Trim(*bank, " ") {
		case "Simple":
			filters = append(filters, "bank=?")
			filterArgs = append(filterArgs, "Simple Bank")
		case "CapOne":
			filters = append(filters, "bank=?")
			filterArgs = append(filterArgs, "CapitalOne")
		}

		switch *tm_TransTypeFilter {
		case "credits":
			results, err = mware.GetCreditsFilterDate(db, strings.Trim(*startDate, " "), strings.Trim(*endDate, " "), filters, filterArgs)
			if *groupTransactions == true {
				groupRes, err = mware.GroupVendorCredits(db, strings.Trim(*startDate, " "), strings.Trim(*endDate, " "), filters, filterArgs)
			}
		case "debits":
			results, err = mware.GetDebitsFilterDate(db, strings.Trim(*startDate, " "), strings.Trim(*endDate, " "), filters, filterArgs)
			if *groupTransactions == true {
				groupRes, err = mware.GroupVendorDebits(db, strings.Trim(*startDate, " "), strings.Trim(*endDate, " "), filters, filterArgs)
			}
		default:
			results, err = mware.GetResultsFilterDate(db, strings.Trim(*startDate, " "), strings.Trim(*endDate, " "), filters, filterArgs)

		}
		if err == nil {
			// Run some aggregation methods
			for i := range results {
				fmt.Print(results[i])
			}

			if *tm_TransTypeFilter != "" {
				total := mware.Total(results)
				fmt.Print("\n--------------------------------------\n")
				fmt.Printf("Total: %.2f", total)

				max := mware.Max(results)
				fmt.Printf("\nLargest Transaction:%v\n", max)
			}

			if *groupTransactions == true {
                //TODO: Sort this list
			    fmt.Print("\n-----------TRANSACTION TOTAL BY DESCRIPTION -------------------\n")
                for i := range groupRes {
                    total := mware.Total(groupRes[i])
                    fmt.Printf("%v: #%v, Total: %.2f\n", i, len(groupRes[i]), total)
                }
			}

		} else {
			log.Println(err)
		}

	}

}
