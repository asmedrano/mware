package main

import (
	"flag"
	"fmt"
	"github.com/asmedrano/mware/mware"
	"log"
	"strings"
)

func main() {

	task := flag.String("t", "show", "What task to run. Options are <import|show>") // task can be used in conjuntion with task modifiers
	tm_TransTypeFilter := flag.String("tt", "", "Transaction type filter, credit|debit")
	// TODO it would be nice to get per task from always being declared
	bank := flag.String("b", "", "Document Source Bank i.e simple | capone | bofa")
	docPath := flag.String("p", "example.csv", "Path to document")
	dbPath := flag.String("d", "transactions.db", "Path to db file")
	startDate := flag.String("start", "", "Start Date, when using <show> task")
	endDate := flag.String("end", "", "End Date, when using <show>task")
	groupTransactions := flag.Bool("gt", false, "Group ALL transactions returned by thier descriptions.")
	fDescription := flag.String("desc", "", "Filter by description")

	flag.Parse()

	// TODO: validate task input
	if *task == "import" {
		db, err := mware.GetDb(*dbPath)
		if err != nil {
			log.Fatal("Could not open db")
		}
		defer db.Close()
		if *bank != "" {
			iT := strings.ToLower(*bank)

			if iT == "simple" {
				log.Println("Importing Simple Bank CSV...")
				i := mware.SimpleImporter{}
				i.Import(*docPath, db)
			} else if iT == "capone" {
				log.Println("Importing CapOne OFX...")

				i := mware.CapOneImporter{}
				i.Import(*docPath, db)
			}else if iT == "bofa" {
				log.Println("Importing BofA QFX...")

				i := mware.BofAImporter{}
				i.Import(*docPath, db)
			}

			log.Print("Done!")
		} else {
			fmt.Print("\nPlease select a bank using `-b` flag\n")
		}

	} else if *task == "show" {
		var results []mware.RowVal
		var filters = []string{}
		var filterArgs = []interface{}{}
		var groupRes = map[string][]mware.RowVal{}
		db, err := mware.GetDb(*dbPath)
		if err != nil {
			log.Fatal("Could not open db")
		}
		defer db.Close()

		// Description filtering. This can be a single item or a | delimited list
		// We need to make sure to group these into a single clause using (clause OR clause). #TODO: This is HACKY
		if strings.Trim(*fDescription, " ") != " " {
			dFilters := strings.Split(*fDescription, "|")
			clause := "OR description like ?"
            dFLen := len(dFilters)
			if dFLen > 1 {
				for i := range dFilters {
					if i == 0 {
						filters = append(filters, "(" + clause)
					}else if i == dFLen - 1 {
						filters = append(filters, clause + ")")
                    }else {
						filters = append(filters, clause)
                    }
				    filterArgs = append(filterArgs, "%"+strings.Trim(dFilters[i], " ")+"%")
				}
			} else {
				filters = append(filters, "AND description like ?")
				filterArgs = append(filterArgs, "%"+strings.Trim(dFilters[0], " ")+"%")

			}

		}
	    // Bank Filtering
		if *bank != "" {
			switch strings.ToLower(strings.Trim(*bank, " ")) {
			case "simple":
				filters = append(filters, "bank=?")
				filterArgs = append(filterArgs, "Simple Bank")
			case "capone":
				filters = append(filters, "bank=?")
				filterArgs = append(filterArgs, "CapitalOne")
			case "bofa":
				filters = append(filters, "bank=?")
				filterArgs = append(filterArgs, "BankOfAmerica")
			}
		}

		// Transaction type filtering
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
				if len(results) != 0 {
					max := mware.Max(results)
					fmt.Printf("\nLargest Transaction:%v\n", max)
					avg := total / float64(len(results))
					fmt.Printf("\nAverage Transaction Amount:%.2f\n", avg)
				}
			}

			if *groupTransactions == true {
				//TODO: Sort this list
				fmt.Print("\n-----------TRANSACTION TOTAL BY DESCRIPTION -------------------\n")
				for i := range groupRes {
					total := mware.Total(groupRes[i])
					fmt.Printf("%v: #%v, Total: %.2f\n", i, len(groupRes[i]), total)
				}
			}

			fmt.Print("\n")

		} else {
			log.Println(err)
		}

	}

}
