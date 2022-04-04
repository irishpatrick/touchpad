package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "html/template"
    "log"
    "net/http"

    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
)

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

        log.Printf("recv %s", message)
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

func main() {
    flag.Parse()
    log.SetFlags(0)

    router := mux.NewRouter()
    router.HandleFunc("/static/{[a-z]+}.js", asset).Methods("GET")
    router.HandleFunc("/static/{[a-z]+}.css", asset).Methods("GET")
    router.HandleFunc("/echo", echo)
    router.HandleFunc("/", home)

    http.Handle("/", router)

    log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.ParseFiles("./static/touchpad.html"))

