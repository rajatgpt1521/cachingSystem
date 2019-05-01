package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rajatgpt1521/cachingSystem/service/models"
	"github.com/rajatgpt1521/cachingSystem/service/pkg/cache_handler"
	"github.com/rajatgpt1521/cachingSystem/service/pkg/database"
	"github.com/rajatgpt1521/cachingSystem/service/server"
	"github.com/rs/zerolog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReadCacheHandler(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	cache_handler.InitializeTestDB()
	path := fmt.Sprintf("/view/page/%s", "123")
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/view/page/{pageno}", server.ReadCachePagination)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"Data":["dog","tree"]}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	database.Instance.Set("gorm:table_options", "CASCADE").DropTableIfExists(&models.Cache{})
}

func TestPutCacheHandler(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	cache_handler.InitializeTestDB()

	path := fmt.Sprintf("/insert/%s", "123")
	req, err := http.NewRequest("PUT", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/insert/{data}", server.PutData)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"Msg":"Successfully added in cache"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	database.Instance.Set("gorm:table_options", "CASCADE").DropTableIfExists(&models.Cache{})
}
func TestReloadHandler(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	path := fmt.Sprintf("/notify/%s", "reload")
	req, err := http.NewRequest("PUT", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/notify/{msg}", server.Reload)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router.ServeHTTP(rr, req)
	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"Msg":"Successfully notified"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
	database.Instance.Set("gorm:table_options", "CASCADE").DropTableIfExists(&models.Cache{})
}
