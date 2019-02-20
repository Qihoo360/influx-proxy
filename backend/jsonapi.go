package backend

import (
    "bytes"
    "compress/gzip"
    "encoding/json"
)

/*
	for the extension, json api is needed.
	TODO: fix string
*/

type seri struct {
    Name    string     `json:"name,omitempty"`
    Columns []string   `json:"columns"`
    Values  [][]string `json:"values"`
}

type statement struct {
    StatementId int    `json:"statement_id"`
    Series      []seri `json:"series,omitempty"`
}

type statementArray struct {
    Results []statement `json:"results"`
}

// GetSerisArray byte转化为seri
func GetSerisArray(sBody []byte) (ss []seri, err error) {
    var tmp statementArray
    err = json.Unmarshal(sBody, &tmp)
    if err == nil {
        if len(tmp.Results) > 0 && len(tmp.Results[0].Series) > 0 {
            ss = tmp.Results[0].Series
        }
    }
    return
}

// GetValuesArray byte转化为[][]string
func GetValuesArray(sBody []byte) (ms [][]string, err error) {
    var tmp statementArray
    err = json.Unmarshal(sBody, &tmp)
    if err == nil {
        if len(tmp.Results) > 0 && len(tmp.Results[0].Series) > 0 {
            ms = tmp.Results[0].Series[0].Values
        }
    }
    return
}

// GetJsonBodyfromSeries seri转化为byte
func GetJsonBodyfromSeries(series []seri) (body []byte, err error) {
    tmpstatement := statement{
        StatementId: 0,
        Series:      series,
    }
    body, err = json.Marshal(statementArray{
        Results: []statement{tmpstatement},
    })
    if err == nil {
        body = append(body, '\n')
    }
    return
}

// GetJsonBodyfromValues [][]string转化为byte
func GetJsonBodyfromValues(name string, columns []string, values [][]string) (body []byte, err error) {
    tmpseri := seri{
        Name:    name,
        Columns: columns,
        Values:  values,
    }
    tmpstatement := statement{
        StatementId: 0,
        Series:      []seri{tmpseri},
    }
    body, err = json.Marshal(statementArray{
        Results: []statement{tmpstatement},
    })
    if err == nil {
        body = append(body, '\n')
    }
    return
}

// GzipEncode 把byte类型压缩
func GzipEncode(body []byte, need bool) (b []byte) {
    if !need {
        b = body
    } else {
        var buf bytes.Buffer
        w := gzip.NewWriter(&buf)
        defer w.Close()
        w.Write(body)
        w.Flush()
        b = buf.Bytes()
    }
    return
}
