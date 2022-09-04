# whistle-pig

A generic HTTP server in go with raw upload, uid, and BLAKE3 features.
This is a template for microservices written in the go language.

The docker image is extremely small at ~ 7.3MB via using `scratch` and a stand-alone go binary.

WARNING: <b>The default settings allow unlimited data transfer over plain HTTP</b>:

## Uploading and BLAKE3 hash functions

Both the /api/v0/uploadbody and /api/v0/hashbody store the entire HTTP BODY in the whistle-pig STDOUT log in hex. The hashbody fuction returns a JSON with a BLAKE3 hash, storing the hash and the body on the server as hex, while return a success message JSON, storing the entire body on the server as hex.

```
$ time curl -X POST --data-binary @/usr/bin/grep https://myserverplace:8088/api/v0/hashbody
"{'status': 'BLAKE3 256 trunc hash', 'type': 'http body', 'b3o256': 'bf592f6a20f5469dec341d1e7d5883c57d3a0978e9030ba3bed7c51ce028a6cc'}"
real    0m0.009s
user    0m0.007s
sys     0m0.000s

```

### UID generation

There is a /uid function that will generate a UUID and return it directly.

```
time curl -k https://localhost:8088/uid
a1451cf5-9ad2-44cd-b06b-a9f8c0c9bc11
real    0m0.009s
user    0m0.004s
sys     0m0.004s
```


#### TLS

There is TLSv1.3 implemented in the template. The files are cert.pem and key.pem, copied in via Dockerfile to the working directory in the container or otherwise in pwd when used directly.
