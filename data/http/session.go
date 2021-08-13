package http

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/simplego/type/stringsutil"
)

type SessionSet struct {
	SessionMap map[string]Session
}

func NewSessionSet() SessionSet {
	return SessionSet{SessionMap: map[string]Session{}}
}

func (ss *SessionSet) AddRequest(r Request) {
	sessionID := strings.TrimSpace(r.SessionID)
	ses, ok := ss.SessionMap[sessionID]
	if !ok {
		ses = NewSession(sessionID)
	}
	ses.RequestSet.Requests = append(ses.RequestSet.Requests, r)
	ss.SessionMap[sessionID] = ses
}

type Session struct {
	SessionId  string
	RequestSet RequestSet
}

func NewSession(sessionId string) Session {
	return Session{
		SessionId:  sessionId,
		RequestSet: NewRequestSet()}
}

func (ses *Session) RequestsByEndpoint() map[string]RequestSet {
	byEndpoint := map[string]RequestSet{}
	for _, req := range ses.RequestSet.Requests {
		epString := req.Endpoint()
		if _, ok := byEndpoint[epString]; !ok {
			byEndpoint[epString] = NewRequestSet()
		}
		rs := byEndpoint[epString]
		rs.Requests = append(rs.Requests, req)
	}
	for ep, rs := range byEndpoint {
		rs.Inflate()
		byEndpoint[ep] = rs
	}
	return byEndpoint
}

type RequestSet struct {
	LastStatusCode int
	Requests       []Request
}

func NewRequestSet() RequestSet {
	return RequestSet{
		LastStatusCode: -1,
		Requests:       []Request{}}
}

func (rs *RequestSet) Inflate() {
	rs.LastStatusCode = rs.BuildLastStatusCode()
}

func (rs *RequestSet) SortByTime() {
	reqs := rs.Requests
	sort.Slice(
		reqs,
		func(i, j int) bool { return reqs[i].Time.Before(reqs[j].Time) })
}

func (rs *RequestSet) BuildLastStatusCode() int {
	if len(rs.Requests) == 0 {
		return -1
	}
	rs.SortByTime()
	lastReq := rs.Requests[len(rs.Requests)-1]
	return lastReq.StatusCode
}

type Request struct {
	Method        string
	URL           string
	URLPattern    string
	StatusCode    int
	SubStatusCode string
	Time          time.Time
	OperationID   string // OpenAPI spec OperationID
	SessionID     string // native SessionID
	RequestID     string // native RequestID, typically UID
	RawData       interface{}
	RawDataString string
}

func (req *Request) Endpoint() string {
	req.Method = strings.TrimSpace(strings.ToUpper(req.Method))
	req.URL = strings.TrimSpace(req.URL)
	parts := []string{}
	if len(req.Method) > 0 {
		parts = append(parts, req.Method)
	}
	if len(req.URL) > 0 {
		parts = append(parts, req.URL)
	}
	return strings.Join(parts, " ")
}

func (req *Request) FullStatus() string {
	return strings.Join(
		stringsutil.SliceCondenseSpace([]string{strconv.Itoa(req.StatusCode), req.SubStatusCode}, false, false),
		" ")
	// join.JoinCondenseTrimSpace([]string{strconv.Itoa(req.StatusCode), req.SubStatusCode}, " ")
}
