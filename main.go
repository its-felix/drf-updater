package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	drfURL = "https://update.drf.rs/drf.dll"
	dest   = "/Users/felix/Downloads/drf.dll"
)

func main() {
	fTemp, err := os.CreateTemp("", "drf_*.dll")
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		_ = fTemp.Close()
		_ = os.Remove(fTemp.Name())
	}()

	res, err := http.Get(drfURL)
	if err != nil {
		log.Fatal(err)
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Fatal("status != 200")
		return
	}

	defer res.Body.Close()
	if _, err = io.Copy(fTemp, res.Body); err != nil {
		log.Fatal(err)
		return
	}

	existingHash, err := sha256File(dest)
	notExists := false

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		if notExists = errors.Is(err, os.ErrNotExist); !notExists {
			log.Fatal(err)
			return
		}
	}

	if !notExists {
		newHash, err := sha256File(fTemp.Name())
		if err != nil {
			log.Fatal(err)
			return
		}

		if bytes.Equal(existingHash, newHash) {
			log.Println("is up to date")
			return
		}
	}

	if err = copyFile(fTemp.Name(), dest); err != nil {
		log.Fatal(err)
		return
	}

	log.Println("updated")
}

func sha256File(name string) ([]byte, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	h := sha256.New()
	if _, err = io.Copy(h, f); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}

func copyFile(src, dst string) error {
	fSrc, err := os.Open(src)
	if err != nil {
		return err
	}

	defer fSrc.Close()

	fDst, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer fDst.Close()

	_, err = io.Copy(fDst, fSrc)
	return err
}
