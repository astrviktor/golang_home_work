package internalhttp

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type ResponseEvent struct {
	Event storage.Event `json:"event"`
	Error string        `json:"error"`
}

type ResponseEventSlice struct {
	Events []storage.Event `json:"events"`
	Error  string          `json:"error"`
}

type ResponseID struct {
	ID    string `json:"id"`
	Error string `json:"error"`
}

type ResponseStatus struct {
	Status bool   `json:"status"`
	Error  string `json:"error"`
}

func (s *Server) WriteResponseEvent(w http.ResponseWriter, resp *ResponseEvent) {
	resBuf, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error(fmt.Sprintf("response marshal error: %s", err))
	}
	_, err = w.Write(resBuf)
	if err != nil {
		s.logger.Error(fmt.Sprintf("response marshal error: %s", err))
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func (s *Server) WriteResponseEventSlice(w http.ResponseWriter, resp *ResponseEventSlice) {
	resBuf, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error(fmt.Sprintf("response marshal error: %s", err))
	}
	_, err = w.Write(resBuf)
	if err != nil {
		s.logger.Error(fmt.Sprintf("response marshal error: %s", err))
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func (s *Server) WriteResponseID(w http.ResponseWriter, resp *ResponseID) {
	resBuf, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error(fmt.Sprintf("response marshal error: %s", err))
	}
	_, err = w.Write(resBuf)
	if err != nil {
		s.logger.Error(fmt.Sprintf("response marshal error: %s", err))
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func (s *Server) WriteResponseStatus(w http.ResponseWriter, resp *ResponseStatus) {
	resBuf, err := json.Marshal(resp)
	if err != nil {
		s.logger.Error(fmt.Sprintf("response marshal error: %s", err))
	}
	_, err = w.Write(resBuf)
	if err != nil {
		s.logger.Error(fmt.Sprintf("response marshal error: %s", err))
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}

func (s *Server) Event(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		s.GetEvent(w, r)
	}

	if r.Method == http.MethodPost {
		s.CreateEvent(w, r)
	}

	if r.Method == http.MethodPut {
		s.UpdateEvent(w, r)
	}

	if r.Method == http.MethodDelete {
		s.DeleteEvent(w, r)
	}
}

// curl --request GET 'http://127.0.0.1:8888/event?id=d9aed75b-3c9a-423b-8455-7ea824e9766e'

func (s *Server) GetEvent(w http.ResponseWriter, r *http.Request) {
	resp := &ResponseEvent{}
	args := r.URL.Query()
	id := args.Get("id")
	if len(id) == 0 {
		resp.Error = "id not found"
		w.WriteHeader(http.StatusBadRequest)
		s.WriteResponseEvent(w, resp)
		return
	}

	event, ok, err := s.storage.Get(id)
	if err != nil {
		resp.Error = fmt.Sprint(err)
		w.WriteHeader(http.StatusInternalServerError)
		s.WriteResponseEvent(w, resp)
		return
	}

	if !ok {
		w.WriteHeader(http.StatusOK)
		s.WriteResponseEvent(w, resp)
		return
	}

	resp.Event = event
	w.WriteHeader(http.StatusOK)
	s.WriteResponseEvent(w, resp)

	s.logger.Info(fmt.Sprintf("get event %#v", event))
}

/*
curl --request POST 'http://127.0.0.1:8888/event' \
--header 'Content-Type: application/json' \
--data-raw '{
"id": "0d59d804-bfe9-427f-ab37-cac59a0fbcd3",
"title": "pLnfgDsc2WD",
"dateStart": "2022-01-31T12:59:47+03:00",
"dateEnd": "2022-01-31T12:59:48+03:00",
"description": "123",
"usedId": 62,
"timeToNotification": 37
}'
*/

func (s *Server) CreateEvent(w http.ResponseWriter, r *http.Request) {
	resp := &ResponseID{}
	buf := make([]byte, r.ContentLength)
	_, err := r.Body.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		resp.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		s.WriteResponseID(w, resp)
		return
	}

	event := storage.Event{}
	err = json.Unmarshal(buf, &event)
	if err != nil {
		resp.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		s.WriteResponseID(w, resp)
		return
	}

	uuid, err := s.storage.Create(event)
	if err != nil {
		resp.Error = fmt.Sprint(err)
		w.WriteHeader(http.StatusInternalServerError)
		s.WriteResponseID(w, resp)
		return
	}

	resp.ID = uuid
	w.WriteHeader(http.StatusOK)
	s.WriteResponseID(w, resp)

	s.logger.Info(fmt.Sprintf("create new event %#v", event))
}

/*
curl --request PUT 'http://127.0.0.1:8888/event' \
--header 'Content-Type: application/json' \
--data-raw '{
"id": "ebcf4276-c189-4b98-a12e-99e7d668deba",
"title": "-----",
"dateStart": "2022-01-24T10:59:47+03:00",
"dateEnd": "2022-01-24T10:59:48+03:00",
"description": "12345",
"usedId": 62,
"timeToNotification": 37
}'
*/

func (s *Server) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	resp := &ResponseStatus{}

	buf := make([]byte, r.ContentLength)
	_, err := r.Body.Read(buf)
	if err != nil && !errors.Is(err, io.EOF) {
		resp.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		s.WriteResponseStatus(w, resp)
		return
	}

	event := storage.Event{}
	err = json.Unmarshal(buf, &event)
	if err != nil {
		resp.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		s.WriteResponseStatus(w, resp)
		return
	}

	ok, err := s.storage.Update(event)
	if err != nil {
		resp.Error = fmt.Sprint(err)
		w.WriteHeader(http.StatusInternalServerError)
		s.WriteResponseStatus(w, resp)
		return
	}

	if !ok {
		w.WriteHeader(http.StatusOK)
		s.WriteResponseStatus(w, resp)
		return
	}

	resp.Status = true
	w.WriteHeader(http.StatusOK)
	s.WriteResponseStatus(w, resp)

	s.logger.Info(fmt.Sprintf("update event %#v", event))
}

// curl --request DELETE 'http://127.0.0.1:8888/event?id=d9aed75b-3c9a-423b-8455-7ea824e9766e'

func (s *Server) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	resp := &ResponseStatus{}
	args := r.URL.Query()
	id := args.Get("id")
	if len(id) == 0 {
		resp.Error = "id not found"
		w.WriteHeader(http.StatusBadRequest)
		s.WriteResponseStatus(w, resp)
		return
	}

	ok, err := s.storage.Delete(id)
	if err != nil {
		resp.Error = fmt.Sprint(err)
		w.WriteHeader(http.StatusInternalServerError)
		s.WriteResponseStatus(w, resp)
		return
	}

	if !ok {
		w.WriteHeader(http.StatusOK)
		s.WriteResponseStatus(w, resp)
		return
	}

	resp.Status = true
	w.WriteHeader(http.StatusOK)
	s.WriteResponseStatus(w, resp)

	s.logger.Info(fmt.Sprintf("delete event %s", id))
}

type Period int32

const (
	PeriodDay   Period = 1
	PeriodWeek  Period = 2
	PeriodMonth Period = 3
)

func (s *Server) GetList(w http.ResponseWriter, r *http.Request, period Period) {
	resp := &ResponseEventSlice{}
	args := r.URL.Query()
	str := args.Get("date")
	date, err := time.Parse("2006-01-02", str)
	if err != nil {
		resp.Error = "date format error"
		w.WriteHeader(http.StatusBadRequest)
		s.WriteResponseEventSlice(w, resp)
		return
	}

	var events []storage.Event

	switch period {
	case PeriodDay:
		events, err = s.storage.EventListDay(date)
	case PeriodWeek:
		events, err = s.storage.EventListWeek(date)
	case PeriodMonth:
		events, err = s.storage.EventListMonth(date)
	}

	if err != nil {
		resp.Error = fmt.Sprint(err)
		w.WriteHeader(http.StatusInternalServerError)
		s.WriteResponseEventSlice(w, resp)
		return
	}

	resp.Events = events
	w.WriteHeader(http.StatusOK)
	s.WriteResponseEventSlice(w, resp)

	s.logger.Info(fmt.Sprintf("get events list %#v", events))
}

// curl --request GET 'http://127.0.0.1:8888/list/day?date=2022-01-25'

func (s *Server) GetListDay(w http.ResponseWriter, r *http.Request) {
	s.GetList(w, r, PeriodDay)
}

// curl --request GET 'http://127.0.0.1:8888/list/week?date=2022-01-24'

func (s *Server) GetListWeek(w http.ResponseWriter, r *http.Request) {
	s.GetList(w, r, PeriodWeek)
}

// curl --request GET 'http://127.0.0.1:8888/list/month?date=2022-01-01'

func (s *Server) GetListMonth(w http.ResponseWriter, r *http.Request) {
	s.GetList(w, r, PeriodMonth)
}
