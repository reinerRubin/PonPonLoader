# /

## usage
### build
```
go build -o bin/pon-pon-loader cmd/loader/main.go
```

### run
```
./bin/pon-pon-loader  http://boards.4chan.org/vr/thread/XXX /path/to/result/folder/
```
or for download all images until the thread dead
```
./bin/pon-pon-loader --watch  http://boards.4chan.org/vr/thread/XXX /path/to/result/folder/
```
