package app

import (
	"crypto/aes"
	"crypto/cipher"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

func AuthCookieMiddleware(cf Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie("user-id")
		if cookie == "" || err != nil {
			encID, err := Encrypt(uuid.NewV4().String(), cf.AppKey)
			if err != nil {
				c.String(http.StatusInternalServerError, "")
				return
			}
			c.SetCookie("user-id", encID, 3600, "/", cf.Address, false, false)
		}
		c.Next()
	}
}

func Encrypt(text string, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", nil
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", nil
	}

	nonce := make([]byte, gcm.NonceSize())
	enc := gcm.Seal(nonce, nonce, []byte(text), nil)

	return string(enc), nil
}

func Decrypt(text []byte, key string) (string, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(text) < nonceSize {
		return "", err
	}

	nonce, ciphertext := text[:nonceSize], text[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
