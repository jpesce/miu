Go is statically typed, but with pointers and metaprogramming via the `reflect` package, it's
possible to mutate variables without knowing their exact structure _a priori_. 

For example, the standard `encoding/json` package, has no idea about the format of the JSON content
it will be processing, but if the user passes a compatible `struct` and JSON content, it will
unmarshal/deserialize/unpickle correctly.

Here's an example of a function that can mutate both integers and strings. In this example, it would
be easy to create one dedicated function for each type, but this can get unfeasible for complex
structures that need to be parsed dynamically (which is the case of JSON strings).

Disclaimer: in production, you should probably do a lot more checks.
```
package main

import (
  "fmt"
  "reflect"
)

func mutate(i any) {
  element := reflect.ValueOf(i).Elem()
  if element.Type().AssignableTo(reflect.TypeOf("")) {
    element.SetString("New string")
  }
  if element.Type().AssignableTo(reflect.TypeOf(0)) {
    element.SetInt(99)
  }
}

func main () {
  i := 0
  fmt.Println(i) // 0
  mutate(&i)
  fmt.Println(i) // 99

  s := "Old string"
  fmt.Println(s) // Old string
  mutate(&s)
  fmt.Println(s) // New string
}
```
