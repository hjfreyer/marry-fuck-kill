package gomfk

import (
	"math/rand"
	"time"
)

type VoteCount []int32

type EntityImage struct {
	SourceUrl   string
	ContentType string
	Data        []byte
}

type Triple struct {
	NameA    string
	ImageIdA int64
	VotesA   VoteCount

	NameB    string
	ImageIdB int64
	VotesB   VoteCount

	NameC    string
	ImageIdC int64
	VotesC   VoteCount

	Creator      UserId
	CreationTime time.Time

	Ordering int32
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

func (t *Triple) AddVote(v Vote) {
	t.VotesA[v.Vote[0]]++
	t.VotesB[v.Vote[1]]++
	t.VotesC[v.Vote[2]]++
}

func (t *Triple) SubtractVote(v Vote) {
	t.VotesA[v.Vote[0]]--
	t.VotesB[v.Vote[1]]--
	t.VotesC[v.Vote[2]]--
}

type Entity struct {
	Name    string
	ImageId int64

	MarryCount int32
	FuckCount  int32
	KillCount  int32
}

func (t *Triple) A() (e Entity) {
	e.Name = t.NameA
	e.ImageId = t.ImageIdA
	e.MarryCount = t.VotesA[0]
	e.FuckCount = t.VotesA[1]
	e.KillCount = t.VotesA[2]
	return
}

func (t *Triple) B() (e Entity) {
	e.Name = t.NameB
	e.ImageId = t.ImageIdB
	e.MarryCount = t.VotesB[0]
	e.FuckCount = t.VotesB[1]
	e.KillCount = t.VotesB[2]
	return
}

func (t *Triple) C() (e Entity) {
	e.Name = t.NameC
	e.ImageId = t.ImageIdC
	e.MarryCount = t.VotesC[0]
	e.FuckCount = t.VotesC[1]
	e.KillCount = t.VotesC[2]
	return
}

const (
	MARRY = 0
	FUCK  = 1
	KILL  = 2
)

type Vote struct {
	Vote []int
}

func VoteFromChar(v byte) int {
	switch v {
	case 'm':
		return MARRY
	case 'f':
		return FUCK
	case 'k':
		return KILL
	}
	panic("Invalid vote string")
}
