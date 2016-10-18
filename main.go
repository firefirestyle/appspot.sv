package firestylesv

import (
	"net/http"
	//	"net/url"

	//	"errors"

	"github.com/firefirestyle/go.miniblob"
	"github.com/firefirestyle/go.minioauth/twitter"
	"github.com/firefirestyle/go.minisession"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	//
	"crypto/rand"
	"encoding/binary"
	"strconv"
)

const (
	UrlApiRoot = "/api/v1"
)

const (
	UrlTwitterTokenUrlRedirect                          = "twitter/tokenurl/redirect"
	UrlTwitterTokenUrlRedirect_callbackUrl              = "cb"
	UrlTwitterTokenUrlRedirect_errorNotFoundCallbackUrl = "1001"
	UrlTwitterTokenUrlRedirect_errorFailedToMakeToken   = "1002"
	UrlTwitterTokenCallback                             = "twitter/tokenurl/callback"
	UrlTwitterTokenCallback_callbackUrl                 = "cb"
)

const (
	UrlBlobRequestUrl = "blob/requesturl"
	UrlBlobCallback   = "blob/callback"
)

var twitterHandlerObj *twitter.TwitterHandler = nil
var blobHandlerObj *miniblob.BlobHandler = nil
var sessionMgrObj *minisession.SessionManager = nil

func GetSessionMgrObj(ctx context.Context) *minisession.SessionManager {
	if sessionMgrObj == nil {
		sessionMgrObj = minisession.NewSessionManager(minisession.SessionManagerConfig{
			ProjectId: "firefirestyle",
			Kind:      "session",
		})
	}
	return sessionMgrObj
}

func GetBlobHandlerObj(ctx context.Context) *miniblob.BlobHandler {
	if blobHandlerObj == nil {
		blobHandlerObj = miniblob.NewBlobHandler(
			UrlApiRoot+"/"+UrlBlobCallback, appengine.VersionID(ctx), //
			miniblob.BlobManagerConfig{
				ProjectId:   "firefirestyle",
				Kind:        "blobstore",
				CallbackUrl: UrlBlobCallback,
			},
			miniblob.BlobHandlerOnEvent{})
	}
	return blobHandlerObj
}

func GetTwitterHandlerObj(ctx context.Context) *twitter.TwitterHandler {
	if twitterHandlerObj == nil {
		twitterHandlerObj = twitter.NewTwitterHandler( //
			"http://"+appengine.DefaultVersionHostname(ctx)+""+UrlApiRoot+"/"+UrlTwitterTokenCallback, twitter.TwitterOAuthConfig{
				ConsumerKey:       TwitterConsumerKey,
				ConsumerSecret:    TwitterConsumerSecret,
				AccessToken:       TwitterAccessToken,
				AccessTokenSecret: TwitterAccessTokenSecret}, twitter.TwitterHundlerOnEvent{
				OnFoundUser: func(w http.ResponseWriter, r *http.Request, handler *twitter.TwitterHandler, accesssToken *twitter.SendAccessTokenResult) map[string]string {
					ctx := appengine.NewContext(r)
					sessionMgrObj := GetSessionMgrObj(ctx)
					tokenObj, err := sessionMgrObj.Login(ctx, twitter.ScreenName, minisession.MakeAccessTokenConfigFromRequest(r))
					if err != nil {
						return map[string]string{"errcode": "1"}
					} else {
						return map[string]string{"token": "" + tokenObj.GetLoginId()}
					}
				},
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
		log.Infof(appengine.NewContext(r), ">> "+makeRandomId())
		w.Write([]byte("Welcome to FireFireStyle!!"))
	})
}

func initApi() {
	http.HandleFunc(UrlApiRoot+"/"+UrlTwitterTokenUrlRedirect, func(w http.ResponseWriter, r *http.Request) {
		//	w.Header().Add("Access-Control-Allow-Origin", "*")
		GetTwitterHandlerObj(appengine.NewContext(r)).TwitterLoginEntry(w, r)
	})
	http.HandleFunc(UrlApiRoot+"/"+UrlTwitterTokenCallback, func(w http.ResponseWriter, r *http.Request) {
		//	w.Header().Add("Access-Control-Allow-Origin", "*")
		GetTwitterHandlerObj(appengine.NewContext(r)).TwitterLoginExit(w, r)
	})
	http.HandleFunc(UrlApiRoot+"/"+UrlBlobRequestUrl, func(w http.ResponseWriter, r *http.Request) {
		//	w.Header().Add("Access-Control-Allow-Origin", "*")
		GetBlobHandlerObj(appengine.NewContext(r)).BlobRequestToken(w, r)
	})
	http.HandleFunc(UrlApiRoot+"/"+UrlBlobCallback, func(w http.ResponseWriter, r *http.Request) {
		//	w.Header().Add("Access-Control-Allow-Origin", "*")
		GetBlobHandlerObj(appengine.NewContext(r)).HandleUploaded(w, r)
	})

}

func makeRandomId() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}
