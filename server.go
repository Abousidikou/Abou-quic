package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"

	//"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	//"github.com/lucas-clemente/quic-go/logging"
	//"github.com/lucas-clemente/quic-go/qlog"
	"github.com/xyproto/sheepcounter"
)

// Size is needed by the /demo/upload handler to determine the size of the uploaded file
type Size interface {
	Size() int64
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
	port := flag.Int("p", 4450 , "Port to listen on")
//	enableQlog := flag.Bool("q", false, "Enable Qlog")
	flag.Parse()

	mux := http.NewServeMux()
	//mux := router.NewRouter().StrictSlash(true)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome\n")
	})

	//serving data directory
	direc, _ := os.Getwd()
	fs := http.FileServer(http.Dir(direc + "/data"))
	mux.Handle("/data/", http.StripPrefix("/data", fs))

	//serving mysite directory
	mysite := http.FileServer(http.Dir(*dir + "/mysite"))
	mux.Handle("/mysite/", http.StripPrefix("/mysite", mysite))

	var downDatalengh int64
	var downSpeed int64
	var durationDown int64
	mux.HandleFunc("/download", func(rw http.ResponseWriter, req *http.Request) {
		fmt.Println("Download Subtest")

		// Guarantee results are written even if function panics.
		defer func() {
			fmt.Println("End")
		}()

		// Run measurement.

		var msgSize = 1 << 27
		msg := generatePRData(int(msgSize))

		//var total int64

		sc := sheepcounter.New(rw)
		var StartTime = time.Now().UTC()
		sc.Write(msg)
		durationDown = time.Since(StartTime).Milliseconds()
		downDatalengh = sc.Counter()
		fmt.Println("COUNTED:", sc.Counter()) // Counts the bytes sent, for this response only

	})
	mux.HandleFunc("/getDownSpeed", func(rw http.ResponseWriter, req *http.Request) {
		downSpeed = int64(((int64(downDatalengh*8) / int64(durationDown)) * 1000) / 1000000)
		fmt.Fprintf(rw, strconv.FormatInt(downSpeed, 10))
		downDatalengh = 0
	})

	var upDatalengh int
	var upSpeed int64
	mux.HandleFunc("/upload", func(rw http.ResponseWriter, req *http.Request) {
		body := &bytes.Buffer{}
		_, err := io.Copy(body, req.Body)
		if err != nil {
			log.Fatal("request", err)
		}
		upDatalengh += body.Len()
	})

	mux.HandleFunc("/getUpSpeed", func(rw http.ResponseWriter, req *http.Request) {
		timeString := req.FormValue("id")
		fmt.Println(reflect.TypeOf(timeString))
		t, _ := strconv.ParseFloat(timeString, 64)
		fmt.Println("t: ", t)
		upSpeed = (((int64(upDatalengh) * 8) / int64(t)) * 1000) / 1000000
		fmt.Println("upSpeed: ", upSpeed)
		fmt.Fprintf(rw, strconv.FormatInt(int64(upDatalengh), 10))
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

	// start server
	 ad := ":"+ strconv.Itoa(*port)
	log.Fatal(http3.ListenAndServe(ad , *cert, *key, mux))

	/*quicConf := &quic.Config{
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
	ad := ":"+ strconv.Itoa(*port)
	server := http3.Server{
		Server:     &http.Server{Handler: mux, Addr: ad},
		QuicConfig: quicConf,
	}
	log.Fatal(server.ListenAndServeTLS(*cert, *key))*/
}
