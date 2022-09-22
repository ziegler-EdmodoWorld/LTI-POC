package main

import (
	"flag"
	"fmt"
	"lti-plat/pkg"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := registerRoutes()

	flag.Parse()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8000),
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("start listen : %s\n", err)
	}
}

func registerRoutes() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.html")
	r.GET("certs", certs)
	r.GET("token", token)
	r.GET("auth", auth)
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})
	return r
}

///auth?client_id=clientid&login_hint=pirlo&nonce=fc33a239-9210-428c-af17-985df22612bf&prompt=none&redirect_uri=http%3A%2F%2Flocalhost%3A9000%2Flaunch&response_mode=form_post&response_type=id_token&scope=openid&state=state-2e03afc9-a4f5-41ac-a326-28810e35df7b
func auth(ctx *gin.Context) {
	scope := ctx.Query("scope")
	if scope != "openid" {
		ctx.Status(400)
		return
	}
	clientId := ctx.Query("client_id")
	username := ctx.Query("login_hint")
	nonce := ctx.Query("nonce")
	state := ctx.Query("state")
	redirectUri := ctx.Query("redirect_uri")
	resId := ctx.Query("lti_message_hint")

	ctx.HTML(http.StatusOK, "launch.html", gin.H{
		"RedirectUri": redirectUri,
		"Jwt":         pkg.IdToken(clientId, username, nonce, resId),
		"State":       state,
		"Iss":         pkg.ISSUER,
		"TargetLink":  redirectUri,
		"UserId":      username,
		"BookId":      resId,
	})
}

func token(ctx *gin.Context) {
	ctx.Status(200)
}

func certs(ctx *gin.Context) {
	ctx.Data(http.StatusOK, "text", []byte(pkg.PublicKey))
}
