// Usage:
// $ curl -s -F 'img=@foo.png' \
//           -F 'key=00000000000000000000000000000000' \
//           -o out.png \
//     http://localhost:8080
// The 'key' form value is optional. If omitted, the image
// will be encrypted with a random key.
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
	"jdtw.dev/tuxify"
)

var (
	port      = flag.Int("port", 8080, "Port to listen on")
	maxBytes  = flag.Int64("max-bytes", 5<<20, "Request size limit")
	rateLimit = flag.Duration("rate", 1*time.Second, "Rate limit")
)

const (
	// The image file.
	formFile = "img"
	// The key. Must be 16 hex-encoded bytes.
	formKey = "key"
)

func main() {
	flag.Parse()

	limiter := rate.NewLimiter(rate.Every(*rateLimit), 10)

	s := http.NewServeMux()
	s.HandleFunc("/", tuxifyHandler(limiter))

	http.ListenAndServe(fmt.Sprintf(":%d", *port), s)
}

func tuxifyHandler(l *rate.Limiter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Hey, money doesn't grow on trees...
		if err := l.Wait(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}

		// Validate the key before we do any heavy lifting...
		var key []byte
		if hexKey := r.FormValue(formKey); hexKey != "" {
			decoded, err := hex.DecodeString(hexKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			if len(decoded) != 16 {
				http.Error(w, "Key must be 16 bytes long", http.StatusBadRequest)
				return
			}
			key = decoded
		}

		// Get the form file...
		r.Body = http.MaxBytesReader(w, r.Body, *maxBytes)
		f, hdr, err := r.FormFile(formFile)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer f.Close()
		log.Printf("Tuxifying %s, Content-Type: %s, Size: %d", hdr.Filename, hdr.Header.Get("Content-Type"), hdr.Size)

		// Decode the image...
		src, _, err := image.Decode(f)
		if err != nil {
			log.Printf("image.Decode(File:%s, Content-Type:%s, Size:%d) failed: %v", hdr.Filename, hdr.Header.Get("Content-Type"), hdr.Size, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Do the transform...
		dst, err := tuxify.Tuxify(key, src)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// And respond...
		if err := png.Encode(w, dst); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/png")
	}
}
