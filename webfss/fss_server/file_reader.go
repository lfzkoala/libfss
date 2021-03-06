package main

import (
  "os"
  "strconv"
  "bufio"
  "strings"
  "github.com/cathieyun/libfss/libfss"
  "fmt"
)

// Runs FSS on (matches) the 1st element of line (key), returns the 3rd element (val).
func readOneFetchSmall(server *libfss.Fss, serverNum byte, fssKey libfss.FssKeyEq2P, fileName string) string {
  var ans int = 0
  file, _ := os.Open(fileName)
  defer file.Close()
  scanner := bufio.NewScanner(file)

  // Read file line by line, on each line evaluate PF on node id
  for scanner.Scan() {
    line := strings.Split(scanner.Text(), " ")
    key, _ := strconv.Atoi(line[1])
    val, _ := strconv.Atoi(line[3])
    ans += server.EvaluatePF(serverNum, fssKey, uint(key))*val
  }
  return strconv.Itoa(ans)
}

// Runs FSS on (matches) the 0th element of line, returns the 1st element.
// Answer is base64-encoded int array, where each int represents an edge.
func readOneFetchLarge(server *libfss.Fss, serverNum byte, fssKey libfss.FssKeyEq2P, fileName string, fetchSize int) string {
  ans := make([]int, fetchSize)
  file, _ := os.Open(fileName)
  defer file.Close()
  scanner := bufio.NewScanner(file)

  maxBytes := 0

  // Read file line by line, on each line evaluate PF on key
  for scanner.Scan() {
    split := strings.SplitAfterN(scanner.Text(), " ", 2)
    fmt.Println("split: ", split)
    key, _ := strconv.Atoi(strings.TrimSpace(split[0]))

    byteArray := []byte(split[1])
    if len(byteArray) > maxBytes {
      maxBytes = len(byteArray)
    }
    fssVal := server.EvaluatePF(serverNum, fssKey, uint(key))
    for i := range byteArray {
      ans[i] += int(byteArray[i]) * fssVal
    }

  }
  // fmt.Println("answer before transmit: ", ans)
  fmt.Println("maxBytes: ", maxBytes)
  transmit := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ans)), ","), "[]")
  return transmit
}

// Runs FSS on (matches) the 0th element of line, returns the 1st element.
// Answer is base64-encoded int array, where each int represents an edge.
func readTwoFetchLarge(server *libfss.Fss, serverNum byte, fssKey libfss.FssKeyEq2P, fileName string, fetchSize int) string {
  ans := make([]int, fetchSize)
  file, _ := os.Open(fileName)
  defer file.Close()
  scanner := bufio.NewScanner(file)

  maxBytes := 0

  // Read file line by line, on each line evaluate PF on key
  for scanner.Scan() {
    line := strings.Split(scanner.Text(), " ")

    val1, _ := strconv.Atoi(line[0])
    val2, _ := strconv.Atoi(line[1])
    key := PRIME1^val1+PRIME2^val2

    byteArray := []byte(line[2])
    if len(byteArray) > maxBytes {
      maxBytes = len(byteArray)
    }
    fssVal := server.EvaluatePF(serverNum, fssKey, uint(key))
    for i := range byteArray {
      ans[i] += int(byteArray[i]) * fssVal
    }
  }
  // fmt.Println("answer before transmit: ", ans)
  fmt.Println("maxBytes: ", maxBytes)
  transmit := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(ans)), ","), "[]")
  return transmit
}