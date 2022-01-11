# *Description*
This is example of client and server Quic.
It runs on 4448 serving the directory "mysite"(Musical template for demo).
We run *download test*  and *upload test* for 13s.
## Download Test
The server generate bytes with minimum size of 8192 to 16MB which is sent on the wire.
Every 250ms the server display data sent. The client waits untils test finished before displaying speed.
Adapting message size.
if msg.Size() <= (1 << 24) || msg.Size() <= (total / 16) {
			  ==> double msg size
		}

## Upload Test
Here, We start with sending  to the server message sized 8192(1<<13).

# Start server quic

- Clone the repository
```bash
git clone https://github.com/Abousidikou/QuicTest_Mysite.git
```


```bash
go mod init quictest
```

```bash
go mod tidy
```
```bash
go run server -c "your certificate" -k "your key" -d "file to serve(./mysite by default)"  -q true (for Qlog creation, false by default)
```


# Start Client Quic
On the local or client machine:
```bash
go run client.go -u url(https://monitor.uac.bj:4450 default) 
```
we can access through browser.

# Url
- */*  Welcome
- */demo/upload* to upload file
- */data* to see uploaded files
- */mysite* for download and upload test.Two Buttons *Download Test* and *Upload Test* will help lauching test.

## Problem
The function http3.ListenAndServe() start server over TCP AND QUIC, but QUIC configuration is not possible.
The fonction server.ListenAndServeTLS() start only QUIC with QUIC Configuration possible.
Only navigator that knows our server as working on QUIC could contact it on QUIC.
THe command line client works for both functions.

## Answer
Connecting first with http3.ListenAndServe() before server.ListenAndServeTLS()

NB:**Answer will be display on the server side and the client side(console and graphically)**


