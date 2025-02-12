# gosoc
## Websocket library written in Go, but the author only kinda remembers golang
A websocket library written in go. If I were you I wouldn't use this in prod lol. This is mainly to learn go bc i wanna.


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

