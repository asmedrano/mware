package mware

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"regexp"
	"time"
)

type RowVal struct {
	Id          int64
	Date        int64
	Amount      string
	Description string
	Category    string
	Key         string // a compound Key that should uniquely identify this entry
	Bank        string // the Source bank
	AccType     string // creditcard | checking
}

// Return a time.Time from the RowVal.Date int64
func (r *RowVal) GetDate() time.Time {
	return time.Unix(r.Date, 0)
}

func (r RowVal) String() string {
	// TODO: Strip out the stuff we dont want from the Date
	tm := fmt.Sprintf("%v", r.GetDate())
	return fmt.Sprintf("\n%v|%v|%v|%v|%v|%v", r.Id, tm, r.Amount, r.Description, r.Category, r.Bank)
}

func GetDb(dbname string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbname)
	// Try to create table
	// SQLite does not have a storage class set aside for storing dates and/or times. Instead, the built-in Date And Time Functions of SQLite are capable of storing dates and times as TEXT, REAL, or INTEGER values https://www.sqlite.org/lang_datefunc.html
	sql := `create table if not exists transactions (id integer not null primary key, date integer, amount real, description text, category text, key text unique, bank text, account_type text)`
	_, err = db.Exec(sql)
	if err != nil {
		log.Printf("%q: %s\n", err, sql)
	}
	// TODO: create more indexes
	sql = `create index if not exists keyidx on transactions (key)`
	_, err = db.Exec(sql)
	if err != nil {
		log.Printf("%q: %s\n", err, sql)
	}
	return db, err
}

// this is only here so i dont have to go back and change all the references to this
// TODO: Change this
func getDb(dbname string) (*sql.DB, error) {
	return GetDb(dbname)
}

// insert a list of rowvals into db in a single db transaction. The table happens to also be called transactions
// Lets assume the same transaction will never be duplicated in a single import
// TODO: POSSIBLE BUG HERE (see last comment)
func insertRows(db *sql.DB, rv []RowVal) (in int, ign int) {
	inserted := 0 // how many recored where actually inserted
	ignored := 0  // how many where ignored
	exists := false
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	sql := "insert into transactions (date, amount, description, category, key, bank, account_type) values(?,?,?,?,?,?,?)"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := range rv {
		// each record should have a unique key, so that we dont insert the same transaction in twice. We need to check for it first
		exists, _ = keyExists(db, rv[i].Key)
		if !exists {
			_, err = stmt.Exec(rv[i].Date, rv[i].Amount, rv[i].Description, rv[i].Category, rv[i].Key, rv[i].Bank, rv[i].AccType)
			if err != nil {
				log.Fatal(err)
			}
			inserted += 1
		} else {
			ignored += 1
		}
	}
	tx.Commit()
	log.Printf("Inserted %v records, %v duplicates skipped", inserted, ignored)
	return inserted, ignored
}

// Return all rows wrapped as []RowVal
// TODO: This should return []RowVal, err
func getRows(db *sql.DB) []RowVal {
	results := []RowVal{}
	rows, err := db.Query("select * from transactions order by date")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		r := RowVal{}
		rows.Scan(&r.Id, &r.Date, &r.Amount, &r.Description, &r.Category, &r.Key, &r.Bank, &r.AccType)
		results = append(results, r)
	}
	return results
}

// Return rows as []RowVal
func getRowsWhere(db *sql.DB, where []string, args []interface{}) ([]RowVal, error) {
	results := []RowVal{}
	query := "select * from transactions"
	or := regexp.MustCompile("\\bor\\b|\\bOR\\b")
    cleanQueryOps := regexp.MustCompile("OR|or|AND|and") // we are gonna strip these out in a sec
	if len(where) > 0 {
		w := ""
		for i := range where {
			if i > 0 {
				if or.FindString(where[i]) != "" {
					// Implicit AND, Ors when required
					w += fmt.Sprintf(" %v %v", "OR", cleanQueryOps.ReplaceAllString(where[i], ""))
				} else {
					w += fmt.Sprintf(" %v %v", "AND", cleanQueryOps.ReplaceAllString(where[i], ""))
				}
			} else {
				w += cleanQueryOps.ReplaceAllString(where[i], "")
			}
		}

		query += " WHERE " + w
	}
	query += " Order By date"

	stmt, err := db.Prepare(query)
	if err != nil {
		return results, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)

	if err != nil {
		return results, err
	}

	for rows.Next() {
		r := RowVal{}
		rows.Scan(&r.Id, &r.Date, &r.Amount, &r.Description, &r.Category, &r.Key, &r.Bank, &r.AccType)
		results = append(results, r)
	}
	return results, nil
}

func keyExists(db *sql.DB, key string) (bool, error) {

	stmt, err := db.Prepare("select key from transactions where key = ?")
	if err != nil {
		log.Println(err)
		return false, err
	}
	defer stmt.Close()
	var k string
	err = stmt.QueryRow(key).Scan(&k)
	if err != nil {
		// this will most likely be "sql: no rows in result set"
		return false, err
	}
	if k == key {
		return true, nil
	} else {
		return false, nil
	}
}

// Return distinct descriptions
func getDistinctDescriptions(db *sql.DB) ([]string, error) {
	results := []string{}
	rows, err := db.Query("select distinct(description) from transactions")
	if err != nil {
		return results, err
	}
	defer rows.Close()
	for rows.Next() {
		var d string
		rows.Scan(&d)
		results = append(results, d)
	}
	return results, err

}
