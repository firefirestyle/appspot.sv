package firestylesv

import (
	"net/http"
	//	"net/url"

	//	"errors"

	"github.com/firefirestyle/gominioauth/twitter"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
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

var twitterHandlerObj *twitter.TwitterHandler = nil

func GetTwitterHandlerObj(ctx context.Context) *twitter.TwitterHandler {
	if twitterHandlerObj == nil {
		twitterHandlerObj = twitter.NewTwitterHandler( //
			"http://"+appengine.DefaultVersionHostname(ctx)+""+UrlApiRoot+"/"+UrlTwitterTokenCallback, twitter.TwitterOAuthConfig{
				ConsumerKey:       TwitterConsumerKey,
				ConsumerSecret:    TwitterConsumerSecret,
				AccessToken:       TwitterAccessToken,
				AccessTokenSecret: TwitterAccessTokenSecret}, nil, nil)
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
	http.HandleFunc(UrlApiRoot+"/"+UrlTwitterTokenUrlRedirect, func(w http.ResponseWriter, r *http.Request) {
		//	w.Header().Add("Access-Control-Allow-Origin", "*")
		GetTwitterHandlerObj(appengine.NewContext(r)).TwitterLoginEntry(w, r)
	})
	http.HandleFunc(UrlApiRoot+"/"+UrlTwitterTokenCallback, func(w http.ResponseWriter, r *http.Request) {
		//	w.Header().Add("Access-Control-Allow-Origin", "*")
		GetTwitterHandlerObj(appengine.NewContext(r)).TwitterLoginExit(w, r)
	})
}
