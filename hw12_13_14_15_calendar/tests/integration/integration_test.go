//go:build integration
// +build integration

package integration_test

import (
	"testing"
	"time"

	internalgrpc "github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/server/http"
	generate "github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/suite"
)

const (
	HTTPServiceHost = "127.0.0.1"
	HTTPServicePort = "8888"
	GRPCServiceHost = "127.0.0.1"
	GRPCServicePort = "9999"
)

type CalendarSuite struct {
	suite.Suite
	httpClient *internalhttp.Client
	grpcClient *internalgrpc.Client
}

func (s *CalendarSuite) SetupSuite() {
	time.Sleep(30 * time.Second)
	httpClient := internalhttp.NewClient(HTTPServiceHost, HTTPServicePort, time.Second)
	grpcClient, err := internalgrpc.NewClient(GRPCServiceHost, GRPCServicePort)
	s.Require().NoError(err)

	s.httpClient = httpClient
	s.grpcClient = grpcClient
}

func (s *CalendarSuite) SetupTest() {
}

func (s *CalendarSuite) TestHTTPCreateEventAndGetEvent() {
	date, err := time.Parse("2006-01-02", "2022-01-01")
	s.Require().NoError(err)

	newEvent := generate.GenerateEventDate(date, date.Add(time.Second))

	id, err := s.httpClient.CreateEvent(newEvent)
	s.Require().NoError(err)
	s.Require().Len(id, 36)
	newEvent.ID = id

	getEvent, err := s.httpClient.GetEvent(id)
	s.Require().NoError(err)

	s.Require().Equal(newEvent, getEvent)

	_, err = s.httpClient.DeleteEvent(id)
	s.Require().NoError(err)
}

func (s *CalendarSuite) TestGRPCCreateEventAndGetEvent() {
	date, err := time.Parse("2006-01-02", "2022-01-02")
	s.Require().NoError(err)

	newEvent := generate.GenerateEventDate(date, date.Add(time.Second))

	id, err := s.grpcClient.CreateEvent(newEvent)
	s.Require().NoError(err)
	s.Require().Len(id, 36)
	newEvent.ID = id

	getEvent, err := s.grpcClient.GetEvent(id)
	s.Require().NoError(err)

	s.Require().Equal(newEvent, getEvent)

	_, err = s.grpcClient.DeleteEvent(id)
	s.Require().NoError(err)
}

func (s *CalendarSuite) TestHTTPCreateEventAndUpdateEvent() {
	date, err := time.Parse("2006-01-02", "2022-01-03")
	s.Require().NoError(err)

	newEvent := generate.GenerateEventDate(date, date.Add(time.Second))

	id, err := s.httpClient.CreateEvent(newEvent)
	s.Require().NoError(err)
	s.Require().Len(id, 36)
	newEvent.ID = id

	updEvent := generate.GenerateEventDate(date.Add(2*time.Second), date.Add(3*time.Second))
	updEvent.ID = id
	ok, err := s.httpClient.UpdateEvent(updEvent)
	s.Require().NoError(err)
	s.Require().True(ok)

	getEvent, err := s.httpClient.GetEvent(id)
	s.Require().NoError(err)

	s.Require().Equal(updEvent, getEvent)

	_, err = s.httpClient.DeleteEvent(id)
	s.Require().NoError(err)
}

func (s *CalendarSuite) TestGRPCCreateEventAndUpdateEvent() {
	date, err := time.Parse("2006-01-02", "2022-01-04")
	s.Require().NoError(err)

	newEvent := generate.GenerateEventDate(date, date.Add(time.Second))

	id, err := s.grpcClient.CreateEvent(newEvent)
	s.Require().NoError(err)
	s.Require().Len(id, 36)
	newEvent.ID = id

	updEvent := generate.GenerateEventDate(date.Add(2*time.Second), date.Add(3*time.Second))
	updEvent.ID = id
	ok, err := s.grpcClient.UpdateEvent(updEvent)
	s.Require().NoError(err)
	s.Require().True(ok)

	getEvent, err := s.grpcClient.GetEvent(id)
	s.Require().NoError(err)

	s.Require().Equal(updEvent, getEvent)

	_, err = s.grpcClient.DeleteEvent(id)
	s.Require().NoError(err)
}

func (s *CalendarSuite) TestHTTPCreateEventAndDeleteEvent() {
	date, err := time.Parse("2006-01-02", "2022-01-05")
	s.Require().NoError(err)

	newEvent := generate.GenerateEventDate(date, date.Add(time.Second))

	id, err := s.httpClient.CreateEvent(newEvent)
	s.Require().NoError(err)
	s.Require().Len(id, 36)
	newEvent.ID = id

	_, err = s.httpClient.DeleteEvent(id)
	s.Require().NoError(err)

	_, err = s.httpClient.GetEvent(id)
	s.Require().Error(err)
}

func (s *CalendarSuite) TestGRPCCreateEventAndDeleteEvent() {
	date, err := time.Parse("2006-01-02", "2022-01-06")
	s.Require().NoError(err)

	newEvent := generate.GenerateEventDate(date, date.Add(time.Second))

	id, err := s.grpcClient.CreateEvent(newEvent)
	s.Require().NoError(err)
	s.Require().Len(id, 36)
	newEvent.ID = id

	_, err = s.grpcClient.DeleteEvent(id)
	s.Require().NoError(err)

	_, err = s.grpcClient.GetEvent(id)
	s.Require().Error(err)
}

func (s *CalendarSuite) TestHTTPCreateAndListDayWeekMonth() {
	dateA, err := time.Parse("2006-01-02", "2022-02-07")
	s.Require().NoError(err)
	dateB, err := time.Parse("2006-01-02", "2022-02-08")
	s.Require().NoError(err)
	dateC, err := time.Parse("2006-01-02", "2022-02-14")
	s.Require().NoError(err)

	eventA := generate.GenerateEventDate(dateA, dateA.Add(time.Hour))
	id, err := s.httpClient.CreateEvent(eventA)
	s.Require().NoError(err)
	s.Require().Len(id, 36)
	eventA.ID = id

	eventB := generate.GenerateEventDate(dateB, dateB.Add(time.Hour))
	id, err = s.httpClient.CreateEvent(eventB)
	s.Require().NoError(err)
	s.Require().Len(id, 36)
	eventB.ID = id

	eventC := generate.GenerateEventDate(dateC, dateC.Add(time.Hour))
	id, err = s.httpClient.CreateEvent(eventC)
	s.Require().NoError(err)
	s.Require().Len(id, 36)
	eventC.ID = id

	events, err := s.httpClient.GetListDay("2022-02-08")
	s.Require().NoError(err)
	s.Require().Equal(1, len(events))

	events, err = s.httpClient.GetListWeek("2022-02-07")
	s.Require().NoError(err)
	s.Require().Equal(2, len(events))

	events, err = s.httpClient.GetListMonth("2022-02-01")
	s.Require().NoError(err)
	s.Require().Equal(3, len(events))

	_, err = s.httpClient.DeleteEvent(eventA.ID)
	s.Require().NoError(err)

	_, err = s.httpClient.DeleteEvent(eventB.ID)
	s.Require().NoError(err)

	_, err = s.httpClient.DeleteEvent(eventC.ID)
	s.Require().NoError(err)
}

func (s *CalendarSuite) TestGRPCCreateAndListDayWeekMonth() {
	dateA, err := time.Parse("2006-01-02", "2022-03-07")
	s.Require().NoError(err)
	dateB, err := time.Parse("2006-01-02", "2022-03-08")
	s.Require().NoError(err)
	dateC, err := time.Parse("2006-01-02", "2022-03-14")
	s.Require().NoError(err)

	eventA := generate.GenerateEventDate(dateA, dateA.Add(time.Hour))
	id, err := s.grpcClient.CreateEvent(eventA)
	s.Require().NoError(err)
	s.Require().Len(id, 36)
	eventA.ID = id

	eventB := generate.GenerateEventDate(dateB, dateB.Add(time.Hour))
	id, err = s.grpcClient.CreateEvent(eventB)
	s.Require().NoError(err)
	s.Require().Len(id, 36)
	eventB.ID = id

	eventC := generate.GenerateEventDate(dateC, dateC.Add(time.Hour))
	id, err = s.grpcClient.CreateEvent(eventC)
	s.Require().NoError(err)
	s.Require().Len(id, 36)
	eventC.ID = id

	events, err := s.grpcClient.GetListDay("2022-03-08")
	s.Require().NoError(err)
	s.Require().Equal(1, len(events))

	events, err = s.grpcClient.GetListWeek("2022-03-07")
	s.Require().NoError(err)
	s.Require().Equal(2, len(events))

	events, err = s.grpcClient.GetListMonth("2022-03-01")
	s.Require().NoError(err)
	s.Require().Equal(3, len(events))

	_, err = s.grpcClient.DeleteEvent(eventA.ID)
	s.Require().NoError(err)

	_, err = s.grpcClient.DeleteEvent(eventB.ID)
	s.Require().NoError(err)

	_, err = s.grpcClient.DeleteEvent(eventC.ID)
	s.Require().NoError(err)
}

func (s *CalendarSuite) TestHTTPCreate2EventAndGetError() {
	date, err := time.Parse("2006-01-02", "2022-01-07")
	s.Require().NoError(err)

	newEvent := generate.GenerateEventDate(date, date.Add(time.Second))

	id, err := s.httpClient.CreateEvent(newEvent)
	s.Require().NoError(err)
	s.Require().Len(id, 36)
	newEvent.ID = id

	_, err = s.httpClient.CreateEvent(newEvent)
	s.Require().Error(err)

	_, err = s.httpClient.DeleteEvent(id)
	s.Require().NoError(err)
}

func (s *CalendarSuite) TestGRPCCreate2EventAndGetError() {
	date, err := time.Parse("2006-01-02", "2022-01-08")
	s.Require().NoError(err)

	newEvent := generate.GenerateEventDate(date, date.Add(time.Second))

	id, err := s.grpcClient.CreateEvent(newEvent)
	s.Require().NoError(err)
	s.Require().Len(id, 36)
	newEvent.ID = id

	_, err = s.grpcClient.CreateEvent(newEvent)
	s.Require().Error(err)

	_, err = s.grpcClient.DeleteEvent(id)
	s.Require().NoError(err)
}

func (s *CalendarSuite) TestHTTPCreateEventAndSendNotification() {
	date, err := time.Parse("2006-01-02", "2022-01-09")
	s.Require().NoError(err)

	newEvent := generate.GenerateEventDate(date, date.Add(time.Second))

	id, err := s.httpClient.CreateEvent(newEvent)
	s.Require().NoError(err)
	s.Require().Len(id, 36)
	newEvent.ID = id

	getEventBefore, err := s.httpClient.GetEvent(id)
	s.Require().NoError(err)
	s.Require().Equal("no", getEventBefore.Notified)

	time.Sleep(30 * time.Second)

	getEventAfter, err := s.httpClient.GetEvent(id)
	s.Require().NoError(err)
	s.Require().Equal("yes", getEventAfter.Notified)

	_, err = s.httpClient.DeleteEvent(id)
	s.Require().NoError(err)
}

func (s *CalendarSuite) TestGRPCCreateEventAndSendNotification() {
	date, err := time.Parse("2006-01-02", "2022-01-10")
	s.Require().NoError(err)

	newEvent := generate.GenerateEventDate(date, date.Add(time.Second))

	id, err := s.grpcClient.CreateEvent(newEvent)
	s.Require().NoError(err)
	s.Require().Len(id, 36)
	newEvent.ID = id

	getEventBefore, err := s.grpcClient.GetEvent(id)
	s.Require().NoError(err)
	s.Require().Equal("no", getEventBefore.Notified)

	time.Sleep(30 * time.Second)

	getEventAfter, err := s.grpcClient.GetEvent(id)
	s.Require().NoError(err)
	s.Require().Equal("yes", getEventAfter.Notified)

	_, err = s.grpcClient.DeleteEvent(id)
	s.Require().NoError(err)
}

func (s *CalendarSuite) TearDownTest() {
}

func (s *CalendarSuite) TearDownSuite() {
}

func TestCalendarSuite(t *testing.T) {
	suite.Run(t, new(CalendarSuite))
}
