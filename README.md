[![Go Report Card](https://goreportcard.com/badge/github.com/pwed/disttube)](https://goreportcard.com/report/github.com/pwed/disttube)

# DistTube
**A decentralized video hosting platform.**

**_Currently in development don't expect everything to work!_**

## Installation
*As of 08/03/2017*

Dependencies: GoLang (1.8 or newer), ffmpeg


### Windows

Download Golang from the official website and install

Download ffmpeg and add ffmpeg\bin to your %PATH%

```
go get github.com/pwed/disttube
```

### Linux

```
sudo apt install ffmpeg golang -y
go get github.com/pwed/disttube
```

*If ffmpeg complains about aac being experimental then you need a newer version*

*If using an earlier version of GoLang you will need to set up a GOPATH and GOROOT*

### Mac

Look at the linux instructions and try your luck.

## Use

From $GOPATH/src/github.com/pwed/disttube, run `go run main.go` in your terminal.

Open `localhost:8080` in your browser.

Have fun.

## TODO

- [x] Ingestion _(is mostly working :D )_
- [ ] Users
- [ ] Channels
- [ ] Video Player
- [ ] Video Browser
- [ ] Admin Tools
- [ ] Config File Setup
- [ ] Docker Container
- [ ] Aggregation Server
- [ ] MPEG-DASH support