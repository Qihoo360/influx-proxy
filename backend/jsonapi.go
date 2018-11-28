package backend

import (
	"encoding/json"
	"fmt"
)

/*
	for the extension, json api is needed.
	todo fix string
 */

type seri struct {
	Name string
	Columns []string
	Values [][]string
}

type statement struct {
	Statement_id int
	Series []seri
}

type statement_array struct {
	Results []statement
}

func GetMeasurementsArray(s_body []byte) (ms [][]string, err error) {
	var tmp statement_array
	fmt.Println(string(s_body))
	err = json.Unmarshal(s_body, &tmp)
	ms = tmp.Results[0].Series[0].Values
	return
}

func GetJsonBody(values [][]string) (body []byte, err error) {
	tmpseri := seri {
		Name: 	"measurements",
		Columns: 	[]string{"name"},
		Values:	values,
	}
	tmpstatement := statement{
		Statement_id: 0,
		Series: []seri{tmpseri},
	}

	body, err = json.Marshal(statement_array{
		Results: []statement{tmpstatement},
	})
	body = append(body, '\n')

	return
}