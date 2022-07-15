package service

import (
	"fmt"
)

// We could keep partners in database and wrap it arount with REST
// but for the demo sake I kept the partners in config file.
type Partner struct {
	// ID         int    `gorm:"primaryKey" mapstructure:"-"`
	Domain     string `mapstructure:"domain"`
	Redirect   string `mapstructure:"redirect"`
	Pixel      string `mapstructure:"pixel"`
	CookieName string `mapstructure:"cookieName"`
}

func (p *Partner) PixelUrl(userID string, caller Partner, redirect bool) string {
	url := fmt.Sprintf("%s%s?user_id=%s&domain=%s", p.Domain, p.Pixel, userID, caller.Domain)
	if redirect {
		url += "&redirect=" + caller.RedirectUrl()
	}
	return url
}

func (p *Partner) RedirectUrl() string {
	return fmt.Sprintf("%s%s", p.Domain, p.Redirect)
}
