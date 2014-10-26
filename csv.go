package mware

import (
	"encoding/csv"
	"os"
)

type CSVData struct {
	Src        string // the source of the csv data
	Results    [][]string
	Header     []string
	fieldIndex map[string]int // a mapping to what index a result is in
}

// return field f from fieldIndex if it exist, otherwise return index, false
func (c *CSVData) GetFieldIndex(f string) (int, bool) {
	val, ok := c.fieldIndex[f]
	return val, ok
}

// return value from row, by getting the proper field index first
func (c *CSVData) GetVal(fieldName string, row []string) (string, bool) {
	i, exists := c.GetFieldIndex(fieldName)
	if !exists {
		return "", false
	}
	// otherwise the fieldname is valid, return the result
	return row[i], true
}

// Read is a simple abstraction over file opening and returning results in a csv
func Read(path string) (records CSVData, err error) {
	data := CSVData{}
	file, err := os.Open(path)
	if err != nil {
		return data, err
	}
	// automatically call Close() at the end of current method
	defer file.Close()
	//
	reader := csv.NewReader(file)
	results, err := reader.ReadAll()
	if err != nil {
		return data, err
	}
	data.Results = results
	data.Src = path
	data.Header = results[0] // I guess we are assuming there is always a header
	data.fieldIndex = makeFieldIndex(data.Header)
	return data, nil
}

// given a []string, return a map that has string as key and index as val
func makeFieldIndex(s []string) map[string]int {
	m := map[string]int{}
	for i := range s {
		m[s[i]] = i
	}
	return m
}
