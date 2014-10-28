package mware

import (
	"os"
	"testing"
)


// Dont like to test both but it makes sense to me
func TestInsertAndQuery(t *testing.T) {
	db, _ := getDb("/tmp/test.db")
	defer db.Close()
	defer rmDB("/tmp/test.db")
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

func rmDB(path string) {
	os.Remove(path)
}
