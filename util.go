package main

import (
	json "github.com/json-iterator/go"
	"github.com/rs/xid"
)

func UUID() string {
	return xid.New().String()
}

var (
	jsoniter = json.ConfigCompatibleWithStandardLibrary
)

func Encode(data interface{}) ([]byte, error) {
	return jsoniter.Marshal(data)
}

func EncodeString(data interface{}) string {
	if binData, err := jsoniter.Marshal(data); err == nil {
		return string(binData)
	}
	return ""
}

func Decode(str string, data interface{}) error {
	return jsoniter.UnmarshalFromString(str, data)
}

func DecodeByte(b []byte, data interface{}) error {
	return jsoniter.Unmarshal(b, data)
}
