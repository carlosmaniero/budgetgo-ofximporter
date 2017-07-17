package importer

import (
	"strings"
	"testing"
	"time"

	"github.com/carlosmaniero/budgetgo-ofximporter/domain"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSpecTransactionImporter(t *testing.T) {
	Convey("Scenario: Importing the ofx file", t, func() {
		Convey("Given I've a valid ofx file", func() {
			reader := strings.NewReader(ofxExample)
			importer := TransactionOfxImporter{
				File: reader,
			}
			Convey("When I import the file", func() {
				iterator := importer.Parse()

				Convey("Then I can see the total of transactions", func() {
					So(iterator.Count(), ShouldEqual, 3)
				})

				Convey("And I can iterate over the transactions", func() {
					expectedDescription := []string{
						"automatic deposit",
						"Transfer from checking",
						"John Hancock",
					}
					expectedAmount := []float64{
						200.0,
						150,
						-99.39,
					}
					location, _ := time.LoadLocation("UTC")
					expectedDate := []time.Time{
						time.Date(2007, 3, 15, 0, 0, 0, 0, location),
						time.Date(2007, 3, 29, 0, 0, 0, 0, location),
						time.Date(2007, 7, 9, 0, 0, 0, 0, location),
					}
					transaction := domain.Transaction{}
					current := 0

					for iterator.HasNext() {
						iterator.Next(&transaction)

						So(transaction.Description, ShouldEqual, expectedDescription[current])
						So(transaction.Amount, ShouldEqual, expectedAmount[current])
						So(transaction.Date.Format(time.RFC3339Nano), ShouldEqual, expectedDate[current].Format(time.RFC3339Nano))
						current++
					}

					So(current, ShouldEqual, 3)
				})
			})
		})
	})
}

var ofxExample = `OFXHEADER:100
DATA:OFXSGML
VERSION:103
SECURITY:NONE
ENCODING:USASCII
CHARSET:1252
COMPRESSION:NONE
OLDFILEUID:NONE
NEWFILEUID:NONE

<OFX>
  <SIGNONMSGSRSV1>
    <SONRS>
      <STATUS>
        <CODE>0
        <SEVERITY>INFO
      </STATUS>
      <DTSERVER>20071015021529.000[-8:PST]
      <LANGUAGE>ENG
      <DTACCTUP>19900101000000
      <FI>
        <ORG>MYBANK
        <FID>01234
      </FI>
    </SONRS>
  </SIGNONMSGSRSV1>
  <BANKMSGSRSV1>
      <STMTTRNRS>
        <TRNUID>23382938
        <STATUS>
          <CODE>0
          <SEVERITY>INFO
        </STATUS>
        <STMTRS>
          <CURDEF>USD
          <BANKACCTFROM>
            <BANKID>987654321
            <ACCTID>098-121
            <ACCTTYPE>SAVINGS
          </BANKACCTFROM>
          <BANKTRANLIST>
            <DTSTART>20070101
            <DTEND>20071015
            <STMTTRN>
              <TRNTYPE>CREDIT
              <DTPOSTED>20070315
              <DTUSER>20070315
              <TRNAMT>200.00
              <FITID>980315001
              <NAME>DEPOSIT
              <MEMO>automatic deposit
            </STMTTRN>
            <STMTTRN>
              <TRNTYPE>CREDIT
              <DTPOSTED>20070329
              <DTUSER>20070329
              <TRNAMT>150.00
              <FITID>980310001
              <NAME>TRANSFER
              <MEMO>Transfer from checking
            </STMTTRN>
            <STMTTRN>
              <TRNTYPE>PAYMENT
              <DTPOSTED>20070709
              <DTUSER>20070709
              <TRNAMT>-99.39
              <FITID>980309001
                <CHECKNUM>1025
              <NAME>John Hancock
            </STMTTRN>
          </BANKTRANLIST>
          <LEDGERBAL>
            <BALAMT>5250.00
            <DTASOF>20071015021529.000[-8:PST]
          </LEDGERBAL>
          <AVAILBAL>
            <BALAMT>5250.00
            <DTASOF>20071015021529.000[-8:PST]
          </AVAILBAL>
        </STMTRS>
      </STMTTRNRS>
  </BANKMSGSRSV1>
</OFX>`
