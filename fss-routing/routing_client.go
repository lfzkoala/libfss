package main

import (
  "net/http"
  "io/ioutil"
  "fmt"
  "strconv"
  "../libfss"
  "encoding/json"
  "encoding/base64"
  "time"
  "strings"
)

const (
  CONN_HOST = "localhost"
  CONN_START_PORT = 8000
)

func main() {
  // makeQuery(0, 191397, 20)
  makeQuery(1, 2, 20)
}

func makeQuery(queryType int, lookup uint, size uint) {
  t0 := time.Now()

  // Initialize client and generate keys based on query
  client := libfss.ClientInitialize(size)
  fssKeys := client.GenerateTreePF(lookup, 1) 

  chan0 := make(chan string)
  chan1 := make(chan string)
  go queryServer(chan0, strconv.Itoa(queryType), packageKeys(fssKeys[0]), packageKeys(client.PrfKeys), strconv.Itoa(int(client.NumBits)), 0)
  go queryServer(chan1, strconv.Itoa(queryType), packageKeys(fssKeys[1]), packageKeys(client.PrfKeys), strconv.Itoa(int(client.NumBits)), 1)
  ans0 := <-chan0
  ans1 := <-chan1

  t1 := time.Now()
  fmt.Printf("The call took %v to run.\n", t1.Sub(t0))
  if queryType == 0 {
    int0, _ := strconv.Atoi(ans0)
    int1, _ := strconv.Atoi(ans1)
    fmt.Println("combined answer: ", int0 + int1)
  } else if queryType == 1 {
    fmt.Println("\n\nans0: ", ans0, "\n")
    fmt.Println("\n\nans1: ", ans1, "\n")
    received0 := strings.Split(ans0,",")
    received1 := strings.Split(ans1,",")
    parsed := make([]byte, len(received0))
    for i := range received0 {
      num0, _ := strconv.Atoi(received0[i])
      num1, _ := strconv.Atoi(received1[i])
      parsed[i] = byte(num0 + num1)
    }
    fmt.Println("parsed string: \n",string(parsed))
  }
}

func queryServer(c chan string, queryType, fssKey, prfKeys, numBits string, serverNum int) {
  port := strconv.Itoa(CONN_START_PORT+serverNum)
  resp, _ := http.Get("http://"+CONN_HOST+":"+port+"/"+queryType+"/"+fssKey+"/"+prfKeys+"/"+numBits)
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  var answer map[string]string
  _ = json.Unmarshal(body, &answer)
  c <- answer["ans"]
}

func packageKeys(key interface{}) string {
  marshalledKey, _ := json.Marshal(key)
  return base64.StdEncoding.EncodeToString(marshalledKey)
}