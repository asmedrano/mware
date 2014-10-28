/*
Turns these
<STMTTRN>
<TRNTYPE>DEBIT
<DTPOSTED>20140930170000.000
<TRNAMT>-49.99
<FITID>201410011247118
<NAME>MOVEITX.COCA95032
<MEMO>SOME Memoi things
</STMTTRN>

into csv:
TRNTYPE, DTPOSTED, TRNAMT, FITID, NAME, MEMO
*/

package ofx

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strings"
	"log"
)

var IN_TAG bool

func ConvertToCSV(inputPath string, outputPath string) {
	f, err := os.Open(inputPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	scanner := bufio.NewScanner(f)
	rows := [][]string{}
	row := []string{}
	for scanner.Scan() {
		l := scanner.Text()

		if strings.Contains(l, "<STMTTRN") {
			IN_TAG = true
		}
		if strings.Contains(l, "</STMTTRN") {
			IN_TAG = false
			rows = append(rows, row)
			row = []string{} // reset row
		}
		if IN_TAG && !strings.Contains(l, "<STMTTRN") {
			row = append(row, cleanTag(l))
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	writeRows(rows, outputPath)
}

func writeRows(rows [][]string, outputPath string) {
	file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0777)
	defer file.Close()
	if err != nil {
		log.Fatalf("Could not open %v to for writing", outputPath)
	}
	csvWriter := csv.NewWriter(file) // file here satisfies the Writer Interface
	// write a header row
	csvWriter.Write([]string{"TRNTYPE", "DTPOSTED", "TRNAMT", "FITID", "NAME", "MEMO"})
	csvWriter.Flush()
	csvWriter.WriteAll(rows)
}

func cleanTag(s string) string {
	re := regexp.MustCompile("(<\\w+>|\\s{2,})")
	return re.ReplaceAllString(s, "")
}
