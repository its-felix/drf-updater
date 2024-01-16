package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	drfURL = "https://update.drf.rs/drf.dll"
	dest   = `C:\Program Files (x86)\Steam\steamapps\common\Guild Wars 2\addons\arcdps\drf.dll`
)

func main() {
	err := run()
	if err != nil {
		log.Println(err)
	}

	for i := 0; i < 3; i++ {
		print(".")
		time.Sleep(time.Second)
	}

	if err != nil {
		os.Exit(1)
	}
}

func run() error {
	fTemp, err := os.CreateTemp("", "drf_*.dll")
	if err != nil {
		return err
	}

	defer func() {
		_ = fTemp.Close()
		_ = os.Remove(fTemp.Name())
	}()

	res, err := http.Get(drfURL)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New("status != 200")
	}

	defer res.Body.Close()
	if _, err = io.Copy(fTemp, res.Body); err != nil {
		return err
	}

	existingHash, err := sha256File(dest)
	notExists := false

	if err != nil && !errors.Is(err, os.ErrNotExist) {
		if notExists = errors.Is(err, os.ErrNotExist); !notExists {
			return err
		}
	}

	if !notExists {
		newHash, err := sha256File(fTemp.Name())
		if err != nil {
			return err
		}

		if bytes.Equal(existingHash, newHash) {
			log.Println("is up to date")
			return nil
		}
	}

	if err = copyFile(fTemp.Name(), dest); err != nil {
		return err
	}

	log.Println("updated")
	return nil
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
