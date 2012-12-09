package impl

import (
//	"code.google.com/p/goprotobuf/proto"
	"appengine"
	"appengine/datastore"

	"github.com/hjfreyer/marry-fuck-kill/go/mfklib"
	"github.com/hjfreyer/marry-fuck-kill/go/third_party/proto"
)

func NewDb(c appengine.Context) mfklib.Database {
	return mfkDb{c}
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
			Id: int64(tripleId),
			Err: err,
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

func (db mfkDb) UpdateStats(tripleId mfklib.TripleId, updater mfklib.TripleStatsUpdater) error {
	return nil
}

type dbTriple struct {
	Proto []byte
}