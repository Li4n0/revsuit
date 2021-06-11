package qqwry

import (
	"bytes"
	"compress/zlib"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	log "unknwon.dev/clog/v2"
)

const (
	CopyWriteUrl = "https://qqwry.mirror.noc.one/copywrite.rar"
	Url          = "https://qqwry.mirror.noc.one/qqwry.rar"
)

func get(url string) (b []byte, err error) {
	client := http.Client{
		Timeout: 90 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // disable verify
		}}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Add("User-Agent", "Nali/2.1.2 (Nali CLI, https://nali.skk.moe)")

	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	b, err = ioutil.ReadAll(resp.Body)
	return b, err
}

func getKey(b []byte) (key uint32, err error) {
	if len(b) != 280 {
		return 0, errors.New("copywrite.rar is corrupt")
	}
	key = binary.LittleEndian.Uint32(b[20:])
	return key, err
}

func decrypt(b []byte, key uint32) (_ []byte, err error) {
	for i := 0; i < 0x200; i++ {
		key *= uint32(0x805)
		key++
		key &= uint32(0xff)
		b[i] = b[i] ^ byte(key)
	}
	rc, err := zlib.NewReader(bytes.NewBuffer(b))
	if err != nil {
		return
	}
	defer rc.Close()
	return ioutil.ReadAll(rc)
}

func download() (err error) {
	var (
		copyWriteData, qqwryData []byte
		wg                       sync.WaitGroup
	)
	log.Info("Downloading qqwry.dat...")
	wg.Add(2)
	go func() {
		defer wg.Done()
		if copyWriteData, err = get(CopyWriteUrl); err != nil {
			return
		}
	}()

	go func() {
		defer wg.Done()
		if qqwryData, err = get(Url); err != nil {
			return
		}
	}()
	wg.Wait()
	if err != nil {
		return err
	}
	var key uint32
	if key, err = getKey(copyWriteData); err != nil {
		return err
	}
	b, err := decrypt(qqwryData, key)
	if err != nil {
		return err
	}
	_ = ioutil.WriteFile("qqwry.dat", b, 0644)

	return nil
}
