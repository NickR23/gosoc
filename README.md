# gosoc


<img src="https://github.com/user-attachments/assets/e93e817e-7ba9-4c07-ab0a-4ad3544cae22" width="400">


## Websocket library written in Go, but the author only kinda remembers golang
A websocket library written in go. If I were you I wouldn't use this in prod.

## References
- [RFC 6455](https://www.rfc-editor.org/rfc/rfc6455.html)

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

### For the `go tests`
1. Start [server](https://github.com/jmalloc/echo-server): `sudo docker run --detach -p 10000:8080 jmalloc/echo-server`
2. Run tests: `go test -v ./client`
3. Profit????
