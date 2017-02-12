package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/FactomProject/factom"
	"github.com/dbHackathon2017/hackathon/common"
	"github.com/dbHackathon2017/hackathon/common/constants"
	"github.com/dbHackathon2017/hackathon/common/primitives"
	"github.com/dbHackathon2017/hackathon/company"
	"github.com/dbHackathon2017/hackathon/factom-read"
)

var (
	FILES_PATH string = "../hackathon2"
	templates  *template.Template

	MainCompany *company.FakeCompany

	cacheLock    sync.RWMutex
	PensionCache map[string]common.Pension

	mux           *http.ServeMux
	TemplateMutex sync.Mutex

	templateDelims = []string{"{--{", "}--}"}
)

func AddToPensionCache(penid string, pen common.Pension) {
	if !ready_to_disp {
		return
	}

	cacheLock.Lock()
	PensionCache[penid] = pen
	cacheLock.Unlock()
}

func GetFromPensionCache(penid string) *common.Pension {
	if !ready_to_disp {
		return nil
	}
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	if pp, ok := PensionCache[penid]; ok {
		p := new(common.Pension)
		*p = pp
		return p
	}
	return nil
}

func GetCacheList() []common.Pension {
	if !ready_to_disp {
		return nil
	}
	cacheLock.RLock()
	defer cacheLock.RUnlock()
	list := make([]common.Pension, 0)
	for _, p := range PensionCache {
		list = append(list, p)
	}

	return list
}

var ready_to_disp = true

func firsLoadLocal() {
	time.Sleep(5 * time.Second)
	ready_to_disp = false
	loading = true
	//time.Sleep(3 * time.Second)
	list := MainCompany.LoadPenCacheFromDB()
	cacheLock.Lock()
	for _, p := range list {
		PensionCache[p.PensionID.String()] = p
		// AddToPensionCache(p.PensionID.String(), p)
	}
	cacheLock.Unlock()
	loading = false
	ready_to_disp = true
	fmt.Printf("Done Loading from cache, loaded %d pensions\n", len(list))
}

// loadCache loads pensions into the cache
var loading bool = false

func loadCache(time.Time) {
	if loading {
		return
	}
	loading = true
	//fmt.Println("Adding to cache")
	for i, p := range MainCompany.Pensions {
		//fmt.Printf("- #%d -", i)
		var _ = i
		fpen, err := read.GetPensionFromFactom(p.PensionID)
		if err != nil {
			fmt.Println("ERR in load: " + err.Error())
			continue
		}
		AddToPensionCache(fpen.PensionID.String(), *fpen)
	}
	loading = false

	MainCompany.Save(GetCacheList(), FULL_CACHE)
}

func InitTemplate() {
	TemplateMutex.Lock()
	// Put function into templates
	templates = template.New("main")
	// templates = template.Must(templates.ParseGlob(FILES_PATH + "/html/*"))
	templates = templates.Delims(templateDelims[0], templateDelims[1])
	TemplateMutex.Unlock()

	PensionCache = make(map[string]common.Pension)
}

func delayed() {
	err := MainCompany.Pensions[0].MoveChainTo(MainCompany.Pensions[1], "StevenMod", *primitives.RandomFileList(10))
	if err != nil {
		panic(err)
	}
}

func ServeFrontEnd(port int) {
	factom.SetFactomdServer(constants.REMOTE_HOST)
	// Templates
	InitTemplate()

	MainCompany = company.RandomFakeCompay()
	if USE_DB {
		MainCompany.LoadFromDB()
		if FULL_CACHE { // Load from DB before factom
			go func() {
				firsLoadLocal()
			}()
		} else { // Load from factom
			ready_to_disp = true // yea, bad name
		}
	}

	if MAKE_TRANS {
		fmt.Println()
		amt := 5
		for i := 0; i < amt; i++ {
			penId, err := MainCompany.CreateRandomPension()
			if err != nil {
				panic(err)
			}

			fmt.Println("Chain made, can be found here: " +
				"http://altcoin.host:8090/search?input=" + penId.String() + "&type=chainhead")
		}

		for i := 0; i < amt; i++ {
			MainCompany.Pensions[i].AddValue(10000, "Steven WOOT!", *primitives.RandomFileList(10), true)
			MainCompany.Pensions[i].AddValue(2500, "Steven WOOT!", *primitives.RandomFileList(10), true)
			MainCompany.Pensions[i].AddValue(2500, "Steven WOOT!", *primitives.RandomFileList(10), true)
			MainCompany.Pensions[i].AddValue(5025, "Steven WOOT!", *primitives.RandomFileList(10), true)
			MainCompany.Pensions[i].AddValue(6025, "Steven WOOT!", *primitives.RandomFileList(10), true)
		}

		go func() {
			time.Sleep(15 * time.Second)
			delayed()
		}()
	}

	go func() {
		time.Sleep(5 * time.Second)
		loadCache(time.Now())
	}()
	go doEvery(10*time.Second, loadCache)

	mux = http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(FILES_PATH)))

	http.HandleFunc("/", static(pageHandler))
	http.HandleFunc("/GET", HandleGETRequests)
	http.HandleFunc("/POST", HandlePOSTRequests)

	portStr := ":" + strconv.Itoa(port)

	fmt.Println("Starting GUI on http://locahost" + portStr + "/")
	http.ListenAndServe(portStr, nil)
}

// doEvery
// For go routines. Calls function once each duration.
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
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
	templates.Delims(templateDelims[0], templateDelims[1]).ParseGlob(FILES_PATH + "/html/*.html")
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
		w.Write(jsonError("Cannot do a POST request as /GET"))
		return
	}
	req := r.FormValue("request")
	switch req {
	case "on":
		w.Write(jsonResp(true))
	default:
		w.Write(jsonError("Not a valid GET request"))
	}
}

type POSTRequest struct {
	Request string          `json:"request"`
	Params  json.RawMessage `json:"params,omitempty"`
}

func HandlePOSTRequests(w http.ResponseWriter, r *http.Request) {
	// Only handles POST
	if r.Method != "POST" {
		w.Write(jsonError("Cannot do a GET request as /POST"))
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.Write(jsonError(err.Error()))
	}

	p := new(POSTRequest)
	err = json.Unmarshal(data, p)
	if err != nil {
		w.Write(jsonError(err.Error()))
	}
	// Form:
	//	request -- Request Function
	//	json	-- json object

	//req := r.FormValue("request")
	fmt.Println(p.Request)
	switch p.Request {
	case "on":
		w.Write(jsonResp(true))
	case "all-pensions":
		err := handleAllPensionsCompany(w, r)
		if err != nil {
			w.Write(jsonError(err.Error()))
		}
	case "all-pensions-user":
		err := handleAllPensionsUser(w, r)
		if err != nil {
			w.Write(jsonError(err.Error()))
		}
	case "pension":
		err := handlePension(w, r, data)
		if err != nil {
			w.Write(jsonError(err.Error()))
		}
	case "transaction":
		err := handleTransaction(w, r, data)
		if err != nil {
			w.Write(jsonError(err.Error()))
		}
	case "company-stats":
		err := handleCompanyStats(w, r, data)
		if err != nil {
			w.Write(jsonError(err.Error()))
		}
	case "makepension":
		err := handleMakePension(w, r, data)
		if err != nil {
			w.Write(jsonError(err.Error()))
		}
	case "addvalue":
		err := handleAddValue(w, r, data)
		if err != nil {
			w.Write(jsonError(err.Error()))
		}
	case "transfer":
		err := handleTransfer(w, r, data)
		if err != nil {
			w.Write(jsonError(err.Error()))
		}
	default:
		w.Write(jsonError("Not a valid POST request"))
	}

}
