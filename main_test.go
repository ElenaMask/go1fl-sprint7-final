package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	table := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}
	for _, v := range table {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?search=%s&city=moscow", url.QueryEscape(v.search)), nil)
		handler.ServeHTTP(response, req)
		cafes := strings.Split(strings.TrimSpace(response.Body.String()), ",")
		if len(cafes) == 1 && cafes[0] == "" {
			cafes = []string{}
		}
		assert.Len(t, cafes, v.wantCount)
		for _, cafe := range cafes {
			cafe = strings.ToLower(cafe)
			v.search = strings.ToLower(v.search)
			assert.Contains(t, cafe, v.search)
		}
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		city  string
		count int
		want  int
	}{
		{"tula", 0, 0},
		{"moscow", 1, 1},
		{"tula", 2, 2},
		{"moscow", 100, min(100, len(cafeList["moscow"]))},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=%s&count=%d", url.QueryEscape(v.city), v.count), nil)
		handler.ServeHTTP(response, req)
		cafes := strings.Split(strings.TrimSpace(response.Body.String()), ",")
		if len(cafes) == 1 && cafes[0] == "" {
			cafes = []string{}
		}
		assert.Len(t, cafes, v.want)
	}
}

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}
