package track

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"inet.af/netaddr"
)

const (
	maxSessionsPerIP          = 10
	minUpdateInterval         = 1 * time.Second
	maxTimeAbsDiff    float64 = 10.0
	maxTimeRelDiff    float64 = 0.01
)

var (
	ErrIPTooManySessions      = errors.New("IP address has too many sessions registered")
	ErrSessionNotRegistered   = errors.New("Session is not registered")
	ErrTooSoon                = errors.New("Update is too soon after the last")
	ErrCumulativeTimeTooSmall = errors.New("Cumulative time is smaller than what is already stored")
	ErrCumulativeTimeTooLarge = errors.New("Cumulative time is unrealistally large")
)

type session struct {
	mu             sync.Mutex
	cumulativeTime float64
	ip             netaddr.IP
	lastSeen       time.Time
	registered     time.Time
}

type Datastore struct {
	mu   sync.Mutex
	byIP map[netaddr.IP][]*session
	byID map[uuid.UUID]*session
}

func NewDatastore() *Datastore {
	return &Datastore{
		byIP: make(map[netaddr.IP][]*session),
		byID: make(map[uuid.UUID]*session),
	}
}

func (d *Datastore) Register(ip netaddr.IP) (uuid.UUID, error) {
	now := time.Now()

	d.mu.Lock()
	defer d.mu.Unlock()

	if _, found := d.byIP[ip]; !found {
		d.byIP[ip] = make([]*session, 0)
	}

	if len(d.byIP[ip]) >= maxSessionsPerIP {
		return uuid.Nil, ErrIPTooManySessions
	}

	sid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	if _, found := d.byID[sid]; found {
		panic(errors.New("UUID collision"))
	}

	s := &session{
		cumulativeTime: 0,
		ip:             ip,
		lastSeen:       now,
		registered:     now,
	}

	d.byIP[ip] = append(d.byIP[ip], s)
	d.byID[sid] = s

	return sid, nil
}

func (d *Datastore) Update(sid uuid.UUID, ct float64) error {
	now := time.Now()

	s, found := d.byID[sid]
	if !found {
		return ErrSessionNotRegistered
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cumulativeTime > 0 && now.Before(s.lastSeen.Add(minUpdateInterval)) {
		return ErrTooSoon
	}

	if ct < s.cumulativeTime {
		return ErrCumulativeTimeTooSmall
	}

	claimedDelta := ct - s.cumulativeTime
	observedDelta := now.Sub(s.lastSeen).Seconds()
	absDiff := claimedDelta - observedDelta
	relDiff := absDiff / observedDelta
	if absDiff > maxTimeAbsDiff || relDiff > maxTimeRelDiff {
		return ErrCumulativeTimeTooLarge
	}

	s.cumulativeTime = ct
	s.lastSeen = now

	return nil
}
