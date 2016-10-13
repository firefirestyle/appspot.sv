package firestylesv

import (
	"net/http"
	//	"net/url"

	//	"errors"

	"github.com/firefirestyle/gominioauth/twitter"
	//	"google.golang.org/appengine"
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

var twitterHandlerObj = twitter.NewTwitterHandler("https://firefirestyle.appspot.com"+UrlApiRoot+"/"+UrlTwitterTokenCallback, twitter.TwitterOAuthConfig{
	ConsumerKey:       TwitterConsumerKey,
	ConsumerSecret:    TwitterConsumerSecret,
	AccessToken:       TwitterAccessToken,
	AccessTokenSecret: TwitterAccessTokenSecret}, nil, nil)

var l = map[string]func(http.ResponseWriter, *http.Request){
	UrlTwitterTokenUrlRedirect: twitterHandlerObj.TwitterLoginEntry,
	UrlTwitterTokenCallback:    twitterHandlerObj.TwitterLoginExit,
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
	for k, v := range l {
		http.HandleFunc(UrlApiRoot+"/"+k, v)
	}
}
