// Handles basic connection primatives.
//
// Unlike a simple connection, awsconn carries a second URL
// used for proxy-server data.  It is not a connection, but
// but a set of helpers to the "net" class of functionality.
//
// It is the encouraged way for goaws
// utlities to establish and maintain connection details so
// that proxy configuration data is available to all 
// callers.
package awsconn
/* 
  Copyright (c) 2010, Abneptis LLC.
  See COPYRIGHT and LICENSE for details.
*/

import "com.abneptis.oss/urltools"

import "http"
import "net"
import "os"
import "bufio"
import "xml"

// Sends a request and returns an error
// A simple helper exposed as it is often
// useful.
func SendRequest(cc *http.ClientConn, req *http.Request)(resp *http.Response, err os.Error){
  err = cc.Write(req)
  if err == nil {
    resp, err = cc.Read()
  }
  return
}

// Handles xml unmarshalling to a generic object from
// an http.Response
func ParseResponse(resp *http.Response, o interface{})(err os.Error){
  if resp.Body == nil {
    err = os.NewError("Response body is empty")
  } else {
    parser := xml.NewParser(resp.Body)
    err = parser.Unmarshal(o, nil)
  }
  return
}

// Construct a new Endpoint. 
func NewEndpoint(u, pu *http.URL)(*Endpoint){
  return &Endpoint {
   URL: u,
   ProxyURL: pu,
  }
}

// An endpoint has two URL's, the "URL", or rather, the actual service
// address, and a "ProxyURL", which is used for low-level connections if
// it is present -- else the URL is connected to directly.
type Endpoint struct {
  URL *http.URL
  ProxyURL *http.URL
}

// Return the URL to be used for connection purposes.
// While not expected to be needed by external users,
// this is considered the "correct" way to make that
// determination, and shorter to import and use than write.
func (self *Endpoint)ConnectionURL()(out *http.URL){
  if self.ProxyURL != nil {
    out = self.ProxyURL
  } else {
    out = self.URL
  }
  return
}

// Return the URL to be used for request generation.
// This is far more likely to be useful to end users
// than the above.
func (self *Endpoint)GetURL()(out *http.URL){
  return self.URL
}


// Return a new net.Conn using netname and local as net.Dial does.
// NewConn does not explicitly check to ensure that you are using
// a stream protocol, so if you accept this from a user source, it
// is the callers responsibility to verify.
func (self *Endpoint)NewConn(netname, local string)(rawc net.Conn, err os.Error){
  hps, err := urltools.ExtractURLHostPort(self.ConnectionURL())
  if err == nil {
    rawc, err = net.Dial(netname, local, hps)
  }
  return
}

// Returns a new HTTP connection;  As with NewConn, netname is not
// checked, and the behaviour of an HTTP client over a non stream
// protocol is undefined - but probably interesting to watch.
func (self *Endpoint)NewHTTPClientConn(netname, local string, r *bufio.Reader)(hc *http.ClientConn, err os.Error){
  rawc, err := self.NewConn(netname, local)
  if err == nil {
    hc = http.NewClientConn(rawc, r)
  }
  return
}
