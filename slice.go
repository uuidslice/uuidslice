// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuidslice

import (
	"sync"

	"github.com/google/uuid"
)

// UniqueCopy fills the destination array with UUIDs from the source array, preserving order.
func UniqueCopy(dst, src []uuid.UUID) int {
	if len(dst) == 0 || len(src) == 0 {
		return 0
	}

	uniqueness := pool.Get().(map[uuid.UUID]struct{})
	defer pool.Put(uniqueness)

	for number := range uniqueness {
		delete(uniqueness, number)
	}

	var n int

	for i := 0; i < len(src); i++ {
		if _, ok := uniqueness[src[i]]; ok {
			continue
		}
		uniqueness[src[i]] = struct{}{}

		dst[n] = src[i]

		n++
		if len(dst) < n+1 {
			break
		}
	}

	return n
}

var pool = sync.Pool{New: func() interface{} { return make(map[uuid.UUID]struct{}) }}
