package utilCaptcha

import (
	"encoding/json"
	"fmt"
	"github.com/hilaoyu/go-utils/utilCache"
	"github.com/hilaoyu/go-utils/utilRandom"
	"github.com/hilaoyu/go-utils/utilUuid"
	"github.com/mojocn/base64Captcha"
	"image/color"
	"strings"
	"time"
)

type CaptchaCode struct {
	Id     string        `json:"id,omitempty"`
	Val    string        `json:"val,omitempty"`
	Image  string        `json:"image,omitempty"`
	SendTo string        `json:"send_to,omitempty"`
	Expire time.Duration `json:"expire,omitempty"`
}

type CaptchaService struct {
	cache        *utilCache.Cache
	expire       time.Duration
	cachePrefix  string
	fontsStorage base64Captcha.FontsStorage
	fonts        []string
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
	captchaService.fonts = []string{"RitaSmith.ttf", "wqy-microhei.ttc"}
	return captchaService
}

func (c *CaptchaService) SetFontsStorage(storage base64Captcha.FontsStorage) *CaptchaService {
	c.fontsStorage = storage
	return c
}
func (c *CaptchaService) SetFonts(fonts []string) *CaptchaService {
	c.fonts = fonts
	return c
}
func (c *CaptchaService) Verify(id string, val string, sendTo string, clear bool) (err error) {
	id = strings.TrimSpace(id)
	val = strings.TrimSpace(val)
	sendTo = strings.TrimSpace(sendTo)
	if "" == id || "" == val {
		err = fmt.Errorf("id和val不能为空")
		return
	}
	code, err := c.cacheGetCode(id)

	if nil != err {
		return
	}

	if clear {
		c.cacheDeleteCode(id)
	}

	if "" != sendTo && code.SendTo != sendTo {
		err = fmt.Errorf("验证码错误.")
		return
	}

	if code.Val == val {
		err = nil
		return
	} else {
		err = fmt.Errorf("验证码错误")
	}
	//err = fmt.Errorf("验证码错误")
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
	imageDraw := base64Captcha.NewDriverString(height, width, 30, base64Captcha.OptionShowHollowLine, 0, codeVal, &color.RGBA{R: 96, G: 96, B: 96, A: 128}, c.fontsStorage, c.fonts)

	image, err := imageDraw.DrawCaptcha(codeVal)
	if nil != err {
		return
	}

	code, err = c.cacheSaveCode(codeVal, "", ttl...)
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
	mathDraw := base64Captcha.NewDriverMath(height, width, 30, base64Captcha.OptionShowHollowLine, &color.RGBA{R: 96, G: 96, B: 96, A: 128}, c.fontsStorage, c.fonts)

	_, question, codeVal := mathDraw.GenerateIdQuestionAnswer()
	image, err := mathDraw.DrawCaptcha(question)
	if nil != err {
		return
	}

	code, err = c.cacheSaveCode(codeVal, "", ttl...)
	if nil != err {
		return
	}

	code.Image = image.EncodeB64string()
	return
}

func (c *CaptchaService) buildCacheKey(codeId string) string {
	return c.cachePrefix + codeId
}

func (c *CaptchaService) cacheSaveCode(codeVal string, sendTo string, ttl ...time.Duration) (code *CaptchaCode, err error) {
	code = &CaptchaCode{
		Id:     utilUuid.UuidGenerate(),
		SendTo: sendTo,
		Val:    codeVal,
	}
	expire := c.expire
	if len(ttl) > 0 && ttl[0] > 0 {
		expire = ttl[0]
	}
	code.Expire = expire

	codeByte, err := json.Marshal(code)
	if nil != err {
		code = nil
		return
	}
	err = c.cache.Set(c.buildCacheKey(code.Id), codeByte, expire)
	if nil != err {
		code = nil
	}
	return
}
func (c *CaptchaService) cacheGetCode(id string) (code *CaptchaCode, err error) {
	id = strings.TrimSpace(id)
	if "" == id {
		err = fmt.Errorf("id不能为空")
		return
	}
	cacheCode := c.cache.Get(c.buildCacheKey(id))

	if nil == cacheCode {
		err = fmt.Errorf("验证码已过期")
		return
	}

	cacheCodeByte, ok := cacheCode.([]byte)
	if !ok {
		err = fmt.Errorf("验证码解析失败")
		return
	}

	code = &CaptchaCode{}
	err = json.Unmarshal(cacheCodeByte, code)
	return
}
func (c *CaptchaService) cacheDeleteCode(id string) {
	c.cache.Del(c.buildCacheKey(id))
}
