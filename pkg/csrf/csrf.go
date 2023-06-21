package csrf

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"

	"github.com/dchest/uniuri"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	csrfSecret = "csrfSecret"
	csrfSalt   = "csrfSalt"
	csrfToken  = "csrfToken"
)

var defaultIgnoreMethods = []string{"GET", "HEAD", "OPTIONS"}

var defaultErrorFunc = func(c *gin.Context) {
	panic(errors.New("CSRF token mismatch"))
}

var defaultTokenGetter = func(c *gin.Context) string {
	r := c.Request

	if t := r.FormValue("_csrf"); len(t) > 0 {
		return t
	} else if t := r.URL.Query().Get("_csrf"); len(t) > 0 {
		return t
	} else if t := r.Header.Get("X-CSRF-TOKEN"); len(t) > 0 {
		return t
	} else if t := r.Header.Get("X-XSRF-TOKEN"); len(t) > 0 {
		return t
	}

	return ""
}

var defaultStatus = true

// Options stores configurations for a CSRF middleware.
type Options struct {
	Status        bool
	Secret        string
	IgnoreMethods []string
	ErrorFunc     gin.HandlerFunc
	TokenGetter   func(c *gin.Context) string
	IgnoreRoutes  []string
}

var csrfOption Options

func tokenize(secret, salt string) string {
	h := sha1.New()
	io.WriteString(h, salt+"-"+secret)
	hash := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return hash
}

func inArray(arr []string, value string) bool {
	inarr := false

	for _, v := range arr {
		if v == value {
			inarr = true
			break
		}
	}
	return inarr
}

// Middleware validates CSRF token.
func Middleware(options Options) gin.HandlerFunc {

	if options.IgnoreMethods == nil {
		options.IgnoreMethods = defaultIgnoreMethods
	}

	if options.ErrorFunc == nil {
		options.ErrorFunc = defaultErrorFunc
	}

	if options.TokenGetter == nil {
		options.TokenGetter = defaultTokenGetter
	}

	csrfOption = options

	return func(c *gin.Context) {
		session := sessions.Default(c)
		c.Set(csrfSecret, options.Secret)

		if csrfOption.Status {
			if inArray(options.IgnoreMethods, c.Request.Method) {
				c.Next()
				return
			}

			if inArray(options.IgnoreRoutes, c.Request.URL.Path) {
				c.Next()
				return
			}

			salt, ok := session.Get(csrfSalt).(string)

			if !ok || len(salt) == 0 {
				options.ErrorFunc(c)
				return
			}

			token := options.TokenGetter(c)
			if tokenize(options.Secret, salt) != token {
				options.ErrorFunc(c)
				return
			}
			c.Next()
		}
	}
}

// GetToken returns a CSRF token.
func VaildToken(c *gin.Context) {
	session := sessions.Default(c)
	c.Set(csrfSecret, csrfOption.Secret)

	salt, ok := session.Get(csrfSalt).(string)

	if !ok || len(salt) == 0 {
		csrfOption.ErrorFunc(c)
		return
	}

	token := csrfOption.TokenGetter(c)
	if tokenize(csrfOption.Secret, salt) != token {
		csrfOption.ErrorFunc(c)
		return
	}
	c.Next()
}

// GetToken returns a CSRF token.
func GetToken(c *gin.Context) string {
	session := sessions.Default(c)
	secret := c.MustGet(csrfSecret).(string)

	if t, ok := c.Get(csrfToken); ok {
		return t.(string)
	}

	salt, ok := session.Get(csrfSalt).(string)
	if !ok {
		salt = uniuri.New()
		session.Set(csrfSalt, salt)
		session.Save()
	}
	token := tokenize(secret, salt)
	c.Set(csrfToken, token)

	return token
}
