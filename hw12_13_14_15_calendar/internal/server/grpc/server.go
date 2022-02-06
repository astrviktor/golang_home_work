package internalgrpc

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/app"
	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	logger  app.Logger
	app     Application
	storage storage.Storage
	addr    string
	pb.UnimplementedCalendarServer
}

type Application interface { // TODO
}

func NewServer(logger app.Logger, app Application, storage storage.Storage, host string, port string) *Server {
	return &Server{
		logger:  logger,
		app:     app,
		storage: storage,
		addr:    net.JoinHostPort(host, port),
	}
}

func pbEventToStorageEvent(pbEvent *pb.Event) storage.Event {
	event := storage.Event{
		ID:                 pbEvent.Id,
		Title:              pbEvent.Title,
		DateStart:          pbEvent.DateStart.AsTime(),
		DateEnd:            pbEvent.DateEnd.AsTime(),
		Description:        pbEvent.Description,
		UserID:             int(pbEvent.UsedId),
		TimeToNotification: int(pbEvent.TimeToNotification),
		Notified:           pbEvent.Notified,
	}

	return event
}

func storageEventToPbEvent(event storage.Event) *pb.Event {
	pbEvent := pb.Event{}

	pbEvent.Id = event.ID
	pbEvent.Title = event.Title
	pbEvent.DateStart = timestamppb.New(event.DateStart)
	pbEvent.DateEnd = timestamppb.New(event.DateEnd)
	pbEvent.Description = event.Description
	pbEvent.UsedId = uint32(event.UserID)
	pbEvent.TimeToNotification = uint32(event.TimeToNotification)
	pbEvent.Notified = event.Notified

	return &pbEvent
}

func (s *Server) CreateEvent(ctx context.Context, req *pb.Event) (*pb.ResponseID, error) {
	event := pbEventToStorageEvent(req)

	uuid, err := s.storage.Create(event)
	if err != nil {
		return nil, status.Error(codes.Internal, "error on create event")
	}

	return &pb.ResponseID{Id: uuid}, nil
}

func (s *Server) GetEvent(ctx context.Context, req *pb.ID) (*pb.ResponseEvent, error) {
	event, ok, err := s.storage.Get(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "error on get event")
	}

	if !ok {
		return &pb.ResponseEvent{Event: &pb.Event{}}, errors.New("not found")
	}

	return &pb.ResponseEvent{Event: storageEventToPbEvent(event)}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *pb.Event) (*pb.ResponseStatus, error) {
	event := pbEventToStorageEvent(req)

	ok, err := s.storage.Update(event)
	if err != nil {
		return nil, status.Error(codes.Internal, "error on update event")
	}

	if !ok {
		return &pb.ResponseStatus{Status: false}, errors.New("not found")
	}

	return &pb.ResponseStatus{Status: true}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *pb.ID) (*pb.ResponseStatus, error) {
	ok, err := s.storage.Delete(req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, "error on delete event")
	}

	if !ok {
		return &pb.ResponseStatus{Status: false}, errors.New("not found")
	}

	return &pb.ResponseStatus{Status: true}, nil
}

func (s *Server) GetList(ctx context.Context, req *pb.Date) (*pb.ResponseEventSlice, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "date format error")
	}

	var events []storage.Event

	switch req.Period {
	case pb.Period_PERIOD_UNSPECIFIED:
		return nil, status.Error(codes.InvalidArgument, "unspecified period")
	case pb.Period_PERIOD_DAY:
		events, err = s.storage.EventListDay(date)
	case pb.Period_PERIOD_WEEK:
		events, err = s.storage.EventListWeek(date)
	case pb.Period_PERIOD_MONTH:
		events, err = s.storage.EventListMonth(date)
	}

	if err != nil {
		return nil, status.Error(codes.Internal, "error on get events")
	}

	pbEvents := make([]*pb.Event, len(events))
	for idx, event := range events {
		pbEvents[idx] = storageEventToPbEvent(event)
	}

	return &pb.ResponseEventSlice{Events: pbEvents}, nil
}

func (s *Server) Start(ctx context.Context) {
	lsn, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(s.WithLoggingInterceptor),
	)

	pb.RegisterCalendarServer(grpcServer, s)
	s.logger.Info("gprs server starting on address " + s.addr)

	go func() {
		s.logger.Fatal(grpcServer.Serve(lsn).Error())
	}()

	<-ctx.Done()
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

// evans --host 0.0.0.0 --port 9999 --path ./api --proto=calendar.proto
