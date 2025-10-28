package util

import (
	"cambridge-hit.com/gin-base/activateserver/resource/embed"
	"errors"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"strconv"
)

var I18nBundle *i18n.Bundle
var Localizer map[string]*i18n.Localizer

var Locales = []string{"en", "zh-hans"} //供localizer初始化使用
var GlobalLocale = Locales[0]

var LanguageMap = map[int]string{
	0: "en",
	1: "zh-hans",
}

func init() {
	I18nBundle = i18n.NewBundle(language.English)
	I18nBundle.RegisterUnmarshalFunc("json", jsoniter.Unmarshal)
	Localizer = make(map[string]*i18n.Localizer)
	for _, v := range Locales {
		//if _, err := I18nBundle.LoadMessageFile("locales/active." + v + ".json"); err != nil {
		if _, err := I18nBundle.LoadMessageFileFS(embed.FsLocales, "locales/"+v+".json"); err != nil {
			// 在init阶段不记录错误，避免logger未初始化的问题
			// logger.Error("loading locale file error", zap.String("filename", "locales/"+v+".json"))
			continue
		}
		Localizer[v] = i18n.NewLocalizer(I18nBundle, v)
	}
	
	// 确保至少有一个本地化器
	if len(Localizer) == 0 {
		// 如果没有加载任何本地化文件，创建一个默认的本地化器
		Localizer["en"] = i18n.NewLocalizer(I18nBundle, "en")
		GlobalLocale = "en"
	}
}

func T(messageID string, locale ...string) string {
	localizer := Localizer[GlobalLocale] // 默认使用全局locale

	if len(locale) > 0 && Localizer[locale[0]] != nil {
		localizer = Localizer[locale[0]] // 如果指定的locale有效，则使用之
	}

	// 如果localizer为nil，直接返回消息ID
	if localizer == nil {
		return messageID
	}

	if localizedMessage, err := localizer.Localize(&i18n.LocalizeConfig{MessageID: messageID}); err == nil {
		return localizedMessage
	}
	// 未找到本地化消息，返回消息ID作为回退
	return messageID
}

// Translate 必须指定locale
func Translate(messageID string, locale string) string {
	return T(messageID, locale)
}

func GetMessageWithTemplateData(locale string, messageID string, templateData map[string]interface{}) (string, error) {
	localizer := Localizer[locale]
	if localizer != nil {
		return localizer.Localize(&i18n.LocalizeConfig{
			MessageID:    messageID,
			TemplateData: templateData,
		})
	}
	return "", errors.New("cannot find the target locale file")
}

func Unmarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}

func GetLanguage(context *gin.Context) string {
	//_language := request.GetJsonParamStr(context, "Language")
	_language := context.GetHeader("Language")
	languageNum, err := strconv.Atoi(_language)
	if err != nil || languageNum >= len(LanguageMap) || languageNum < 0 {
		_language = GlobalLocale
	} else {
		_language = LanguageMap[languageNum]
	}
	return _language
}
