package main

import (
  "os"
  "fmt"
  "net/http"
  "io/ioutil"
  "regexp"
  "html"
  "strings"
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
  findPayload(bodyStr, url, messages)

  if (depth < maxDepth) {
    followLinks(bodyStr, depth+1, maxDepth, messages)
  }
  messages <- "-close"
  return

}

func findPayload(input string, baseUrl string, messages chan<- string) {
  parseStringRegexp := regexp.MustCompile("(?i)<img\\s([^<>]+\\s)?src=[\"']([^\"'<>\\s]+)")
  protocolRegexp := regexp.MustCompile("(?i)^([a-zA-Z]+:)(//[^/]*/?)(.*)([^/]*)$")

  protoMatch := protocolRegexp.FindStringSubmatch(baseUrl);

  tmp := parseStringRegexp.FindAllStringSubmatch(input, -1)
  if (tmp != nil) {
    for i :=0; i < len(tmp); i++ {
      img := tmp[i][2]
      img = html.UnescapeString(img)

      // Really crude
      if strings.Index(img, "://") == -1 {
        // No protocol => relative url
        if strings.Index(img, "//") == 0 {
          // absolute url sans protocol
          img = protoMatch[1] + img;
        } else if (img[0] == '/') {
          // Absolute wrt baseUrl
          img = protoMatch[1] + protoMatch[2] + img
        } else {
          // Just append away
          img = protoMatch[1] + protoMatch[2] + protoMatch[3] + img
        }
      }
      messages <- img
    }
  } else {
    if debug { fmt.Println("no payload matches") }
  }
}

func followLinks(input string, depth int, maxDepth int, messages chan<- string) {
  parseStringRegexp := regexp.MustCompile("https?://[^\"\\s]+")
  matches := parseStringRegexp.FindAllString(input, -1)
  if (matches != nil) {
    for i := 0; i < len(matches); i++ {
      messages <- "+open"
      go getUrlAsString(matches[i], depth, maxDepth, messages)
    }
  }

}

func main() {
  fmt.Println(os.Args)
  if len(os.Args) < 2 {
    fmt.Println("usage: scrape url1 [url2 url3 ..]")
    return
  }

  messages := make(chan string, len(os.Args)-1)
  counter := 0

  for i := 1; i < len(os.Args); i++ {
    messages <- "+open"
    go getUrlAsString(os.Args[i], 0, 1, messages)
  }
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
