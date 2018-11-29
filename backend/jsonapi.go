package backend

import (
	"encoding/json"

	"log"
)

/*
	for the extension, json api is needed.
	todo fix string
 */

type seri struct {
	Name string `json:"name"`
	Columns []string `json:"columns"`
	Values [][]string `json:"values"`
}

type statement struct {
	Statement_id int `json:"statement_id"`
	Series []seri `json:"series"`
}

type statement_array struct {
	Results []statement `json:"results"`
}

func GetMeasurementsArray(s_body []byte) (ms [][]string, err error) {
	var tmp statement_array
	log.Println("*******************************************")
	log.Println(string(s_body))
	log.Println("*******************************************")
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