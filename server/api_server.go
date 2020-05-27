package server

import (
	"github.com/gorilla/mux"
	"github.com/senseyman/image-media-processor/service"
	"github.com/sirupsen/logrus"
	"net/http"
)

// - process/resize image
// - list prev resized images with inputted params
// - all api can return errors
// - return obj - links to orig and resized images
// - if requested prev image with prev params - just return already processed img link

// APIServer process http requests
// Using api version mechanism

type APIServer struct {
	logger           *logrus.Logger
	address          string
	router           *mux.Router
	imgProcessor     service.MediaProcessor
	cloudStore       service.CloudStore
	requestProcessor *ApiServerRequestProcessor
}

var (
	NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
)

// create new instance of APIServer
func NewAPIServer(bindAddr string, logger *logrus.Logger, imgProcessor service.MediaProcessor, cloudStore service.CloudStore) *APIServer {
	return &APIServer{
		logger:           logger,
		address:          bindAddr,
		router:           mux.NewRouter(),
		requestProcessor: NewApiServerRequestProcessor(logger, imgProcessor, cloudStore),
	}
}

// Starting APIServer using port from config
func (s *APIServer) Start() error {
	s.logger.Infof("Starting api server. Port %s ", s.address)
	s.registerRouters()
	return http.ListenAndServe(s.address, s.router)
}

func (s *APIServer) registerRouters() {
	s.logger.Info("Registering api routers...")
	api := s.router.PathPrefix("/api").Subrouter()
	api.NotFoundHandler = NotFoundHandler

	s.registerRouteV1(api)
}

func (s *APIServer) registerRouteV1(parentRouter *mux.Router) {
	apiV1 := parentRouter.PathPrefix("/v1").Subrouter()

	apiV1.NotFoundHandler = NotFoundHandler

	apiV1.HandleFunc("/resize", s.requestProcessor.HandleResizeRequest).Methods(http.MethodPost)
	apiV1.HandleFunc("/list", s.requestProcessor.HandleListHistoryRequest).Methods(http.MethodGet)
}

func (s *APIServer) GetRouter() *mux.Router {
	return s.router
}
