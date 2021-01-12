package handler

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(handler *APIHandler) http.Handler {
	router := mux.NewRouter().StrictSlash(true)

	var publicRoutes = Routes{
		{
			"ListBuilds",
			"GET",
			"/apis/builds",
			handler.ListBuilds,
		},
		{
			"GetBuild",
			"GET",
			"/apis/build",
			handler.GetBuild,
		},
		{
			"CreateBuild",
			"POST",
			"/apis/build",
			handler.CreateBuild,
		},
		{
			"StopBuild",
			"DELETE",
			"/apis/build/{build}",
			handler.StopBuild,
		},
		{
			"ListRepos",
			"GET",
			"/apis/repo",
			handler.ListRepos,
		},
		{
			"GetRepo",
			"GET",
			"/apis/repo",
			handler.GetRepo,
		},
		{
			"EnableRepo",
			"POST",
			"/apis/repo",
			handler.EnableRepo,
		},
		{
			"DisableRepo",
			"DELETE",
			"/apis/repo/{repo}",
			handler.DisableRepo,
		},
		{
			"UpdateRepo",
			"PATCH",
			"/apis/repo/{repo}",
			handler.UpdateRepo,
		},
	}

	// The public route is always accessible
	for _, route := range publicRoutes {
		router.Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"Authorization"},
	})

	return c.Handler(router)
}
