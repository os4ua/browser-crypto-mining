package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"inet.af/netaddr"
)

type trackRequest struct {
	SessionID         uuid.UUID
	CumulativeSeconds int
}

func trackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s := http.StatusMethodNotAllowed
		respondJSON(w, s, newErrorResponse(http.StatusText(s)))
		return
	}

	defer r.Body.Close()

	req := trackRequest{}
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

	if req.CumulativeSeconds == 0 {
		respondJSON(w, http.StatusBadRequest, newErrorResponse("Missing cumulativeSeconds"))
		return
	}

	// Get IP address from X-Forwarded-For, and fall back to request host.
	var ip netaddr.IP
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		orig := strings.SplitN(xff, ",", 1)[0]
		ip, err = netaddr.ParseIP(orig)
	} else {
		var ipPort netaddr.IPPort
		ipPort, err = netaddr.ParseIPPort(r.RemoteAddr)
		ip = ipPort.IP()
	}
	if err != nil {
		panic(err)
	}

	log.Printf("/track %s %+v\n", ip, req)

	// TODO: actually track :D
}
