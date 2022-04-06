package main

// #cgo CFLAGS: -Idriver
// #cgo LDFLAGS: -Ldriver/build -ldriver
// #include "driver.h"
import "C"

import (
    "flag"
    "fmt"
    "io/ioutil"
    "html/template"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "strings"
    "syscall"

    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
)

var homeTemplate = template.Must(template.ParseFiles("./static/touchpad.html"))
var addr = flag.String("addr", "0.0.0.0:8080", "http service address")
var upgrader = websocket.Upgrader{}

func echo (w http.ResponseWriter, r *http.Request) {
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade: ", err)
        return
    }
    defer c.Close()

    for {
        mt, message, err := c.ReadMessage()
        if err != nil {
            log.Println("read: ", err)
            break
        }

        //log.Printf("recv %s", message)
        processCommand(message)
        err = c.WriteMessage(mt, message)
        if err != nil {
            log.Println("write: ", err)
            break
        }
    }
}

func home(w http.ResponseWriter, r *http.Request) {
    homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func asset(w http.ResponseWriter, r *http.Request) {
    fn := "." + r.URL.Path
    buf, err := ioutil.ReadFile(fn)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintln(w, "404 not found")
        return
    }

    w.Write(buf)
}

func processCommand(msg []byte) {
    if len(msg) < 2 {
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
            continue
        }

        dx, err := strconv.Atoi(numbers[1])
        if err != nil {
            panic(err)
        }
        dy, err := strconv.Atoi(numbers[2])
        if err != nil {
            panic(err)
        }

        C.driver_mouse_rel(C.int(dx), C.int(dy))
    }
}

func setupHandlers() {
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        fmt.Println("shutting down...")
        C.driver_destroy_device()
        os.Exit(0)
    }()
}

func main() {
    flag.Parse()
    log.SetFlags(0)

    setupHandlers()

    errno := C.driver_create_device()
    if errno != 0 {
        log.Fatal("cannot create device! err=", errno)
    }

    router := mux.NewRouter()
    router.HandleFunc("/static/{[a-z]+}.js", asset).Methods("GET")
    router.HandleFunc("/static/{[a-z]+}.css", asset).Methods("GET")
    router.HandleFunc("/echo", echo)
    router.HandleFunc("/", home)

    http.Handle("/", router)

    log.Fatal(http.ListenAndServe(*addr, nil))
}

