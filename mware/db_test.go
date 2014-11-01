package mware

import (
	"os"
	"testing"
	"time"
)

func TestKeyExists(t *testing.T) {
	db, _ := getDb("/tmp/test.db")
	defer db.Close()
	defer rmDB("/tmp/test.db")
	r := RowVal{
		Date:        time.Now().Unix(),
		Amount:      "100",
		Description: "TestTrans",
		Category:    "Test",
		Key:         "ABC",
		Bank:        "test",
		AccType:     "creditcard",
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
		Date:        time.Now().Unix(),
		Amount:      "100",
		Description: "TestTrans",
		Category:    "Test",
		Key:         "123",
		Bank:        "Test",
		AccType:     "checking",
	}

	v := []RowVal{r}
	insertRows(db, v)
	rows := getRows(db)

	if rows[0].Id != 1 {
		t.Error("Expected id to be 1")
	}

	//rows[0].GetDate() // TODO: this needs an assertion

	fRows, err := getRowsWhere(db, []string{"amount=?", "category=?"}, []interface{}{"100", "Test"})
	if err != nil {
		t.Error(err)
	}
	if fRows[0].Description != "TestTrans" {
		t.Error("Expected Description to be 'TestTrans'")
	}
}

func TestSameKeyInsert(t *testing.T) {
	db, _ := getDb("/tmp/test.db")
	defer db.Close()
	defer rmDB("/tmp/test.db")
	r := RowVal{
		Date:        time.Now().Unix(),
		Amount:      "100",
		Description: "TestTrans",
		Category:    "Test",
		Key:         "123",
		Bank:        "Test",
		AccType:     "creditcard",
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
