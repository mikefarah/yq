package main

import (
  "fmt"
  "gopkg.in/yaml.v2"
  "log"
  "io/ioutil"
)

func main() {
  var raw_data, read_error = ioutil.ReadFile("sample.yaml")

  if read_error != nil {
    log.Fatalf("error: %v", read_error)
  }

  var parsed_data interface{}

  err := yaml.Unmarshal([]byte(raw_data), &parsed_data)
  if err != nil {
    log.Fatalf("error: %v", err)
  }

  fmt.Println(parsed_data)
}
