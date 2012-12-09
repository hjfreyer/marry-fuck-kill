package gomfk

// import (
// 	_ "math/rand"
// 	"time"
// )

// type TripleId int64
// type ImageId struct {
// 	parent TripleId
// 	idx    int64
// }

// type TripleCursor int64

// type Triple struct {
// 	Id TripleId

// 	Entities [3]Entity

// 	UserVoted bool
// 	UserVote  Vote

// 	Creator      UserId
// 	CreationTime time.Time

// 	Ordering int64
// 	Disabled bool
// }

// type Entity struct {
// 	Name string

// 	MarryCount int32
// 	FuckCount  int32
// 	KillCount  int32
// }

// type TripleCreation struct {
// 	A, B, C EntityCreation
// 	Creator UserId
// }

// type EntityCreation struct {
// 	Name  string
// 	Image *FetchedImage
// }

// type Vote string

// func (v Vote) IsValid() bool {
// 	return v == "mfk" ||
// 		v == "mkf" ||
// 		v == "fmk" ||
// 		v == "fkm" ||
// 		v == "kmf" ||
// 		v == "kfm"
// }

// type DataAccessor interface {
// 	GetTripleIds(int) ([]TripleId, error)
// 	GetTriples(UserId, []TripleId) ([]Triple, error)
// 	MakeTriple(TripleCreation) (TripleId, error)

// 	GetImage(ImageId) (string, []byte, error)

// 	UpdateVote(TripleId, UserId, Vote) error
// }

// // // A triple from a specific user's perspective.
// // type TripleView struct {
// // 	A, B, C Entity
// // 	UserVoted bool
// // }

// // func (t *Triple) A() (e Entity) {
// // 	e.Name = t.NameA
// // 	e.ImageId = t.ImageIdA
// // 	e.MarryCount = t.VotesA[0]
// // 	e.FuckCount = t.VotesA[1]
// // 	e.KillCount = t.VotesA[2]
// // 	return
// // }

// // func (t *Triple) B() (e Entity) {
// // 	e.Name = t.NameB
// // 	e.ImageId = t.ImageIdB
// // 	e.MarryCount = t.VotesB[0]
// // 	e.FuckCount = t.VotesB[1]
// // 	e.KillCount = t.VotesB[2]
// // 	return
// // }

// // func (t *Triple) C() (e Entity) {
// // 	e.Name = t.NameC
// // 	e.ImageId = t.ImageIdC
// // 	e.MarryCount = t.VotesC[0]
// // 	e.FuckCount = t.VotesC[1]
// // 	e.KillCount = t.VotesC[2]
// // 	return
// // }
