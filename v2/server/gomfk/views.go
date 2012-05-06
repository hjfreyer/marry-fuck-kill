package gomfk

import (
"bytes"
	"regexp"
_"strings"
"strconv"
	"appengine"
_	"appengine/datastore"
	_ "appengine/urlfetch"
	_ "appengine/user"
	_ "encoding/json"
	_ "fmt"
	"html/template"
	_ "io/ioutil"
	"net/http"
	_ "net/url"
)

type TripleView Triple

func (tv *TripleView) MaxVote() int32 {
	votes := []int32{
		tv.Entities[0].MarryCount,
		tv.Entities[0].FuckCount,
		tv.Entities[0].KillCount,
		tv.Entities[1].MarryCount,
		tv.Entities[1].FuckCount,
		tv.Entities[1].KillCount,
		tv.Entities[2].MarryCount,
		tv.Entities[2].FuckCount,
		tv.Entities[2].KillCount,
	}
	max := int32(0)
	for _, c := range votes {
		if c > max {
			max = c
		}
	}
	return max
}

func (tv *TripleView) Vote(idx int) string {
	switch tv.UserVote[idx] {
	case 'm': return "marry"
	case 'f': return "fuck"
	case 'k': return "kill"
	}
	panic("invalid vote")
}

const URL_TEMPL = "http://chart.apis.google.com/chart" +
	"?chxr=0,0,{{.Max}}" +
	"&chxt=y" +
	"&chbh=a" +
	"&chs={{.Width}}x{{.Height}}" +
	"&cht=bvg" +
	"&chco=9911BB,C76FDD,63067A" +
	"&chds=0,{{.Max}},0,{{.Max}},0,{{.Max}}" +
	"&chd=t:{{.MarryCount}}|{{.FuckCount}}|{{.KillCount}}" +
	"&chdl=Marry|Fuck|Kill" +
	"&chdlp=r"

func (tv TripleView) ChartUrl(entity int,
	marryAdd, fuckAdd, killAdd int32) template.URL {
	marryCount := tv.Entities[entity].MarryCount + marryAdd
	fuckCount := tv.Entities[entity].FuckCount + fuckAdd
	killCount := tv.Entities[entity].KillCount + killAdd

	t := template.Must(template.New("charturl").Parse(URL_TEMPL))

	var buffer bytes.Buffer
	err := t.Execute(&buffer, map[string]int32{
		"Width" : 202,
		"Height" : 90,
		"Max" : tv.MaxVote() + 1,
		"MarryCount" : marryCount,
		"FuckCount" : fuckCount,
		"KillCount" : killCount,
	})

	if err != nil {
		panic(err)
	}

	return template.URL(buffer.String())
}

func Templates() *template.Template {
	return template.Must(template.ParseFiles(
		"gomfk/templates/header.html",
		"gomfk/templates/make.html",
		"gomfk/templates/triple.html",
		"gomfk/templates/vote.html",
	))
}

func NotFound(w http.ResponseWriter) {
	http.Error(w, "404 not found!!!1", 404)
}

var IMAGE_RE,_ = regexp.Compile("^/i/([0-9]+)/([012])$")
func ImageHandler(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)

	match := IMAGE_RE.FindStringSubmatch(r.URL.Path)
	if match == nil {
		http.Error(w, "Invalid path", 400)
		return
	}

	parentId, err := strconv.ParseInt(match[1], 10, 64)
	if err != nil {
		http.Error(w, "Invalid path", 400)
		return
	}
	idx, err := strconv.ParseInt(match[2], 10, 64)
	if err != nil {
		http.Error(w, "Invalid path", 400)
		return
	}

	imageId := ImageId{TripleId(parentId), idx}
	cxt.Infof("Looking up image: %d", imageId)

 	db := NewAppengineDataAccessor(cxt)
	contentType, data, err := db.GetImage(imageId)

	if err != nil {
		// TODO(hjfreyer): standardize
		http.Error(w, "Image not found", 404)
		return
	}

	w.Header().Set("content-type", contentType)
	w.Write(data)
}

func ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		NotFound(w)
		return
	}

	cxt := appengine.NewContext(r)
 	db := NewAppengineDataAccessor(cxt)

	user := UserIdFromContext(r)

	tripleIds, err := db.GetTripleIds(10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	triples, err := db.GetTriples(user, tripleIds)
	if err != nil {
		cxt.Errorf("%v", []error(err.(appengine.MultiError)))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tripleViews := make([]TripleView, len(triples))
	for i, t := range triples {
		tripleViews[i] = TripleView(t)
	}

	templates := Templates()
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "tripleList", tripleViews); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	buffer.WriteTo(w)
}

func MakeHandler(w http.ResponseWriter, r *http.Request) {
	t := Templates()

	t.ExecuteTemplate(w, "make", nil)
}

func VoteHandler(w http.ResponseWriter, r *http.Request) {
	t := Templates()

	t.ExecuteTemplate(w, "vote", nil)
}

	// fetcher := urlfetch.Client(cxt)

	// response, err := fetcher.Get("https://www.googleapis.com/customsearch/v1?q=fox&cx=017343173679326196998%3Aomutomvh_wi&safe=medium&searchType=image&fields=items(image(thumbnailHeight%2CthumbnailLink%2CthumbnailWidth))&pp=1&key=AIzaSyDbjy0CKTMV5DoJR07ZYF5w-KL7Ey5lyGY")

	// if err != nil {
	//   cxt.Infof("%s", err)
	// 	return
	// } else {
	//   defer response.Body.Close()
	//   contents, err := ioutil.ReadAll(response.Body)
	//   if err != nil {
	//     cxt.Infof("%s", err)
	// 		return
	//   }
	//   cxt.Infof("%s\n", string(contents))
	// }
