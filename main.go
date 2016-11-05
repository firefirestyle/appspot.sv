package firestylesv

import (
	"net/http"

	userTmp "github.com/firefirestyle/go.miniuser/template"

	artTmp "github.com/firefirestyle/go.miniarticle/template"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

/*
const (
	UrlTwitterTokenUrlRedirect_callbackUrl              = "cb"
	UrlTwitterTokenUrlRedirect_errorNotFoundCallbackUrl = "1001"
	UrlTwitterTokenUrlRedirect_errorFailedToMakeToken   = "1002"
	UrlTwitterTokenCallback_callbackUrl                 = "cb"
)*/

const (

	//
	UrlArtNew             = "/api/v1/art/new"
	UrlArtUpdate          = "/api/v1/art/update"
	UrlArtFind            = "/api/v1/art/find"
	UrlArtGet             = "/api/v1/art/get"
	UrlArtBlobGet         = "/api/v1/art/getblob"
	UrlArtRequestBlobUrl  = "/api/v1/art/requestbloburl"
	UrlArtCallbackBlobUrl = "/api/v1/art/callbackbloburl"

	// blob
	UrlBlobRequestUrl = "/api/v1/blob/requesturl"
	UrlBlobCallback   = "/api/v1/blob/callback"
	UrlBlobGet        = "/api/v1/blob/get"
)

var userTemplate = userTmp.NewUserTemplate(userConfig)
var artTemplate *artTmp.ArtTemplate = artTmp.NewArtTemplate(artTmp.ArtTemplateConfig{
	GroupName:    "Main",
	KindBaseName: "FFArt",
}, userTemplate.GetUserHundlerObj)

func init() {
	initHomepage()
	userTemplate.InitUserApi()
	artTemplate.InitArtApi()
}

func initHomepage() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to FireFireStyle!!"))
	})
}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
