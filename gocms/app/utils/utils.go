package utils

import (
	"log"

	"github.com/davidddw/gopj/gocms/app/conf"
	"github.com/gin-gonic/gin"

	ut "github.com/go-playground/universal-translator"
	"github.com/noxue/locales/zh"
	zh_translations "github.com/noxue/validator/translations/zh"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"gopkg.in/go-playground/validator.v9"
)

func init() {
	initValidator()
	initQiniu()
}

var Validate *validator.Validate
var Trans ut.Translator

func initValidator() {
	Validate = validator.New()

	zh1 := zh.New()
	uni := ut.New(zh1, zh1)
	Trans, _ = uni.GetTranslator("zh")

	zh_translations.RegisterDefaultTranslations(Validate, Trans)
}

func ValidateStruct(t interface{}) (errs []gin.H) {
	err := Validate.Struct(t)
	if err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			log.Println(err)
			return
		}
		for _, err := range err.(validator.ValidationErrors) {
			errs = append(errs, gin.H{err.Field(): err.Translate(Trans)})
		}

		return
	}
	return
}

var QiNiuMac *qbox.Mac
var BucketManager *storage.BucketManager

func initQiniu() {
	QiNiuMac = qbox.NewMac(conf.Conf.Upload.QiNiu.Id, conf.Conf.Upload.QiNiu.Key)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	BucketManager = storage.NewBucketManager(QiNiuMac, &cfg)
}
