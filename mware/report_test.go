package mware

import "testing"
import "time"

func TestTransactionFilters(t *testing.T) {
	db, _ := getDb("/tmp/test.db")
	defer db.Close()
	defer rmDB("/tmp/test.db")
	a := RowVal{
		Date:        time.Date(2014, time.July, 1, 0, 0, 0, 0, time.Local).Unix(),
		Amount:      "-100",
		Description: "TestTrans",
		Category:    "Test",
		Key:         "123",
		Bank:        "test",
		AccType:     "checking",
	}

	b := RowVal{
		Date:        time.Date(2014, time.July, 6, 0, 0, 0, 0, time.Local).Unix(),
		Amount:      "100",
		Description: "Income",
		Category:    "Test",
		Key:         "ABC",
		Bank:        "test",
		AccType:     "checking",
	}

	c := RowVal{
		Date:        time.Date(2014, time.July, 10, 0, 0, 0, 0, time.Local).Unix(),
		Amount:      "100",
		Description: "Credit Payment",
		Category:    "Test",
		Key:         "333",
		Bank:        "test",
		AccType:     "creditcard",
	}

	v := []RowVal{a, b, c}
	insertRows(db, v)

	res, err := GetResultsFilterDate(db, "07-02-2014", "07-10-2014", []string{}, []interface{}{})
	if err != nil {
		t.Error(err)
	}
	if len(res) != 1 {
		t.Error("res should be == 1")
	}

	res, err = GetCreditsFilterDate(db, "07-02-2014", "07-10-2014", []string{}, []interface{}{})
	if err != nil {
		t.Error(err)
	}
	
	if res[0].Id != 2 {
		t.Error("res.Id should be == 2")
	}


}
