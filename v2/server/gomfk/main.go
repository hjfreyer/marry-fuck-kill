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
	http.HandleFunc("/api/v1", ApiHandler)
	http.HandleFunc("/vote", VoteHandler)
	http.HandleFunc("/vote.do", VoteApiHandler)
	http.HandleFunc("/i/", ImageHandler)
	http.HandleFunc("/", ListHandler)
}
