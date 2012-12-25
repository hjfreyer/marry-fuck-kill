package gomfk

import (
	_ "appengine"
	_ "appengine/datastore"
	_ "appengine/user"
	cryptorand "crypto/rand"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net/http"
	"github.com/hjfreyer/marry-fuck-kill/go/handlers"
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
	http.HandleFunc("/make", handlers.MakeHandler)
	http.Handle("/api/v1/imagesearch", handlers.ImageSearchApiHandler)
	http.Handle("/api/v1/make", handlers.MakeTripleApiHandler)
	http.Handle("/api/v1/vote", handlers.VoteApiHandler)
	http.Handle("/i/", handlers.GetImageHandler)
	http.Handle("/vote/", handlers.SingleTripleHandler)
}
