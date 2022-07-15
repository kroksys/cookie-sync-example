package service

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Server struct {
	LocalPort   string    `mapstructure:"localPort"`
	Host        Partner   `mapstructure:"host"`
	DSN         string    `mapstructure:"connectionString"`
	TablePrefix string    `mapstructure:"tablePrefix"`
	Partners    []Partner `mapstructure:"partners"`

	// Gin logger on/off
	ShowRequests bool `mapstructure:"showRequests"`
}

var iframeLayout = template.Must(template.New("t1").Parse(
	"<html><body>{{range .PixelUrls}}<img src='{{.}}' />{{end}}</body></html>",
))

const pixel = "\x47\x49\x46\x38\x39\x61\x01\x00\x01\x00\x80\x00\x00\x00\x00\x00\x00\x00\x00\x21\xF9\x04\x01\x00\x00\x00\x00\x2C\x00\x00\x00\x00\x01\x00\x01\x00\x00\x02\x02\x44\x01\x00\x3B"

type iframeData struct {
	PixelUrls []string
}

func (s *Server) Start() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	if s.ShowRequests {
		r.Use(gin.LoggerWithWriter(gin.DefaultWriter))
	}

	r.GET("/", s.homeHandler)
	r.GET("/iframe", s.iframeHandler)
	r.GET("/pixel", s.pixelHanlder)
	r.GET("/redirect", s.redirectHandler)
	r.GET("/user/:id", s.getUserInfo)

	fmt.Printf("Started hosting %s localy on :%s\n", s.Host.Domain, s.LocalPort)
	err := r.Run(fmt.Sprintf(":%s", s.LocalPort))
	if err != nil {
		log.Printf("Server stopped with error: %v\n", err)
	}
}

// Shows info about partner. Should not exist in production system.
func (s *Server) homeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, fmt.Sprintf("%+v", s.Host))
}

func (s *Server) iframeHandler(ctx *gin.Context) {
	hostUserId, err := ctx.Cookie(s.Host.CookieName)
	if err == http.ErrNoCookie {
		hostUserId = newUserId()
		s.setCookie(ctx.Writer, hostUserId)
	}
	iframeLayout.Execute(ctx.Writer, iframeData{
		PixelUrls: s.partnerUrls(hostUserId),
	})
}

// Notifies partners about this user
func (s *Server) pixelHanlder(ctx *gin.Context) {
	hostUserId, err := ctx.Cookie(s.Host.CookieName)
	if err == http.ErrNoCookie {
		hostUserId = newUserId()
		s.setCookie(ctx.Writer, hostUserId)
	}

	senderUserId := ctx.Query("user_id")
	senderDomain := ctx.Query("domain")
	senderRecirectUrl := ctx.Query("redirect")

	err = SaveMatchingTable(senderDomain, senderUserId, hostUserId)
	if err != nil {
		log.Printf("Error createing marging table record: %s\n", err.Error())
	}

	log.Printf("%s userID: %s | %s userID: %s", s.Host.Domain, hostUserId, senderDomain, senderUserId)

	if senderRecirectUrl != "" {
		log.Println("Redirecting back to", senderRecirectUrl)
		ctx.Redirect(http.StatusTemporaryRedirect,
			fmt.Sprintf("%s?user_id=%s&domain=%s", senderRecirectUrl, hostUserId, s.Host.Domain))
	}
}

// Receive partner user id back from pixelHandler
func (s *Server) redirectHandler(ctx *gin.Context) {
	hostUserId, err := ctx.Cookie(s.Host.CookieName)
	if err == http.ErrNoCookie {
		hostUserId = newUserId()
		s.setCookie(ctx.Writer, hostUserId)
	}

	senderUserId := ctx.Query("user_id")
	senderDomain := ctx.Query("domain")

	err = SaveMatchingTable(senderDomain, senderUserId, hostUserId)
	if err != nil {
		log.Printf("Error createing marging table record: %s\n", err.Error())
	}

	ctx.Writer.Header().Set("Content-Type", "image/gif")
	fmt.Fprint(ctx.Writer, pixel)
}

func (s *Server) getUserInfo(ctx *gin.Context) {
	mtRecord, err := ReadMatchingTableByLocalID(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, nil)
		return
	}
	ctx.JSON(http.StatusOK, mtRecord)
}

// Sets cookie to response writer
func (s *Server) setCookie(w gin.ResponseWriter, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     s.Host.CookieName,
		Value:    value,
		HttpOnly: false,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})
}

// Returns list of urls to be used for img pixels at iframe to sync cookies
func (s *Server) partnerUrls(hostUserId string) []string {
	urls := []string{}
	for _, p := range s.Partners {
		urls = append(urls, p.PixelUrl(hostUserId, s.Host, true))
	}
	return urls
}

// Just a wrapper to generate new user id.
func newUserId() string {
	return uuid.NewString()
}
