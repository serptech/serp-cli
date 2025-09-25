package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/pretty"
)

func prettyPrintJSON(b []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "    ")
	return out.Bytes(), err
}

func GetPretty(e interface{}) ([]byte, error) {
	buf, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}
	out, err := prettyPrintJSON(buf)
	return out, err
}

func PrettyPrint(e interface{}) error {
	out, err := GetPretty(e)
	if err != nil {
		return err
	}
	result := pretty.Color(out, nil)
	fmt.Println(string(result))
	return nil
}
