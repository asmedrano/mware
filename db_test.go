package mware

import (
	"os"
	"testing"
)

func TestInitDB(t *testing.T) {
	defer rmDB()
	initDB("/tmp/test.db")
	// TODO: This needs an assertion
}

// Dont like to test both but it makes sense to me
func TestInsertAndQuery(t *testing.T) {
	initDB("/tmp/test.db")
	db, _ := getDb("/tmp/test.db")
	defer rmDB()
	defer db.Close()
	r := RowVal{
		Date:        "7/15/2014",
		Amount:      "100",
		Description: "TestTrans",
		Category:    "Test",
	}

	v := []RowVal{r}
	insertRows(db, v)

	rows := getRows(db)
	if rows[0].Id != "1" {
		t.Error("Expected id to be 1")
	}

}

func rmDB() {
	os.Remove("/tmp/test.db")
}
