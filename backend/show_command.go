package backend

import (
	"fmt"
	"net/http"
	"strings"
)

func (ic *InfluxCluster) QueryAll(req *http.Request) (s_header http.Header, bodys [][]byte, err error) {
	bodys = make([][]byte, 0)
	for _, v := range ic.m2bs {
		need := false
		actu := false

		for _, api := range v {
			if api.GetZone() != ic.Zone {
				continue
			}
			if (!api.IsActive() || api.IsWriteOnly()) {
				continue
			}
			need = true

			_header, _, s_body, _err := api.JustQuery(req)
			if _err != nil {
				err = _err
				continue
			}

			s_header = _header
			bodys = append(bodys, s_body)
			actu = true
			break
		}

		if need && !actu {
			s_header = nil
			bodys = nil
			return
		}
	}
	err = nil
	return
}

func (ic *InfluxCluster) ShowQuery(w http.ResponseWriter, req *http.Request) (err error) {
	fmt.Println("here to do SHOW *")
	f_header, bodys, _err := ic.QueryAll(req)
	err = _err
	if _err != nil {
		err = _err
		return
	}

	var f_body []byte

	q := strings.TrimSpace(req.FormValue("q"))
	if strings.Contains(q, "field") || strings.Contains(q, "tag") {
		// TODO: combine series

	} else {
		name := ""
		columns := []string{"key"}
		if strings.Contains(q, "measurements") {
			name = "measurements"
			columns = []string{"name"}
		}

		occur := make(map[string]bool)
		values := make([][]string, 0)
		for _, s_body := range bodys {
			s_ms, _err := GetValuesArray(s_body)
			if _err != nil {
				err = _err
				return
			}
			for _, s := range s_ms {
				if strings.Contains(s[0], "influxdb.cluster") {
					continue
				}
				if !occur[s[0]] {
					values = append(values, s)
					occur[s[0]] = true
				}
			}
		}

		f_body, err = GetJsonBodyfromValues(name, columns, values)
		if err != nil {
			return
		}
	}

	copyHeader(w.Header(), f_header)
	w.WriteHeader(200)
	w.Write(GzipEncode(f_body, f_header.Get("Content-Encoding") == "gzip"))
	err = nil
	return
}
