package main
import (
  "fmt"
  "io/ioutil"
  "log"
  "strings"
  "net"
)

func main() {
  data, err := ioutil.ReadFile("/etc/hostname")
  if err != nil {
    log.Fatal("Unable to identify the hostname, please grant required permissions")
  }
  hostname := strings.TrimSpace(string(data)) + ".cs.binghamton.edu"
  portNum := "8090"
  fmt.Printf("Hostname = %s, Port = %s\nStarting server...\n", hostname, portNum)
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
  fmt.Printf("The address of the client is as follows: ", string(conn.RemoteAddr().String()))
}
