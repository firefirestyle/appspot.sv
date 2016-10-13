package firestylesv

import (
	"net/http"
	"net/url"

	"github.com/firefirestyle/gominioauth/twitter"
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

var twitterHandlerObj = NewTwitterHandler(TwitterConsumerKey, TwitterConsumerSecret, TwitterAccessToken, TwitterAccessTokenSecret)
var l = map[string]func(http.ResponseWriter, *http.Request){
	UrlTwitterTokenUrlRedirect: twitterHandlerObj.twitterLoginEntry, //()func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Welcome to!!")) },
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
		http.HandleFunc(UrlApiRoot+k, v)
	}
}

// ---
//  twitter
// ---
type TwitterHandler struct {
	twitterManager *twitter.TwitterManager
}

func NewTwitterHandler(consumerKey string, consumerSecret string, accessToken string, accessTokenSecret string) *TwitterHandler {
	twitterHandlerObj := new(TwitterHandler)
	twitterHandlerObj.twitterManager = twitter.NewTwitterManager( //
		consumerKey, consumerSecret, accessToken, accessTokenSecret)
	return twitterHandlerObj
}

func (obj *TwitterHandler) MakeUrlNotFoundCallbackError(baseAddr string) (string, error) {
	urlObj, err := url.Parse(baseAddr)
	if err != nil {
		return "", err
	}
	query := urlObj.Query()
	query.Add("error", UrlTwitterTokenUrlRedirect_errorNotFoundCallbackUrl)
	urlObj.RawQuery = query.Encode()
	return urlObj.String(), nil
}

func (obj *TwitterHandler) MakeUrlFailedToMakeToken(baseAddr string) (string, error) {
	urlObj, err := url.Parse(baseAddr)
	if err != nil {
		return "", err
	}
	query := urlObj.Query()
	query.Add("error", UrlTwitterTokenUrlRedirect_errorFailedToMakeToken)
	urlObj.RawQuery = query.Encode()
	return urlObj.String(), nil
}

func (obj *TwitterHandler) twitterLoginEntry(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	values := r.URL.Query()

	callbackUrl := values.Get(UrlTwitterTokenUrlRedirect_callbackUrl)
	redirectUrl := ""
	if callbackUrl == "" {
		redirectUrl, _ = obj.MakeUrlNotFoundCallbackError(r.RemoteAddr)
	} else {
		twitterObj := obj.twitterManager.NewTwitter()
		twitterCallback := "http://" + appengine.DefaultVersionHostname(ctx) + UrlApiRoot + UrlTwitterTokenCallback + "?" + UrlTwitterTokenCallback_callbackUrl + "=" + url.QueryEscape(callbackUrl)
		oauthResult, err := twitterObj.SendRequestToken(ctx, twitterCallback)
		if err != nil {
			urlPattern1, errPattern1 := obj.MakeUrlFailedToMakeToken(callbackUrl)
			if errPattern1 != nil {
				redirectUrl, _ = obj.MakeUrlNotFoundCallbackError(r.RemoteAddr)
			} else {
				redirectUrl = urlPattern1
			}
		} else {
			redirectUrl = oauthResult.GetOAuthTokenUrl()
		}
	}
	http.Redirect(w, r, redirectUrl, http.StatusFound)
}
