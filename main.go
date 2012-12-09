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
	http.HandleFunc("/make", MakeHandler)
	http.Handle("/api/v1/imagesearch", ImageSearchApiHandler)
	http.Handle("/api/v1/make", MakeTripleApiHandler)
	// http.HandleFunc("/api/v1/vote", WrapHandler(ApiVoteHandler))
	http.Handle("/i/", GetImageHandler)
	http.Handle("/vote/", SingleTripleHandler)
}
