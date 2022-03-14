package internalgrpc

import (
	"context"
	"net"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	client pb.CalendarClient
}

func NewClient(host string, port string) (*Client, error) {
	conn, err := grpc.Dial(net.JoinHostPort(host, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewCalendarClient(conn)

	return &Client{client: client}, nil
}

func (c *Client) CreateEvent(event storage.Event) (string, error) {
	pbEvent := storageEventToPbEvent(event)

	responseID, err := c.client.CreateEvent(context.Background(), pbEvent)
	if err != nil {
		return "", err
	}

	return responseID.Id, nil
}

func (c *Client) GetEvent(id string) (storage.Event, error) {
	responseEvent, err := c.client.GetEvent(context.Background(), &pb.ID{Id: id})
	if err != nil {
		return storage.Event{}, err
	}

	event := pbEventToStorageEvent(responseEvent.Event)

	return event, nil
}

func (c *Client) UpdateEvent(event storage.Event) (bool, error) {
	pbEvent := storageEventToPbEvent(event)

	responseStatus, err := c.client.UpdateEvent(context.Background(), pbEvent)
	if err != nil {
		return false, err
	}

	return responseStatus.Status, nil
}

func (c *Client) DeleteEvent(id string) (bool, error) {
	responseStatus, err := c.client.DeleteEvent(context.Background(), &pb.ID{Id: id})
	if err != nil {
		return false, err
	}

	return responseStatus.Status, nil
}

func (c *Client) GetListDay(date string) ([]storage.Event, error) {
	return c.GetListPeriod(&pb.Date{Date: date, Period: pb.Period_PERIOD_DAY})
}

func (c *Client) GetListWeek(date string) ([]storage.Event, error) {
	return c.GetListPeriod(&pb.Date{Date: date, Period: pb.Period_PERIOD_WEEK})
}

func (c *Client) GetListMonth(date string) ([]storage.Event, error) {
	return c.GetListPeriod(&pb.Date{Date: date, Period: pb.Period_PERIOD_MONTH})
}

func (c *Client) GetListPeriod(req *pb.Date) ([]storage.Event, error) {
	responseEventSlice, err := c.client.GetList(context.Background(), req)
	if err != nil {
		return nil, err
	}

	events := make([]storage.Event, len(responseEventSlice.Events))

	for idx, pbEvent := range responseEventSlice.Events {
		events[idx] = pbEventToStorageEvent(pbEvent)
	}

	return events, nil
}
