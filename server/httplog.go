// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// AccessLogRecord struct for holding access log data.
type LogHttpRequest struct {
	RemoteAddr    string `json:"remote_addr"`
	RequestMethod string `json:"request_method"`
	//	Request        string `json:"request"`
	ServerProtocol string `json:"server_protocol"`
	Host           string `json:"host"`
	HTTPReferrer   string `json:"http_referrer"`
	HTTPUserAgent  string `json:"http_user_agent"`
	RemoteUser     string `json:"remote_user"`
}

type LogResponse struct {
	Status string `json:"status"`    // e.g. "200 OK"
	Proto  string `json:"respproto"` // e.g. "HTTP/1.0"
}

func (r *LogHttpRequest) json() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)

	err := encoder.Encode(r)
	return buffer.Bytes(), err
}

func (r *LogResponse) json() ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)

	err := encoder.Encode(r)
	return buffer.Bytes(), err
}

// ---- Logging

func WrapHandlerWithLogging(wrappedHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("--> %s %s\n", req.Method, req.URL.Path)

		lrw := NewLoggingResponseWriter(w)
		wrappedHandler.ServeHTTP(lrw, req)

		statusCode := lrw.statusCode
		fmt.Printf("<-- %d %s\n", statusCode, http.StatusText(statusCode))
	})
}

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	// WriteHeader(int) is not called if our response implicitly returns 200 OK, so
	// we default to that status code.
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func LogResp(w http.ResponseWriter) (msg string) {

	logResponseWriter := &LogResponse{

		Status: "STATUS OK",
		Proto:  "PROTO HTTP 1.1",
	}

	if jsonData, err := logResponseWriter.json(); err != nil {
		msg = fmt.Sprintf(`{"Error": "%s"}`, err)
	} else {

		msg = string(jsonData)

	}

	return msg

}

func RequestToLog(r *http.Request) string {

	var msg string

	logRequestData := &LogHttpRequest{
		RemoteAddr:    r.RemoteAddr, //w.Context.Input.IP(),
		RequestMethod: r.Method,
		//	Request:        fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto),
		ServerProtocol: r.Proto,
		Host:           r.Host,
		HTTPReferrer:   r.Header.Get("Referer"),
		HTTPUserAgent:  r.Header.Get("User-Agent"),
		RemoteUser:     r.Header.Get("Remote-User"),
	}

	if jsonData, err := logRequestData.json(); err != nil {
		msg = fmt.Sprintf(`{"Error": "%s"}`, err)
	} else {

		msg = string(jsonData)

	}

	return msg

}
