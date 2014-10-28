package ofx

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCreateCSVFromOFX(t *testing.T) {
	writeSampleCapOneOFX(t)
	ConvertToCSV("/tmp/testdir/capone.ofx", "/tmp/testdir/capone.ofx.cvs")
	f, err := ioutil.ReadFile("/tmp/testdir/capone.ofx.cvs")
	if err != nil {
		t.Error(err)
	}

	output := string(f)
    // TODO: needs assertion. Just gonna print it for now
    t.Log(output)

	defer cleanupTestFiles()
}

func writeSampleCapOneOFX(t *testing.T) {
	// Sucks to put this here.
	s := `<OFX>
        <SIGNONMSGSRSV1>
                <SONRS>
                        <STATUS>
                                <CODE>0
                                <SEVERITY>INFO
                        </STATUS>
                        <DTSERVER>20141028161034.051
                        <LANGUAGE>ENG
                        <DTPROFUP>20050531040000.000
                        <FI>
                                <ORG>C1
                                <FID>1001
                        </FI>
                        <INTU.BID>0000
                        <INTU.USERID>nobodymcgee
                </SONRS>
        </SIGNONMSGSRSV1>
        <CREDITCARDMSGSRSV1>
                <CCSTMTTRNRS>
                        <TRNUID>0
                        <STATUS>
                                <CODE>0
                                <SEVERITY>INFO
                        </STATUS>
                        <CCSTMTRS>
                                <CURDEF>USD
                                <CCACCTFROM>
                                        <ACCTID>1234567892567
                                </CCACCTFROM>
                                <BANKTRANLIST>
                                        <DTSTART>20141001170000.000
                                        <DTEND>20141028170000.000
                                        <STMTTRN>
                                                <TRNTYPE>DEBIT
                                                <DTPOSTED>20140930170000.000
                                                <TRNAMT>-7.99
                                                <FITID>201410011247118
                                                <NAME>N.COM N.COM CA 95032
                                                <MEMO>2349: N.COM N.COM CA 95032 US
                                        </STMTTRN>
                                        <STMTTRN>
                                                <TRNTYPE>DEBIT
                                                <DTPOSTED>20140930170000.000
                                                <TRNAMT>-2.50
                                                <FITID>201410011247128
                                                <NAME>BILLS MARKET PROVIDE RI
                                                <MEMO>1234: BILL's MARKET PROVIDE RI 02904 US
                                        </STMTTRN>
                                </BANKTRANLIST>
                                <LEDGERBAL>
                                        <BALAMT>-242.34
                                        <DTASOF>20141028161034.051
                                </LEDGERBAL>
                                <AVAILBAL>
                                        <BALAMT>3839.34
                                        <DTASOF>20141028161034.051
                                </AVAILBAL>
                        </CCSTMTRS>
                </CCSTMTTRNRS>
        </CREDITCARDMSGSRSV1>
</OFX>`
	err := writeTestFile([]byte(s), "capone.ofx")
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
