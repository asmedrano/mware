package mware

import "testing"
import "time"

func TestTransactionFilter(t *testing.T) {
	db, _ := getDb("/tmp/test.db")
	defer db.Close()
	defer rmDB("/tmp/test.db")
	a := RowVal{
		Date:        time.Date(2014, time.July, 1, 0, 0, 0, 0, time.UTC).Unix(),
		Amount:      "100",
		Description: "TestTrans",
		Category:    "Test",
		Key:         "123",
		Bank:        "test",
	}
	b := RowVal{
		Date:        time.Date(2014, time.July, 5, 0, 0, 0, 0, time.UTC).Unix(),
		Amount:      "100",
		Description: "TestTrans",
		Category:    "Test",
		Key:         "ABC",
		Bank:        "test",
	}

	c := RowVal{
		Date:        time.Date(2014, time.July, 10, 0, 0, 0, 0, time.UTC).Unix(),
		Amount:      "100",
		Description: "TestTrans",
		Category:    "Test",
		Key:         "333",
		Bank:        "test",
	}

	v := []RowVal{a, b, c}
	insertRows(db, v)

	res, err := GetResultsFilterDate(db, "07-01-2014", "07-11-2014")
	t.Log(err)
	if err == nil {
		t.Log(res)
	}

}
