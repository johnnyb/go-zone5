package zone5

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

func mustJsonBody(data map[string]interface{}) io.ReadCloser {
	bindata, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	reader := bytes.NewReader(bindata)
	return ioutil.NopCloser(reader)
}

func unmarshalJsonFromReader(reader io.Reader) (map[string]interface{}, error) {
	decoder := json.NewDecoder(reader)
	result := map[string]interface{}{}
	err := decoder.Decode(&result)
	return result, err
}

func mustUnmarshalJsonFromReader(reader io.Reader) (map[string]interface{}) {
	data, err := unmarshalJsonFromReader(reader)
	if err != nil {
		panic(err)
	}
	return data
}
