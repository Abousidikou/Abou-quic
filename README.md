## Description
This is example of client and server Quic.
It runs on 4448 serving the directory "mysite"(Musical template for demo).
We run *download test*  and *upload test* for 13s.
# *Download Test*
The server generate bytes with minimum size of 8192 to 16MB which is sent on the wire.
Every 250ms the server display data sent.
# *Upload Test*
 
## Start server quic

---*go mod init* 

---*go mod tidy*

---**go run server -c "your certificate" -k "your key"**


Now you can Access Enro web site template.

Two Buttons *Download Test* and *Upload Test* will help lauching test.

NB:**Answer will be display on the server side and the client side(console and graphically)**


