package str

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
)

// 返回一个32位md5加密后的字符串
func Md532(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func SHA1(s string) string {
	o := sha1.New()
	o.Write([]byte(s))
	return hex.EncodeToString(o.Sum(nil))
}

// 返回一个32位md5加密后的字符串,sha1
func Md5_32_sha1(str string) string {
	h := md5.New()
	h.Write([]byte(SHA1(str)))
	s := hex.EncodeToString(h.Sum(nil))
	return s
}

// 返回一个16位md5加密后字符串
func Md5To16(str string) string {
	return Md532(str)[8:24]
}

// GetFileMD5 计算文件的 MD5 值
// @param: path string文件路径 /data/file1.png
func GetFileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := md5.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	md5sum := hex.EncodeToString(hasher.Sum(nil))

	return md5sum, nil
}
