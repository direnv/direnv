package main

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func marshal(obj interface{}) string {
	jsonData, err := json.Marshal(obj)

	if err != nil {
		panic(fmt.Errorf("marshal(): %v", err))
	}

	zlibData := bytes.NewBuffer([]byte{})
	w := zlib.NewWriter(zlibData)
	w.Write(jsonData)
	w.Close()

	base64Data := base64.URLEncoding.EncodeToString(zlibData.Bytes())

	return base64Data
}

func unmarshal(base64env string, obj interface{}) error {
	base64env = strings.TrimSpace(base64env)

	data, err := base64.URLEncoding.DecodeString(base64env)
	if err != nil {
		return fmt.Errorf("unmarshal() base64 decoding: %v", err)
	}

	zlibReader := bytes.NewReader(data)
	w, err := zlib.NewReader(zlibReader)
	if err != nil {
		return fmt.Errorf("unmarshal() zlib opening: %v", err)
	}

	envData := bytes.NewBuffer([]byte{})
	_, err = io.Copy(envData, w)
	if err != nil {
		return fmt.Errorf("unmarshal() zlib decoding: %v", err)
	}
	w.Close()

	err = json.Unmarshal(envData.Bytes(), &obj)
	if err != nil {
		return fmt.Errorf("unmarshal() json parsing: %v", err)
	}

	return nil
}
