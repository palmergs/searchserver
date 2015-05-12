package main

import (
	"fmt"
	"io/ioutil"
	"encoding/json"
	"github.com/palmergs/tokensearch"
)

func AppendTokens(pathToFile string, root *tokensearch.TokenNode) (int, error) {

	body, err := ioutil.ReadFile(pathToFile)
	if err != nil {
		return 0, err
	}

	count := 0
	var importJson []interface{}
	json.Unmarshal(body, &importJson)

	for _, mapJson := range importJson {
		m := mapJson.(map[string]interface{})
		idStr := fmt.Sprintf("%.f", m["id"].(float64))
		_, err := root.Insert(tokensearch.NewToken(
				idStr,
				m["label"].(string),
				m["category"].(string)))
		if err == nil {
			count++
		}
	}
	return count, nil
}