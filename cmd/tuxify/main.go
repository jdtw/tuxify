package main

import (
	"encoding/hex"
	"flag"
	"image"
	_ "image/jpeg"
	"image/png"
	"io"
	"log"
	"os"

	"jdtw.dev/tuxify"
)

var (
	in     = flag.String("in", "", "Path to input image")
	out    = flag.String("out", "", "Output file (defaults to stdout)")
	keyHex = flag.String("key", "", "Hex-encoded key")
)

func main() {
	flag.Parse()

	f, err := os.Open(*in)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var key []byte
	if *keyHex != "" {
		decoded, err := hex.DecodeString(*keyHex)
		if err != nil {
			log.Fatal(err)
		}
		key = decoded
	}

	src, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	dst, err := tuxify.Tuxify(key, src)
	if err != nil {
		log.Fatal(err)
	}

	var w io.Writer = os.Stdout
	if *out != "" {
		of, err := os.Create(*out)
		if err != nil {
			log.Fatal(err)
		}
		defer of.Close()
		w = of
	}

	if err := png.Encode(w, dst); err != nil {
		log.Fatal(err)
	}
}
