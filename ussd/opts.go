/*
 * MIT License
 *
 * Copyright (c) 2021 TechCraft Technologies Co. Ltd
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package ussd

import (
	"io"
	"net/http"
)

// ClientOption is a setter func to set Client details like
// Timeout, context, Http and Logger
type ClientOption func(client *Client)

// WithDebugMode set debug mode to true or false
func WithDebugMode(debugMode bool) ClientOption {
	return func(client *Client) {
		client.base.DebugMode = debugMode

	}
}

// WithLogger set a Logger of user preference but of type io.Writer
// that will be used for debugging use cases. A default value is os.Stderr
// it can be replaced by any io.Writer unless its nil which in that case
// it will be ignored
func WithLogger(out io.Writer) ClientOption {
	return func(client *Client) {
		if out == nil {
			return
		}
		client.base.Logger = out
	}
}

// WithHTTPClient when called unset the present http.Client and replace it
// with c. In case user tries to pass a nil value referencing the pkg
// i.e WithHTTPClient(nil), it will be ignored and the pkg wont be replaced
// Note: the new pkg Transport will be modified. It will be wrapped by another
// middleware that enables pkg to
func WithHTTPClient(httpClient *http.Client) ClientOption {

	// TODO check if its really necessary to set the default Timeout to 1 minute

	return func(client *Client) {
		if httpClient == nil {
			return
		}

		client.base.Http = httpClient
	}
}
