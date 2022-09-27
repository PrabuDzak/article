package idgenerator

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"regexp"
	"time"
)

type ShortUrlIdGenerator struct {
}

func NewShortUrlIdGenerator() *ShortUrlIdGenerator {
	return &ShortUrlIdGenerator{}
}

func (s *ShortUrlIdGenerator) Generate(ctx context.Context) (string, error) {
	now := time.Now().String()

	h := sha1.New()
	h.Write([]byte(now))

	regex, _ := regexp.Compile(`[^\w]`)
	res := base64.URLEncoding.EncodeToString(h.Sum(nil))
	res = regex.ReplaceAllString(res, "")

	return res[:6], nil
}
