package utilCaptcha

import (
	"fmt"
	"github.com/hilaoyu/go-utils/utilCache"
	"github.com/hilaoyu/go-utils/utilRandom"
	"github.com/hilaoyu/go-utils/utilUuid"
	"github.com/mojocn/base64Captcha"
	"image/color"
	"time"
)

type CaptchaCode struct {
	Id     string
	Val    string
	Image  string
	Expire time.Duration
}

type CaptchaService struct {
	cache       *utilCache.Cache
	expire      time.Duration
	cachePrefix string
}

func NewCaptchaService(cache *utilCache.Cache, prefix string, ttl ...time.Duration) *CaptchaService {
	captchaService := &CaptchaService{
		cache:       cache,
		expire:      0,
		cachePrefix: prefix,
	}
	if len(ttl) > 0 && ttl[0] > 0 {
		captchaService.expire = ttl[0]
	}
	return captchaService
}

func (c *CaptchaService) Verify(id string, val string, clear bool) (err error) {
	cacheVal := c.cache.Get(c.buildCacheKey(id))

	if nil == cacheVal {
		err = fmt.Errorf("验证码已过期")
		return
	}

	if clear {
		c.cache.Del(c.buildCacheKey(id))
	}

	if codeVal, ok := cacheVal.(string); ok {
		if codeVal == val {
			return nil
		}
	}
	err = fmt.Errorf("验证码错误")
	return
}

func (c *CaptchaService) GenerateImageString(width int, height int, codeLen int, ttl ...time.Duration) (code *CaptchaCode, err error) {
	codeVal := utilRandom.RandString(codeLen)
	if width <= 0 {
		width = 100
	}
	if width > 300 {
		width = 300
	}
	if height <= 0 {
		height = 36
	}
	if height > 300 {
		height = 300
	}
	imageDraw := base64Captcha.NewDriverString(height, width, 30, base64Captcha.OptionShowHollowLine, 0, codeVal, &color.RGBA{R: 0, G: 0, B: 0, A: 60}, nil, []string{"RitaSmith.ttf", "chromohv.ttf", "wqy-microhei.ttc"})

	image, err := imageDraw.DrawCaptcha(codeVal)
	if nil != err {
		return
	}

	code, err = c.saveCode(codeVal, ttl...)
	if nil != err {
		return
	}

	code.Image = image.EncodeB64string()
	return
}
func (c *CaptchaService) GenerateImageMath(width int, height int, ttl ...time.Duration) (code *CaptchaCode, err error) {
	if width <= 0 {
		width = 100
	}
	if width > 300 {
		width = 300
	}
	if height <= 0 {
		height = 36
	}
	if height > 300 {
		height = 300
	}
	mathDraw := base64Captcha.NewDriverMath(height, width, 30, base64Captcha.OptionShowHollowLine, &color.RGBA{R: 0, G: 0, B: 0, A: 60}, nil, []string{"RitaSmith.ttf", "chromohv.ttf", "wqy-microhei.ttc"})

	_, question, codeVal := mathDraw.GenerateIdQuestionAnswer()
	image, err := mathDraw.DrawCaptcha(question)
	if nil != err {
		return
	}

	code, err = c.saveCode(codeVal, ttl...)
	if nil != err {
		return
	}

	code.Image = image.EncodeB64string()
	return
}

func (c *CaptchaService) buildCacheKey(codeId string) string {
	return c.cachePrefix + codeId
}

func (c *CaptchaService) saveCode(codeVal string, ttl ...time.Duration) (code *CaptchaCode, err error) {
	code = &CaptchaCode{
		Id:  utilUuid.UuidGenerate(),
		Val: codeVal,
	}
	expire := c.expire
	if len(ttl) > 0 && ttl[0] > 0 {
		expire = ttl[0]
	}
	code.Expire = expire
	err = c.cache.Set(c.buildCacheKey(code.Id), code.Val, expire)
	if nil != err {
		code = nil
	}
	return
}
