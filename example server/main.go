// dheilema 2018
// cached server example

package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/Nexinto/webcache"
)

var (
	apiCache webcache.CachedPage
	htmlt    *template.Template
)

const (
	// simple page with auto refresh
	testpage = `<html><head><meta http-equiv="Refresh" content="3">
<title>Webcache</title></head><body><h1>{{.Content}}</h1>
<a href="clear/">Clear cache</a>
</body></html>`
	port = "8080"
)

type tpTempl struct {
	Content string
}

func main() {
	htmlt = template.Must(template.New("testpage").Parse(testpage))

	// cache content for 10 seconds
	apiCache = webcache.NewCachedPage(time.Second * 10)

	go showStats(&apiCache)

	// the normal http handler
	http.HandleFunc("/clear/", clearHandler)
	http.HandleFunc("/", apiHandler)
	fmt.Println("HTTP Listening on port", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("http.ListenAndServe() failed: %s\n", err)
	}
}

// simple handler with a expensive backend function
func apiHandler(w http.ResponseWriter, r *http.Request) {
	if !apiCache.Valid() {
		if apiCache.StartUpdate() == nil {
			result := complexBackendFunction()
			htmlt.ExecuteTemplate(&apiCache, "testpage", tpTempl{Content: result})
			apiCache.EndUpdate()
		}
	}
	w.Write(apiCache.Get())
}

// clear the cache and go back to /
func clearHandler(w http.ResponseWriter, r *http.Request) {
	apiCache.Clear()
	http.Redirect(w, r, "/", http.StatusFound)
}

// simulate a long running backend function
func complexBackendFunction() string {
	time.Sleep(time.Second * 5)
	return "Last update at " + fmt.Sprint(time.Now().Format("15:04:05"))
}

// dump the stats to the console
func showStats(p *webcache.CachedPage) {
	for {
		time.Sleep(time.Minute)
		r, u := apiCache.GetStatistics()
		apiCache.ClearStatistics()
		fmt.Printf(time.Now().Format("15:04:05")+" requests=%d updates=%d\n", r, u)
	}
}
