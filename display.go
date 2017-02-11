package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"

	"github.com/FactomProject/factom"
	"github.com/dbHackathon2017/hackathon/common/constants"
	"github.com/dbHackathon2017/hackathon/common/primitives"
	"github.com/dbHackathon2017/hackathon/company"
)

var (
	FILES_PATH string = "../hackathon2"
	templates  *template.Template

	MainCompany *company.FakeCompany

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
	factom.SetFactomdServer(constants.REMOTE_HOST)
	// Templates
	InitTemplate()

	MainCompany = company.RandomFakeCompay()
	if MAKE_TRANS {
		for i := 0; i < 3; i++ {
			penId, err := MainCompany.CreatePension()
			if err != nil {
				panic(err)
			}

			MainCompany.Pensions[i].AddValue(100, "Steven WOOT!", *primitives.RandomFileList(10))
			MainCompany.Pensions[i].AddValue(25, "Steven WOOT!", *primitives.RandomFileList(10))
			MainCompany.Pensions[i].AddValue(25, "Steven WOOT!", *primitives.RandomFileList(10))
			MainCompany.Pensions[i].AddValue(25, "Steven WOOT!", *primitives.RandomFileList(10))
			MainCompany.Pensions[i].AddValue(25, "Steven WOOT!", *primitives.RandomFileList(10))

			fmt.Println("Chain made, can be found here: " +
				"http://altcoin.host:8090/search?input=" + penId.String() + "&type=chainhead")
		}
	}

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
	Error   string      `json:"error"`
	Content interface{} `json:"content"`
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
	case "on":
		w.Write(jsonResp(true))
	case "all-pensions":
		err := handleAllPensions(w, r)
		if err != nil {
			jsonError(err.Error())
		}
	case "pension":
	default:
		w.Write(jsonError("Not a post valid request"))
	}

}
