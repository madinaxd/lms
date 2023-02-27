package client

import (
	"context"
	"encoding/json"
	"fmt"
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

type Student struct {
	ID          int64     `json:"id,omitempty"`
	Fullname    string    `json:"fullname"`
	DateOfBirth time.Time `json:"date_of_birth"`
	Grade       int       `json:"grade,omitempty"`
	Phone       int64     `json:"phone"`
}

type StudentServiceClient struct {
	GetCourseStudentsEndpoint endpoint.Endpoint
}

func NewHTTPClient(instance string, logger log.Logger) (StudentServiceClient, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return StudentServiceClient{}, err
	}

	limiter := ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 100))
	var getCourseStudentsEndpoint endpoint.Endpoint
	{
		getCourseStudentsEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/"),
			encodeGetCourseStudentsRequest,
			decodeGetCourseStudentsResponse,
		).Endpoint()
		getCourseStudentsEndpoint = limiter(getCourseStudentsEndpoint)
		getCourseStudentsEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(getCourseStudentsEndpoint)
	}

	return StudentServiceClient{
		GetCourseStudentsEndpoint: getCourseStudentsEndpoint,
	}, nil
}

func (c *StudentServiceClient) GetCourseStudents(ctx context.Context, id string) ([]Student, error) {
	response, err := c.GetCourseStudentsEndpoint(ctx, getCourseStudentsRequest{CourseID: id})
	if err != nil {
		return []Student{}, err
	}
	resp := response.(getCourseStudentsResponse)
	if resp.Err != nil {
		return []Student{}, resp.Err
	}
	return resp.Students, nil
}

type getCourseStudentsRequest struct {
	CourseID string
}

type getCourseStudentsResponse struct {
	Students []Student
	Err      error `json:"error,omitempty"`
}

func decodeGetCourseStudentsResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getCourseStudentsResponse
	fmt.Println(resp.Body)
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func encodeGetCourseStudentsRequest(_ context.Context, req *http.Request, request interface{}) error {
	r := request.(getCourseStudentsRequest)

	req.URL.Path = path.Join(req.URL.Path, "courses", url.PathEscape(r.CourseID), "students")
	return nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}
