package gomfk

import (
	"math/rand"
	"time"
)

type VoteCount []int32

type Triple struct {
	NameA           string
	ImageIdA        string
	ImageSourceUrlA string
	VotesA          VoteCount

	NameB           string
	ImageIdB        string
	ImageSourceUrlB string
	VotesB          VoteCount

	NameC           string
	ImageIdC        string
	ImageSourceUrlC string
	VotesC          VoteCount

	Creator      UserId
	CreationTime time.Time
	Ordering     int32

	Disabled bool
}

func (t *Triple) Init(r *rand.Rand) {
	t.VotesA = []int32{0, 0, 0}
	t.VotesB = []int32{0, 0, 0}
	t.VotesC = []int32{0, 0, 0}

	t.Disabled = false
	t.Ordering = int32(r.Uint32())
	t.CreationTime = time.Now()
}

type Entity struct {
	Name           string
	ImageId        string
	ImageSourceUrl string

	MarryCount int32
	FuckCount  int32
	KillCount  int32
}

func (t *Triple) A() (e Entity) {
	e.Name = t.NameA
	e.ImageId = t.ImageIdA
	e.ImageSourceUrl = t.ImageSourceUrlA
	e.MarryCount = t.VotesA[0]
	e.FuckCount = t.VotesA[1]
	e.KillCount = t.VotesA[2]
	return
}

func (t *Triple) B() (e Entity) {
	e.Name = t.NameB
	e.ImageId = t.ImageIdB
	e.ImageSourceUrl = t.ImageSourceUrlB
	e.MarryCount = t.VotesB[0]
	e.FuckCount = t.VotesB[1]
	e.KillCount = t.VotesB[2]
	return
}

func (t *Triple) C() (e Entity) {
	e.Name = t.NameC
	e.ImageId = t.ImageIdC
	e.ImageSourceUrl = t.ImageSourceUrlC
	e.MarryCount = t.VotesC[0]
	e.FuckCount = t.VotesC[1]
	e.KillCount = t.VotesC[2]
	return
}

const (
	MARRY = 1
	FUCK  = 2
	KILL  = 3
)

type Vote struct {
	Vote []int
}
