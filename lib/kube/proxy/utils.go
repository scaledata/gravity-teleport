package proxy

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/klog/v2"
	"net"
	"net/http"
	"net/url"
)

// ConnectWithRedirects uses dialer to send req, following up to 10 redirects (relative to
// originalLocation). It returns the opened net.Conn and the raw response bytes.
// If requireSameHostRedirects is true, only redirects to the same host are permitted.
func ConnectWithRedirects(originalMethod string, originalLocation *url.URL, header http.Header, originalBody io.Reader, dialer utilnet.Dialer, requireSameHostRedirects bool) (net.Conn, []byte, error) {
	const (
		maxRedirects    = 9     // Fail on the 10th redirect
		maxResponseSize = 16384 // play it safe to allow the potential for lots of / large headers
	)

	var (
		location         = originalLocation
		method           = originalMethod
		intermediateConn net.Conn
		rawResponse      = bytes.NewBuffer(make([]byte, 0, 256))
		body             = originalBody
	)

	defer func() {
		if intermediateConn != nil {
			intermediateConn.Close()
		}
	}()

redirectLoop:
	for redirects := 0; ; redirects++ {
		if redirects > maxRedirects {
			return nil, nil, fmt.Errorf("too many redirects (%d)", redirects)
		}

		req, err := http.NewRequest(method, location.String(), body)
		if err != nil {
			return nil, nil, err
		}

		req.Header = header

		intermediateConn, err = dialer.Dial(req)
		if err != nil {
			return nil, nil, err
		}

		// Peek at the backend response.
		rawResponse.Reset()
		respReader := bufio.NewReader(io.TeeReader(
			io.LimitReader(intermediateConn, maxResponseSize), // Don't read more than maxResponseSize bytes.
			rawResponse)) // Save the raw response.
		resp, err := http.ReadResponse(respReader, nil)
		if err != nil {
			// Unable to read the backend response; let the client handle it.
			klog.Warningf("Error reading backend response: %v", err)
			break redirectLoop
		}

		switch resp.StatusCode {
		case http.StatusFound:
			// Redirect, continue.
		default:
			// Don't redirect.
			break redirectLoop
		}

		// Redirected requests switch to "GET" according to the HTTP spec:
		// https://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.3
		method = "GET"
		// don't send a body when following redirects
		body = nil

		resp.Body.Close() // not used

		// Prepare to follow the redirect.
		redirectStr := resp.Header.Get("Location")
		if redirectStr == "" {
			return nil, nil, fmt.Errorf("%d response missing Location header", resp.StatusCode)
		}
		// We have to parse relative to the current location, NOT originalLocation. For example,
		// if we request http://foo.com/a and get back "http://bar.com/b", the result should be
		// http://bar.com/b. If we then make that request and get back a redirect to "/c", the result
		// should be http://bar.com/c, not http://foo.com/c.
		location, err = location.Parse(redirectStr)
		if err != nil {
			return nil, nil, fmt.Errorf("malformed Location header: %v", err)
		}

		// Only follow redirects to the same host. Otherwise, propagate the redirect response back.
		if requireSameHostRedirects && location.Hostname() != originalLocation.Hostname() {
			return nil, nil, fmt.Errorf("hostname mismatch: expected %s, found %s", originalLocation.Hostname(), location.Hostname())
		}

		// Reset the connection.
		intermediateConn.Close()
		intermediateConn = nil
	}

	connToReturn := intermediateConn
	intermediateConn = nil // Don't close the connection when we return it.
	return connToReturn, rawResponse.Bytes(), nil
}
