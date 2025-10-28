package str

import (
	"math/rand"
	"time"

	"cambridge-hit.com/gin-base/activateserver/pkg/util/cache"
)

var randSeedMap = make(map[string]bool)

// 生成随机的 num 位大小写字母
// param: num生成code的位数
func GenerateCode(num int) string {
	code := ""
	letters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

	if _, ok := randSeedMap["randSeed"]; !ok {
		rand.Seed(time.Now().UnixNano())
		randSeedMap["randSeed"] = true
	}

	for i := 0; i < num; i++ {
		randomIndex := rand.Intn(len(letters))
		code += string(letters[randomIndex])
	}

	return code
}

func GenerateRandomNumberString(length int) string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano())) // 使用独立的随机数生成器
	digits := "0123456789"
	var result string

	for i := 0; i < length; i++ {
		randomIndex := rng.Intn(len(digits))  // 生成一个随机索引
		result += string(digits[randomIndex]) // 拼接随机数字字符
	}

	return result
}

func GetEmailActivateKey(code string) string {
	return "Email_Activate:" + code
}

func GetCaptchaExpireKey(captchaID string) string {
	return "Captcha_Expire:" + captchaID
}

// 成功验证邮箱后设置一个redis key
func SetEmailActivateSuccess(code string, email string) {
	cache.MyRedis.Set("Email_Activated_"+code, email, time.Minute*30)
}

// 调用注册或其他需要验证邮箱激活成功时使用此函数获取对应的key以验证是否成功激活过邮箱
func GetEmailActivateSuccessKey(code string) string {
	return "Email_Activated_" + code
}
