// Copyright 2016 Qiang Xue. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package routing

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestNewHttpError(t *testing.T) {
	e := NewHTTPError(http.StatusNotFound)
	if e.StatusCode() != http.StatusNotFound {
		t.Errorf("e.StatusCode() = %d, want %d", e.StatusCode(), http.StatusNotFound)
	}
	if e.Error() != http.StatusText(http.StatusNotFound) {
		t.Errorf("e.Error() = %s, want %s", e.Error(), http.StatusText(http.StatusNotFound))
	}

	e = NewHTTPError(http.StatusNotFound, "abc")
	if e.StatusCode() != http.StatusNotFound {
		t.Errorf("e.StatusCode() = %d, want %d", e.StatusCode(), http.StatusNotFound)
	}
	if e.Error() != "abc" {
		t.Errorf("e.Error() = %s, want %s", e.Error(), "abc")
	}

	s, _ := json.Marshal(e)
	expected := `{"status":404,"message":"abc"}`
	if string(s) != expected {
		t.Errorf("json.Marshal(e) = %s, want %s", string(s), expected)
	}
}
