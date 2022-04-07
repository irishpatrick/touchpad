package main

// #cgo CFLAGS: -Idriver
// #cgo LDFLAGS: -Ldriver/build -ldriver
// #include "driver.h"
import "C"

import (
    "encoding/json"
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
    "time"

    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
    "github.com/dgrijalva/jwt-go"
    //"golang.org/x/oath2"
    //"golang.org/x/oath2/google"
)

var homeTemplate = template.Must(template.ParseFiles("./static/touchpad.html"))
var addr = flag.String("addr", "0.0.0.0:8080", "http service address")
var certFile = flag.String("cert", "", "tls cert file")
var keyFile = flag.String("key", "", "tlk key file")
var upgrader = websocket.Upgrader{}

var isAlive = false
var aliveTimer = time.Now()

/*var googleOauthConfig = &oauth2.Config{
    RedirectURL:    "http://" + addr + "/auth/google/callback",
    ClientID:       os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
    ClientSecret:   os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
    Scopes:         []string{"https://www.googleapis.com/auth/userinfo.email"},
    Endpoint:       google.Endpoint,
}*/

var users = map[string]string{
    "user1": "password1",
    "user2": "password2",
}

var jwtKey = []byte("my_secret_key")

type Creds struct {
    Password string `json:"password"`
    Username string `json:"username"`
}

type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

func ResetAliveTimer() {
    isAlive = true
    aliveTimer = time.Now().Add(1 * time.Minute)
}

func echo (w http.ResponseWriter, r *http.Request) {
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade: ", err)
        return
    }
    defer c.Close()

    for {
        if time.Now().Unix() > aliveTimer.Unix() {
            isAlive = false;
            break
        }

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

func bind(w http.ResponseWriter, r *http.Request) {
    var creds Creds

    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    expectedPwd, ok := users[creds.Username]
    if !ok || expectedPwd != creds.Password {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &Claims{
        Username: creds.Username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenStr, err := token.SignedString(jwtKey)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name: "token",
        Value: tokenStr,
        Expires: expirationTime,
    })
}

func alive(w http.ResponseWriter, r *http.Request) {
    c, err := r.Cookie("token")
    if err != nil {
        if err == http.ErrNoCookie {
            w.WriteHeader(http.StatusUnauthorized)
            isAlive = false;
            return
        }

        w.WriteHeader(http.StatusBadRequest)
        isAlive = false;
        return
    }

    tokenStr := c.Value
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            w.WriteHeader(http.StatusUnauthorized)
            isAlive = false;
            return
        }

        w.WriteHeader(http.StatusBadRequest)
        isAlive = false;
        return
    }

    if !token.Valid {
        w.WriteHeader(http.StatusUnauthorized)
        isAlive = false;
        return
    }

    ResetAliveTimer()
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
        fmt.Println("\nshutting down...")
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
    router.HandleFunc("/bind", bind).Methods("POST")
    router.HandleFunc("/alive", alive).Methods("POST")
    router.HandleFunc("/echo", echo)
    router.HandleFunc("/", home)

    http.Handle("/", router)

    if len(*certFile) > 0 && len(*keyFile) > 0 {
        log.Fatal(http.ListenAndServeTLS(*addr, *certFile, *keyFile, nil))
    } else {
        log.Fatal(http.ListenAndServe(*addr, nil))
    }
}

