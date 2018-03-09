// Copyright (c) 2018
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package kongclient

// https://getkong.org/docs/0.12.x/admin-api/

// KongGetResp defines the returned json format for "GET /apis/"
type KongGetResp struct {
	Total int           `json:"total,omitempty"`
	Data  []KongAPISpec `json:"data,omitempty"`
}

// "POST /apis/"

// TODO: there are more fields in recent releases

// KongAPISpec defines the request format for "POST /apis/"
type KongAPISpec struct {
	ID           string `json:"id,omitempty"`
	Name         string `json:"name"`
	CreatedAt    int64  `json:"-"`
	UpstreamURL  string `json:"upstream_url,omitempty"`
	PreserveHost bool   `json:"preserve_host"`
	// kong 0.9.x and earlier
	RequestPath      string `json:"request_path,omitempty"`
	RequestHost      string `json:"request_host,omitempty"`
	StripRequestPath bool   `json:"strip_request_path,omitempty"`
	// kong 0.10.x and later
	Hosts    []string `json:"hosts,omitempty"`
	Uris     []string `json:"uris,omitempty"`
	StripURI bool     `json:"strip_uri,omitempty"`
}
