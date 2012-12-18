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

func getStatsAndStatus(cxt appengine.Context,
	statsKey, userStatusKey *datastore.Key,
	stats *mfklib.TripleStats, status *mfklib.TripleUserStatus) error {
	err := datastore.GetMulti(cxt, []*datastore.Key{statsKey, userStatusKey},
		[]interface{}{stats, status})
	if err != nil {
		merr := err.(appengine.MultiError)
		for _, e := range merr {
			if e != nil && e != datastore.ErrNoSuchEntity {
				return merr
			}
		}
		if merr[0] == datastore.ErrNoSuchEntity {
			*stats = mfklib.TripleStats{}
		}
		if merr[1] == datastore.ErrNoSuchEntity {
			*status = mfklib.TripleUserStatus{}
		}
	}
	return nil
}

func (db mfkDb) UpdateStats(
	tripleId mfklib.TripleId, userId mfklib.UserId,
	stats *mfklib.TripleStats, status *mfklib.TripleUserStatus,
	updater mfklib.Updater) error {
	statsKey := datastore.NewKey(db, "dbTripleStats", "", int64(tripleId), nil)
	userStatusKey := datastore.NewKey(
		db, "dbTripleUserStatus", string(userId), 0, statsKey)

	updateFunc := func(c appengine.Context) error {
		if err := getStatsAndStatus(c, statsKey, userStatusKey, stats, status); err != nil {
			return err
		}

		store, err := updater()
		if err != nil {
			return err
		}

		if store {
			_, err := datastore.PutMulti(c, []*datastore.Key{statsKey, userStatusKey},
				[]interface{}{stats, status})
			if err != nil {
				return err
			}
		}
		return nil
	}

	if err := datastore.RunInTransaction(
		db, updateFunc, &datastore.TransactionOptions{XG: false}); err != nil {
		return err
	}
	return nil
}
