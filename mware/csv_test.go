package mware

import (
	"io/ioutil"
	"os"
	"testing"
)

// Test Read method
func TestRead(t *testing.T) {
	writeSampleSimple(t)
	defer cleanupTestFiles()
	data, _ := Read("/tmp/testdir/simple.csv")
	if len(data.Header) == 0 {
		t.Error("Header should not be 0")
	}

	if i, _ := data.GetFieldIndex("Amount"); i != 3 {
		t.Error("Amount should be at index 3")
	}

	if _, exists := data.GetFieldIndex("boogieboogieboo"); exists != false {
		t.Error("This field does not exist and should return false")
	}

	val, _ := data.GetVal("Date", data.Results[0])

	if val != "7/28/14" {
		t.Error("Date value is not correct")
	}

}

// The Date format doesnt match the real thing wich looks like this 2014/08/04. Its ok cause tests wont fail and ultimately it gets converted properly where it counts, this can be considered some dummy text
func writeSampleSimple(t *testing.T) {
	s := `Date,Recorded at,Scheduled for,Amount,Activity,Pending,Raw description,Description,Category folder,Category,Street address,City,     *State,Zip,Latitude,Longitude,Memo
7/28/14,7/28/14 7:41,,-100,ACH,FALSE,Electronic Funds Transfer ,Transfer,Financial,Credit Card Payment,,,,,,,
7/29/14,7/29/14 7:38,,1200,ACH,FALSE,Electronic Funds Transfer ,Direct Dep,Inome,Other,,,,,,,
7/30/14,7/30/14 13:25,,-50,Signature purchase,FALSE,Some Store With drew money,DR. Zoidbergs Store,Food & Drink,Groceries,,,,,,,
7/31/14,7/31/14 7:33,,-76.75,ACH,FALSE,A bill,Bill Company,Utilities,Electricity,,,,,,,`
	err := writeTestFile([]byte(s), "simple.csv")
	if err != nil {
		t.Log(err)
	}
}

func writeSampleCapOne(t *testing.T) {
	s := `"Date","No.","Description","Debit","Credit"
"10/1/2014","4739","NETFLIX.COM NETFLIX.COM CA 95032 US","7.99",""
"10/1/2014","4739","SHORE'S MARKET NORTH PROVIDE RI 02904 US","2.50",""
"10/2/2014","4739","CVS 1026 PARK 18207225 WOONSOCKET RI 02895 US","2.15",""
"10/2/2014","4739","CVS 1026 PARK 18207225 WOONSOCKET RI 02895 US","3.88",""
"10/2/2014","4739","TARGET 00014043 SMITHFIELD RI 02917 US","32.71",""`
	err := writeTestFile([]byte(s), "cap.csv")
	if err != nil {
		t.Log(err)
	}

}

func writeTestFile(data []byte, filename string) error {
	// create file that declares functions
	os.MkdirAll("/tmp/testdir", 0777)
	return ioutil.WriteFile("/tmp/testdir/"+filename, data, 0644)
}

// Clean up test garbage files
func cleanupTestFiles() {
	os.RemoveAll("/tmp/testdir")
}
