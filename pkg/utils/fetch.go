package utils

import (
	"crypto/rand"
	"github.com/parnurzeal/gorequest"
	"golang.org/x/xerrors"
	"log"
	"math"
	"math/big"
	"time"
)

func RandInt() int {
	seed, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	return int(seed.Int64())
}

func FetchURL(url, apikey string, retry int) (res []byte, err error) {
	for i := 0; i <= retry; i++ {
		if i > 0 {
			wait := math.Pow(float64(i), 2) + float64(RandInt()%10)
			log.Printf("retry after %f seconds\n", wait)
			time.Sleep(time.Duration(time.Duration(wait) * time.Second))
		}
		res, err = fetchURL(url, map[string]string{"api-key": apikey})
		if err == nil {
			return res, nil
		}
	}
	return nil, xerrors.Errorf("failed to fetch URL: %w", err)
}
func fetchURL(url string, headers map[string]string) ([]byte, error) {
	req := gorequest.New().Get(url)
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	resp, body, errs := req.Type("text").EndBytes()
	if len(errs) > 0 {
		return nil, xerrors.Errorf("HTTP error. url: %s, err: %w", url, errs[0])
	}
	if resp.StatusCode != 200 {
		return nil, xerrors.Errorf("HTTP error. status code: %d, url: %s", resp.StatusCode, url)
	}
	return body, nil
}
