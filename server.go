package main

import (
	//	"bytes"
	//	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
//	"strconv"
	"time"

	"github.com/gorilla/websocket"
	//"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func helloJson(w http.ResponseWriter, r *http.Request) {

	/*// Origin check
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	//upgrding connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connection successfully Upgraded ")

	fmt.Fprintln(w, "your are connected via ws: "+ws.Subprotocol())
	*/
	if r.URL.Path == "/nf" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(`{"hello": "Worl"}`))
}

func alt(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/nf" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Alt-Svc", "h3=\":4448\"")
	w.Write([]byte("Alt-Svc"))
}

// Generate data Byte from interger(lengh)
func generatePRData(l int) []byte {
	res := make([]byte, l)
	seed := uint64(1)
	for i := 0; i < l; i++ {
		seed = seed * 48271 % 2147483647
		res[i] = byte(seed)
	}
	return res
}

func main() {
	fmt.Println("Hello A server @ quic:4448")

	cert := flag.String("c", "fullchain.pem", "Enter the certificate file")
	key := flag.String("k", "privkey.pem", "Enter the key file")
	dir := flag.String("d", "", "Directory to serve")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome\n")
	})
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "You are on Api\n")
	})
	mux.HandleFunc("/json", helloJson)

	mux.HandleFunc("/alt", alt)

	mux.HandleFunc("/download", func(rw http.ResponseWriter, req *http.Request) {
		// Create ultimate result.

		var StartTime = time.Now().UTC()
		fmt.Println("StartTime: ", StartTime.String())
		start := StartTime

		// Guarantee results are written even if function panics.
		defer func() {
			var EndTime = time.Now().UTC()
			fmt.Println("EndTime: ", EndTime.String())
			//h.writeResult(data.UUID, kind, result)
		}()

		// Run measurement.

		fmt.Println("Download Subtest")
		var msgSize = 1 << 13
		const every time.Duration = 250 * time.Millisecond

		msg := generatePRData(int(msgSize))

		var total int64
		for {
			//fmt.Println("Time Since: ",time.Since(StartTime).String() )
			if time.Since(StartTime) >= 13*time.Second {
				//fmt.Println("fin" , time.Since(StartTime).String())
				//rw.WriteHeader(400)
				return
			}
			rw.Write(msg)

			// The following block of code implements the scaling of message size
			// as recommended in the spec's appendix. We're not accounting for the
			// size of JSON messages because that is small compared to the bulk
			// message size. The net effect is slightly slowing down the scaling,
			// but this is currently fine. We need to gather data from large
			// scale deployments of this algorithm anyway, so there's no point
			// in engaging in fine grained calibration before knowing.
			total += int64(msgSize)
			if time.Now().UTC().Sub(start) >= every {
				//fmt.Println("Total: ", strconv.FormatInt(total, 10), " bytes")
				//fmt.Println("Time: ", strconv.FormatInt(int64(time.Since(StartTime).Milliseconds()), 10))
				bitSentTillNow := total * 8 // bytes * 8
				//fmt.Println(bitSentTillNow)
				// time.Since(StartTime)/time.Millisecond)  time since start in millisecond
				//(bitSentTillNow/int64(time.Since(StartTime).Milliseconds()))*1000 nombre de bit en 1 s( nmbre de bit en ms * 1000)
				fmt.Println("Speed : ", ((bitSentTillNow/int64(time.Since(StartTime).Milliseconds()))*1000)/1000000, " Mbits/s")
				start = time.Now().UTC()
			}
			if int64(msgSize) >= 1<<10 {
				continue // No further scaling is required
			}
			if int64(msgSize) > total/16 {
				continue // message size still too big compared to sent data
			}

			if int64(msgSize) > 1<<24 {
				continue // message size still too big compared
			}
			msgSize *= 2
			msg = generatePRData(int(msgSize))
		}

	})

	mux.HandleFunc("/upload", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Println("Upload Subtest")

		//fmt.Println("StartTime: ", StartTime.String())
		//start := StartTime

		/*body := &bytes.Buffer{}
		_, err := io.Copy(body, req.Body)
		if err != nil {
			log.Fatal(err)
		}*/
		var buf []byte
		var err error

		/*go func(buf *[]byte) {
		/*ticker := time.NewTicker(1 * time.Second)
		for range ticker.C {
			//fmt.Println(len(buf))
			fmt.Println(binary.Size(buf))
		}*/
		/*timer := time.NewTimer(time.Second * 13)
			<-timer.C
			fmt.Println("Buf Size: ", binary.Size(buf))
		}(&buf)*/
		var StartTime = time.Now().UTC()
		buf, err = ioutil.ReadAll(req.Body)
		var EndTime = time.Since(StartTime).Milliseconds()
		if err != nil {
			log.Fatal("request", err)
		}
		//Ticker

		//fmt.Println(int64(EndTime)) // do whatever you want with the binary file buf
		//fmt.Println(strconv.FormatInt(int64(len(buf)), 10))
		fmt.Println("Speed: ", ((int64(len(buf)*8)/int64(EndTime))*1000)/1000000, " Mbits/s")
		fmt.Fprintf(rw, "Success")
		//fmt.Println(body)
		//fmt.Println(time.Since(StartTime).Milliseconds())
		/*var total int64
		logger.Infof("Request Body: %d bytes", body.Len())
		logger.Infof("Request Body:")
			logger.Infof("%s", body.Bytes())*/
	})

	mux.HandleFunc("/time", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Println(req.Body)
	})
	mux.Handle("/mysite/", http.StripPrefix("/mysite", http.FileServer(http.Dir(*dir+"/mysite"))))

	/*quicConf := &quic.Config{}
	server := http3.server{
		Server: &http.Server{Handler:mux, Addr: ":4448" },
		QuicConfig: quicConf,
	}*/
	log.Fatal(http3.ListenAndServe(":4448", *cert, *key, mux))
	//log.Fatal(http3.ListenAndServe(":4448", *cert, *key, mux))

	//log.Fatal(http.ListenAndServe(":4449", mux))

	/*done := make(chan int,1)
		go func(channel chan <- int){
		fmt.Println("in goroutine")
		quicConf := &quic.Config{}
		server := http3.Server{
			Server:     &http.Server{Handler: mux, Addr: ":4448" },
			QuicConfig: quicConf,
		}
		log.Fatal(server.ListenAndServeTLS(*cert,*key))

		fmt.Println("on channel")
	 	channel <- 5
		}(done)
		fmt.Println("Waiting for go routine end")
		s := <- done
		fmt.Println(s)*/

	/*	quicConf := &quic.Config{}
		server := http3.Server{
			Server:     &http.Server{Handler: mux, Addr: ":4448" },
			QuicConfig: quicConf,
		}
		log.Fatal(server.ListenAndServeTLS(*cert, *key))*/

}

