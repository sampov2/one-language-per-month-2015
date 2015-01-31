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
const MESSAGE_TYPE_START = 0
const MESSAGE_TYPE_END = 1
const MESSAGE_TYPE_PAYLOAD = 2

type Message struct {
  msgType int
  payload string
}

func sendClose(messages chan<- Message) {
  messages <- Message{msgType:MESSAGE_TYPE_END}
}

func getUrlAsString(url string, depth int, maxDepth int, messages chan<- Message) {
  var resp *http.Response

  // Send close message once this function exits
  defer sendClose(messages)

  if debug { fmt.Println("Retrieving", url)}
  resp, err := http.Get(url)
  if err != nil {
    return
  }

  defer resp.Body.Close()
  if debug { fmt.Println("Reading response", url)}
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return
  }

  if debug { fmt.Println("Read", len(body), "bytes from", url) }

  bodyStr := string(body)
  // Process payload
  findPayload(bodyStr, url, messages)

  if (depth < maxDepth) {
    followLinks(bodyStr, depth+1, maxDepth, messages)
  }
}

func findPayload(input string, baseUrl string, messages chan<- Message) {
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
      messages <- Message{msgType:MESSAGE_TYPE_PAYLOAD, payload:img}
    }
  } else {
    if debug { fmt.Println("no payload matches") }
  }
}

func followLinks(input string, depth int, maxDepth int, messages chan<- Message) {
  parseStringRegexp := regexp.MustCompile("https?://[^\"\\s]+")
  matches := parseStringRegexp.FindAllString(input, -1)
  if (matches != nil) {
    for i := 0; i < len(matches); i++ {
      messages <- Message{msgType:MESSAGE_TYPE_START}
      go getUrlAsString(matches[i], depth, maxDepth, messages)
    }
  }

}

func main() {
  if len(os.Args) < 2 {
    fmt.Println("usage: scrape url1 [url2 url3 ..]")
    return
  }

  messages := make(chan Message, len(os.Args)-1)
  counter := 0
  running := true

  for i := 1; i < len(os.Args); i++ {
    messages <- Message{msgType:MESSAGE_TYPE_START}
    go getUrlAsString(os.Args[i], 0, 1, messages)
  }
  for running {
    msg := <-messages
    switch {
    case msg.msgType == MESSAGE_TYPE_START:
      if debug { fmt.Println("open!") }
      counter++;
      break;
    case msg.msgType == MESSAGE_TYPE_END:
      if debug { fmt.Println("close..") }
      counter--;
      if (counter == 0) {
        if debug { fmt.Println("done") }
        running = false
      }
      break
    case msg.msgType == MESSAGE_TYPE_PAYLOAD:
      fmt.Println(msg.payload)
      break
    }
  }
}
