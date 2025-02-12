# gosoc


<img src="https://github.com/user-attachments/assets/e93e817e-7ba9-4c07-ab0a-4ad3544cae22" width="400">


## Websocket library written in Go, but the author only kinda remembers golang
A websocket library written in go. If I were you I wouldn't use this in prod. This is mainly to learn go bc i wanna.

## References
- [RFC 6455](https://www.rfc-editor.org/rfc/rfc6455.html#section-5)

## Testing
Goal is to make https://github.com/crossbario/autobahn-testsuite happy (i.e. make this gosoc rfc compliant)
Start autobahn via:
```
docker run -it --rm \
    -v "${PWD}/config:/config" \
    -v "${PWD}/reports:/reports" \
    -p 9001:9001 \
    --name fuzzingserver \
    crossbario/autobahn-testsuite
```

Some unit tests use the fuzzing server. Once the server has been started, test with `go test -v ./client`
