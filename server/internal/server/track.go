package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/os4ua/browser-crypto-mining/server/internal/track"
	"inet.af/netaddr"
)

type trackRegisterResponse struct {
	SessionID uuid.UUID `json:"sessionId"`
}

type trackUpdateRequest struct {
	SessionID      uuid.UUID `json:"sessionId"`
	CumulativeTime float64   `json:"cumulativeTime"`
}

func newTrackHandler(tr *track.Datastore) http.Handler {
	mux := http.NewServeMux()

	handleFuncExact(mux, http.MethodPost, "/register", func(w http.ResponseWriter, r *http.Request) {
		ip, err := requestIP(r)
		if err != nil {
			panic(err)
		}

		sid, err := tr.Register(ip)
		if err == track.ErrIPTooManySessions {
			respondJSON(w, http.StatusForbidden, newErrorResponse(err.Error()))
			return
		}

		resp := trackRegisterResponse{
			SessionID: sid,
		}

		log.Printf("register %s %s\n", ip, sid)

		respondJSON(w, http.StatusOK, &resp)
	})

	handleFuncExact(mux, http.MethodPost, "/update", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		req := trackUpdateRequest{}
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&req)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, newErrorResponse("Malformed JSON"))
			return
		}

		if req.SessionID == uuid.Nil {
			respondJSON(w, http.StatusBadRequest, newErrorResponse("Missing sessionId"))
			return
		}

		if req.CumulativeTime == 0 {
			respondJSON(w, http.StatusBadRequest, newErrorResponse("Missing cumulativeTime"))
			return
		}

		ip, err := requestIP(r)
		if err != nil {
			panic(err)
		}

		err = tr.Update(req.SessionID, req.CumulativeTime)
		if err != nil {
			respondJSON(w, http.StatusBadGateway, newErrorResponse(err.Error()))
			return
		}

		log.Printf("update %s %v\n", ip, req)
	})

	// Register default 404.
	mux.HandleFunc("/", notFoundHandler)

	return mux
}

func requestIP(r *http.Request) (netaddr.IP, error) {
	// Get IP address from X-Forwarded-For, and fall back to request host.

	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		orig := strings.SplitN(xff, ",", 1)[0]
		return netaddr.ParseIP(orig)
	}

	ipPort, err := netaddr.ParseIPPort(r.RemoteAddr)
	return ipPort.IP(), err
}
