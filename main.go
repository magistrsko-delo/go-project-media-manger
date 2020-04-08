package main

import (
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go-project-media-manger/Models"
	"go-project-media-manger/router"
	"log"
	"net/http"
)

func init()  {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	Models.InitEnv()
}

func main()  {
	r := mux.NewRouter()

	api := r.PathPrefix("/v1").Subrouter()
	projectMediaRouter := &router.ProjectMediaRouter{Router:api}
	projectMediaRouter.RegisterHandlers()

	r.NotFoundHandler = http.HandlerFunc(NotFound)

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{Models.GetEnvStruct().OriginAllowed},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"X-Requested-With",
		},
	})

	log.Fatal(http.ListenAndServe(":" + Models.GetEnvStruct().Port, corsOpts.Handler(r)))
}


func NotFound(w http.ResponseWriter, r *http.Request) {
	rsp := "route not found: " + r.URL.Path
	w.Write([]byte(rsp))
}
