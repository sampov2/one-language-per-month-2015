package main

// This file is used for fooling around with the syntax

import (
  //"net/http"
  "fmt"
  "time"
  "math/rand"
  "math/cmplx"
)

/**
 * Two assignment operators
 *    x := foo()    and    x = foo()
 *
 * Where := defines a new variable x and assigns the return value of foo()
 * to it. The operator = assigns the return value of foo() to the existing
 * variable x.
 */

const (
  my_number = 1234
  )

var (
  j string = "initialized"
  woo complex128 = cmplx.Sqrt(-5 + 12i)
)

func foo(x string) (ret string) {
  //ret = x + "Fish" + string(time.Now()) + "!"
  defer fmt.Println("two")
  fmt.Println("one")
  ret = x + " Fish "
  j = "foo ran"
  return
}

func main() {
  var x, y, z = "Foo", false, -10
  var jjj int
  var tmp string
  //tmp string
  //tmp := "Rumble"
  fmt.Println(j)
  for jjj = 0; jjj < 2; jjj++ {
    tmp = foo("Rumble")
  }
  fmt.Println(j)
  fmt.Println(tmp)
  rand.Seed(time.Now().Unix())
  fmt.Println("Number", rand.Intn(100))
  fmt.Println(x, y, z)
  //resp, err := http.Get("http://www.spatineo.com/")
  fmt.Println(woo)
  for jjj = 0; jjj < 10; jjj++ {
    fmt.Println(jjj)
  }
  switch {
    case jjj < 5:
      fmt.Println("Never??")
    case jjj < 15:
      fmt.Println("Yeah!")
    case jjj < 25:
      fmt.Println("Interesting!")
  }
}
