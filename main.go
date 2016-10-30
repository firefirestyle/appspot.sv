package firestylesv

import (
	"net/http"
	//	"net/url"

	//	"errors"

	arthundler "github.com/firefirestyle/go.miniarticle/hundler"
	//	miniblob "github.com/firefirestyle/go.miniblob/blob"
	//	blobhandler "github.com/firefirestyle/go.miniblob/handler"
	"github.com/firefirestyle/go.minioauth/twitter"
	"github.com/firefirestyle/go.miniprop"
	"github.com/firefirestyle/go.minisession"
	//	"github.com/firefirestyle/go.miniuser"
	userhundler "github.com/firefirestyle/go.miniuser/handler"
	//
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	//
	//"crypto/rand"
	//"encoding/binary"
	//"strconv"
	"io/ioutil"

	//	"google.golang.org/appengine/blobstore"
)

const (
	UrlTwitterTokenUrlRedirect_callbackUrl              = "cb"
	UrlTwitterTokenUrlRedirect_errorNotFoundCallbackUrl = "1001"
	UrlTwitterTokenUrlRedirect_errorFailedToMakeToken   = "1002"
	UrlTwitterTokenCallback_callbackUrl                 = "cb"
)

const (
	UrlTwitterTokenUrlRedirect = "/api/v1/twitter/tokenurl/redirect"
	UrlTwitterTokenCallback    = "/api/v1/twitter/tokenurl/callback"

	UrlBlobRequestUrl      = "/api/v1/blob/requesturl"
	UrlBlobCallback        = "/api/v1/blob/callback"
	UrlBlobGet             = "/api/v1/blob/get"
	UrlUserGet             = "/api/v1/user/get"
	UrlUserFind            = "/api/v1/user/find"
	UrlUserBlobGet         = "/api/v1/user/getblob"
	UrlUserRequestBlobUrl  = "/api/v1/user/requestbloburl"
	UrlUserCallbackBlobUrl = "/api/v1/user/callbackbloburl"
	UrlMeLogout            = "/api/v1/me/logout"
	UrlArtNew              = "/api/v1/art/new"
	UrlArtUpdate           = "/api/v1/art/update"
	UrlArtFind             = "/api/v1/art/find"
	UrlArtGet              = "/api/v1/art/get"
	UrlArtBlobGet          = "/api/v1/art/getblob"
	UrlArtRequestBlobUrl   = "/api/v1/art/requestbloburl"
	UrlArtCallbackBlobUrl  = "/api/v1/art/callbackbloburl"
)

var twitterHandlerObj *twitter.TwitterHandler = nil

//var blobHandlerObj *blobhandler.BlobHandler = nil
var sessionMgrObj *minisession.SessionManager = nil
var userHandlerObj *userhundler.UserHandler = nil
var artHandlerObj *arthundler.ArticleHandler = nil

func GetArtHundlerObj(ctx context.Context) *arthundler.ArticleHandler {
	if artHandlerObj == nil {
		artHandlerObj = arthundler.NewArtHandler(
			arthundler.ArticleHandlerManagerConfig{
				ProjectId:       "firefirestyle",
				ArticleKind:     "article",
				BlobKind:        "artblob",
				PointerKind:     "artpointer",
				BlobCallbackUrl: UrlArtCallbackBlobUrl,
				BlobSign:        appengine.VersionID(ctx),
			}, //
			arthundler.ArticleHandlerOnEvent{})
	}
	return artHandlerObj
}

func GetUserHundlerObj(ctx context.Context) *userhundler.UserHandler {
	if userHandlerObj == nil {
		userHandlerObj = userhundler.NewUserHandler(UrlUserCallbackBlobUrl,
			userhundler.UserHandlerManagerConfig{ //
				ProjectId:   "firefirestyle",
				UserKind:    "user",
				RelayIdKind: "relayId",
				SessionKind: "session",
			}, userhundler.UserHandlerOnEvent{})
	}
	return userHandlerObj
}

func GetTwitterHandlerObj(ctx context.Context) *twitter.TwitterHandler {
	if twitterHandlerObj == nil {
		twitterHandlerObj = GetUserHundlerObj(ctx).GetTwitterHandlerObj(ctx, //
			"http://"+appengine.DefaultVersionHostname(ctx)+""+UrlTwitterTokenCallback, twitter.TwitterOAuthConfig{
				ConsumerKey:       TwitterConsumerKey,
				ConsumerSecret:    TwitterConsumerSecret,
				AccessToken:       TwitterAccessToken,
				AccessTokenSecret: TwitterAccessTokenSecret,
			})
	}
	return twitterHandlerObj
}

func init() {
	initApi()
	initHomepage()
}

func initHomepage() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to FireFireStyle!!"))
	})
}

func initApi() {
	// twitter
	http.HandleFunc(UrlTwitterTokenUrlRedirect, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetTwitterHandlerObj(appengine.NewContext(r)).TwitterLoginEntry(w, r)
	})
	http.HandleFunc(UrlTwitterTokenCallback, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetTwitterHandlerObj(appengine.NewContext(r)).TwitterLoginExit(w, r)
	})

	// user
	http.HandleFunc(UrlUserGet, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetUserHundlerObj(appengine.NewContext(r)).HandleGet(w, r)
	})
	http.HandleFunc(UrlUserFind, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetUserHundlerObj(appengine.NewContext(r)).HandleFind(w, r)
	})
	http.HandleFunc(UrlUserRequestBlobUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetUserHundlerObj(appengine.NewContext(r)).HandleBlobRequestToken(w, r)
	})
	http.HandleFunc(UrlUserCallbackBlobUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetUserHundlerObj(appengine.NewContext(r)).HandleBlobUpdated(w, r)
	})
	http.HandleFunc(UrlUserBlobGet, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetUserHundlerObj(appengine.NewContext(r)).HandleBlobGet(w, r)
	})

	// me
	http.HandleFunc(UrlMeLogout, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		bodyBytes, _ := ioutil.ReadAll(r.Body)
		propObj := miniprop.NewMiniPropFromJson(bodyBytes)
		token := propObj.GetString("token", "")
		ctx := appengine.NewContext(r)
		GetUserHundlerObj(ctx).GetSessionMgr().Logout(ctx, token, minisession.MakeAccessTokenConfigFromRequest(r))
	})

	// art
	// UrlArtNew
	http.HandleFunc(UrlArtNew, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		GetArtHundlerObj(ctx).HandleNew(w, r)
	})

	// art
	// UrlArtNew
	http.HandleFunc(UrlArtUpdate, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		GetArtHundlerObj(ctx).HandleUpdate(w, r)
	})

	http.HandleFunc(UrlArtFind, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		GetArtHundlerObj(ctx).HandleFind(w, r)
	})

	http.HandleFunc(UrlArtGet, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		GetArtHundlerObj(ctx).HandleGet(w, r)
	})
	//UrlArtGet

	http.HandleFunc(UrlArtRequestBlobUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		GetArtHundlerObj(ctx).HandleBlobRequestToken(w, r)
	})

	http.HandleFunc(UrlArtCallbackBlobUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		Debug(ctx, "asdfasdfasdf")

		GetArtHundlerObj(ctx).HandleBlobUpdated(w, r)
	})

	http.HandleFunc(UrlArtBlobGet, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		ctx := appengine.NewContext(r)
		Debug(ctx, "asdfasdfasdf")

		GetArtHundlerObj(ctx).HandleBlobGet(w, r)
	})

}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
