package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rajatgpt1521/cachingSystem/service/pkg/cache_handler"
	"github.com/rajatgpt1521/cachingSystem/service/pkg/redis_handler"
	"net/http"
	"strconv"
)

type readResponse struct {
	Data []string
}
type putResponse struct {
	Msg string
}

var PAGESIZE = 2

func PutData(w http.ResponseWriter, r *http.Request) {
	data := mux.Vars(r)

	err, msg := cache_handler.Cacheing.Put(data["data"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

	} else {
		w.WriteHeader(http.StatusOK)
	}
	res := putResponse{msg}
	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	w.Write(js)
}

func Reload(w http.ResponseWriter, r *http.Request) {
	data := mux.Vars(r)
	msg := redis_handler.AddMessage(data["msg"])
	res := putResponse{msg}
	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)


}

func ReadCachePagination(w http.ResponseWriter, r *http.Request) {
	var cache_data readResponse
	skip := mux.Vars(r)
	err, cache := cache_handler.Cacheing.ReadAll()

	if err != nil {

		w.WriteHeader(http.StatusNotFound)
	}
	i1, err := strconv.Atoi(skip["pageno"])
	cache = paginate(cache, i1, PAGESIZE)
	for _, data := range cache {
		cache_data.Data = append(cache_data.Data, data)

	}

	js, err := json.Marshal(cache_data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)

}

func paginate(x []string, skip int, size int) []string {
	start := func() int {
		if skip*size > len(x) {
			return len(x)
		} else {
			return skip * size
		}
	}()

	limit := func() int {
		if start+size > len(x) {
			return len(x)
		} else {
			return start + size
		}
	}()

	return x[start:limit]
}
