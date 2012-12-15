package impl

import (
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
}

type dbTripleStats struct {
	Proto []byte
}

type dbTripleUserStatus struct {
	Proto []byte
}

type mfkDb struct {
	appengine.Context
}

func (db mfkDb) AddTriple(triple *mfklib.Triple) (mfklib.TripleId, error) {
	tripleStr, err := proto.Marshal(triple)
	if err != nil {
		panic(err)
	}

	t := dbTriple{tripleStr}
	tripleKey := datastore.NewIncompleteKey(db, "dbTriple", nil)
	tripleKey, err = datastore.Put(db, tripleKey, &t)
	if err != nil {
		return 0, err
	}

	return mfklib.TripleId(tripleKey.IntID()), nil
}

func (db mfkDb) GetTriple(tripleId mfklib.TripleId) (*mfklib.Triple, error) {
	key := datastore.NewKey(db, "dbTriple", "", int64(tripleId), nil)

	triple := dbTriple{}
	if err := datastore.Get(db, key, &triple); err == datastore.ErrNoSuchEntity {
		return nil, mfklib.EntityNotFoundError{
			Type: "Triple",
			Id:   int64(tripleId),
			Err:  err,
		}
	} else if err != nil {
		return nil, err
	}

	result := &mfklib.Triple{}
	if err := proto.Unmarshal(triple.Proto, result); err != nil {
		panic(err)
	}

	return result, nil
}

func getStatsAndStatus(statsKey, statusKey *datastore.Key, 

func (db mfkDb) UpdateStats(
	tripleId mfklib.TripleId, userId mfklib.UserId,
	stats *mfklib.TripleStats, status *mfklib.TripleUserStatus,
	updater mfklib.Updater) error {

	statsKey := datastore.NewKey(db, "dbTripleStats", "", int64(tripleId), nil)
	userStatusKey := datastore.NewKey(
		db, "dbTripleUserStatus", string(userId), 0, statsKey)

	err := datastore.RunInTransaction(db, func(c appengine.Context) error {
		if err := datastore.GetMulti(c,
			[]*datastore.Key{statsKey, userStatusKey}, 
			[]interface{}{stats, status}); err != nil {
			merr := err.(appengine.MultiError)
			if 
		}

		if err := datastore.Get(c, statusKey, status); err == datastore.ErrNoSuchEntity {
			*status := mfklib.TripleUserStatus{}
		} else if err != nil {
			return err
		}

		store, err := updater()
		
		if store {
			
		}
	})
}
