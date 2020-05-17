package BTCMarkets

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	Endpoint          = "https://api.btcmarkets.net/v3"
	httpClientTimeout = 15 * time.Second
)

var (
	ErrUnexpectedResponse = errors.New("the BTCMarkets API is currently unavailable")
)

type Client struct {
	AccessKey  string       // The API access key.
	PrivateKey string       // The API private key.
	HTTPClient *http.Client // The HTTP client to send requests on.
}

type contentType string

const (
	contentTypeEmpty contentType = ""
	contentTypeJSON  contentType = "application/json"
)

func New(apiKey string, privateKey string) *Client {
	return &Client{
		AccessKey:  apiKey,
		PrivateKey: privateKey,
		HTTPClient: &http.Client{
			Timeout: httpClientTimeout,
		},
	}
}

func paginationQuery(endpointSpecificOptions url.Values, options *ListOptions) string {
	if options == nil {
		return ""
	}

	query := endpointSpecificOptions
	if options.Limit != 0 {
		query.Set("limit", strconv.FormatUint(options.Limit, 10))
	}
	if options.Before != 0 {
		query.Set("before", strconv.FormatUint(options.Before, 10))
	}
	if options.After != 0 {
		query.Set("after", strconv.FormatUint(options.After, 10))

	}

	return query.Encode()
}

func (c *Client) Do(v interface{}, method string, path string, data interface{}) error {
	u := &url.URL{
		Path: path,
	}

	if err := c.Request(&v, method, u.String(), data, nil); err != nil {
		return err
	}

	return nil
}

func (c *Client) DoWithPagination(v interface{}, method string, path string, endpointSpecificOptions url.Values, listOptions *ListOptions) error {
	query := paginationQuery(endpointSpecificOptions, listOptions)

	u := &url.URL{
		Path:     path,
		RawQuery: query,
	}

	if err := c.Request(&v, method, u.String(), nil, listOptions); err != nil {
		return err
	}

	return nil
}

func (c *Client) Request(v interface{}, method string, path string, data interface{}, listOptions *ListOptions) error {
	request, err := buildRequest(c, method, path, data)
	if err != nil {
		return err
	}

	response, err := c.HTTPClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated:
		// Status codes 200 and 201 are indicative of being able to convert the
		// response body to the struct that was specified.
		if err := json.Unmarshal(responseBody, &v); err != nil {
			return fmt.Errorf("could not decode response JSON, %s: %v", string(responseBody), err)
		}

		// Set the pagination values based upon headers received in response
		bmBefore := response.Header.Get("Bm-Before")
		bmAfter := response.Header.Get("Bm-After")
		if listOptions != nil && bmBefore != "" && bmAfter != "" { // Ugly guard statement
			listOptions.PrevBefore, err = strconv.ParseUint(bmBefore, 10, 0)
			listOptions.PrevAfter, err = strconv.ParseUint(bmAfter, 10, 0)
		}

		return nil
	case http.StatusInternalServerError:
		// Status code 500 is a server error and means nothing can be done at this
		// point.
		return ErrUnexpectedResponse
	default:
		var errorResponse Error
		if err := json.Unmarshal(responseBody, &errorResponse); err != nil {
			return err
		}

		return errorResponse
	}
}

func buildRequest(client *Client, method string, path string, data interface{}) (*http.Request, error) {
	if !strings.HasPrefix(path, "https://") && !strings.HasPrefix(path, "http://") {
		path = fmt.Sprintf("%s/%s", Endpoint, path)
	}
	uri, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	body, contentType, err := prepareRequestBody(data)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, uri.String(), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	if contentType != contentTypeEmpty {
		request.Header.Set("Content-Type", string(contentType))
	}

	BuildAuthHeaders(request, client.PrivateKey, method, uri.Path, string(body))
	request.Header.Set("Accept", "application/json")
	request.Header.Set("BM-AUTH-APIKEY", client.AccessKey)

	return request, nil
}

func prepareRequestBody(data interface{}) ([]byte, contentType, error) {
	switch data := data.(type) {
	case nil:
		// Nil bodies are accepted by `net/http`, so this is not an error.
		return nil, contentTypeEmpty, nil
	default:
		b, err := json.Marshal(data)
		if err != nil {
			return nil, contentType(""), err
		}

		return b, contentTypeJSON, nil
	}
}

func BuildAuthHeaders(request *http.Request, privateKey string, method string, path string, body string) {
	//getting now() in milliseconds
	nowMs := strconv.FormatInt(time.Now().UTC().UnixNano()/1000000, 10)

	stringToSign := method + path + nowMs

	if body != "" {
		stringToSign += body
	}

	request.Header.Set("BM-AUTH-TIMESTAMP", nowMs)
	request.Header.Set("BM-AUTH-SIGNATURE", signMessage(privateKey, stringToSign))
}

func signMessage(key string, message string) string {
	encodedKey, _ := base64.StdEncoding.DecodeString(key)
	mac := hmac.New(sha512.New, encodedKey)
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}
