package firestylesv

import (
	"net/http"

	userTmp "github.com/firefirestyle/go.miniuser/template"

	artTmp "github.com/firefirestyle/go.miniarticle/template"
	"golang.org/x/net/context"
	//	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
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
