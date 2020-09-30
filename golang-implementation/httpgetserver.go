package main
import (
  "fmt"
  "io/ioutil"
  "log"
  "strings"
  "net"
  "os"
  "bufio"
  "sync"
  "strconv"
)

type AccessCounter struct {
  counts map[string]int
  mux sync.Mutex
}

var accessCounter = AccessCounter{counts: make(map[string]int)}
var printMutex sync.Mutex

func main() {
  data, err := ioutil.ReadFile("/etc/hostname")
  if err != nil {
    log.Fatal("Unable to identify the hostname, please grant required permissions")
  }
  hostname := strings.TrimSpace(string(data)) + ".cs.binghamton.edu"
  portNum := "8090"
  fmt.Printf("Hostname = %s, Port = %s\nStarting server...\n", hostname, portNum)
  if _, err := os.Stat("www"); err != nil {
        if os.IsNotExist(err) {
	  log.Fatal("www directory not present. Cannot instantiate the http server(as per the requirements)")
	}
  }
  ln, err := net.Listen("tcp", ":"+portNum)
  if err != nil {
    log.Print("Unable to create a listening server. Please check if anything is wrong!")
    log.Fatal(err)
  }
  fmt.Println("Server created successfully. Now accepting connections via TCP...")
  for {
    conn, err := ln.Accept()
    if err != nil {
	    log.Print("Error while accepting the request. The details of the of the client are as follows: ", conn.RemoteAddr().String())
    }
    go handleConnection(conn)
  }
}

func handleConnection(conn net.Conn) {
  clientAddr := string(conn.RemoteAddr().String())
  reader := bufio.NewReader(conn)
  var data string
  var err error
  for i := 0; i < 2; i++{ //Because the second word points to the resource requested
    data, err = reader.ReadString(' ');
    data = strings.TrimSpace(data)
    if err != nil {
      log.Print("Error in HTTP header, no resource specified")
      conn.Close()
      return
    }
  }
  filename := "www"+data
  count := accessCounter.updateCount(filename)
  log.Print("Access count: ", count)
  response := ""
  log.Print("Searching for Requested resource: ", data)
  if _, err := os.Stat(filename); err != nil {
        if os.IsNotExist(err) {
	  response = "404 Not Found"
	} else {
	  response = "200 OK"
	}
  }
  log.Print(response)
  // TODO: build HTTP repsonse 
  // logging the required output in the required format
  printLog(clientAddr, data, count)
}

func printLog(clientAddr string, resource string, count int) {
  printMutex.Lock()
  defer printMutex.Unlock()
  clientAddrParts := strings.Split(clientAddr, ":")
  log.Print(resource + "|" + clientAddrParts[0] + "|" + clientAddrParts[1] + "|" + strconv.Itoa(count))
}

func (accessCounter *AccessCounter) updateCount(filename string) int {
  accessCounter.mux.Lock()
  defer accessCounter.mux.Unlock()
  value, ok := accessCounter.counts[filename]
  if(ok) {
    accessCounter.counts[filename] = value + 1
  }else {
    accessCounter.counts[filename] = 1
  }
  return accessCounter.counts[filename]
}


