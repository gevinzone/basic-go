// Copyright 2023 igevin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package week1

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDeleteAt(t *testing.T) {
	testCases := []struct {
		name        string
		src         []int
		index       int
		expectedRes []int
		expectedErr error
	}{
		{
			name:        "nil",
			src:         nil,
			expectedErr: ErrIndexOutOfRange,
		},
		{
			name:        "empty",
			src:         []int{},
			index:       0,
			expectedErr: ErrIndexOutOfRange,
		},
		{
			name:        "boundary",
			src:         []int{1, 2, 3},
			index:       3,
			expectedErr: ErrIndexOutOfRange,
		},
		{
			name:        "out of range",
			src:         []int{1, 2, 3},
			index:       4,
			expectedErr: ErrIndexOutOfRange,
		},
		{
			name:        "first",
			src:         []int{1, 2, 3},
			index:       0,
			expectedRes: []int{2, 3},
		},
		{
			name:        "last",
			src:         []int{1, 2, 3},
			index:       2,
			expectedRes: []int{1, 2},
		},
		{
			name:        "normal",
			src:         []int{1, 2, 3},
			index:       1,
			expectedRes: []int{1, 3},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := DeleteAt(tc.src, tc.index)
			require.Equal(t, tc.expectedErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.expectedRes, res)
		})
	}
}

func TestShrink(t *testing.T) {
	testCases := []struct {
		name        string
		originCap   int
		enqueueLoop int
		expectCap   int
	}{
		{
			name:        "小于64",
			originCap:   32,
			enqueueLoop: 6,
			expectCap:   32,
		},
		{
			name:        "小于2048, 不足1/4",
			originCap:   1000,
			enqueueLoop: 20,
			expectCap:   500,
		},
		{
			name:        "小于2048, 超过1/4",
			originCap:   1000,
			enqueueLoop: 400,
			expectCap:   1000,
		},
		{
			name:        "大于2048，不足一半",
			originCap:   3000,
			enqueueLoop: 60,
			expectCap:   1875,
		},
		{
			name:        "大于2048，大于一半",
			originCap:   3000,
			enqueueLoop: 2000,
			expectCap:   3000,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			l := make([]int, 0, tc.originCap)

			for i := 0; i < tc.enqueueLoop; i++ {
				l = append(l, i)
			}
			l = Shrink[int](l)
			assert.Equal(t, tc.expectCap, cap(l))
		})
	}
}
