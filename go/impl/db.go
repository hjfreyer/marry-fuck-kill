package impl

import (
"math/rand"
"time"
	"appengine"
	"appengine/datastore"
	"github.com/hjfreyer/marry-fuck-kill/go/mfklib"
	"github.com/hjfreyer/marry-fuck-kill/go/third_party/proto"
)

func NewDb(c appengine.Context) mfklib.Database {
	return mfkDb{c}
}

type dbTriple struct {
	Proto []byte
	Random int64  // Non-negative random number.
}

type dbTripleStats struct {
	Proto []byte
}

func dbTripleStatsFromStats(in mfklib.TripleStats) *dbTripleStats {
	copyTally := func(in mfklib.Tally) *mfklib.TripleStatsProto_Tally {
		return &mfklib.TripleStatsProto_Tally{
			Marry: proto.Uint64(in.Marry),
			Fuck: proto.Uint64(in.Fuck),
			Kill: proto.Uint64(in.Kill),
		}
	}

	statsProto := mfklib.TripleStatsProto{
		A: copyTally(in.A),
		B: copyTally(in.B),
		C: copyTally(in.C),
	}

	data, err := proto.Marshal(&statsProto)
	panicOnError(err)
	return &dbTripleStats{data}
}

func (s *dbTripleStats) ToStats() mfklib.TripleStats {
	var statsProto mfklib.TripleStatsProto
	panicOnError(proto.Unmarshal(s.Proto, &statsProto))

	copyTally := func(in *mfklib.TripleStatsProto_Tally) mfklib.Tally {
		return mfklib.Tally{
			Marry: in.GetMarry(),
			Fuck: in.GetFuck(),
			Kill: in.GetKill(),
		}
	}

	return mfklib.TripleStats{
		A: copyTally(statsProto.A),
		B: copyTally(statsProto.B),
		C: copyTally(statsProto.C),
	}
}

type dbTripleUserStatus struct {
	Proto []byte
}

func dbTripleUserStatusFromVote(vote mfklib.VoteStatus) *dbTripleUserStatus {
	protoStatus := mfklib.TripleUserStatus{
		Vote: vote.Enum(),
	}

	data, err := proto.Marshal(&protoStatus)
	panicOnError(err)

	return &dbTripleUserStatus{data}
}

func (s *dbTripleUserStatus) ToVote() mfklib.VoteStatus {
	var statusProto mfklib.TripleUserStatus
	panicOnError(proto.Unmarshal(s.Proto, &statusProto))

	return statusProto.GetVote()
}

type mfkDb struct {
	appengine.Context
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func (db mfkDb) AddTriple(triple *mfklib.Triple) (mfklib.TripleId, error) {
	tripleStr, err := proto.Marshal(triple)
	panicOnError(err)

	// Pick an int64 uniformly, based on the time.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	t := dbTriple{
		Proto: tripleStr,
		Random: r.Int63(),
	}
	tripleKey := datastore.NewIncompleteKey(db, "dbTriple", nil)
	tripleKey, err = datastore.Put(db, tripleKey, &t)
	if err != nil {
		return 0, mfklib.NewLowerLevelError(err, "adding new Triple: %s", *triple)
	}

	return mfklib.TripleId(tripleKey.IntID()), nil
}

func (db mfkDb) GetTriple(tripleId mfklib.TripleId) (*mfklib.Triple, error) {
	// Appengine datastore doesn't like keys with id 0.
	if int64(tripleId) == 0 {
		return nil, &mfklib.TripleNotFoundError{tripleId}
	}

	key := datastore.NewKey(db, "dbTriple", "", int64(tripleId), nil)

	triple := dbTriple{}
	switch err := datastore.Get(db, key, &triple); err {
	case nil:
		break
	case datastore.ErrNoSuchEntity:
		return nil, &mfklib.TripleNotFoundError{tripleId}
	default:
		return nil, mfklib.NewLowerLevelError(err, "get Triple %d", tripleId)
	}

	result := &mfklib.Triple{}
	err := proto.Unmarshal(triple.Proto, result)
	panicOnError(err)

	return result, nil
}

func (db mfkDb) UpdateStats(
	tripleId mfklib.TripleId, userId mfklib.UserId,
	stats *mfklib.TripleStats, vote *mfklib.VoteStatus,
	updater mfklib.Updater) error {
	if tripleId == 0 {
		panic("TripleId must not be zero")
	}
	if userId == "" {
		panic("userId must not be empty")
	}

	statsKey := datastore.NewKey(db, "dbTripleStats", "", int64(tripleId), nil)
	voteKey := datastore.NewKey(db, "dbTripleUserStatus", string(userId), 0, statsKey)

	var dbStats dbTripleStats
	var dbStatus dbTripleUserStatus

	updateFunc := func(c appengine.Context) error {
		err := datastore.GetMulti(c, []*datastore.Key{statsKey, voteKey},
			[]interface{}{&dbStats, &dbStatus})
		if err != nil {
			for _, e := range err.(appengine.MultiError) {
				if e != nil && e != datastore.ErrNoSuchEntity {
					return err
				}
			}
		}

		*stats = dbStats.ToStats()
		*vote = dbStatus.ToVote()

		store := updater()
		if store {
			dbStats = *dbTripleStatsFromStats(*stats)
			dbStatus = *dbTripleUserStatusFromVote(*vote)

			_, err := datastore.PutMulti(c, []*datastore.Key{statsKey, voteKey},
				[]interface{}{&dbStats, &dbStatus})
			if err != nil {
				return err
			}
		}
		return nil
	}


	err := datastore.RunInTransaction(db, updateFunc, &datastore.TransactionOptions{XG: false})
	if err != nil {
		return mfklib.NewLowerLevelError(err, "Updating stats for Triple %d and User %s",
			tripleId, userId)
	}
	return nil
}
