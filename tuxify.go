package tuxify

import (
	"bytes"
	"crypto/aes"
	"crypto/rand"
	"image"

	"image/color"
)

func Tuxify(key []byte, src image.Image) (image.Image, error) {
	rect := src.Bounds()

	// Put all raw RGB values into a buffer...
	buffy := &bytes.Buffer{}
	pixel := make([]byte, 3)
	for y := 0; y < rect.Dy(); y++ {
		for x := 0; x < rect.Dx(); x++ {
			rgba := color.RGBAModel.Convert(src.At(x, y))
			r, g, b, _ := rgba.RGBA()
			pixel[0], pixel[1], pixel[2] = byte(r), byte(g), byte(b)
			buffy.Write(pixel)
		}
	}

	// Encrypt RGB values with the given key (or a random key if the key parameter is nil)
	ciphertext, err := encrypt(key, buffy.Bytes())
	if err != nil {
		return nil, err
	}

	// Put the encrypted RGB values into a new image...
	r := bytes.NewReader(ciphertext)
	dst := image.NewRGBA(rect)
	for y := 0; y < rect.Dy(); y++ {
		for x := 0; x < rect.Dx(); x++ {
			r.Read(pixel)
			dst.Set(x, y, color.RGBA{pixel[0], pixel[1], pixel[2], 255})
		}
	}

	return dst, nil
}

func encrypt(key []byte, data []byte) ([]byte, error) {
	if key == nil {
		key = make([]byte, 16)
		rand.Read(key)
	}
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := cipher.BlockSize()
	if mod := len(data) % bs; mod != 0 {
		zeros := make([]byte, bs-mod)
		data = append(data, zeros...)
	}
	for block := data; len(block) > 0; block = block[bs:] {
		cipher.Encrypt(block, block)
	}
	return data, nil
}
