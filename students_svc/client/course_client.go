package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/ratelimit"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type Course struct {
	ID   int64  `json:"id,omitempty"`
	Name string `json:"name"`
}

type CourseServiceClient struct {
	GetCourseEndpoint endpoint.Endpoint
}

func NewHTTPClient(instance string, logger log.Logger) (CourseServiceClient, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return CourseServiceClient{}, err
	}

	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))
	var getCourseEndpoint endpoint.Endpoint
	{
		getCourseEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/courses"),
			encodeGetCourseRequest,
			decodeGetCourseResponse,
		).Endpoint()
		getCourseEndpoint = limiter(getCourseEndpoint)
		getCourseEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getCourseEndpoint)
	}

	return CourseServiceClient{
		GetCourseEndpoint: getCourseEndpoint,
	}, nil
}

func (c *CourseServiceClient) GetCourse(ctx context.Context, id string) (Course, error) {
	response, err := c.GetCourseEndpoint(ctx, getCourseRequest{ID: id})
	if err != nil {
		return Course{}, err
	}
	resp := response.(getCourseResponse)
	if resp.Err != nil {
		return Course{}, resp.Err
	}
	return resp.Course, nil
}

type getCourseRequest struct {
	ID string
}

type getCourseResponse struct {
	Course Course
	Err    error `json:"error,omitempty"`
}

func decodeGetCourseResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getCourseResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func encodeGetCourseRequest(_ context.Context, req *http.Request, request interface{}) error {
	r := request.(getCourseRequest)

	req.URL.Path = path.Join(req.URL.Path, url.PathEscape(r.ID))

	return nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}
