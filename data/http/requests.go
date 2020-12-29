package http

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/simplego/encoding/csvutil"
	"github.com/grokify/simplego/type/stringsutil/join"
)

// Endpoints writes a CSV with request data. Use Endpoints.Add(),
// Endpoints.Inflate() and then Endpoints.WriteCSV()
type Endpoints struct {
	EndpointsMap map[string]Endpoint
}

func NewEndpoints() Endpoints {
	return Endpoints{
		EndpointsMap: map[string]Endpoint{}}
}

// Add. time is optional.
func (eps *Endpoints) Add(method, url string, status int, subStatus string, dt time.Time) {
	method = strings.ToUpper(strings.TrimSpace(method))
	url = strings.TrimSpace(url)
	endpoint := method + " " + url
	ep, ok := eps.EndpointsMap[endpoint]
	if !ok {
		ep = Endpoint{
			Method:   method,
			URL:      url,
			Statuses: Statuses{StatusMap: map[string]StatusInfo{}},
		}
	}
	ep.AddStatus(status, subStatus, dt)
	eps.EndpointsMap[endpoint] = ep
}

func EndpointFromMethodAndURL(method, url string) string {
	parts := []string{}
	parts = append(parts, strings.ToUpper(strings.TrimSpace(method)))
	parts = append(parts, strings.TrimSpace(url))
	return strings.Join(parts, " ")
}

func (eps *Endpoints) Inflate() {
	for endpointUID, endpoint := range eps.EndpointsMap {
		endpoint.Statuses.Inflate()
		eps.EndpointsMap[endpointUID] = endpoint
	}
}

func (eps *Endpoints) AllFullStatusCodes() []string {
	codes := map[string]int{}
	for _, ep := range eps.EndpointsMap {
		for code := range ep.Statuses.StatusMap {
			if _, ok := codes[code]; !ok {
				codes[code] = 0
			}
			codes[code] += 1
		}
	}
	codesSlice := []string{}
	for code := range codes {
		codesSlice = append(codesSlice, code)
	}
	sort.Strings(codesSlice)
	return codesSlice
}

// WriteCSV outputs a summary status table with a response status
// distribution for URL requests.
func (eps *Endpoints) WriteCSV(filename string) error {
	codes := eps.AllFullStatusCodes()
	codesStr := codes
	//codesStr := SliceIntToString(codes)
	w, f, err := csvutil.NewWriterFile(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	defer w.Flush()
	header := []string{"http_method", "request_uri", "request_count"}
	for _, codeStr := range codesStr {
		header = append(header, codeStr+"_%")
	}
	header = append(header, codesStr...)
	err = w.Write(header)
	if err != nil {
		return err
	}
	for _, epInfo := range eps.EndpointsMap {
		row := []string{
			epInfo.Method,
			epInfo.URL,
			strconv.Itoa(epInfo.Statuses.AllRequestCount())}
		for _, code := range codes {
			if statusCode, ok := epInfo.Statuses.StatusMap[code]; ok {
				row = append(row, strconv.FormatFloat(statusCode.StatusDistribution, 'E', -1, 64))
			} else {
				row = append(row, "0")
			}
		}
		for _, code := range codes {
			if statusCode, ok := epInfo.Statuses.StatusMap[code]; ok {
				row = append(row, strconv.Itoa(statusCode.RequestCount))
			} else {
				row = append(row, "0")
			}
		}
		err := w.Write(row)
		if err != nil {
			return err
		}
	}
	w.Flush()
	f.Close()
	return nil
}

type Endpoint struct {
	Method   string
	URL      string
	Statuses Statuses
	//StatusTimes StatusTimes
}

func (ep *Endpoint) AddStatus(status int, subStatus string, dt time.Time) {
	fullStatus := StatusPartsToFullStatus(status, subStatus)
	if ep.Statuses.StatusMap == nil {
		ep.Statuses = NewStatuses()
	}
	if _, ok := ep.Statuses.StatusMap[fullStatus]; !ok {
		ep.Statuses.StatusMap[fullStatus] = StatusInfo{Times: []time.Time{}}
	}
	statusInfo := ep.Statuses.StatusMap[fullStatus]
	statusInfo.Status = status
	statusInfo.SubStatus = subStatus
	statusInfo.RequestCount += 1
	statusInfo.Times = append(statusInfo.Times, dt)
	ep.Statuses.StatusMap[fullStatus] = statusInfo
}

func StatusPartsToFullStatus(status int, subStatus string) string {
	parts := []string{strconv.Itoa(status)}
	subStatus = strings.TrimSpace(subStatus)
	if len(subStatus) > 0 {
		parts = append(parts, subStatus)
	}
	return strings.Join(parts, " ")
}

func SliceIntToString(s []int) []string {
	s2 := []string{}
	for _, i := range s {
		s2 = append(s2, strconv.Itoa(i))
	}
	return s2
}

type Statuses struct {
	StatusMap map[string]StatusInfo
}

func NewStatuses() Statuses {
	return Statuses{StatusMap: map[string]StatusInfo{}}
}

func (st *Statuses) AllRequestCount() int {
	allCount := 0
	for _, status := range st.StatusMap {
		allCount += status.RequestCount
	}
	return allCount
}

func (st *Statuses) Inflate() {
	allCount := st.AllRequestCount()
	for fullStatus, info := range st.StatusMap {
		thisCount := info.RequestCount
		dist := float64(thisCount) / float64(allCount)
		info.StatusDistribution = dist
		//info.StatusText = http.StatusText(info.Status)
		st.StatusMap[fullStatus] = info
	}
}

type StatusInfo struct {
	Status             int
	SubStatus          string
	RequestCount       int
	StatusDistribution float64     // This status vs. all statuses
	Times              []time.Time // This is optional. Can be used by manually adding to RequestCount
}

func (si *StatusInfo) Inflate() {
	// Set Request Count to times only if times is populated.
	// This can be used without Times being populated.
	if len(si.Times) != 0 && si.RequestCount != len(si.Times) {
		si.RequestCount = len(si.Times)
	}
}

func (si *StatusInfo) StatusText() string {
	return http.StatusText(si.Status)
}

func (si *StatusInfo) FullStatus() string {
	return join.JoinCondenseTrimSpace([]string{strconv.Itoa(si.Status), si.SubStatus}, " ")
}

func StatusMapToTimesArray(statuses map[string]StatusInfo) []StatusTime {
	//times := []StatusTime{}
	times := StatusTimeSlice{}
	for _, si := range statuses {
		for _, t := range si.Times {
			st := StatusTime{
				Status: si.Status,
				Time:   t,
			}
			times = append(times, st)
		}
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i].Time.Before(times[j].Time)
	})
	return times
}

// Endpoint Chart that shows errors by hour or day

//func EndpointTo

type StatusTimeSlice []StatusTime

type StatusTime struct {
	Method     string
	RequestURL string
	Status     int
	Time       time.Time
}
