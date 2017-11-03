package main

// go run main.go [access key]

import (
  "encoding/csv"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "strconv"
  "strings"
  // "net/http"
  // "oauth2sample/config"
  // "oauth2sample/handlers"
  quickbooks "github.com/jinmatt/go-quickbooks.v2"
)

// creditCardBill - Associate employee to their statement amount.
type creditCardBill struct {
  Employee     string  `json:"employee"`
  StatementAmt float64 `json:"amount"`
}

//Account selection based on employee. Must add employees here when issued card.
func employeeAccountLookup(employee string) (accountID string, accountName string) {
  switch employee {
  case "Luke Johnson":
    return "231", "L Johnson Credit Card Chase"
  case "Tom Furey":
    return "232", "T Furey Credit Card Chase"
  case "Steve Herscleb":
    return "233", "S Herschleb Credit Card Chase"
  case "Nicole Farrar":
    return "230", "N Farrar Credit Card Chase"
  case "Ryan Brennan":
    return "242", "Ryan Brennan"
  }

  return "00", "Employee not found."

}

func main() {

  files, err := ioutil.ReadDir("./csv/")
  if err != nil {
    log.Fatal(err)
  }

  csvCombinedContent := []string{}
  creditsCSV := []string{}
  total := 0.00
  aggregateStatements := []creditCardBill{}

  //Loops thru csv files in the ./csv/ folder, reads each csv file.
  for index, f := range files {
    if strings.Contains(f.Name(), ".csv") {
      content, err := ioutil.ReadFile(fmt.Sprintf("./csv/%s", f.Name()))
      if err != nil {
        log.Fatal(err)
      }
      reader := csv.NewReader(strings.NewReader(string(content)))

      records, err := reader.ReadAll()
      if err != nil {
        log.Fatal(err)
      }

      sumOfPurchases := 0.00
      //Add headers
      if index == 0 {
        csvCombinedContent = append(csvCombinedContent, strings.Join(records[0], ","))
        creditsCSV = append(creditsCSV, strings.Join(records[0], ","))
      }

      employee := ""
      //For each record in file, sum up purchases and append to aggregateStatements array.
      for index, record := range records[1:] {
        if index == 0 {
          employee = record[5]
          fmt.Println(record[5])
        }
        if record[0] == "Sale" {
          csvCombinedContent = append(csvCombinedContent, strings.Join(record, ","))
          floatValue, err := strconv.ParseFloat(record[4], 64)

          if err != nil {
            log.Fatal(err)
          }

          sumOfPurchases = sumOfPurchases + (floatValue * -1.00)

        }
        //Add returns and reversals to creditsCSV array for later use.
        if record[0] == "Reversal" || record[0] == "Return" {
          creditsCSV = append(creditsCSV, strings.Join(record, ","))
        }
      }
      total = sumOfPurchases + total

      currentStatement := creditCardBill{
        Employee:     employee,
        StatementAmt: sumOfPurchases,
      }
      aggregateStatements = append(aggregateStatements, currentStatement)
    }
  }
  //Check for credits that would affect journal entry.
  for _, Credits := range creditsCSV {
    fmt.Println(Credits)
  }
  fmt.Println("Total payment: ", total)
  fmt.Println("====================")
  fmt.Println("Credits")

  // For Loop for each Object in aggregate Statements create debit JE objects and append to array journalEntryLines.
  journalEntryLines := []quickbooks.Line{}
  for index, statement := range aggregateStatements {

    accountID, accountName := employeeAccountLookup(statement.Employee)

    currentStatmentLine := quickbooks.Line{
      LineID:      strconv.Itoa(index),
      Description: statement.Employee,
      Amount:      statement.StatementAmt,
      DetailType:  "JournalEntryLineDetail",
      JournalEntryLineDetail: &quickbooks.JournalEntryLineDetail{
        PostingType: "Debit",
        AccountRef: quickbooks.JournalEntryRef{
          Value: accountID,
          Name:  accountName,
        },
      },
    }

    journalEntryLines = append(journalEntryLines, currentStatmentLine)
  }

  //At this point, array has all Debit JE objects :)

  //Create credit and append to array
  totalStatementLine := quickbooks.Line{
    LineID:      strconv.Itoa(len(aggregateStatements)),
    Description: "To record payment of credit card statement.",
    Amount:      total,
    DetailType:  "JournalEntryLineDetail",
    JournalEntryLineDetail: &quickbooks.JournalEntryLineDetail{
      PostingType: "Credit",
      AccountRef: quickbooks.JournalEntryRef{
        Value: "231",
        Name:  "L Johnson Credit Card Chase",
      },
    },
  }

  journalEntryLines = append(journalEntryLines, totalStatementLine)

  //create quickbooks client
  quickbooksClient := quickbooks.NewClient(os.Getenv("QBO_realmID_production"), os.Args[3], false)

  journalEntry := quickbooks.Journalentry{
    Line: journalEntryLines,
  }

  JournalentryObject, err := quickbooksClient.CreateJE(journalEntry)

  fmt.Println(JournalentryObject, err)

}
