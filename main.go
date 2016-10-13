package firestylesv

import (
	"net/http"
	"net/url"

	"errors"

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

var twitterHandlerObj = NewTwitterHandler(TwitterOAuthConfig{
	ConsumerKey:       TwitterConsumerKey,
	ConsumerSecret:    TwitterConsumerSecret,
	AccessToken:       TwitterAccessToken,
	AccessTokenSecret: TwitterAccessTokenSecret}, nil, nil)

var l = map[string]func(http.ResponseWriter, *http.Request){
	UrlTwitterTokenUrlRedirect: twitterHandlerObj.TwitterLoginEntry, //()func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Welcome to!!")) },
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
type TwitterOAuthConfig struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

type TwitterHandler struct {
	twitterManager *twitter.TwitterManager
	onRequest      func(url.Values) map[string]string
	onFoundUser    func(url.Values, *twitter.SendAccessTokenResult) map[string]string
}

func NewTwitterHandler(config TwitterOAuthConfig, onRequest func(url.Values) map[string]string,
	onFoundUser func(url.Values, *twitter.SendAccessTokenResult) map[string]string) *TwitterHandler {
	twitterHandlerObj := new(TwitterHandler)
	twitterHandlerObj.twitterManager = twitter.NewTwitterManager( //
		config.ConsumerKey, config.ConsumerSecret, config.AccessToken, config.AccessTokenSecret)
	twitterHandlerObj.onFoundUser = onFoundUser
	twitterHandlerObj.onRequest = onRequest
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

func (obj *TwitterHandler) TwitterLoginEntry(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	values := r.URL.Query()

	callbackUrl := values.Get(UrlTwitterTokenUrlRedirect_callbackUrl)
	redirectUrl := ""
	if callbackUrl == "" {
		redirectUrl, _ = obj.MakeUrlNotFoundCallbackError(r.RemoteAddr)
	} else {
		//
		twitterCallback := "http://" + appengine.DefaultVersionHostname(ctx) + UrlApiRoot + UrlTwitterTokenCallback + "?" + UrlTwitterTokenCallback_callbackUrl + "=" + url.QueryEscape(callbackUrl)
		twitterCallbackObj, _ := url.Parse(twitterCallback)
		if obj.onRequest != nil {
			twitterCallbackValues := twitterCallbackObj.Query()
			opts := obj.onRequest(r.URL.Query())
			for k, v := range opts {
				twitterCallbackValues.Add(k, v)
			}
			twitterCallbackObj.RawQuery = twitterCallbackValues.Encode()
		}
		//
		twitterObj := obj.twitterManager.NewTwitter()
		oauthResult, err := twitterObj.SendRequestToken(ctx, twitterCallbackObj.String())
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

func (obj *TwitterHandler) TwitterLoginExit(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	ctx := appengine.NewContext(r)
	//
	//
	callbackUrl := r.URL.Query().Get(UrlTwitterTokenCallback_callbackUrl)
	urlObj, err := url.Parse(callbackUrl)
	if err != nil {
		removeUrlObj, _ := url.Parse(r.RemoteAddr)
		query := removeUrlObj.Query()
		query.Add("error", "error")
		removeUrlObj.RawQuery = query.Encode()
		http.Redirect(w, r, removeUrlObj.String(), http.StatusFound)
		return
	}
	//
	//
	twitterObj := obj.twitterManager.NewTwitter()
	rt, err := twitterObj.OnCallbackSendRequestToken(ctx, r.URL)
	if err != nil || rt.GetScreenName() == "" || rt.GetUserID() == "" {
		rt = nil
		if err == nil && (rt.GetScreenName() == "" || rt.GetUserID() == "") {
			err = errors.New("empty user")
		}
	}

	if obj.onFoundUser != nil {
		values := urlObj.Query()
		opts := obj.onFoundUser(r.URL.Query(), rt)
		for k, v := range opts {
			values.Add(k, v)
		}
		urlObj.RawQuery = values.Encode()
	}
	//

	if err != nil {
		query := urlObj.Query()
		query.Add("error", "oauth")
		urlObj.RawQuery = query.Encode()
		http.Redirect(w, r, urlObj.String(), http.StatusFound)
		return
	} else {
		http.Redirect(w, r, urlObj.String(), http.StatusFound)
	}
}
