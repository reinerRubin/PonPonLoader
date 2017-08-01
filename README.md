# /
pon-pon-loader is an yet another 4chan thread image grabber in go.

The main "feature" (lol) is the watch flag. So you can run downloader with watch flag and all images will be downloaded until 404.

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
