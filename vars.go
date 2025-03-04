package main

import "flag"

var (
	flagServerName     = flag.String("name", "smtp2http", "the server name")
	flagListenAddr     = flag.String("listen", ":smtp", "the smtp address to listen on")
	flagMaxMessageSize = flag.Int64("msglimit", 1024*1024*2, "maximum incoming message size")
	flagReadTimeout    = flag.Int("timeout.read", 5, "the read timeout in seconds")
	flagWriteTimeout   = flag.Int("timeout.write", 5, "the write timeout in seconds")
	flagGotifyURL      = flag.String("gotify", "http://localhost:8080/", "the gotify to send the message to")
)

func init() {
	flag.Parse()
}
