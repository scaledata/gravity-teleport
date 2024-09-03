// package forwarder implements http handler that forwards requests to remote server
// and serves back the response
// websocket proxying support based on https://github.com/yhat/wsutil
package forward

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gravitational/oxy/utils"
)

// ReqRewriter can alter request headers and body
type ReqRewriter interface {
	Rewrite(r *http.Request)
}

type optSetter func(f *Forwarder) error

// PassHostHeader specifies if a client's Host header field should
// be delegated
func PassHostHeader(b bool) optSetter {
	return func(f *Forwarder) error {
		f.passHost = b
		return nil
	}
}

// RoundTripper sets a new http.RoundTripper
// Forwarder will use http.DefaultTransport as a default round tripper
func RoundTripper(r http.RoundTripper) optSetter {
	return func(f *Forwarder) error {
		f.roundTripper = r
		return nil
	}
}

// Dialer mirrors the net.Dial function to be able to define alternate
// implementations
type Dialer func(network, address string) (net.Conn, error)

// WebsocketDial defines a new network dialer to use to dial to remote websocket destination.
// If no dialer has been defined, net.Dial will be used.
func WebsocketDial(dial Dialer) optSetter {
	return func(f *Forwarder) error {
		f.websocketForwarder.dial = dial
		return nil
	}
}

// Rewriter defines a request rewriter for the HTTP forwarder
func Rewriter(r ReqRewriter) optSetter {
	return func(f *Forwarder) error {
		f.httpForwarder.rewriter = r
		return nil
	}
}

// WebsocketRewriter defines a request rewriter for the websocket forwarder
func WebsocketRewriter(r ReqRewriter) optSetter {
	return func(f *Forwarder) error {
		f.websocketForwarder.rewriter = r
		return nil
	}
}

// ErrorHandler is a functional argument that sets error handler of the server
func ErrorHandler(h utils.ErrorHandler) optSetter {
	return func(f *Forwarder) error {
		f.errHandler = h
		return nil
	}
}

// Logger specifies the logger to use.
// Forwarder will default to oxyutils.NullLogger if no logger has been specified
func Logger(l utils.Logger) optSetter {
	return func(f *Forwarder) error {
		f.log = l
		return nil
	}
}

// FlushInterval sets flush interval for streaming response
func FlushInterval(t time.Duration) optSetter {
	return func(f *Forwarder) error {
		f.httpForwarder.flushInterval = t
		return nil
	}
}

// Forwarder wraps two traffic forwarding implementations: HTTP and websockets.
// It decides based on the specified request which implementation to use
type Forwarder struct {
	*httpForwarder
	*websocketForwarder
	*handlerContext
}

// handlerContext defines a handler context for error reporting and logging
type handlerContext struct {
	errHandler utils.ErrorHandler
	log        utils.Logger
}

// httpForwarder is a handler that can reverse proxy
// HTTP traffic
type httpForwarder struct {
	roundTripper  http.RoundTripper
	rewriter      ReqRewriter
	passHost      bool
	flushInterval time.Duration
}

// websocketForwarder is a handler that can reverse proxy
// websocket traffic
type websocketForwarder struct {
	dial            Dialer
	rewriter        ReqRewriter
	TLSClientConfig *tls.Config
}

// New creates an instance of Forwarder based on the provided list of configuration options
func New(setters ...optSetter) (*Forwarder, error) {
	f := &Forwarder{
		httpForwarder:      &httpForwarder{},
		websocketForwarder: &websocketForwarder{},
		handlerContext:     &handlerContext{},
	}
	for _, s := range setters {
		if err := s(f); err != nil {
			return nil, err
		}
	}
	if f.httpForwarder.roundTripper == nil {
		f.httpForwarder.roundTripper = http.DefaultTransport
	}
	if f.websocketForwarder.dial == nil {
		f.websocketForwarder.dial = net.Dial
	}
	if f.httpForwarder.rewriter == nil {
		f.httpForwarder.rewriter = NewHeaderRewriter()
	}
	if f.log == nil {
		f.log = utils.NullLogger
	}
	if f.errHandler == nil {
		f.errHandler = utils.DefaultHandler
	}
	return f, nil
}

// ServeHTTP decides which forwarder to use based on the specified
// request and delegates to the proper implementation
func (f *Forwarder) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if isWebsocketRequest(req) {
		f.websocketForwarder.serveHTTP(w, req, f.handlerContext)
	} else {
		f.httpForwarder.serveHTTP(w, req, f.handlerContext)
	}
}

// serveHTTP forwards HTTP traffic using the configured transport
func (f *httpForwarder) serveHTTP(w http.ResponseWriter, req *http.Request, ctx *handlerContext) {
	start := time.Now().UTC()
	response, err := f.roundTripper.RoundTrip(f.copyRequest(req, req.URL, ctx))
	if err != nil {
		ctx.log.Errorf("Error forwarding to %v, err: %v", req.URL, err)
		ctx.errHandler.ServeHTTP(w, req, err)
		return
	}

	if req.TLS != nil {
		ctx.log.Infof("Round trip: %v %v, code: %v, duration: %v tls:version: %x, tls:resume:%t, tls:csuite:%x, tls:server:%v",
			req.Method, req.URL, response.StatusCode, time.Now().UTC().Sub(start),
			req.TLS.Version,
			req.TLS.DidResume,
			req.TLS.CipherSuite,
			req.TLS.ServerName)
	} else {
		ctx.log.Infof("Round trip: %v %v, code: %v, duration: %v",
			req.Method, req.URL, response.StatusCode, time.Now().UTC().Sub(start))
	}

	utils.CopyHeaders(w.Header(), response.Header)
	w.WriteHeader(response.StatusCode)
	written, err := copyResponse(ctx, f.flushInterval, w, response.Body)
	defer response.Body.Close()

	if err != nil && err != io.EOF {
		ctx.log.Errorf("Error copying upstream response Body: %v", err)
		ctx.errHandler.ServeHTTP(w, req, err)
		return
	}

	if written != 0 {
		w.Header().Set(ContentLength, strconv.FormatInt(written, 10))
	}
}

func (f *httpForwarder) getURLFromRequest(req *http.Request, ctx *handlerContext) *url.URL {
	// If the Request was created by Go via a real HTTP request,  RequestURI will
	// contain the original query string. If the Request was created in code, RequestURI
	// will be empty, and we will use the URL object instead
	u := req.URL
	if req.RequestURI != "" {
		parsedURL, err := url.ParseRequestURI(req.RequestURI)
		if err == nil {
			u = parsedURL
		} else {
			ctx.log.Warningf("gravitational/oxy/forward: error when parsing RequestURI: %s", err)
		}
	}
	return u
}

// copyRequest makes a copy of the specified request to be sent using the configured
// transport
func (f *httpForwarder) copyRequest(req *http.Request, target *url.URL, ctx *handlerContext) *http.Request {
	outReq := new(http.Request)
	*outReq = *req // includes shallow copies of maps, but we handle this below

	outReq.URL = utils.CopyURL(outReq.URL)
	outReq.URL.Scheme = target.Scheme
	outReq.URL.Host = target.Host

	u := f.getURLFromRequest(outReq, ctx)

	outReq.URL.Path = u.Path
	outReq.URL.RawPath = u.RawPath
	outReq.URL.RawQuery = u.RawQuery
	outReq.RequestURI = "" // Outgoing request should not have RequestURI

	outReq.Proto = "HTTP/1.1"
	outReq.ProtoMajor = 1
	outReq.ProtoMinor = 1

	// Overwrite close flag so we can keep persistent connection for the backend servers
	outReq.Close = false

	outReq.Header = make(http.Header)
	utils.CopyHeaders(outReq.Header, req.Header)

	if f.rewriter != nil {
		f.rewriter.Rewrite(outReq)
	}

	// Do not pass client Host header unless optsetter PassHostHeader is set.
	if !f.passHost {
		outReq.Host = target.Host
	}
	return outReq
}

// serveHTTP forwards websocket traffic
func (f *websocketForwarder) serveHTTP(w http.ResponseWriter, req *http.Request, ctx *handlerContext) {
	outReq := f.copyRequest(req)
	host := outReq.URL.Host

	// if host does not specify a port, use the default http port
	if !strings.Contains(host, ":") {
		if outReq.URL.Scheme == "wss" {
			host = host + ":443"
		} else {
			host = host + ":80"
		}
	}

	targetConn, err := f.dial("tcp", host)
	if err != nil {
		ctx.log.Errorf("Error dialing `%v`: %v", host, err)
		ctx.errHandler.ServeHTTP(w, req, err)
		return
	}
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		ctx.log.Errorf("Unable to hijack the connection: does not implement http.Hijacker")
		ctx.errHandler.ServeHTTP(w, req, err)
		return
	}
	underlyingConn, _, err := hijacker.Hijack()
	if err != nil {
		ctx.log.Errorf("Unable to hijack the connection: %v", err)
		ctx.errHandler.ServeHTTP(w, req, err)
		return
	}
	// it is now caller's responsibility to Close the underlying connection
	defer underlyingConn.Close()
	defer targetConn.Close()

	// write the modified incoming request to the dialed connection
	if err = outReq.Write(targetConn); err != nil {
		ctx.log.Errorf("Unable to copy request to target: %v", err)
		ctx.errHandler.ServeHTTP(w, req, err)
		return
	}

	// read response code to make sure connection upgrade succeeded
	respCode, respBody, err := readResponseAndCode(targetConn)
	if err != nil {
		ctx.log.Errorf("Unable to read websocket upgrade response: %v", err)
		return
	}

	// write upgrade response to the client
	if len(respBody) > 0 {
		if _, err := underlyingConn.Write(respBody); err != nil {
			ctx.log.Errorf("Unable to write websocket upgrade response: %v", err)
			return
		}
	}

	// make sure we got 101 before establishing a bidirectional pipe
	if respCode != http.StatusSwitchingProtocols {
		ctx.log.Warningf("Unable to upgrade websocket connection: %v", string(respBody))
		return
	}

	ctx.log.Infof("Websocket upgrade: %v", outReq.URL.String())

	errc := make(chan error, 2)
	replicate := func(dst io.Writer, src io.Reader) {
		_, err := io.Copy(dst, src)
		errc <- err
	}
	go replicate(targetConn, underlyingConn)
	go replicate(underlyingConn, targetConn)
	<-errc
}

// readResponseAndCode reads an HTTP response and its status code from reader.
func readResponseAndCode(r io.Reader) (int, []byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 256))
	resp, err := http.ReadResponse(bufio.NewReader(io.TeeReader(r, buf)), nil)
	if err != nil {
		return 0, nil, err
	}
	return resp.StatusCode, buf.Bytes(), nil
}

// copyRequest makes a copy of the specified request.
func (f *websocketForwarder) copyRequest(req *http.Request) (outReq *http.Request) {
	outReq = new(http.Request)
	*outReq = *req
	outReq.URL = utils.CopyURL(req.URL)
	outReq.URL.Scheme = req.URL.Scheme
	outReq.URL.Host = req.URL.Host
	if f.rewriter != nil {
		f.rewriter.Rewrite(outReq)
	}
	return outReq
}

// isWebsocketRequest determines if the specified HTTP request is a
// websocket handshake request
func isWebsocketRequest(req *http.Request) bool {
	containsHeader := func(name, value string) bool {
		items := strings.Split(req.Header.Get(name), ",")
		for _, item := range items {
			if value == strings.ToLower(strings.TrimSpace(item)) {
				return true
			}
		}
		return false
	}
	return containsHeader(Connection, "upgrade") && containsHeader(Upgrade, "websocket")
}
