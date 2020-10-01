package main
import (
  "fmt"
  "io"
  "io/ioutil"
  "log"
  "strings"
  "net"
  "os"
  "bufio"
  "sync"
  "strconv"
  "time"
  "mime"
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
  response := ""
  var fileLastModified time.Time
  var contentType string
  var contentLength int64
  var httpHeader string
  if stats, err := os.Stat(filename); err != nil {
    if os.IsNotExist(err) {
      response = "404 Not Found"
      httpHeader = "HTTP/1.1 " + response + "\n"
      bytesToTransfer := []byte(httpHeader)
      conn.Write(bytesToTransfer)
    }
  } else {
    response = "200 OK"
    fileLastModified = stats.ModTime()
    resourceParts := strings.Split(data, ".")
    if len(resourceParts) == 1 {
      contentType = "application/octet-stream"
    } else {
      extension := resourceParts[len(resourceParts)-1]
      contentType = mime.TypeByExtension("." + extension)
      if len(contentType) == 0 {
        contentType = "application/octet-stream"
      }
    }
    contentLength = stats.Size()
    httpHeader = buildHTTPHeader(response, fileLastModified, contentType, contentLength)
    bytesToTransfer := []byte(httpHeader)
    conn.Write(bytesToTransfer)
    file, err := os.Open(strings.TrimSpace(filename)) // For read access.
    if err != nil {
      log.Fatal(err)
    }
    defer file.Close() // make sure to close the file even if we panic.
    _, err = io.Copy(conn, file)
    if err != nil {
      log.Print(err)
    }
  }
  conn.Close()
  printLog(clientAddr, data, count)
}

func buildHTTPHeader(response string, fileLastModified time.Time, contentType string, contentLength int64) string {
  header := "HTTP/1.1 " + response + "\n"
  header = header + "Date: " + time.Now().UTC().Format(time.RFC1123) + "\n"
  header = header + "Server: CS557/assignment1/1" + "\n"
  header = header + "Last-Modified: " + fileLastModified.UTC().Format(time.RFC1123) + "\n"
  header = header + "Content-Type: " + contentType + "\n"
  header = header + "Content-Length: " + strconv.FormatInt(contentLength, 10) + "\n\n" //The second new line is to indicate the end of header
  return header
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


