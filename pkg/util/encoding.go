package util

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/traditionalchinese"
	xunicode "golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

/**
	文件编码重写
**/

func ConvertToUTF8(data []byte) ([]byte, error) {
	// 判断data是否是有效的UTF-8编码
	if !utf8.Valid(data) {
		// 检测编码格式
		sourceEncoding, err := detectEncoding(data)
		if err != nil {
			return nil, err
		}

		// 获取编码格式对应的编码器
		encoder := getEncoder(sourceEncoding)
		if encoder == nil {
			return nil, fmt.Errorf("encoder not found for encoding: %s", sourceEncoding)
		}

		// 使用编码器将数据转换为 UTF-8 编码
		utf8Data, err := io.ReadAll(transform.NewReader(bytes.NewReader(data), encoder.NewDecoder()))
		if err != nil {
			return nil, err
		}

		return utf8Data, nil
	}
	// 如果data已经是有效的UTF-8编码，则直接返回
	return data, nil
}

func detectEncoding(data []byte) (string, error) {
	detector := chardet.NewTextDetector()
	result, err := detector.DetectBest(data)
	if err != nil {
		return "", err
	}
	return result.Charset, nil
}

func getEncoder(encodingName string) encoding.Encoding {
	switch encodingName {
	case "UTF-8":
		return encoding.Nop
	case "GBK":
		return simplifiedchinese.GBK
	case "GB-18030":
		return simplifiedchinese.GB18030
	case "HZ-GB-2312":
		return simplifiedchinese.HZGB2312
	case "Big5":
		return traditionalchinese.Big5
	case "UTF-16LE":
		return xunicode.UTF16(xunicode.LittleEndian, xunicode.IgnoreBOM)
	case "UTF-16BE":
		return xunicode.UTF16(xunicode.BigEndian, xunicode.IgnoreBOM)
	case "EUC-JP":
		return japanese.EUCJP
	case "Shift_JIS":
		return japanese.ShiftJIS
	case "ISO-2022-JP":
		return japanese.ISO2022JP
	case "ISO-8859-1":
		return charmap.ISO8859_1
	case "ISO-8859-2":
		return charmap.ISO8859_2
	case "ISO-8859-3":
		return charmap.ISO8859_3
	case "ISO-8859-4":
		return charmap.ISO8859_4
	case "ISO-8859-5":
		return charmap.ISO8859_5
	case "ISO-8859-6":
		return charmap.ISO8859_6
	case "ISO-8859-7":
		return charmap.ISO8859_7
	case "ISO-8859-8":
		return charmap.ISO8859_8
	case "ISO-8859-9":
		return charmap.ISO8859_9
	case "ISO-8859-10":
		return charmap.ISO8859_10
	case "ISO-8859-13":
		return charmap.ISO8859_13
	case "ISO-8859-14":
		return charmap.ISO8859_14
	case "ISO-8859-15":
		return charmap.ISO8859_15
	case "ISO-8859-16":
		return charmap.ISO8859_16
	case "Windows-1250":
		return charmap.Windows1250
	case "Windows-1251":
		return charmap.Windows1251
	case "Windows-1252":
		return charmap.Windows1252
	case "Windows-1253":
		return charmap.Windows1253
	case "Windows-1254":
		return charmap.Windows1254
	case "Windows-1255":
		return charmap.Windows1255
	case "Windows-1256":
		return charmap.Windows1256
	case "Windows-1257":
		return charmap.Windows1257
	case "Windows-1258":
		return charmap.Windows1258
	case "KOI8-R":
		return charmap.KOI8R
	case "KOI8-U":
		return charmap.KOI8U
	default:
		return nil
	}
}
