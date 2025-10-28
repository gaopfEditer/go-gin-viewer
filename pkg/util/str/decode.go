package str

import (
	"bytes"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"os"
)

func ConvertGB2312INItoUTF8(iniFile string) ([]byte, error) {
	content, err := os.ReadFile(iniFile)
	if err != nil {
		return nil, err
	}

	// 将 GB2312 编码的内容转换为 UTF-8
	reader := transform.NewReader(bytes.NewReader(content), simplifiedchinese.GB18030.NewDecoder())
	utf8Content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return utf8Content, nil
}
