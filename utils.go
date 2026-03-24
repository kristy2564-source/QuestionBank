package main

import (
	"bytes"
	"io"
	"strconv"
	"strings"
)

func BuildComposeQuery(req ComposeRequest) string {
	var b strings.Builder
	b.WriteString("组卷需求：")
	if req.Title != "" {
		b.WriteString("标题=" + req.Title + "；")
	}
	if req.Subject != "" {
		b.WriteString("学科=" + req.Subject + "；")
	}
	if req.Grade != "" {
		b.WriteString("学段/年级=" + req.Grade + "；")
	}
	if req.Difficulty != "" {
		b.WriteString("难度=" + req.Difficulty + "；")
	}
	if len(req.Counts) > 0 {
		b.WriteString("题型数量=")
		var parts []string
		for k, v := range req.Counts {
			parts = append(parts, k+":"+strconv.Itoa(v))
		}
		b.WriteString(strings.Join(parts, ","))
		b.WriteString("；")
	}
	if len(req.Tags) > 0 {
		b.WriteString("标签=" + strings.Join(req.Tags, ",") + "；")
	}
	return b.String()
}

func ScanDoubleCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if i := bytes.Index(data, []byte("\r\n\r\n")); i >= 0 {
		return i + 4, data[0:i], nil
	}
	if atEOF && strings.Contains(string(data), "\"end\":true") {
		return len(data), data, nil
	}
	return 0, nil, nil
}

type bytesReader []byte

func (b bytesReader) Read(p []byte) (int, error) {
	n := copy(p, b)
	if n == len(b) {
		return n, io.EOF
	}
	return n, nil
}
