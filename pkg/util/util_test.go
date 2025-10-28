package util

import (
	"cambridge-hit.com/gin-base/activateserver/pkg/util/validate"
	"testing"
)

func TestEmail(t *testing.T) {
	// 测试函数
	emails := []string{
		"example@example.com",
		"invalid-email",
		"user.name+tag+sorting@example.com",
		"username@.com",
		"username@model.c",
		"username@model.co",
	}

	for _, email := range emails {
		if validate.Email(email) {
			println(email, "is a validate email address.")
		} else {
			println(email, "is not a validate email address.")
		}
	}
}
