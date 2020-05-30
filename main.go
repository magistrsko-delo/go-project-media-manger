package main

import (
	"fmt"
	"github.com/heptiolabs/healthcheck"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"go-project-media-manger/Models"
	"go-project-media-manger/router"
	"go-project-media-manger/services"

	// "go-project-media-manger/router"
	// "go-project-media-manger/services"
	"log"
	"net/http"
	"time"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
)

func init()  {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	Models.InitEnv()
}

func main()  {
	env := Models.GetEnvStruct()
	health := healthcheck.NewHandler()
	health.AddLivenessCheck("rabbitmq: " + env.RabbitHost + ":" + env.RabbitPort, healthcheck.TCPDialCheck(env.RabbitHost + ":" + env.RabbitPort, 5*time.Second))
	health.AddLivenessCheck("timeshift: " + env.TimeShiftGrpcServer + ":" + env.TimeShiftGrpcPort, healthcheck.TCPDialCheck(env.TimeShiftGrpcServer + ":" + env.TimeShiftGrpcPort, 5*time.Second))
	health.AddLivenessCheck("media-metadata: " + env.MediaMetadataGrpcServer + ":" + env.MediaMetadataGrpcPort, healthcheck.TCPDialCheck(env.MediaMetadataGrpcServer + ":" + env.MediaMetadataGrpcPort, 5*time.Second))
	health.AddLivenessCheck("media-chunks: " + env.MediaChunkMetadataServer + ":" + env.MediaChunkMetadataPort, healthcheck.TCPDialCheck(env.MediaChunkMetadataServer + ":" + env.MediaChunkMetadataPort, 5*time.Second))

	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	injector := jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, zipkinPropagator)
	extractor := jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, zipkinPropagator)

	zipkinSharedRPCSpan := jaeger.TracerOptions.ZipkinSharedRPCSpan(true)

	sender, err := jaeger.NewUDPTransport(Models.GetEnvStruct().TracingConnection, 0)

	rabbitMQ := services.InitRabbitMQ()

	// defer rabbitMQ.Conn.Close()
	// defer rabbitMQ.Ch.Close()

	r := mux.NewRouter()

	api := r.PathPrefix("/v1").Subrouter()
	projectMediaRouter := &router.ProjectMediaRouter{Router:api}
	projectMediaRouter.RegisterHandlers(rabbitMQ)

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

	go http.ListenAndServe("0.0.0.0:8888", health)

	if err == nil {
		fmt.Println("success: TRACING")
		tracer, closer := jaeger.NewTracer(
			"go-project-media-manager",
			jaeger.NewConstSampler(true),
			jaeger.NewRemoteReporter(
				sender,
				jaeger.ReporterOptions.BufferFlushInterval(1*time.Second)),
			injector,
			extractor,
			zipkinSharedRPCSpan,
		)
		defer closer.Close()
		log.Fatal(http.ListenAndServe(":" + Models.GetEnvStruct().Port, nethttp.Middleware(tracer, corsOpts.Handler(r)))  )
	} else {
		fmt.Println( "err: ", err)
		log.Fatal(http.ListenAndServe(":" + Models.GetEnvStruct().Port, corsOpts.Handler(r)) )
	}
}


func NotFound(w http.ResponseWriter, r *http.Request) {
	rsp := "route not found: " + r.URL.Path
	w.Write([]byte(rsp))
}
