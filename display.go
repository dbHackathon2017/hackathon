package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"
)

var (
	FILES_PATH string = "web/"
	templates  *template.Template

	mux           *http.ServeMux
	TemplateMutex sync.Mutex
)

func InitTemplate() {
	TemplateMutex.Lock()
	// Put function into templates
	templates = template.New("main")
	//templates = template.Must(templates.ParseGlob(FILES_PATH + "/html/*"))

	TemplateMutex.Unlock()
}

func ServeFrontEnd(port int) {
	// Templates
	InitTemplate()

	mux = http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(FILES_PATH)))

	http.HandleFunc("/", static(pageHandler))
	http.HandleFunc("/GET", HandleGETRequests)
	http.HandleFunc("/POST", HandlePOSTRequests)

	portStr := ":" + strconv.Itoa(port)

	fmt.Println("Starting GUI on http://" + portStr + "/")
	http.ListenAndServe(portStr, nil)
}

// mkArray makes an array inside a template
func mkArray(args ...interface{}) []interface{} {
	return args
}

// compareInts is used inside templates to compare ints
func compareInts(a int, b int) bool {
	return (a == b)
}

// compareStrings used inside templates to compare strings
func compareStrings(a string, b string) bool {
	return (a == b)
}

// For all static files. (CSS, JS, IMG, etc...)
func static(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.ContainsRune(r.URL.Path, '.') {
			mux.ServeHTTP(w, r)
			return
		}
		h.ServeHTTP(w, r)
	}
}

// pageHandler redirects all page requests to proper handlers
func pageHandler(w http.ResponseWriter, r *http.Request) {
	TemplateMutex.Lock()
	templates.ParseGlob(FILES_PATH + "/html/*.html")
	TemplateMutex.Unlock()
	request := strings.Split(r.RequestURI, "?")
	var err error
	switch request[0] {
	case "/":
		err = HandleIndexPage(w, r)
	default:
		err = HandleNotFoundError(w, r)
	}

	if err != nil {
		fmt.Printf("An error has occurred")
	}
}

// jsonResponse is used for responding to Post/Get Requests
type jsonResponse struct {
	Error   string      `json:"Error"`
	Content interface{} `json:"Content"`
}

func newJsonResponse(err string, content interface{}) *jsonResponse {
	j := new(jsonResponse)
	j.Error = err
	j.Content = content

	return j
}

func (j *jsonResponse) Bytes() []byte {
	data, err := json.Marshal(j)
	if err != nil {
		return nil
	}

	return data
}

// jsonResp used if request is successful
func jsonResp(content interface{}) []byte {
	e := newJsonResponse("none", content)
	return e.Bytes()
}

// jsonError used if request has an error
func jsonError(err string) []byte {
	e := newJsonResponse(err, "none")
	return e.Bytes()
}

func HandleGETRequests(w http.ResponseWriter, r *http.Request) {
	// Only handles GET
	if r.Method != "GET" {
		return
	}
	req := r.FormValue("request")
	switch req {
	case "on":
		w.Write(jsonResp(true))
	default:
		w.Write(jsonError("Not a valid request"))
	}
}

func HandlePOSTRequests(w http.ResponseWriter, r *http.Request) {
	// Only handles POST
	if r.Method != "POST" {
		return
	}

	// Form:
	//	request -- Request Function
	//	json	-- json object

	req := r.FormValue("request")
	switch req {
	default:
		w.Write(jsonError("Not a post valid request"))
	}

}
