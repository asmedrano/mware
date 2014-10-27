package mware

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type RowVal struct {
	Id          string
	Date        string
	Amount      string
	Description string
	Category    string
}

func getDb(dbname string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbname)
	return db, err
}

// insert a list of rowvals into db in a single db transaction. The table happens to also be called transactions
func insertRows(db *sql.DB, rv []RowVal) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	sql := "insert into transactions (date, amount, description, category) values(?,?,?,?)"
	stmt, err := tx.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := range rv {
		_, err = stmt.Exec(rv[i].Date, rv[i].Amount, rv[i].Description, rv[i].Category)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}

// Return all rows wrapped as []RowVal
func getRows(db *sql.DB) []RowVal {
	results := []RowVal{}
	defer db.Close()
	rows, err := db.Query("select * from transactions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		r := RowVal{}
		rows.Scan(&r.Id, &r.Date, &r.Amount, &r.Description, &r.Category)
		results = append(results, r)
	}
	return results
}

// initialize the DB with its corresponding tables
// SQLite does not have a storage class set aside for storing dates and/or times. Instead, the built-in Date And Time Functions of SQLite are capable of storing dates and times as TEXT, REAL, or INTEGER values https://www.sqlite.org/lang_datefunc.html
func initDB(dbname string) {
	db, err := getDb(dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sql := `create table transactions (id integer not null primary key, date integer, amount real, description text, category text)`
	_, err = db.Exec(sql)
	if err != nil {
		log.Printf("%q: %s\n", err, sql)
		return
	}
}
