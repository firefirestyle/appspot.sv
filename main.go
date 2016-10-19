package firestylesv

import (
	"net/http"
	//	"net/url"

	//	"errors"

	"github.com/firefirestyle/go.miniblob"
	"github.com/firefirestyle/go.minioauth/twitter"
	"github.com/firefirestyle/go.minisession"
	"github.com/firefirestyle/go.miniuser"
	//
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	//
	"crypto/rand"
	"encoding/binary"
	"strconv"
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

	UrlBlobRequestUrl = "/api/v1/blob/requesturl"
	UrlBlobCallback   = "/api/v1/blob/callback"

	UrlUserGetUrl = "/api/v1/user/get"

//	UrlBlobCallback   = "/api/v1/blob/callback"
)

var twitterHandlerObj *twitter.TwitterHandler = nil
var blobHandlerObj *miniblob.BlobHandler = nil
var sessionMgrObj *minisession.SessionManager = nil
var userHandlerObj *miniuser.UserHandler = nil

func GetUserMgrObj(ctx context.Context) *miniuser.UserHandler {
	if userHandlerObj == nil {
		userHandlerObj = miniuser.NewUserHandler(miniuser.UserManagerConfig{
			ProjectId:   "firefirestyle",
			UserKind:    "user",
			RelayIdKind: "relayId",
		}, miniuser.UserHandlerOnEvent{})
	}
	return userHandlerObj
}

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
			UrlBlobCallback, appengine.VersionID(ctx), //
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
			"http://"+appengine.DefaultVersionHostname(ctx)+""+UrlTwitterTokenCallback, twitter.TwitterOAuthConfig{
				ConsumerKey:       TwitterConsumerKey,
				ConsumerSecret:    TwitterConsumerSecret,
				AccessToken:       TwitterAccessToken,
				AccessTokenSecret: TwitterAccessTokenSecret}, twitter.TwitterHundlerOnEvent{
				OnFoundUser: func(w http.ResponseWriter, r *http.Request, handler *twitter.TwitterHandler, accesssToken *twitter.SendAccessTokenResult) map[string]string {
					ctx := appengine.NewContext(r)
					sessionMgrObj := GetSessionMgrObj(ctx)
					tokenObj, err := sessionMgrObj.Login(ctx, //
						accesssToken.GetScreenName(), //
						minisession.MakeAccessTokenConfigFromRequest(r))
					if err != nil {
						return map[string]string{"errcode": "1"}
					}

					userMgrObj := GetUserMgrObj(ctx)
					//_, userSessionObj, userObj. :=
					_, _, userObj, err1 := userMgrObj.GetManager().LoginRegistFromTwitter(ctx, //
						accesssToken.GetScreenName(), //
						accesssToken.GetUserID(),     //
						accesssToken.GetOAuthToken()) //
					if err1 != nil {
						return map[string]string{"errcode": "2", "errindo": err1.Error()}
					} else {
						return map[string]string{"token": "" + tokenObj.GetLoginId(), "userName": userObj.GetUserName()}
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
	http.HandleFunc(UrlTwitterTokenUrlRedirect, func(w http.ResponseWriter, r *http.Request) {
		//	w.Header().Add("Access-Control-Allow-Origin", "*")
		GetTwitterHandlerObj(appengine.NewContext(r)).TwitterLoginEntry(w, r)
	})
	http.HandleFunc(UrlTwitterTokenCallback, func(w http.ResponseWriter, r *http.Request) {
		//	w.Header().Add("Access-Control-Allow-Origin", "*")
		GetTwitterHandlerObj(appengine.NewContext(r)).TwitterLoginExit(w, r)
	})
	http.HandleFunc(UrlBlobRequestUrl, func(w http.ResponseWriter, r *http.Request) {
		//	w.Header().Add("Access-Control-Allow-Origin", "*")
		GetBlobHandlerObj(appengine.NewContext(r)).BlobRequestToken(w, r)
	})
	http.HandleFunc(UrlBlobCallback, func(w http.ResponseWriter, r *http.Request) {
		//	w.Header().Add("Access-Control-Allow-Origin", "*")
		GetBlobHandlerObj(appengine.NewContext(r)).HandleUploaded(w, r)
	})
	http.HandleFunc(UrlUserGetUrl, func(w http.ResponseWriter, r *http.Request) {
		//	w.Header().Add("Access-Control-Allow-Origin", "*")
		GetUserMgrObj(appengine.NewContext(r)).HandleGet(w, r)
	})
}

func makeRandomId() string {
	var n uint64
	binary.Read(rand.Reader, binary.LittleEndian, &n)
	return strconv.FormatUint(n, 36)
}
