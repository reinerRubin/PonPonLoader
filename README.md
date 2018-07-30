# /
pon-pon-loader is an yet another 4chan thread image grabber in go.

The main "feature" (lol) is a watch flag. So you can run the downloader with a watch flag and all images will be downloaded until 404.

## usage
### build
```
make
```

### run
```
./bin/pon-pon-loader  http://boards.4chan.org/vr/thread/XXX /path/to/result/folder/
```
or for downloading all images until the thread dead
```
./bin/pon-pon-loader --watch  http://boards.4chan.org/vr/thread/XXX /path/to/result/folder/
```
