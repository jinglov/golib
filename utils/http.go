package utils

import (
	"audience_api/lib/logger"
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

func DoRequestJson(method string, url string, data []byte, header http.Header, proxy *url.URL) ([]byte, error) {
	body := bytes.NewReader(data)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	if header != nil {
		req.Header = header
	}
	req.Header.Set("Content-Type", "application/json")
	return doRequest(req, proxy)
}

func doRequest(req *http.Request, proxy *url.URL) (b []byte, err error) {
	req.Header.Set("Connection", "Keep-Alive")
	/*
		设置transport，先复制默认的Transpport过来，然后再根据当前情况做变化
	*/
	var ts *http.Transport = new(http.Transport)
	*ts = *http.DefaultTransport.(*http.Transport)
	if req.URL.Scheme == "https" {
		ts.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	/*
		设置代理服务器
	*/
	if proxy != nil {
		ts.Proxy = http.ProxyURL(proxy)
	}
	client := &http.Client{
		Transport: ts,
		//Timeout:   time.Second * 15,
	}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		logger.Error(err)
		return
	}
	defer resp.Body.Close()
	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return
	}
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		err = errors.New("return Status Code is " + strconv.Itoa(resp.StatusCode))
	}
	return
}
