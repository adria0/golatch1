package golatch1

import "fmt"
import "sort"
import "strings"
import "crypto/hmac"
import "crypto/sha1"
import "golang.org/x/text/transform"
import "golang.org/x/text/encoding/charmap"
import "encoding/base64"
import "encoding/json"
import "net/url"
import "time"
import "bytes"
import "net/http"
import "io/ioutil"

const (
	ServiceURL = "https://latch.elevenpaths.com"
)

const (
	httpGet    = "GET"
	httpPut    = "PUT"
	httpPost   = "POST"
	httpDelete = "DELETE"
)

type LatchApp struct {
	appId     string
	secretKey string
	transport *http.Client
}

// Creates a new LatchApp with the supplied appId and secret
func NewLatchApp(appId string, secret string) *LatchApp {
	return &LatchApp{appId, secret, &http.Client{}}
}

// Creates a new LatchApp with the supplied appId and secret, but
// also specifying the transport
func NewLatchAppWithTransport(appId string, secret string, transport *http.Client) *LatchApp {
	return &LatchApp{appId, secret, transport}
}

// Serialize x-11paths-headers, used to calc authentication
func serialHeaders(xheaders map[string]string) []byte {
	if xheaders == nil {
		return []byte{}
	}

	keys := make([]string, 0, len(xheaders))
	for k := range xheaders {
		if strings.Index(strings.ToLower(k), "x-11paths-") == 0 {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	var b bytes.Buffer
	for i, k := range keys {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(strings.ToLower(k))
		b.WriteString(":")
		b.WriteString(strings.Replace(xheaders[k], "\n", " ", -1))
	}
	return b.Bytes()
}

// Serialize URL parameters
func serialParams(params map[string]string) []byte {

	if params == nil {
		return []byte{}
	}
	keys := make([]string, 0, len(params))

	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b bytes.Buffer
	for i, k := range keys {
		if i > 0 {
			b.WriteString("&")
		}
		b.WriteString(url.QueryEscape(k))
		b.WriteString("=")
		b.WriteString(url.QueryEscape(params[k]))
	}

	return b.Bytes()
}

// Calc the HMACSHA1 autentication header
func (app *LatchApp) authHeader(method string, path string, params map[string]string, xheaders map[string]string, utcTime string) (header string, err error) {

	var b bytes.Buffer

	b.WriteString(method)
	b.WriteString("\n")
	b.WriteString(utcTime)
	b.WriteString("\n")
	b.Write(serialHeaders(xheaders))
	b.WriteString("\n")
	b.WriteString(path)

	if len(params) > 0 {
		b.WriteString("\n")
		b.Write(serialParams(params))
	}

	encoder := charmap.Windows1252.NewEncoder()
	encoded, _, err := transform.Bytes(encoder, b.Bytes())
	if err != nil {
		return "", err
	}

	hmacsha1 := hmac.New(sha1.New, []byte(app.secretKey))
	hmacsha1.Write(encoded)
	signature := base64.StdEncoding.EncodeToString(hmacsha1.Sum(nil))

	var authHeader bytes.Buffer
	authHeader.WriteString("11PATHS ")
	authHeader.WriteString(app.appId)
	authHeader.WriteString(" ")
	authHeader.WriteString(signature)

	return authHeader.String(), nil
}

func (app *LatchApp) call(method string, path string, params map[string]string, xheaders map[string]string, jsonStruct interface{}) error {

	body := serialParams(params)
	req, err := http.NewRequest(method, ServiceURL+path, bytes.NewReader(body))

	if err != nil {
		return err
	}

	if xheaders != nil {
		for k, v := range xheaders {
			req.Header.Set(k, v)
		}
	}

	if method == httpPost {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Content-Length", fmt.Sprintf("%v", len(body)))
	}

	utcTime := time.Now().UTC().Format("2006-01-02 15:04:05")
	auth, err := app.authHeader(method, path, xheaders, params, utcTime)
	req.Header.Set("Authorization", auth)
	req.Header.Set("X-11paths-Date", utcTime)

	resp, err := app.transport.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, jsonStruct)
	if err != nil {
		return err
	}

	return nil
}
