package main

import (
	//	"bytes"
	//	"encoding/binary"
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	//"crypto/md5"
	"errors"
	"mime/multipart"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	//"github.com/lucas-clemente/quic-go/internal/utils"
	"github.com/lucas-clemente/quic-go/logging"
	"github.com/lucas-clemente/quic-go/qlog"
)

// Size is needed by the /demo/upload handler to determine the size of the uploaded file
type Size interface {
	Size() int64
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func helloJson(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/nf" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write([]byte(`{"hello": "World"}`))
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

type bufferedWriteCloser struct {
	*bufio.Writer
	io.Closer
}

// NewBufferedWriteCloser creates an io.WriteCloser from a bufio.Writer and an io.Closer
func NewBufferedWriteCloser(writer *bufio.Writer, closer io.Closer) io.WriteCloser {
	return &bufferedWriteCloser{
		Writer: writer,
		Closer: closer,
	}
}

func (h bufferedWriteCloser) Close() error {
	if err := h.Writer.Flush(); err != nil {
		return err
	}
	return h.Closer.Close()
}

func main() {
	fmt.Println("Serving @ quic:4448")

	cert := flag.String("c", "fullchain.pem", "Enter the certificate file")
	key := flag.String("k", "privkey.pem", "Enter the key file")
	dir := flag.String("d", ".", "Directory to serve")
	enableQlog := flag.Bool("q", false, "Enable Qlog")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome\n")
	})

	mux.HandleFunc("/json", helloJson)

	mux.HandleFunc("/download", func(rw http.ResponseWriter, req *http.Request) {
		// Create ultimate result.
		testFile, logerr := os.OpenFile("test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		if logerr != nil {
			return
		}
		defer testFile.Close()
		var StartTime = time.Now().UTC()
		fmt.Println("StartTime: ", StartTime.String())
		testFile.WriteString("Time: " + StartTime.String() + "\n")
		start := StartTime

		// Guarantee results are written even if function panics.
		defer func() {
			var EndTime = time.Now().UTC()
			fmt.Println("EndTime: ", EndTime.String())
			testFile.WriteString("EndTime: " + EndTime.String() + "\n")
			//h.writeResult(data.UUID, kind, result)
		}()

		// Run measurement.

		fmt.Println("Download Subtest")
		testFile.WriteString("Download Subtest...\n")
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
				sped := ((bitSentTillNow / int64(time.Since(StartTime).Milliseconds())) * 1000) / 1000000
				fmt.Println("Speed : ", sped, " Mbits/s")
				testFile.WriteString("Speed: " + strconv.FormatInt(sped, 10) + "\n")
				start = time.Now().UTC()
			}
			fmt.Println("message size:", msgSize)
			if int64(msgSize) >= total/16 {
				continue // message size still too big compared to sent data
			}

			if int64(msgSize) >= 1<<24 {
				continue // message size still too big compared
			}
			msgSize *= 2
			msg = generatePRData(int(msgSize))
		}

	})

	mux.HandleFunc("/upload", func(rw http.ResponseWriter, req *http.Request) {
		testFile, logerr := os.OpenFile("test.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
		if logerr != nil {
			return
		}
		defer testFile.Close()
		fmt.Println("Upload Subtest")
		testFile.WriteString("Upload Subtest...\n")
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
		sped := ((int64(len(buf)*8) / int64(EndTime)) * 1000) / 1000000
		fmt.Println("Speed: ", sped, " Mbits/s")
		fmt.Fprintf(rw, "Success")
		testFile.WriteString("Speed: " + sped)
		//fmt.Println(body)
		//fmt.Println(time.Since(StartTime).Milliseconds())
		/*var total int64
		logger.Infof("Request Body: %d bytes", body.Len())
		logger.Infof("Request Body:")
			logger.Infof("%s", body.Bytes())*/
	})

	mux.HandleFunc("/demo/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			start := time.Now()
			err := r.ParseMultipartForm(1 << 30) // 4 GB
			if err == nil {
				var file multipart.File
				var hand *multipart.FileHeader
				file, hand, err = r.FormFile("uploadfile")
				if err == nil {
					//defer file.Close()
					//var size int64
					if sizeInterface, ok := file.(Size); ok {
						_ = sizeInterface.Size() // _ == size
						//b := make([]byte, size)
						//start := time.Now()
						//i, _ := file.Read(b)
						//fmt.Println(time.Since(start))
						//fmt.Println("Size of file: ", i, " bytes")

						logFile, logerr := os.OpenFile("log.file", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
						f, err := os.OpenFile("./data/"+hand.Filename, os.O_WRONLY|os.O_CREATE, 0666)
						if err != nil || logerr != nil {
							fmt.Println(err)
							return
						}
						defer logFile.Close()
						defer f.Close()
						// Calling Copy method with its parameters
						bytes, erro := io.Copy(f, file)
						// If error is not nil then panics
						if erro != nil {
							panic(erro)
						}

						// Prints output
						fmt.Printf("The number of bytes are: %d\n", bytes)
						fmt.Println(time.Since(start))
						logFile.WriteString("Nom du fichier: " + hand.Filename + "\n")
						logFile.WriteString("Taille du fichier: " + strconv.FormatInt(int64(bytes), 10) + " bytes\n")
						logFile.WriteString("Temps d'envoie: " + time.Since(start).String() + "\n\n")

						//file.Read(b)
						//md5 := md5.Sum(b)
						//fmt.Fprintf(w, "File Received---md5:%x---Header:%v", md5, hand.Header)
						fmt.Fprintf(w, "File Received---%v", hand.Header)
						return
					}
					err = errors.New("couldn't get uploaded file size")
				}
			}
			fmt.Printf("Error receiving upload: %#v", err)
		}
		io.WriteString(w, `<html><body><form action="/demo/upload" method="post" enctype="multipart/form-data">
				<input type="file" name="uploadfile"><br>
				<input type="submit">
			</form></body></html>`)
	})

	// create file server handler
	direc, _ := os.Getwd()
	//fmt.Println("current path :" + direc)
	fs := http.FileServer(http.Dir(direc + "/data"))
	mux.Handle("/data/", http.StripPrefix("/data", fs))

	mysite := http.FileServer(http.Dir(*dir + "/mysite"))
	mux.Handle("/mysite/", http.StripPrefix("/mysite", mysite))

	/*
		This function start both tcp and quic
	*/
	//log.Fatal(http3.ListenAndServe(":4448", *cert, *key, mux))

	/*
		This part start only quic but with quicConfig which enable qlog records.
		It doesn't works if clients haven't already being connected on tcp but it works
		after being connected with the function above.
	*/

	quicConf := &quic.Config{
		MaxIncomingStreams:         1000000,
		MaxIncomingUniStreams:      1000000,
		InitialStreamReceiveWindow: 524288,
		// MaxStreamReceiveWindow is the maximum stream-level flow control window for receiving data.
		// If this value is zero, it will default to 6 MB.
		//MaxStreamReceiveWindow: 6,
		// InitialConnectionReceiveWindow is the initial size of the stream-level flow control window for receiving data.
		// If the application is consuming data quickly enough, the flow control auto-tuning algorithm
		// will increase the window up to MaxConnectionReceiveWindow.
		// If this value is zero, it will default to 512 KB.
		//InitialConnectionReceiveWindow uint64,
		// MaxConnectionReceiveWindow is the connection-level flow control window for receiving data.
		// If this value is zero, it will default to 15 MB.
		//MaxConnectionReceiveWindow uint64,
		// MaxIncomingStreams is the maximum number of concurrent bidirectional streams that a peer is allowed to open.
		// Values above 2^60 are invalid.
		// If not set, it will default to 100.
		// If set to a negative value, it doesn't allow any bidirectional streams.
		//MaxIncomingStreams int64,
		// MaxIncomingUniStreams is the maximum number of concurrent unidirectional streams that a peer is allowed to open.
		// Values above 2^60 are invalid.
		// If not set, it will default to 100.
		// If set to a negative value, it doesn't allow any unidirectional streams.
		//MaxIncomingUniStreams int64,
	}
	if *enableQlog {
		quicConf.Tracer = qlog.NewTracer(func(_ logging.Perspective, connID []byte) io.WriteCloser {
			fmt.Println(connID)
			filename := fmt.Sprintf("server_%s.qlog", time.Now().String())
			f, err := os.Create(filename)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("Creating qlog file %s.\n", filename)
			return NewBufferedWriteCloser(bufio.NewWriter(f), f)
		})
	}
	server := http3.Server{
		Server:     &http.Server{Handler: mux, Addr: ":4448"},
		QuicConfig: quicConf,
	}
	log.Fatal(server.ListenAndServeTLS(*cert, *key))

}
