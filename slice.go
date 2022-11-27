// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuidslice

import (
	"sync"

	"github.com/google/uuid"
)

// StringCopy fills the destination slice with UUIDs parsed from the source slice.
func StringCopy(dst []uuid.UUID, src []string) (int, error) {
	if len(dst) == 0 || len(src) == 0 {
		return 0, nil
	}

	var n int

	for i := 0; i < len(src); i++ {
		id, err := uuid.Parse(src[i])
		if err != nil {
			return 0, err
		}

		dst[n] = id

		n++
		if len(dst) < n+1 {
			break
		}
	}

	return n, nil
}

// UniqueCopy fills the destination slice with unique UUIDs from the source slice, preserving order.
func UniqueCopy(dst, src []uuid.UUID) int {
	if len(dst) == 0 || len(src) == 0 {
		return 0
	}

	uniqueness := pool.Get().(map[uuid.UUID]struct{})
	defer pool.Put(uniqueness)

	for id := range uniqueness {
		delete(uniqueness, id)
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

// ExceptCopy fills the destination slice with UUIDs from first srouce
// excepts UUIDs from another source, preserving order.
func ExceptCopy(dst, src []uuid.UUID, src2 []uuid.UUID) int {
	if len(dst) == 0 || len(src) == 0 {
		return 0
	}

	exceptness := pool.Get().(map[uuid.UUID]struct{})
	defer pool.Put(exceptness)

	for id := range exceptness {
		delete(exceptness, id)
	}

	for i := 0; i < len(src2); i++ {
		exceptness[src2[i]] = struct{}{}
	}

	var n int

	for i := 0; i < len(src); i++ {
		if _, ok := exceptness[src[i]]; ok {
			continue
		}

		dst[n] = src[i]

		n++
		if len(dst) < n+1 {
			break
		}
	}

	return n
}

var pool = sync.Pool{New: func() interface{} { return make(map[uuid.UUID]struct{}) }}
