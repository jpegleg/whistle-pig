package main
import (
    "fmt"
    "io"
    "net/http"
    "crypto/tls"
    "encoding/hex"
    "time"
    "flag"
    "github.com/google/uuid"
    "lukechampine.com/blake3"
)

func verhandler(w http.ResponseWriter, r *http.Request) {
    clienta := r.RemoteAddr
    dt := time.Now()
    fmt.Fprint(w, "API HEALTHY v0")
    fmt.Println(dt.String(), "GENERIC HTTP SERV - health check recv", clienta, r.URL.Path[1:])
}

func verjsonhandler(w http.ResponseWriter, r *http.Request) {
    clienta := r.RemoteAddr
    dt := time.Now()
    fmt.Fprint(w, "\"{'server': 'API HEALTHY v0'}\"")
    fmt.Println(dt.String(), "GENERIC HTTP SERV - health check recv", clienta, r.URL.Path[1:])
}

func uidhandler(w http.ResponseWriter, r *http.Request) {
    uuidW := uuid.New()
    id := uuidW.String()
    clienta := r.RemoteAddr
    dt := time.Now()
    fmt.Fprint(w, "", id)
    fmt.Println(dt.String(), "resource accessed", clienta, r.URL.Path[1:], "refid:", id)
}

func dochandler(w http.ResponseWriter, r *http.Request) {
    clienta := r.RemoteAddr
    dt := time.Now()
    fmt.Fprint(w, "\"{'DOCUMENT': 'my document here'}\"")
    fmt.Println(dt.String(), "resource accessed", clienta, r.URL.Path[1:])
}

func bodyhandler(w http.ResponseWriter, r *http.Request) {
    clienta := r.RemoteAddr
    defer r.Body.Close()
    buf, err := io.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }
    encodedString := hex.EncodeToString(buf)
    dt := time.Now()
    fmt.Fprint(w, "\"{'status': 'uploaded', 'type': 'http body'}\"")
    fmt.Println(dt.String(), "resource accessed", clienta, r.URL.Path[1:], "http body recv as hex: ", encodedString)
}

func hashhandler(w http.ResponseWriter, r *http.Request) {
    clienta := r.RemoteAddr
    defer r.Body.Close()
    buf, err := io.ReadAll(r.Body)
    if err != nil {
        panic(err)
    }
    encodedString := hex.EncodeToString(buf)
    valByte := blake3.Sum256([]byte(encodedString))
    slice32 := valByte[:]
    encodedB3 := hex.EncodeToString(slice32)
    dt := time.Now()
    fmt.Fprint(w, "\"{'status': 'BLAKE3 256 trunc hash', 'type': 'http body', 'b3o256': '", encodedB3, "'}\"")
    fmt.Println(dt.String(), "resource accessed", clienta, r.URL.Path[1:], "http body recv as hex: ", encodedString, "BLAKE3 truncated to 256 bytes as hex: ", encodedB3)
}

func main() {
    certFile := flag.String("certfile", "cert.pem", "certificate PEM file")
    keyFile := flag.String("keyfile", "key.pem", "key PEM file")
    flag.Parse()

    http.HandleFunc("/", verjsonhandler)

    http.HandleFunc("/uid", uidhandler)
    http.HandleFunc("/doc", dochandler)
    http.HandleFunc("/health", verhandler)
    http.HandleFunc("/ver", verjsonhandler)
    http.HandleFunc("/api/v0/sendbody", bodyhandler)
    http.HandleFunc("/api/v0/hashbody", hashhandler)

    s := &http.Server{
        Addr:           ":8088",
        //ReadTimeout:    10 * time.Second,
        //WriteTimeout:   10 * time.Second,
        //MaxHeaderBytes: 1 << 32,
        TLSConfig: &tls.Config{
          MinVersion:               tls.VersionTLS13,
          PreferServerCipherSuites: true,
      },
    }
    s.ListenAndServeTLS(*certFile, *keyFile)
}
