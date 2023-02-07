// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuidslice_test

import (
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
	"github.com/pfmt/pfmt"
	"github.com/pfmt/uuidslice"
)

var stringCopyTests = []struct {
	test    string
	line    string
	src     []string
	dst     []uuid.UUID
	want    []uuid.UUID
	wantErr error
	bench   bool
	skip    bool
	keep    bool
}{
	{
		test:  "UUIDs",
		line:  testline(),
		src:   []string{"f23133ea-e89f-467e-a757-ffa215332e6a", "ef4f8e2b-d723-41d0-a23a-ac74678e06a7"},
		dst:   []uuid.UUID{uuid.UUID{}, uuid.UUID{}},
		want:  []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))},
		bench: true,
	}, {
		test:    "not UUID",
		line:    testline(),
		src:     []string{"Not UUID."},
		dst:     []uuid.UUID{uuid.UUID{}},
		want:    []uuid.UUID{},
		wantErr: errors.New("invalid UUID length: 9"),
	}, {
		test: "without destination",
		line: testline(),
		src:  []string{"f23133ea-e89f-467e-a757-ffa215332e6a", "ef4f8e2b-d723-41d0-a23a-ac74678e06a7"},
		dst:  nil,
		want: nil,
	}, {
		test: "empty destination",
		line: testline(),
		src:  []string{"f23133ea-e89f-467e-a757-ffa215332e6a", "ef4f8e2b-d723-41d0-a23a-ac74678e06a7"},
		dst:  []uuid.UUID{},
		want: []uuid.UUID{},
	}, {
		test: "short destination",
		line: testline(),
		src:  []string{"f23133ea-e89f-467e-a757-ffa215332e6a", "ef4f8e2b-d723-41d0-a23a-ac74678e06a7"},
		dst:  []uuid.UUID{uuid.UUID{}},
		want: []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a"))},
	},
}

func TestStringCopy(t *testing.T) {
	t.Parallel()

	keep := stringCopyTests[:0:0]
	skip := stringCopyTests[:0:0]

	for _, tt := range stringCopyTests {
		if tt.keep {
			keep = append(keep, tt)
		} else {
			skip = append(skip, tt)
		}
	}

	if len(keep) == 0 {
		keep = stringCopyTests

	} else {
		for _, tt := range skip {
			t.Logf("%s/unkeep: %s", tt.line, tt.test)
		}
	}

	for _, tt := range keep {
		if tt.skip {
			t.Logf("%s/skip: %s", tt.line, tt.test)
			continue
		}

		tt := tt

		t.Run(tt.line+"/"+tt.test, func(t *testing.T) {
			t.Parallel()

			n, err := uuidslice.StringCopy(tt.dst, tt.src)
			if !strings.Contains(fmt.Sprint(err), fmt.Sprint(tt.wantErr)) {
				t.Errorf("\nwant error: %s\n got error: %s\ntest: %s", tt.wantErr, err, tt.line)
			}

			got := tt.dst[:n]

			if !cmp.Equal(got, tt.want) {
				t.Errorf("\nwant: %s\n got: %s\ntest: %s", pfmt.Sprint(tt.want), got, tt.line)
			}
		})
	}
}

func BenchmarkStringCopy(b *testing.B) {
	b.ReportAllocs()

	keep := stringCopyTests[:0:0]
	skip := stringCopyTests[:0:0]

	for _, tt := range stringCopyTests {
		if tt.keep {
			keep = append(keep, tt)
		} else {
			skip = append(skip, tt)
		}
	}

	if len(keep) == 0 {
		keep = stringCopyTests

	} else {
		for _, tt := range skip {
			b.Logf("%s/unkeep: %s", tt.line, tt.test)
		}
	}

	for _, tt := range keep {
		if tt.skip {
			b.Logf("%s/skip: %s", tt.line, tt.test)
			continue
		}

		if !tt.bench {
			continue
		}

		b.Run(tt.line, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = uuidslice.StringCopy(tt.dst, tt.src)
			}
		})
	}
}

var uniqueCopyTests = []struct {
	test  string
	line  string
	src   []uuid.UUID
	dst   []uuid.UUID
	want  []uuid.UUID
	bench bool
	skip  bool
	keep  bool
}{
	{
		test:  "not unique",
		line:  testline(),
		src:   []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7")), uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a"))},
		dst:   []uuid.UUID{uuid.UUID{}, uuid.UUID{}},
		want:  []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))},
		bench: true,
	}, {
		test: "already unique",
		line: testline(),
		src:  []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))},
		dst:  []uuid.UUID{uuid.UUID{}, uuid.UUID{}},
		want: []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))},
	}, {
		test: "not unique",
		line: testline(),
		src:  []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a"))},
		dst:  []uuid.UUID{uuid.UUID{}, uuid.UUID{}, uuid.UUID{}},
		want: []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a"))},
	}, {
		test: "without destination",
		line: testline(),
		src:  []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))},
		dst:  nil,
		want: nil,
	}, {
		test: "empty destination",
		line: testline(),
		src:  []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))},
		dst:  []uuid.UUID{},
		want: []uuid.UUID{},
	}, {
		test: "short destination",
		line: testline(),
		src:  []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))},
		dst:  []uuid.UUID{uuid.UUID{}},
		want: []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a"))},
	},
}

func TestUniqueCopy(t *testing.T) {
	t.Parallel()

	keep := uniqueCopyTests[:0:0]
	skip := uniqueCopyTests[:0:0]

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
			t.Logf("%s/unkeep: %s", tt.line, tt.test)
		}
	}

	for _, tt := range keep {
		if tt.skip {
			t.Logf("%s/skip: %s", tt.line, tt.test)
			continue
		}

		tt := tt

		t.Run(tt.line+"/"+tt.test, func(t *testing.T) {
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

	keep := uniqueCopyTests[:0:0]
	skip := uniqueCopyTests[:0:0]
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
			b.Logf("%s/unkeep: %s", tt.line, tt.test)
		}
	}

	for _, tt := range keep {
		if tt.skip {
			b.Logf("%s/skip: %s", tt.line, tt.test)
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
	t.Parallel()

	src := []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7")), uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a"))}
	want := []uuid.UUID{uuid.Must(uuid.Parse("f23133ea-e89f-467e-a757-ffa215332e6a")), uuid.Must(uuid.Parse("ef4f8e2b-d723-41d0-a23a-ac74678e06a7"))}

	n := uuidslice.UniqueCopy(src, src)
	got := src[:n]

	if !cmp.Equal(got, want) {
		t.Errorf("\nwant: %s\n got: %s", pfmt.Sprint(want), got)
	}
}

var exceptCopyTests = []struct {
	test  string
	line  string
	src   []uuid.UUID
	src2  []uuid.UUID
	dst   []uuid.UUID
	want  []uuid.UUID
	bench bool
	skip  bool
	keep  bool
}{
	{
		test:  "foobar",
		line:  testline(),
		src:   []uuid.UUID{uuid.Must(uuid.Parse("68e12387-fafc-4ae5-aa5e-e28b424f3fbe")), uuid.Must(uuid.Parse("91144820-7bdb-4b6a-99b1-b42d1bd7c72a")), uuid.Must(uuid.Parse("1c7ae620-d921-4ac2-964a-4e7cb10d8a01"))},
		src2:  []uuid.UUID{uuid.Must(uuid.Parse("1c7ae620-d921-4ac2-964a-4e7cb10d8a01"))},
		dst:   []uuid.UUID{uuid.UUID{}, uuid.UUID{}},
		want:  []uuid.UUID{uuid.Must(uuid.Parse("68e12387-fafc-4ae5-aa5e-e28b424f3fbe")), uuid.Must(uuid.Parse("91144820-7bdb-4b6a-99b1-b42d1bd7c72a"))},
		bench: true,
	},
}

func TestExceptCopy(t *testing.T) {
	t.Parallel()

	keep := exceptCopyTests[:0:0]
	skip := exceptCopyTests[:0:0]

	for _, tt := range exceptCopyTests {
		if tt.keep {
			keep = append(keep, tt)
		} else {
			skip = append(skip, tt)
		}
	}

	if len(keep) == 0 {
		keep = exceptCopyTests

	} else {
		for _, tt := range skip {
			t.Logf("%s/unkeep: %s", tt.line, tt.test)
		}
	}

	for _, tt := range keep {
		if tt.skip {
			t.Logf("%s/skip: %s", tt.line, tt.test)
			continue
		}

		tt := tt

		t.Run(tt.line+"/"+tt.test, func(t *testing.T) {
			t.Parallel()

			n := uuidslice.ExceptCopy(tt.dst, tt.src, tt.src2)
			got := tt.dst[:n]

			if !cmp.Equal(got, tt.want) {
				t.Errorf("\nwant: %s\n got: %s\ntest: %s", pfmt.Sprint(tt.want), got, tt.line)
			}
		})
	}
}

func BenchmarkExceptCopy(b *testing.B) {
	b.ReportAllocs()

	keep := exceptCopyTests[:0:0]
	skip := exceptCopyTests[:0:0]
	for _, tt := range exceptCopyTests {
		if tt.keep {
			keep = append(keep, tt)
		} else {
			skip = append(skip, tt)
		}
	}

	if len(keep) == 0 {
		keep = exceptCopyTests
	} else {
		for _, tt := range skip {
			b.Logf("%s/unkeep: %s", tt.line, tt.test)
		}
	}

	for _, tt := range keep {
		if tt.skip {
			b.Logf("%s/skip: %s", tt.line, tt.test)
			continue
		}

		if !tt.bench {
			continue
		}

		b.Run(tt.line, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = uuidslice.ExceptCopy(tt.dst, tt.src, tt.src2)
			}
		})
	}
}

func testline() string {
	_, file, line, ok := runtime.Caller(1)
	if ok {
		return fmt.Sprintf("%s:%d", filepath.Base(file), line)
	}
	return "it was not possible to recover file and line number information about function invocations"
}
