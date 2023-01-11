package main

// #cgo CFLAGS: -I../driver
// #cgo LDFLAGS: -L../driver/artifacts -ldriver
// #include "driver.h"
import "C"

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"touchpad/security"
	"touchpad/server"
	"touchpad/util"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

/** Token Parameters **/
var TEMPLATE_HTML = "index.html"
var TOKEN_COOKIE_NAME = "token"
var TOKEN_TIME_VALID = 5 * time.Minute

/** Flags **/
var port = flag.String("port", "8080", "http service port")
var addr = flag.String("addr", "0.0.0.0", "http service address")
var siteDir = flag.String("site", "./static/dist/", "static site assets")
var certFile = flag.String("cert", "", "tls cert file")
var keyFile = flag.String("key", "", "tlk key file")
var HOST = *addr + ":" + *port

var upgrader = websocket.Upgrader{}

var isAlive = false

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade error: ", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read error: ", err)
			break
		}

		processCommand(message)

		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write error: ", err)
			break
		}
	}
}

func processCommand(msg []byte) {
	if len(msg) < 3 {
		log.Println("bad message format")
		return
	}

	str := string(msg)
	fingers := strings.Split(str, ")")
	for _, finger := range fingers {
		if len(finger) < 1 {
			continue
		}

		tmp := finger[1:] // trim open paren

		numbers := strings.Split(tmp, ",")
		if len(numbers) != 3 {
			log.Printf("error parsing: %v\n", numbers)
			continue
		}

		fingerNum, err := strconv.ParseUint(numbers[0], 10, 8) // base 10, 8 bits
		if err != nil {
			log.Printf("error parsing: %v\n", numbers)
			panic(err)
		}

		dx, err := strconv.ParseFloat(numbers[1], 32)
		if err != nil {
			panic(err)
		}
		dy, err := strconv.ParseFloat(numbers[2], 32)
		if err != nil {
			panic(err)
		}

		if fingerNum == 0 {
			err := C.driver_mouse_rel(C.int(dx), C.int(dy))
			if err != 0 {
				log.Printf("driver_mouse_rel error: %v\n", err)
			}
		}
	}
}

func setupHandlers() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\nshutting down...")
		C.driver_destroy_device()
		os.Exit(0)
	}()
}

func url(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(util.GetURL(*port)))
}

func main() {
	flag.Parse()
	log.SetFlags(0)

	log.Println("generating keys...")
	security.GenerateKeys()
	security.GenerateSessionID()

	log.Println("creating signal handlers...")
	setupHandlers()

	log.Println("creating virtual device...")
	errno := C.driver_create_device()
	if errno != 0 {
		log.Fatal("cannot create device! err=", errno)
	}

	log.Println("building routes...")
	router := mux.NewRouter()
	router.HandleFunc("/api/url", url).Methods("GET")
	router.HandleFunc("/api/echo", echo)
	router.HandleFunc("/api/auth/alive", server.AuthAliveHandler).Methods("GET")
	router.HandleFunc("/api/auth/challenge", server.AuthLoginChallengeHandler).Methods("GET")
	router.HandleFunc("/api/auth/response", server.AuthLoginResponseHandler).Methods("POST")
	router.PathPrefix("/").Handler(server.NewFileServer(*siteDir))

	router.Use(server.NewLoggerMiddleware)
	router.Use(server.NewAuthMiddleware)

	http.Handle("/", router)

	if len(*certFile) > 0 && len(*keyFile) > 0 {
		log.Printf("serving at https://%s\n", util.GetURL(*port))
		log.Fatal(http.ListenAndServeTLS(HOST, *certFile, *keyFile, nil))
	} else {
		log.Printf("serving at http://%s\n", util.GetURL(*port))
		log.Fatal(http.ListenAndServe(HOST, nil))
	}
}
