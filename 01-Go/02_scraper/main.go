package main

import (
  "fmt"
  "net/http"
  "io/ioutil"
  "regexp"
  "html"
)

const debug = false

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

  bodyStr := string(body)
  // Process payload
  payload := findPayload(bodyStr)
  for i := 0; i < len(payload); i++ {
    messages <- payload[i]
  }

  urls := parseString(bodyStr)
  if (depth < maxDepth) {
    for i := 0; i < len(urls); i++ {
      messages <- "+open"
      go getUrlAsString(urls[i], depth+1, maxDepth, messages)
    }
  }
  messages <- "-close"
  return

}

func findPayload(input string) []string {
  parseStringRegexp := regexp.MustCompile("<img\\s([^<>]+\\s)?src=[\"']([^\"'<>\\s]+)")
  tmp := parseStringRegexp.FindAllStringSubmatch(input, -1)
  matches := []string{}
  if (tmp != nil) {
    for i :=0; i < len(tmp); i++ {
      img := tmp[i][2]
      img = html.UnescapeString(img)
      matches = append(matches, img)
    }
  } else {
    if debug { fmt.Println("no payload matches") }
  }


  return matches
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
    } else {
      fmt.Println(msg)
    }
  }

}
