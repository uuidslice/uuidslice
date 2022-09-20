// Copyright 2022 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuidslice_test

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/pfmt/pfmt"
	"github.com/pfmt/uuidslice"
)

var uniqueCopyTests = []struct {
	name  string
	line  string
	src   []uuid.UUID
	dst   []uuid.UUID
	want  []uuid.UUID
	bench bool
	skip  bool
	keep  bool
}{
	{
		name:  "non unique",
		line:  testline(),
		src:   []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7")), uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a"))},
		dst:   []uuid.UUID{uuid.UUID{}, uuid.UUID{}},
		want:  []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))},
		bench: true,
	}, {
		name: "already unique",
		line: testline(),
		src:  []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))},
		dst:  []uuid.UUID{uuid.UUID{}, uuid.UUID{}},
		want: []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))},
	}, {
		name: "non unique",
		line: testline(),
		src:  []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a"))},
		dst:  []uuid.UUID{uuid.UUID{}, uuid.UUID{}, uuid.UUID{}},
		want: []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a"))},
	}, {
		name: "without destination",
		line: testline(),
		src:  []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))},
		dst:  nil,
		want: nil,
	}, {
		name: "empty destination",
		line: testline(),
		src:  []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))},
		dst:  []uuid.UUID{},
		want: []uuid.UUID{},
	}, {
		name: "short destination",
		line: testline(),
		src:  []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))},
		dst:  []uuid.UUID{uuid.UUID{}, uuid.UUID{}},
		want: []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))},
	}, {
		name: "very short destination",
		line: testline(),
		src:  []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))},
		dst:  []uuid.UUID{uuid.UUID{}},
		want: []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a"))},
	},
}

func TestUniqueCopy(t *testing.T) {
	keep := uniqueCopyTests[:0]
	skip := uniqueCopyTests[:0]
	for _, tt := range uniqueCopyTests {
		if tt.keep {
			keep = append(keep, tt)
		} else {
			skip = append(skip, tt)
		}
	}

	if len(keep) == 0 {
		keep = uniqueCopyTests
	} else {
		for _, tt := range skip {
			t.Logf("%s/unkeep: %s", tt.line, tt.name)
		}
	}

	for _, tt := range keep {
		if tt.skip {
			t.Logf("%s/skip: %s", tt.line, tt.name)
			continue
		}

		tt := tt

		t.Run(tt.line+"/"+tt.name, func(t *testing.T) {
			t.Parallel()

			n := uuidslice.UniqueCopy(tt.dst, tt.src)
			got := tt.dst[:n]

			if !cmp.Equal(got, tt.want) {
				t.Errorf("\nwant: %s\n got: %s\ntest: %s", pfmt.Sprint(tt.want), got, tt.line)
			}
		})
	}
}

func BenchmarkUniqueCopy(b *testing.B) {
	b.ReportAllocs()

	keep := uniqueCopyTests[:0]
	skip := uniqueCopyTests[:0]
	for _, tt := range uniqueCopyTests {
		if tt.keep {
			keep = append(keep, tt)
		} else {
			skip = append(skip, tt)
		}
	}

	if len(keep) == 0 {
		keep = uniqueCopyTests
	} else {
		for _, tt := range skip {
			b.Logf("%s/unkeep: %s", tt.line, tt.name)
		}
	}

	for _, tt := range keep {
		if tt.skip {
			b.Logf("%s/skip: %s", tt.line, tt.name)
			continue
		}

		if !tt.bench {
			continue
		}

		b.Run(tt.line, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = uuidslice.UniqueCopy(tt.dst, tt.src)
			}
		})
	}
}

func TestUniqueCopyToSelf(t *testing.T) {
	src := []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7")), uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a"))}
	want := []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))}

	n := uuidslice.UniqueCopy(src, src)
	got := src[:n]

	if !cmp.Equal(got, want) {
		t.Errorf("\nwant: %s\n got: %s", pfmt.Sprint(want), got)
	}
}

func testline() string {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		return fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	return "it was not possible to recover file and line number information about function invocations"
}
