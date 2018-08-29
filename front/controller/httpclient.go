package controller

import (
	"github.com/valyala/fasthttp"
	"time"
	"errors"
	"strconv"
)

func HttpGet200ok(url string) ([]byte, error){
	code, body, err := fasthttp.GetTimeout(nil, url, 3*time.Second)
	if err != nil {
		return nil, err
	}

	if code != 200 {
		return nil, errors.New(strconv.Itoa(code))
	}

	return body, nil
}
