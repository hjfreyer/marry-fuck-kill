package gomfk

// import (
// 	"appengine"
// 	"appengine/datastore"
// 	"errors"
// 	"fmt"
// 	_ "math/rand"
// 	"time"
// )

// func NewAppengineDataAccessor(cxt appengine.Context) DataAccessor {
// 	return appEngineDataAccessor{cxt}
// }

// type dbEntityImage FetchedImage

// type dbTriple struct {
// 	NameA  string
// 	VotesA dbVoteCount

// 	NameB  string
// 	VotesB dbVoteCount

// 	NameC  string
// 	VotesC dbVoteCount

// 	Creator      UserId
// 	CreationTime time.Time

// 	Ordering int64
// 	Disabled bool
// }

// type dbVoteCount []int32

// func voteIdx(vote byte) int {
// 	switch vote {
// 	case 'm':
// 		return 0
// 	case 'f':
// 		return 1
// 	case 'k':
// 		return 2
// 	}
// 	panic(fmt.Sprintf("Invalid vote: %c", vote))
// }

// func (t *dbTriple) addVote(vote dbVote) {
// 	if vote.Vote == "" {
// 		return
// 	}
// 	t.VotesA[voteIdx(vote.Vote[0])]++
// 	t.VotesB[voteIdx(vote.Vote[1])]++
// 	t.VotesC[voteIdx(vote.Vote[2])]++
// }

// func (t *dbTriple) subtractVote(vote dbVote) {
// 	if vote.Vote == "" {
// 		return
// 	}
// 	t.VotesA[voteIdx(vote.Vote[0])]--
// 	t.VotesB[voteIdx(vote.Vote[1])]--
// 	t.VotesC[voteIdx(vote.Vote[2])]--
// }

// type appEngineDataAccessor struct {
// 	cxt appengine.Context
// }

// func (db appEngineDataAccessor) GetImage(id ImageId) (string, []byte, error) {
// 	tripleKey := datastore.NewKey(db.cxt, "dbTriple", "", int64(id.parent), nil)

// 	key := datastore.NewKey(db.cxt, "dbEntityImage",
// 		fmt.Sprintf("%d", id.idx), 0, tripleKey)

// 	var image dbEntityImage
// 	if err := datastore.Get(db.cxt, key, &image); err != nil {
// 		return "", []byte{}, err
// 	}

// 	return image.ContentType, image.Data, nil
// }

// func (db appEngineDataAccessor) GetTripleIds(count int) ([]TripleId, error) {
// 	q := datastore.NewQuery("dbTriple").Limit(count).KeysOnly()

// 	keys, err := q.GetAll(db.cxt, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ids := make([]TripleId, len(keys))
// 	for i, key := range keys {
// 		ids[i] = TripleId(key.IntID())
// 	}

// 	return ids, nil
// }

// func (db appEngineDataAccessor) GetTriples(userId UserId, tripleIds []TripleId) (
// 	[]Triple, error) {
// 	keys := make([]*datastore.Key, 0, 2*len(tripleIds))
// 	triplesAndVotes := make([]interface{}, 0, 2*len(tripleIds))

// 	for _, tripleId := range tripleIds {
// 		tripleKey := datastore.NewKey(db.cxt, "dbTriple", "", int64(tripleId), nil)
// 		voteKey := datastore.NewKey(db.cxt, "dbVote", string(userId), 0, tripleKey)
// 		keys = append(keys, tripleKey, voteKey)
// 		triplesAndVotes = append(triplesAndVotes, &dbTriple{}, &dbVote{})
// 	}

// 	err := datastore.GetMulti(db.cxt, keys, triplesAndVotes)
// 	var errs appengine.MultiError
// 	if err != nil {
// 		errs = err.(appengine.MultiError)
// 	} else {
// 		errs = appengine.MultiError(make([]error, len(keys)))
// 	}

// 	triples := make([]Triple, len(tripleIds))
// 	for i, _ := range triples {
// 		tripleIdx := 2 * i
// 		voteIdx := tripleIdx + 1

// 		if errs[tripleIdx] != nil {
// 			return nil, errs[tripleIdx]
// 		}

// 		dbT := triplesAndVotes[tripleIdx].(*dbTriple)
// 		dbV := triplesAndVotes[voteIdx].(*dbVote)

// 		userVoted := (errs[voteIdx] == nil)
// 		if userVoted {
// 			dbT.subtractVote(*dbV)
// 		}

// 		triples[i] = Triple{
// 			Id: TripleId(keys[tripleIdx].IntID()),
// 			Entities: [3]Entity{
// 				Entity{
// 					Name:       dbT.NameA,
// 					MarryCount: dbT.VotesA[0],
// 					FuckCount:  dbT.VotesA[1],
// 					KillCount:  dbT.VotesA[2],
// 				},
// 				Entity{
// 					Name:       dbT.NameB,
// 					MarryCount: dbT.VotesB[0],
// 					FuckCount:  dbT.VotesB[1],
// 					KillCount:  dbT.VotesB[2],
// 				},
// 				Entity{
// 					Name:       dbT.NameC,
// 					MarryCount: dbT.VotesC[0],
// 					FuckCount:  dbT.VotesC[1],
// 					KillCount:  dbT.VotesC[2],
// 				},
// 			},
// 			UserVoted: userVoted,
// 			UserVote:  Vote(dbV.Vote),
// 		}
// 	}

// 	return triples, nil
// }

// func (db appEngineDataAccessor) MakeTriple(request TripleCreation) (
// 	TripleId, error) {
// 	r := NewRandom()

// 	t := dbTriple{
// 		NameA:        request.A.Name,
// 		VotesA:       []int32{0, 0, 0},
// 		NameB:        request.B.Name,
// 		VotesB:       []int32{0, 0, 0},
// 		NameC:        request.C.Name,
// 		VotesC:       []int32{0, 0, 0},
// 		Creator:      request.Creator,
// 		CreationTime: time.Now(),
// 		Ordering:     r.Int63(),
// 		Disabled:     false,
// 	}

// 	images := []*dbEntityImage{
// 		&dbEntityImage{
// 			SourceUrl:   request.A.Image.SourceUrl,
// 			ContentType: request.A.Image.ContentType,
// 			Data:        request.A.Image.Data,
// 		},
// 		&dbEntityImage{
// 			SourceUrl:   request.B.Image.SourceUrl,
// 			ContentType: request.B.Image.ContentType,
// 			Data:        request.B.Image.Data,
// 		},
// 		&dbEntityImage{
// 			SourceUrl:   request.C.Image.SourceUrl,
// 			ContentType: request.C.Image.ContentType,
// 			Data:        request.C.Image.Data,
// 		},
// 	}

// 	tripleKey := datastore.NewIncompleteKey(db.cxt, "dbTriple", nil)
// 	err := datastore.RunInTransaction(db.cxt, func(cxt appengine.Context) error {
// 		var err error
// 		tripleKey, err = datastore.Put(cxt, tripleKey, &t)
// 		if err != nil {
// 			return err
// 		}

// 		// Add images.
// 		imageKeys := []*datastore.Key{
// 			datastore.NewKey(db.cxt, "dbEntityImage", "0", 0, tripleKey),
// 			datastore.NewKey(db.cxt, "dbEntityImage", "1", 0, tripleKey),
// 			datastore.NewKey(db.cxt, "dbEntityImage", "2", 0, tripleKey),
// 		}

// 		if _, err := datastore.PutMulti(cxt, imageKeys, images); err != nil {
// 			return err
// 		}

// 		return nil
// 	}, nil)

// 	if err != nil {
// 		return 0, err
// 	}
// 	return TripleId(tripleKey.IntID()), nil
// }

// type dbVote struct {
// 	Vote string
// }

// func (db appEngineDataAccessor) GetVote(TripleId, UserId) (v Vote, e error) {
// 	return
// }

// func (db appEngineDataAccessor) UpdateVote(tripleId TripleId,
// 	userId UserId, vote Vote) error {
// 	tripleKey := datastore.NewKey(db.cxt, "dbTriple", "", int64(tripleId), nil)
// 	voteKey := datastore.NewKey(db.cxt, "dbVote", string(userId), 0, tripleKey)

// 	newVote := dbVote{string(vote)}
// 	if !vote.IsValid() {
// 		panic("invalid vote")
// 	}

// 	for tryTime := 0; tryTime < RETRY_COUNT; tryTime++ {
// 		err := datastore.RunInTransaction(db.cxt,
// 			func(cxt appengine.Context) error {
// 				var triple dbTriple
// 				if err := datastore.Get(cxt, tripleKey, &triple); err != nil {
// 					return err
// 				}

// 				var oldVote dbVote
// 				err := datastore.Get(cxt, voteKey, &oldVote)
// 				if err != nil && err != datastore.ErrNoSuchEntity {
// 					return err
// 				}
// 				if err != datastore.ErrNoSuchEntity {
// 					triple.subtractVote(oldVote)
// 				}
// 				triple.addVote(newVote)

// 				if _, err := datastore.PutMulti(cxt, []*datastore.Key{tripleKey, voteKey},
// 					[]interface{}{&triple, &newVote}); err != nil {
// 					return err
// 				}

// 				return nil
// 			}, nil)

// 		if err == nil {
// 			return nil
// 		} else if err == datastore.ErrNoSuchEntity {
// 			return err
// 		} else {
// 			db.cxt.Errorf("Transaction failed: %v", err)
// 		}
// 	}
// 	return errors.New("Transaction failed too many times")
// }

// // func (t *Triple) Init(r *rand.Rand) {
// // }

// // func (t *Triple) AddVote(v Vote) {
// // 	t.VotesA[v.Vote[0]]++
// // 	t.VotesB[v.Vote[1]]++
// // 	t.VotesC[v.Vote[2]]++
// // }

// // func (t *Triple) SubtractVote(v Vote) {
// // 	t.VotesA[v.Vote[0]]--
// // 	t.VotesB[v.Vote[1]]--
// // 	t.VotesC[v.Vote[2]]--
// // }
