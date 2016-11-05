package firestylesv

import (
	"net/http"

	//	"errors"

	//	"github.com/firefirestyle/go.miniarticle/article"
	arthundler "github.com/firefirestyle/go.miniarticle/hundler"
	//	blobhandler "github.com/firefirestyle/go.miniblob/handler"
	//	"github.com/firefirestyle/go.miniprop"
	"github.com/firefirestyle/go.minisession"
	userTmp "github.com/firefirestyle/go.miniuser/template"

	userhundler "github.com/firefirestyle/go.miniuser/handler"
	//
	"golang.org/x/net/context"
	//	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	//
	artTmp "github.com/firefirestyle/go.miniarticle/template"
	//"io/ioutil"
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

var sessionMgrObj *minisession.SessionManager = nil
var userHandlerObj *userhundler.UserHandler = nil
var artHandlerObj *arthundler.ArticleHandler = nil
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
