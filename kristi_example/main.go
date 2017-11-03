package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "strings"
)

func main() {

  files, err := ioutil.ReadDir("./")
  if err != nil {
      log.Fatal(err)
  }

  csvCombinedContent := []string{}

  for _, f := range files {
      if strings.Contains(f.Name(), ".csv") {
        content, err := ioutil.ReadFile(f.Name())
      	if err != nil {
      		log.Fatal(err)
      	}
        csvCombinedContent = append(csvCombinedContent, string(content))
      }
  }

  // fmt.Println(csvCombinedContent)
  for index, record := range csvCombinedContent {
      fmt.Println(record)
  }


}
