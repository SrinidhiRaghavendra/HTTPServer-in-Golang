Summary:
This project is written using golang. 
This project implements basic  http server which serves only HTTP get requests for files.
The server must have a directory called "www" where all the requested resources are searched. If this www directory is not present, then the server quits.
If the requested resource is not present in www directory, 404 error code is returned as status in the HTTP header as response.
If the requested resource is present, the following headers are added in addition to the file contents.(NOTE: there is one empty line between headers and the body of the http response.):
1. Status line
2. Date
3. Server
4. Last-Modified
5. Content-Type
6. Content-Length

For each request, the server logs the following after servicing the request and closing the socket connection:
<requested resource name>|<client ip address>|<client port number>|<access count of the requested resource>

Each new request is handled in a separate thread and hence the main thread listening for the connections isn't blocked.

Implementation Description:
The implementation details are as follows:
1. Finding the hostname using /etc/hostname
2. Generating the FQDN for the host by appending ".cs.binghamton.edu"
3. Port name is selected as 8090.
4. Host name an port information is logged onto stdout.
5. Check if www directory is present, if not, quit, else, go on to start a TCP listener(server).
6. With each request, create a new thread to handle the connection
7. In the induvidual connection thread, synchronously increment the access count of the requeste resource in a map which maps from a string(resource) to an int(access count).
8. Now, if the resource is present, create all the above mentioned headers and write the bytes of the header as well as the file contents into the connection, else write the bytes of just the header(with status as 404 Not Found) into the connection and close the connection. 
9. Log the required output(format mentioned above for each request) onto stdout.

Instructions:
From the base directory, move to golang-implementation directory(This has the source file as well as the www diretory).
>cd golang-implementation

Compile:
>make
The above creates an executable called server in the same directory

Run:
>./server
Starts the http server if www directory is present, else quits with  an error message.

Sample Input/Output:
If the host where the server runs is remote02, then the FQDN for the server will be remote02.cs.binghamton.edu ands the port is 8090.
For a file named objects.zip in the www directory, the following is a sample request and sample outputs on server and the client side.

Sample request using wget:
wget http://remote02.cs.binghamton.edu:8090/emptyFile.txt

Sample output on the server:
/emptyFile.txt|128.226.114.203|54326|1

Sample output on the client:
Resolving remote02.cs.binghamton.edu (remote02.cs.binghamton.edu)... 128.226.114.202
Connecting to remote02.cs.binghamton.edu (remote02.cs.binghamton.edu)|128.226.114.202|:8090... connected.
HTTP request sent, awaiting response... 200 OK
Length: 0 [text/plain]
Saving to: ‘emptyFile.txt’

emptyFile.txt                                    [ <=>                                                                                         ]       0  --.-KB/s    in 0s      

2020-10-01 16:02:59 (0.00 B/s) - ‘emptyFile.txt’ saved [0/0] 
