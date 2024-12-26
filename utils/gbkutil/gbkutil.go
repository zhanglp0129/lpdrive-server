package gbkutil

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"strings"
)

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	return io.ReadAll(reader)
}

func StrToGbk(s string) ([]byte, error) {
	reader := transform.NewReader(strings.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	return io.ReadAll(reader)
}

func Utf8ReaderToGbkReader(s io.Reader) io.Reader {
	return transform.NewReader(s, simplifiedchinese.GBK.NewEncoder())
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	return io.ReadAll(reader)
}

func GbkToStr(s []byte) (string, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	res, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func GbkReaderToUtf8Reader(s io.Reader) io.Reader {
	return transform.NewReader(s, simplifiedchinese.GBK.NewDecoder())
}
