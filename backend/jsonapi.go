package backend

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
)

/*
	for the extension, json api is needed.
	TODO: fix string
 */

type seri struct {
	Name string `json:"name,omitempty"`
	Columns []string `json:"columns"`
	Values [][]string `json:"values"`
}

type statement struct {
	Statement_id int `json:"statement_id"`
	Series []seri `json:"series,omitempty"`
}

type statement_array struct {
	Results []statement `json:"results"`
}

func fmtlog(title []byte, content []byte, err error) {
	fmt.Println(string(title) + ":")
	fmt.Println("\terr:" +err.Error())
	fmt.Println("\tcontent:"+string(content))
}

func GetSerisArray(s_body []byte) (ss []seri, err error) {
	var tmp statement_array
	err = json.Unmarshal(s_body, &tmp)
	if err == nil {
		if len(tmp.Results) > 0 {
			if len(tmp.Results[0].Series) > 0 {
				ss = tmp.Results[0].Series
				return
			}
		}
		ss = make([]seri, 0)
		return
	}
	fmtlog([]byte("json unmarshal"), s_body, err)
	return
}

func GetValuesArray(s_body []byte) (ms [][]string, err error) {
	var tmp statement_array
	err = json.Unmarshal(s_body, &tmp)
	if err == nil {
		if len(tmp.Results) > 0 {
			if len(tmp.Results[0].Series) > 0 {
				ms = tmp.Results[0].Series[0].Values
				return
			}
		}
		ms = make([][]string, 0)
		return
	}
	fmtlog([]byte("json unmarshal"), s_body, err)
	return
}

func GetJsonBodyfromSeries(series []seri) (body []byte, err error){
	tmpstatement := statement {
		Statement_id: 0,
		Series: series,
	}
	body, err = json.Marshal(statement_array{
		Results: []statement{tmpstatement},
	})
	if err == nil {
		body = append(body, '\n')
	}
	return
}

func GetJsonBodyfromValues(name string, columns []string, values [][]string) (body []byte, err error) {
	tmpseri := seri {
		Name: 	name,
		Columns: 	columns,
		Values:	values,
	}
	tmpstatement := statement{
		Statement_id: 0,
		Series: []seri{tmpseri},
	}
	body, err = json.Marshal(statement_array{
		Results: []statement{tmpstatement},
	})
	if err == nil {
		body = append(body, '\n')
	}
	return
}

func GzipEncode(body []byte, need bool) (b []byte) {
	if !need {
		b = body
	} else {
		var buf bytes.Buffer
		w := gzip.NewWriter(&buf)
		w.Write(body)
		defer w.Close()
		w.Flush()
		b = buf.Bytes()
	}
	return
}