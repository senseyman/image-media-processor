package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

// TODO add routes registration
// TODO add api version mechanism
// TODO add check user requests params

// TODO routes:
// - process/resize image
// - list prev resized images with inputted params
// - all api can return errors
// - return obj - links to orig and resized images
// - if requested prev image with prev params - just return already processed img link

// APIServer process http requests
// Using api version mechanism
type APIServer struct {
	logger  *logrus.Logger
	address string
	router  *mux.Router
}

var (
	NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})
)

// create new instance of APIServer
func NewAPIServer(bindAddr string, logger *logrus.Logger) *APIServer {
	return &APIServer{
		logger:  logger,
		address: bindAddr,
		router:  mux.NewRouter(),
	}
}

// Starting APIServer using port from config
func (s *APIServer) Start() error {
	s.logger.Infof("Starting api server. Port %s ", s.address)
	s.registerRouters()
	return http.ListenAndServe(s.address, s.router)
}

func (s *APIServer) registerRouters() {
	api := s.router.PathPrefix("/api").Subrouter()
	api.NotFoundHandler = NotFoundHandler

	s.registerRouteV1(api)
}

func (s *APIServer) registerRouteV1(parentRouter *mux.Router) {
	apiV1 := parentRouter.PathPrefix("/v1").Subrouter()

	apiV1.NotFoundHandler = NotFoundHandler

	// TODO put handle func
	apiV1.HandleFunc("/resize", s.handleResizeRequest).Methods(http.MethodPost)
	apiV1.HandleFunc("/list", nil).Methods(http.MethodGet)
}

func (s *APIServer) marshalDtoToJson(data interface{}) []byte {
	v, err := json.Marshal(data)
	if err != nil {
		s.logger.Errorf("Error while marshaling dto to byte array: %v", err)
		return nil
	}
	return v
}
