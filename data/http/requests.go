package http

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/grokify/gotilla/encoding/csvutil"
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

func (eps *Endpoints) Add(method, url string, status int, subStatus string) {
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
	ep.AddStatus(status, subStatus)
	eps.EndpointsMap[endpoint] = ep
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
		for code, _ := range ep.Statuses.StatusMap {
			if _, ok := codes[code]; !ok {
				codes[code] = 0
			}
			codes[code] += 1
		}
	}
	codesSlice := []string{}
	for code, _ := range codes {
		codesSlice = append(codesSlice, code)
	}
	sort.Strings(codesSlice)
	return codesSlice
}

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
}

func (ep *Endpoint) AddStatus(status int, subStatus string) {
	fullStatus := StatusPartsToFullStatus(status, subStatus)
	if _, ok := ep.Statuses.StatusMap[fullStatus]; !ok {
		ep.Statuses.StatusMap[fullStatus] = StatusInfo{}
	}
	statusInfo := ep.Statuses.StatusMap[fullStatus]
	statusInfo.Status = status
	statusInfo.SubStatus = subStatus
	statusInfo.RequestCount += 1
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
	StatusDistribution float64
}

func (si *StatusInfo) StatusText() string {
	return http.StatusText(si.Status)
}

func (si *StatusInfo) FullStatus() string {
	return StatusPartsToFullStatus(si.Status, si.SubStatus)
}
