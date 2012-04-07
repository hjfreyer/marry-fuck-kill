package gomfk

import (
	"appengine"
	"appengine/datastore"
	_ "appengine/user"
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net/http"
)

func NewRandom() *rand.Rand {
	var seed int64
	err := binary.Read(cryptorand.Reader, binary.LittleEndian, &seed)

	if err != nil {
		panic("binary read fail")
		fmt.Println("binary.Read failed:", err)
	}
	fmt.Println("seed", seed)
	return rand.New(rand.NewSource(seed))
}

func init() {
	http.HandleFunc("/make", makeHandler)
	http.HandleFunc("/make.do", makeDoHandler)
	http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	q := datastore.NewQuery("Triple").Limit(10)
	templates := Templates()

	triples := make([]Triple, 0, 10)
	if _, err := q.GetAll(cxt, &triples); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	cxt.Infof("%v", triples)

	if err := templates.ExecuteTemplate(w, "tripleList", triples); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//	fmt.Fprintf(w, "<a href='/make'>Hello</a>, %v!", "you")
}
