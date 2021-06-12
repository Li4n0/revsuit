package geoip

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	log "unknwon.dev/clog/v2"
)

const (
	Url = "https://download.maxmind.com/app/geoip_download?edition_id=GeoLite2-City&license_key=%s&suffix=tar.gz"
)

func get(url string) (b []byte, err error) {
	client := http.Client{
		Timeout: 90 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // disable verify
		}}
	request, _ := http.NewRequest(http.MethodGet, url, nil)
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.77 Safari/537.36")

	resp, err := client.Do(request)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

func extractTarGz(gzipStream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeReg:
			log.Trace(header.Name)
			if !strings.HasSuffix(header.Name, "GeoLite2-City.mmdb") {
				continue
			}
			outFile, err := os.Create("GeoLite2-City.mmdb")
			if err != nil {
				return err
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
		}
	}
	return nil
}

func download(licenseKey string) (err error) {
	var GeoLite2TarGz []byte
	log.Info("Downloading GeoLite2-City.mmdb...")
	if GeoLite2TarGz, err = get(fmt.Sprintf(Url, licenseKey)); err != nil {
		return err
	}

	return extractTarGz(bytes.NewReader(GeoLite2TarGz))
}
