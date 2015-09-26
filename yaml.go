package main

import (
  "fmt"
  "gopkg.in/yaml.v2"
  "log"
)

var data = `
a: Easy!
b:
  c: 2
  d: [3, 4]
`

func main() {
  var m interface{}

  err := yaml.Unmarshal([]byte(data), &m)
  if err != nil {
          log.Fatalf("error: %v", err)
  }

  fmt.Println("Hello, 世界")
  fmt.Println(m)
}
