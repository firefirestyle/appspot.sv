package firestylesv

import (
	"net/http"
	//	"net/url"

	//	"errors"

	arthundler "github.com/firefirestyle/go.miniarticle/hundler"
	miniblob "github.com/firefirestyle/go.miniblob/blob"
	blobhandler "github.com/firefirestyle/go.miniblob/handler"
	"github.com/firefirestyle/go.minioauth/twitter"
	"github.com/firefirestyle/go.miniprop"
	"github.com/firefirestyle/go.minisession"
	//	"github.com/firefirestyle/go.miniuser"
	userhundler "github.com/firefirestyle/go.miniuser/hundler"
	//
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	//
	//"crypto/rand"
	//"encoding/binary"
	//"strconv"
	"errors"
	"io/ioutil"
	"strings"

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

	UrlBlobRequestUrl     = "/api/v1/blob/requesturl"
	UrlBlobCallback       = "/api/v1/blob/callback"
	UrlBlobGet            = "/api/v1/blob/get"
	UrlUserGet            = "/api/v1/user/get"
	UrlUserFind           = "/api/v1/user/find"
	UrlMeLogout           = "/api/v1/me/logout"
	UrlArtNew             = "/api/v1/art/new"
	UrlArtUpdate          = "/api/v1/art/update"
	UrlArtFind            = "/api/v1/art/find"
	UrlArtGet             = "/api/v1/art/get"
	UrlArtRequestBlobUrl  = "/api/v1/art/requestbloburl"
	UrlArtCallbackBlobUrl = "/api/v1/art/callbackbloburl"
)

var twitterHandlerObj *twitter.TwitterHandler = nil
var blobHandlerObj *blobhandler.BlobHandler = nil
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
		userHandlerObj = userhundler.NewUserHandler(userhundler.UserHandlerManagerConfig{
			ProjectId:   "firefirestyle",
			UserKind:    "user",
			RelayIdKind: "relayId",
			SessionKind: "session",
		}, userhundler.UserHandlerOnEvent{})
	}
	return userHandlerObj
}

func GetBlobHandlerObj(ctx context.Context) *blobhandler.BlobHandler {
	if blobHandlerObj == nil {
		blobHandlerObj = blobhandler.NewBlobHandler(
			UrlBlobCallback, appengine.VersionID(ctx), //
			miniblob.BlobManagerConfig{
				ProjectId:   "firefirestyle",
				Kind:        "blobstore",
				CallbackUrl: UrlBlobCallback,
			},
			blobhandler.BlobHandlerOnEvent{
				OnBlobRequest: func(w http.ResponseWriter, r *http.Request, outputProp *miniprop.MiniProp, blobHandlerObj *blobhandler.BlobHandler) (string, map[string]string, error) {
					//
					// login check
					bodyBytes, _ := ioutil.ReadAll(r.Body)
					propObj := miniprop.NewMiniPropFromJson(bodyBytes)
					token := propObj.GetString("token", "")
					ctx := appengine.NewContext(r)

					loginCheckInfo := GetUserHundlerObj(ctx).GetSessionMgr().CheckLoginId(ctx, token, minisession.MakeAccessTokenConfigFromRequest(r))
					if loginCheckInfo.IsLogin == false {
						return "", nil, errors.New("failed to wrong token : (1)")
					}
					//
					// path check
					dir := r.URL.Query().Get("dir")
					if true == strings.HasPrefix(dir, "/user") {
						Debug(ctx, ">>>>>> AAA <<<<<<"+loginCheckInfo.AccessTokenObj.GetUserName())
						if false == strings.HasPrefix(dir, "/user/"+loginCheckInfo.AccessTokenObj.GetUserName()) {
							return "", nil, errors.New("failed to wrong token : (2)")
						}
					} else {
						return "", nil, errors.New("failed to wrong token : (3)")
					}
					return loginCheckInfo.AccessTokenObj.GetLoginId(), map[string]string{}, nil
				},
				OnBlobComplete: func(w http.ResponseWriter, r *http.Request, outputProp *miniprop.MiniProp, blobHandlerObj *blobhandler.BlobHandler, blobObj *miniblob.BlobItem) error {
					dir := r.URL.Query().Get("dir")
					if true == strings.HasPrefix(dir, "/user") {
						ctx := appengine.NewContext(r)
						userName := strings.Replace(dir, "/user/", "", -1)
						userMgrObj := GetUserHundlerObj(ctx)
						userObj, userErr := userMgrObj.GetUserFromUserNameAndRelayId(ctx, userName)
						if userErr != nil {
							outputProp.SetString("error", "not found user")
							return userErr
						}
						userObj.SetIconUrl("key://" + blobObj.GetBlobKey())
						GetUserHundlerObj(ctx).SaveUserWithImmutable(ctx, userObj)
						return nil
					} else {
						return errors.New("unsupport")
					}
				},
			})
	}
	return blobHandlerObj
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
	// blob
	http.HandleFunc(UrlBlobRequestUrl, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetBlobHandlerObj(appengine.NewContext(r)).HandleBlobRequestToken(w, r)
	})

	http.HandleFunc(UrlBlobCallback, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetBlobHandlerObj(appengine.NewContext(r)).HandleUploaded(w, r)
	})

	http.HandleFunc(UrlBlobGet, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		GetBlobHandlerObj(appengine.NewContext(r)).HandleGet(w, r)
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

}

func Debug(ctx context.Context, message string) {
	log.Infof(ctx, message)
}
