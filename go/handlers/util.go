package handlers

import (
	"appengine"
	"appengine/memcache"
	"net/http"
)

type Error struct {
	StatusCode int
	Message    string
	Error      error
}

type ErrorHandler interface {
	Handle(w http.ResponseWriter, r *http.Request) *Error
}

func NewErrorHandler(h ErrorHandler) http.Handler {
	return errorHandler{h}
}

type CachedHandler interface {
	GetKey(*http.Request) string
	Handle(*http.Request) ([]byte, *Error)
}

func NewCachedHandler(h CachedHandler) http.Handler {
	return NewErrorHandler(cachedHandler{h})
}

type errorHandler struct {
	ErrorHandler
}

func (handler errorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	err := handler.Handle(w, r)
	if err != nil {
		cxt.Errorf("Error: %s", err.Error)
		http.Error(w, err.Message, err.StatusCode)
	}
}

type cachedHandler struct {
	CachedHandler
}

func (c cachedHandler) Handle(w http.ResponseWriter, r *http.Request) *Error {
	cacheKey := c.GetKey(r)

	if cacheKey == "" {
		res, err := c.CachedHandler.Handle(r)
		if err != nil {
			return err
		}
		w.Write(res)
		return nil
	}

	cxt := appengine.NewContext(r)
	if item, err := memcache.Get(cxt, cacheKey); err == nil {
		cxt.Infof("Memcache hit: %q", cacheKey)
		w.Write(item.Value)
		return nil
	} else if err != memcache.ErrCacheMiss {
		panic(err)
	}

	cxt.Infof("Memcache miss: %q", cacheKey)

	result, err := c.CachedHandler.Handle(r)

	if err != nil {
		return err
	}

	w.Write(result)

	item := &memcache.Item{
		Key:   cacheKey,
		Value: result,
	}

	if err := memcache.Add(cxt, item); err == memcache.ErrNotStored {
		cxt.Warningf("Memcache item already stored: %q", cacheKey)
	} else if err != nil {
		panic(err)
	} else {
		cxt.Infof("Cached successfully: %q", cacheKey)
	}
	return nil
}
