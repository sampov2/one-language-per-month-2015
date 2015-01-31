package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "regexp"
)

const debug = true

func getUrlAsString(url string, depth int, maxDepth int, messages chan<- string) {
  var resp *http.Response


  if debug { fmt.Println("Retrieving", url)}
  resp, err := http.Get(url)
  if err != nil {
    messages <- "-close"
    return
  }

  defer resp.Body.Close()
  if debug { fmt.Println("Reading response", url)}
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    messages <- "-close"
    return
  }

  if debug { fmt.Println("Read", len(body), "bytes from", url) }

  urls := parseString(string(body))
  //fmt.Println("urls", urls)
  if (depth < maxDepth) {
    for i:= 0; i < len(urls); i++ {
      messages <- "+open"
      go getUrlAsString(urls[i], depth+1, maxDepth, messages)
    }
  }
  messages <- "-close"
  return

}

func parseString(input string) []string {
  parseStringRegexp := regexp.MustCompile("https?://[^\"\\s]+")
  matches := parseStringRegexp.FindAllString(input, -1)
  if (matches == nil) {
    matches = []string{}
  }

  return matches
}

func main() {
  messages := make(chan string, 1)
  counter := 0

  messages <- "+open"
  go getUrlAsString("http://spatineo.com", 0, 1, messages)
  for true {
    msg := <-messages
    if (msg == "+open") {
      if debug { fmt.Println("open!") }
      counter++;
    } else if (msg == "-close") {
      if debug { fmt.Println("close..") }
      counter--;
      if (counter == 0) {
        if debug { fmt.Println("done") }
        break
      }
    }
  }

}
