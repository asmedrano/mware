package mware

import (
	"os"
	"testing"
)

func TestKeyExists(t *testing.T) {
	db, _ := getDb("/tmp/test.db")
	defer db.Close()
	//defer rmDB("/tmp/test.db")
	r := RowVal{
		Date:        "7/15/2014",
		Amount:      "100",
		Description: "TestTrans",
		Category:    "Test",
		Key:         "ABC",
	}

	v := []RowVal{r}
	insertRows(db, v)
	k, _ := keyExists(db, "ABC")
	if k != true {
        t.Error()
    }

	nk, _ := keyExists(db, "999")
	if nk != false {
        t.Error()
    }
}

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
		Key:         "123",
	}

	v := []RowVal{r}
	insertRows(db, v)
	rows := getRows(db)
	if rows[0].Id != "1" {
		t.Error("Expected id to be 1")
	}

}


func TestSameKeyInsert(t *testing.T) {
	db, _ := getDb("/tmp/test.db")
	defer db.Close()
	defer rmDB("/tmp/test.db")
	r := RowVal{
		Date:        "7/15/2014",
		Amount:      "100",
		Description: "TestTrans",
		Category:    "Test",
		Key:         "123",
	}

	v := []RowVal{r}
    i, _ := insertRows(db, v)
    if i != 1 {
        t.Error()
    }
    _, o := insertRows(db, v)
    if o != 1 {
        t.Error()
    }
}


func rmDB(path string) {
	os.Remove(path)
}
