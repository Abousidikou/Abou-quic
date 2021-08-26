package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/lucas-clemente/quic-go/http3"
)

func main() {
	fmt.Println("Hello client")
	url := flag.String("u", "https://monitor.uac.bj:4448/", "Enter url")
	flag.Parse()
	hclient := http.Client{
		Transport: &http3.RoundTripper{},
	}

	res, err := hclient.Get(*url)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println("Got response", res)

	body := &bytes.Buffer{}
	_, err = io.Copy(body, res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Request Body: ")
	fmt.Printf("Lenght: %d bytes", body.Len())
	fmt.Printf("Total space allocated: %d bytes", body.Cap())
	fmt.Println("\n", body.String())

}
