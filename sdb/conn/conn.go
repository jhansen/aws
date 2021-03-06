// Copyright 2012 Aaron Jacobs. All Rights Reserved.
// Author: aaronjjacobs@gmail.com (Aaron Jacobs)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package conn

import (
	"fmt"
	"github.com/jhansen/aws"
	"github.com/jhansen/aws/time"
)

const (
	iso8601Format = "2006-01-02T15:04:05Z"
)

// A connection to the SimpleDB service.
type Conn interface {
	// Send the supplied request to the service, taking care of adding
	// appropriate authentication info.
	SendRequest(req Request) (resp []byte, err error)
}

// Create a connection using the supplied dependencies.
func NewConn(
	key aws.AccessKey,
	httpConn HttpConn,
	signer Signer,
	clock time.Clock) (Conn, error) {
	return &conn{key, httpConn, signer, clock}, nil
}

type conn struct {
	key      aws.AccessKey
	httpConn HttpConn
	signer   Signer
	clock    time.Clock
}

func (c *conn) SendRequest(req Request) (resp []byte, err error) {
	// Make a copy of the request that we can modify below.
	originalReq := req
	req = Request{}
	for key, val := range originalReq {
		req[key] = val
	}

	// Set authentication info.
	req["AWSAccessKeyId"] = c.key.Id
	req["Timestamp"] = c.clock.Now().UTC().Format(iso8601Format)
	req["SignatureVersion"] = "2"
	req["SignatureMethod"] = "HmacSHA1"

	// Sign the request.
	if err = c.signer.SignRequest(req); err != nil {
		err = fmt.Errorf("SignRequest: %v", err)
		return
	}

	// Send the request.
	httpResp, err := c.httpConn.SendRequest(req)
	if err != nil {
		err = fmt.Errorf("SendRequest: %v", err)
		return
	}

	// Did the server return an error?
	if httpResp.StatusCode != 200 {
		err = fmt.Errorf(
			"Error from server (%d): %s",
			httpResp.StatusCode,
			httpResp.Body)

		return
	}

	return httpResp.Body, nil
}
