package cjson

import (
	"encoding/json"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
)

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func RegisterFuzzyDecoders() {
	extra.RegisterFuzzyDecoders()
}

func Stringify(v any) string {
	if v == nil {
		return "null"
	}
	raw, _ := JSON.MarshalToString(v)
	return raw
}

func Prettify(v any) string {
	if v == nil {
		return "null"
	}
	// 因为 json-iterator/go 的缩进有 bug，没有继承 stream.indention 给 subStream.indention，
	// 所以这里用原生的 encoding/json 处理缩进
	bs, _ := json.MarshalIndent(v, "", "    ")
	return string(bs)
}
