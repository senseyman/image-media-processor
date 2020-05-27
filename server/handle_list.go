package server

import (
	"net/http"
)

// Function to handle user request for getting list of all his request history
// - requested image and resized results
// - requested resize params
func (s *ApiServerRequestProcessor) HandleListHistoryRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

}
