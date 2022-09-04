# whistle-pig

A generic HTTP server in go with raw upload, uid, and BLAKE3 features.
This is a template for microservices written in the go language.

The docker image is extremely small at ~ 7.3MB via using `scratch` and a stand-alone go binary.

WARNING: <b>The default settings allow unlimited data transfer over plain HTTP</b>:

## Uploading and BLAKE3 hash functions

Both the /api/v0/uploadbody and /api/v0/hashbody store the entire HTTP BODY in the whistle-pig STDOUT log in hex. The hashbody fuction returns a JSON with a BLAKE3 hash, storing the hash and the body on the server as hex, while return a success message JSON, storing the entire body on the server as hex.

```
$ time curl -X POST --data-binary @/usr/bin/grep myserverplace:8088/api/v0/hashbody
"{'status': 'BLAKE3 256 trunc hash', 'type': 'http body', 'b3o256': 'bf592f6a20f5469dec341d1e7d5883c57d3a0978e9030ba3bed7c51ce028a6cc'}"
real    0m0.009s
user    0m0.007s
sys     0m0.000s

```

### UID generation

There is a /uid function that will generate a UUID and return it directly.

```
time curl localhost:8088/uid
a1451cf5-9ad2-44cd-b06b-a9f8c0c9bc11
real    0m0.009s
user    0m0.004s
sys     0m0.004s
```


#### TLS usage

The default does not have TLS implemented, there is not security at all in the template.

Adding TLS to this natively causes issues for the minimal Docker image build,
however the binary works fine outside of the minimal container.

Add key.pem and cert.pem to pwd and use this code reference instead:

```
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
```
