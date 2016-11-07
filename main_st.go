package firestylesv

import (
	"net/http"

	"os"

	"time"
)

//
//
//
type EasyFileReader struct {
	fp *os.File
}

func NewEasyFileReaderFromFilePath(path string) (*EasyFileReader, error) {
	fp, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	rfrObj := new(EasyFileReader)
	rfrObj.fp = fp
	return rfrObj, nil
}

func (obj *EasyFileReader) Read(p []byte) (n int, err error) {
	return obj.fp.Read(p)
}

func (obj *EasyFileReader) Seek(offset int64, whence int) (int64, error) {
	return obj.fp.Seek(offset, whence)
}

func initTwitterCard() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//		ctx := appengine.NewContext(r)

		obj, err := NewEasyFileReaderFromFilePath("web" + r.URL.Path)
		if err == nil {
			w.Header().Set("Cache-Control", "public, max-age=2592000")
			http.ServeContent(w, r, r.URL.Path, time.Now().Add(2592000), obj)
		}
	})
}
