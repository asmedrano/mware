package mware

import (
	"testing"
)

func TestSimpleImport(t *testing.T) {
	writeSampleSimple(t)
	db, _ := getDb("/tmp/transactions.db")
	i := SimpleImporter{}
	i.Import("/tmp/testdir/simple.csv")
    rows := getRows(db)
    if rows[0].Amount != "-100" {
        t.Error("Amount != -100")
    }
	defer db.Close()
	defer rmDB("/tmp/transactions.db")
	defer cleanupTestFiles()

}

func TestCapImport(t *testing.T) {
	writeSampleCapOne(t)
	db, _ := getDb("/tmp/transactions.db")
	/*
    i := CapOneImporter{}
    //i.Import("/tmp/testdir/cap.csv")
    rows := getRows(db)
    if rows[0].Amount != "-7.99" {
        t.Error("Amount != -7.99") 
    }*/
	defer db.Close()
	defer rmDB("/tmp/transactions.db")
	defer cleanupTestFiles()

}
