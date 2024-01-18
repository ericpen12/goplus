package http

import (
	"bytes"
	"encoding/json"
	"io"
)

type Parser interface {
	Obj() any
}

type Result interface {
	Buffer() *bytes.Buffer
	Bind(data interface{}) error
	UnescapeString() Result
}

func Read(body io.Reader, format any) (Result, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	_, err := buf.ReadFrom(body)
	if err != nil {
		return nil, err
	}
	return &result{
		originBuf: buf,
		buf:       getBuffer(buf, format),
	}, nil
}

func getBuffer(r *bytes.Buffer, format any) *bytes.Buffer {
	if format == nil {
		return r
	}
	_ = json.Unmarshal(r.Bytes(), format)
	var b []byte
	if p, ok := format.(Parser); ok {
		b, _ = json.Marshal(p.Obj())
	} else {
		b, _ = json.Marshal(format)
	}
	return bytes.NewBuffer(b)
}

type result struct {
	originBuf, buf *bytes.Buffer
}

func (r *result) Buffer() *bytes.Buffer {
	return r.buf
}

func (r *result) String() string {
	return r.buf.String()
}

func (r *result) Bind(data interface{}) error {
	return json.Unmarshal(r.buf.Bytes(), data)
}

func (r *result) UnescapeString() Result {
	m := make(map[string]interface{})
	_ = json.Unmarshal([]byte(UnescapeString(r.buf.Bytes())), &m)
	b, _ := json.Marshal(m)
	r.buf = bytes.NewBuffer(b)
	return r
}

func UnescapeString(input []byte) string {
	var s string
	_ = json.Unmarshal(input, &s)
	return s
}
