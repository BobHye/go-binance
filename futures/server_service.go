package futures

import (
	"context"
	"net/http"
)

// PingService ping server
type PingService struct {
	c *Client
}

// Do send request
func (s *PingService) Do(ctx context.Context, opts ...RequestOption) (err error) {
	// GET /fapi/v1/ping | 测试能否联通
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/ping",
	}
	_, _, err = s.c.callAPI(ctx, r, opts...)
	return err
}

// ServerTimeService get server time
type ServerTimeService struct {
	c *Client
}

// Do send request
func (s *ServerTimeService) Do(ctx context.Context, opts ...RequestOption) (int64, error) {
	// GET /fapi/v1/time | 获取服务器时间
	r := &request{
		method:   http.MethodGet,
		endpoint: "/fapi/v1/time",
	}
	data, _, err := s.c.callAPI(ctx, r, opts...)
	if err != nil {
		return 0, err
	}
	var res ServerTime
	err = json.Unmarshal(data, &res)
	return res.ServerTime, err
}

type ServerTime struct {
	ServerTime int64 `json:"serverTime"` // 服务器当前的系统时间
}

// SetServerTimeService set server time
type SetServerTimeService struct {
	c *Client
}

// Do send request
func (s *SetServerTimeService) Do(ctx context.Context, opts ...RequestOption) (timeOffset int64, err error) {
	serverTime, err := s.c.NewServerTimeService().Do(ctx, opts...)
	if err != nil {
		return 0, err
	}
	timeOffset = currentTimestamp() - serverTime
	s.c.TimeOffset = timeOffset
	return timeOffset, nil
}
