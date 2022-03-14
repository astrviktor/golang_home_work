package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/astrviktor/golang_home_work/hw12_13_14_15_calendar/internal/storage"
)

type Client struct {
	addr    string
	timeout time.Duration
}

func NewClient(host string, port string, timeout time.Duration) *Client {
	return &Client{net.JoinHostPort(host, port), timeout}
}

func (c *Client) GetEvent(id string) (storage.Event, error) {
	url := "http://" + c.addr + "/event?id=" + id
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return storage.Event{}, err
	}

	client := &http.Client{Timeout: c.timeout}

	resp, err := client.Do(req)
	if err != nil {
		return storage.Event{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return storage.Event{}, err
	}

	respEvent := ResponseEvent{}
	err = json.Unmarshal(respBody, &respEvent)
	if err != nil {
		return storage.Event{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return storage.Event{}, errors.New("error while getting event")
	}
	return respEvent.Event, nil
}

func (c *Client) CreateEvent(event storage.Event) (string, error) {
	b, err := json.Marshal(event)
	if err != nil {
		return "", err
	}

	body := bytes.NewReader(b)

	url := "http://" + c.addr + "/event"

	req, err := http.NewRequestWithContext(context.Background(), "POST", url, body)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: c.timeout}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	responseID := ResponseID{}
	err = json.Unmarshal(respBody, &responseID)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("error while creating event")
	}
	return responseID.ID, nil
}

func (c *Client) UpdateEvent(event storage.Event) (bool, error) {
	b, err := json.Marshal(event)
	if err != nil {
		return false, err
	}

	body := bytes.NewReader(b)

	url := "http://" + c.addr + "/event"

	req, err := http.NewRequestWithContext(context.Background(), "PUT", url, body)
	if err != nil {
		return false, err
	}

	client := &http.Client{Timeout: c.timeout}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	responseStatus := ResponseStatus{}
	err = json.Unmarshal(respBody, &responseStatus)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		return false, errors.New("error while updating event")
	}
	return responseStatus.Status, nil
}

func (c *Client) DeleteEvent(id string) (bool, error) {
	url := "http://" + c.addr + "/event?id=" + id
	req, err := http.NewRequestWithContext(context.Background(), "DELETE", url, nil)
	if err != nil {
		return false, err
	}

	client := &http.Client{Timeout: c.timeout}

	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	responseStatus := ResponseStatus{}
	err = json.Unmarshal(respBody, &responseStatus)
	if err != nil {
		return false, err
	}

	if resp.StatusCode != http.StatusOK {
		return false, errors.New("error while deleting event")
	}
	return responseStatus.Status, nil
}

func (c *Client) GetListDay(date string) ([]storage.Event, error) {
	return c.GetListPeriod(date, PeriodDay)
}

func (c *Client) GetListWeek(date string) ([]storage.Event, error) {
	return c.GetListPeriod(date, PeriodWeek)
}

func (c *Client) GetListMonth(date string) ([]storage.Event, error) {
	return c.GetListPeriod(date, PeriodMonth)
}

func (c *Client) GetListPeriod(date string, period Period) ([]storage.Event, error) {
	var url string

	switch period {
	case PeriodDay:
		url = "http://" + c.addr + "/list/day?date=" + date
	case PeriodWeek:
		url = "http://" + c.addr + "/list/week?date=" + date
	case PeriodMonth:
		url = "http://" + c.addr + "/list/month?date=" + date
	}

	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: c.timeout}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	respEventSlice := ResponseEventSlice{}
	err = json.Unmarshal(respBody, &respEventSlice)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error while getting list")
	}
	return respEventSlice.Events, nil
}
