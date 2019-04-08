// Copyright 2016 Eleme. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package backend

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/wilhelmguo/influx-proxy/logs"
)

var (
	ErrBadRequest = errors.New("Bad Request\n")
	ErrNotFound   = errors.New("Not Found\n")
	ErrInternal   = errors.New("Internal Error")
	ErrUnknown    = errors.New("Unknown Error\n")
)

func Compress(buf *bytes.Buffer, p []byte) (err error) {
	zip := gzip.NewWriter(buf)
	n, err := zip.Write(p)
	if err != nil {
		return
	}
	if n != len(p) {
		err = io.ErrShortWrite
		return
	}
	err = zip.Close()
	return
}

type HttpBackend struct {
	BasicAuth *BasicAuth
	client    *http.Client
	transport http.Transport
	Interval  int
	URL       string
	DB        string
	Zone      string
	Active    bool
	running   bool
	WriteOnly int
}

func NewHttpBackend(cfg *BackendConfig) (hb *HttpBackend) {
	hb = &HttpBackend{
		client: &http.Client{
			Timeout: time.Millisecond * time.Duration(cfg.Timeout),
		},
		// TODO: query timeout? use req.Cancel
		// client_query: &http.Client{
		//     Timeout: time.Millisecond * time.Duration(cfg.TimeoutQuery),
		// },
		BasicAuth: cfg.BasicAuth,
		Interval:  cfg.CheckInterval,
		URL:       cfg.URL,
		DB:        cfg.DB,
		Zone:      cfg.Zone,
		Active:    true,
		running:   true,
		WriteOnly: cfg.WriteOnly,
	}
	go hb.CheckActive()
	return
}

// TODO: update active when calling successed or failed.

func (hb *HttpBackend) CheckActive() {
	var err error
	for hb.running {
		_, err = hb.Ping()
		hb.Active = (err == nil)
		time.Sleep(time.Millisecond * time.Duration(hb.Interval))
	}
}

func (hb *HttpBackend) IsWriteOnly() bool {
	if hb.WriteOnly == 0 {
		return false
	}
	return true
}

func (hb *HttpBackend) IsActive() bool {
	return hb.Active
}

func (hb *HttpBackend) Ping() (version string, err error) {
	resp, err := hb.client.Get(hb.URL + "/ping")
	if err != nil {
		logs.Error("http error: ", err)
		return
	}
	defer resp.Body.Close()

	version = resp.Header.Get("X-Influxdb-Version")

	if resp.StatusCode == 204 {
		return
	}
	logs.Errorf("write status code: %d, the backend is %s\n", resp.StatusCode, hb.URL)

	respbuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("readall error: ", err)
		return
	}
	logs.Errorf("error response: %s\n", respbuf)
	return
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func (hb *HttpBackend) GetZone() (zone string) {
	return hb.Zone
}

func (hb *HttpBackend) QueryResp(req *http.Request) (header http.Header, status int, body []byte, err error) {
	if len(req.Form) == 0 {
		req.Form = url.Values{}
	}
	req.Form.Set("db", hb.DB)
	req.ContentLength = 0

	req.URL, err = url.Parse(hb.URL + "/query?" + req.Form.Encode())
	if err != nil {
		logs.Error("internal url parse error: ", err)
		return
	}

	q := strings.TrimSpace(req.FormValue("q"))

	resp, err := hb.transport.RoundTrip(req)
	defer resp.Body.Close()
	if err != nil {
		logs.Errorf("query error: %s,the query is %s\n", err, q)
		hb.Active = false
		return
	}

	respDody := resp.Body
	if resp.Header.Get("Content-Encoding") == "gzip" {
		respDody, err = gzip.NewReader(resp.Body)
		defer respDody.Close()
		if err != nil {
			logs.Errorf("unable to decode gzip body")
			return
		}
	}

	body, err = ioutil.ReadAll(respDody)
	if err != nil {
		logs.Errorf("read body error: %s,the query is %s\n", err, q)
		return
	}

	header = resp.Header
	status = resp.StatusCode
	return
}

// Don't setup Accept-Encoding: gzip. Let real client do so.
// If real client don't support gzip and we setted, it will be a mistake.
func (hb *HttpBackend) Query(w http.ResponseWriter, req *http.Request) (err error) {
	if len(req.Form) == 0 {
		req.Form = url.Values{}
	}
	req.Form.Set("db", hb.DB)
	req.ContentLength = 0

	// Add basic auth
	hb.basicAuth(req)

	req.URL, err = url.Parse(hb.URL + "/query?" + req.Form.Encode())
	if err != nil {
		logs.Error("internal url parse error: ", err)
		return
	}

	q := strings.TrimSpace(req.FormValue("q"))
	resp, err := hb.transport.RoundTrip(req)
	if err != nil {
		logs.Errorf("query error: %s,the query is %s\n", err, q)
		hb.Active = false
		return
	}
	defer resp.Body.Close()

	copyHeader(w.Header(), resp.Header)

	p, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Errorf("read body error: %s,the query is %s\n", err, q)
		return
	}

	w.WriteHeader(resp.StatusCode)
	w.Write(p)
	return
}

func (hb *HttpBackend) basicAuth(req *http.Request) {
	// Add basic auth
	if hb.BasicAuth != nil {
		req.Header.Add("Authorization", fmt.Sprintf("Basic %s",
			base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", hb.BasicAuth.Username, hb.BasicAuth.Password)))))
	}
}

func (hb *HttpBackend) Write(p []byte) (err error) {
	var buf bytes.Buffer
	err = Compress(&buf, p)
	if err != nil {
		logs.Error("compress error: ", err)
		return
	}

	logs.Errorf("http backend write %s", hb.DB)
	err = hb.WriteStream(&buf, true)
	return
}

func (hb *HttpBackend) WriteCompressed(p []byte) (err error) {
	buf := bytes.NewBuffer(p)
	err = hb.WriteStream(buf, true)
	return
}

func (hb *HttpBackend) WriteStream(stream io.Reader, compressed bool) (err error) {
	q := url.Values{}
	q.Set("db", hb.DB)

	req, err := http.NewRequest("POST", hb.URL+"/write?"+q.Encode(), stream)
	if compressed {
		req.Header.Add("Content-Encoding", "gzip")
	}

	// Add basic auth
	hb.basicAuth(req)

	resp, err := hb.client.Do(req)
	if err != nil {
		logs.Error("http error: ", err)
		hb.Active = false
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 {
		return
	}
	logs.Error("write status code: ", resp.StatusCode)

	respbuf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("readall error: ", err)
		return
	}
	logs.Errorf("error response: %s\n", respbuf)

	// translate code to error
	// https://docs.influxdata.com/influxdb/v1.1/tools/api/#write
	switch resp.StatusCode {
	case 400:
		err = ErrBadRequest
	case 404:
		err = ErrNotFound
	default: // mostly tcp connection timeout
		logs.Errorf("status: %d", resp.StatusCode)
		err = ErrUnknown
	}
	return
}

func (hb *HttpBackend) Close() (err error) {
	hb.running = false
	hb.transport.CloseIdleConnections()
	return
}
