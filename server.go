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
    "net"
    "net/http"
    "os"
    "os/signal"
    "path"
    "strconv"
    "strings"
    "syscall"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/gorilla/mux"
    "github.com/gorilla/websocket"
)

/** Token Parameters **/
var TEMPLATE_HTML = "index.html"
var TOKEN_COOKIE_NAME = "token"
var TOKEN_TIME_VALID = 5 * time.Minute;

/** Flags **/
var port = flag.String("port", "8080", "http service port")
var addr = flag.String("addr", "0.0.0.0", "http service address")
var siteDir = flag.String("site", "./static/dist/", "static site assets")
var certFile = flag.String("cert", "", "tls cert file")
var keyFile = flag.String("key", "", "tlk key file")
var HOST = *addr + ":" + *port

var homeTemplate = template.Must(template.ParseFiles(path.Join(*siteDir, TEMPLATE_HTML)))
var upgrader = websocket.Upgrader{}

var isAlive = false
var aliveTimer = time.Now()

// TODO remove
var users = map[string]string{
    "user1": "password1",
    "user2": "password2",
}

// TODO remove
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

func GetOutboundIP() net.IP {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }

    defer conn.Close()

    localAddr := conn.LocalAddr().(*net.UDPAddr)

    return localAddr.IP
}

func url(w http.ResponseWriter, r *http.Request) {
    ipStr := GetOutboundIP().String()
    w.Write([]byte("http://" + ipStr + ":" + *port))
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
            log.Printf("token expired: current_time=%d, expiration_time=%d\n", time.Now().Unix(), aliveTimer.Unix())
            isAlive = false;
            break
        }

        mt, message, err := c.ReadMessage()
        if err != nil {
            log.Println("read: ", err)
            break
        }

        processCommand(message)
        err = c.WriteMessage(mt, message)
        if err != nil {
            log.Println("write: ", err)
            break
        }
    }
}

func home(w http.ResponseWriter, r *http.Request) {
    homeTemplate.Execute(w, r.Host)
}

func asset(w http.ResponseWriter, r *http.Request) {
    fn := path.Join(*siteDir, r.URL.Path)
    mime := "text/plain"
    if strings.HasSuffix(fn, ".js") {
        mime = "text/javascript"
    } else if strings.HasSuffix(fn, ".css") {
        mime = "text/css"
    }

    w.Header().Set("Content-Type", mime)

    buf, err := ioutil.ReadFile(fn)
    if err != nil {
        w.WriteHeader(http.StatusNotFound)
        fmt.Fprintln(w, "404 not found")
        return
    }

    w.Write(buf)
}

func bind(w http.ResponseWriter, r *http.Request) {
    r.Body = http.MaxBytesReader(w, r.Body, 32 << 20 + 512)
    r.ParseMultipartForm(32 << 20) // 32 Mb

    uname := r.FormValue("username")
    pass := r.FormValue("password")
    fmt.Printf("%s,%s\n", uname, pass)

    expectedPwd, ok := users[uname]
    if !ok || expectedPwd != pass {
        w.WriteHeader(http.StatusUnauthorized)
        return
    }

    expirationTime := time.Now().Add(TOKEN_TIME_VALID)
    claims := &Claims{
        Username: uname,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    ResetAliveTimer()
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenStr, err := token.SignedString(jwtKey)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name: TOKEN_COOKIE_NAME,
        Value: tokenStr,
        Expires: expirationTime,
    })
}

func alive(w http.ResponseWriter, r *http.Request) {
    c, err := r.Cookie(TOKEN_COOKIE_NAME)
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

func renew(w http.ResponseWriter, r *http.Request) {
    c, err := r.Cookie(TOKEN_COOKIE_NAME)
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

    expirationTime := time.Now().Add(TOKEN_TIME_VALID)
    claims.ExpiresAt = expirationTime.Unix()
    token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenStr, err = token.SignedString(jwtKey)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name: TOKEN_COOKIE_NAME,
        Value: tokenStr,
        Expires: expirationTime,
    })
}

func processCommand(msg []byte) {
    if len(msg) < 2 {
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
            C.driver_mouse_rel(C.int(dx), C.int(dy))
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

func main() {
    flag.Parse()
    log.SetFlags(0)

    log.Println("creating signal handlers...")
    setupHandlers()

    log.Println("creating virtual device...")
    errno := C.driver_create_device()
    if errno != 0 {
        log.Fatal("cannot create device! err=", errno)
    }

    log.Println("building routes...")
    router := mux.NewRouter()
    router.HandleFunc("/{[a-z]+}.js", asset).Methods("GET")
    router.HandleFunc("/{[a-z]+}.css", asset).Methods("GET")
    router.HandleFunc("/{[a-z]+}.ico", asset).Methods("GET")
    router.HandleFunc("/url", url).Methods("GET")
    router.HandleFunc("/bind", bind).Methods("POST")
    router.HandleFunc("/alive", alive).Methods("POST")
    router.HandleFunc("/renew", renew).Methods("POST")
    router.HandleFunc("/echo", echo)
    router.HandleFunc("/", home)

    http.Handle("/", router)

    if len(*certFile) > 0 && len(*keyFile) > 0 {
        log.Printf("serving at https://%s\n", HOST)
        log.Fatal(http.ListenAndServeTLS(HOST, *certFile, *keyFile, nil))
    } else {
        log.Printf("serving at http://%s\n", HOST)
        log.Fatal(http.ListenAndServe(HOST, nil))
    }
}

