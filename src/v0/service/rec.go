package service

import (
	"encoding/json"
	"errors"
	"github.com/valyala/fasthttp"
	"os"
	"time"
)

const (
	RecApiAddr     = "https://presentio-gorse-master.herokuapp.com/"
	ReqContentType = "application/json"
)

var client *fasthttp.Client
var apiKey string

func init() {
	client = &fasthttp.Client{
		ReadTimeout:                   time.Second,
		WriteTimeout:                  time.Second,
		MaxIdleConnDuration:           time.Hour,
		NoDefaultUserAgentHeader:      false,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        false,
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}

	apiKey = os.Getenv("GORSE_SERVER_API_KEY")
}

type ItemEntity struct {
	Categories []string
	IsHidden   bool
	ItemId     string
	Labels     []string
	Timestamp  string
}

func CreateOrUpdateRecItem(entity *ItemEntity) error {
	req := fasthttp.AcquireRequest()

	body, err := json.Marshal(entity)

	if err != nil {
		return err
	}

	req.SetRequestURI(RecApiAddr + "/item")
	req.Header.SetContentType(ReqContentType)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set("X-API-Key", apiKey)
	req.SetBodyRaw(body)

	resp := fasthttp.AcquireResponse()

	err = client.Do(req, resp)
	fasthttp.ReleaseRequest(req)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New("unable to create item")
	}

	fasthttp.ReleaseResponse(resp)

	return nil
}

type FeedbackEntity struct {
	FeedbackType string
	ItemId       string
	Timestamp    string
	UserId       string
}

func AddFeedback(entity *FeedbackEntity) error {
	req := fasthttp.AcquireRequest()

	body, err := json.Marshal(entity)

	if err != nil {
		return err
	}

	req.SetRequestURI(RecApiAddr + "/api/feedback")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentType(ReqContentType)
	req.Header.Set("X-API-Key", apiKey)
	req.SetBodyRaw(body)

	resp := fasthttp.AcquireResponse()

	err = client.Do(req, resp)
	fasthttp.ReleaseRequest(req)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New("unable to insert feedback")
	}

	fasthttp.ReleaseResponse(resp)

	return nil
}

func RemoveFeedback(entity *FeedbackEntity) error {
	req := fasthttp.AcquireRequest()

	req.SetRequestURI(RecApiAddr + "/api/feedback/" + entity.FeedbackType + "/" + entity.UserId + "/" + entity.ItemId)
	req.Header.SetContentType(ReqContentType)
	req.Header.SetMethod(fasthttp.MethodDelete)
	req.Header.Set("X-API-Key", apiKey)

	resp := fasthttp.AcquireResponse()

	err := client.Do(req, resp)
	fasthttp.ReleaseRequest(req)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.New("unable to delete feedback")
	}

	fasthttp.ReleaseResponse(resp)

	return nil
}
