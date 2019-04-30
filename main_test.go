package main

import (
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
	//models.AutoMigrateSQL()
	//cache_handler.InitializeTest()
	req, err := http.NewRequest("GET", "/view/page/0", nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.ReadCachePagination)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

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

	req, err := http.NewRequest("PUT", "/insert/page", nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.PutData)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

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

	req, err := http.NewRequest("PUT", "/notify/reload", nil)
	if err != nil {
		t.Fatal(err)
	}
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.Reload)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

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
