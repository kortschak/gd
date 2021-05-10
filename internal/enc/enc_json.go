// Copyright ©2020 Dan Kortschak. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build gd

package enc

import (
	"encoding/json"
	"os"
	"sync"
)

var (
	mu     sync.Mutex
	render func(interface{}) error
)

func init() {
	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	render = enc.Encode
}
