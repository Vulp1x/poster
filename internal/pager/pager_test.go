package pager

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Pager(t *testing.T) {
	type args struct {
		name       string
		page       uint32
		perPage    uint32
		perPageMax uint32
		wantLimit  uint32
		wantOffset uint32
	}

	tests := []args{
		{
			name:       "when all is default",
			wantLimit:  DefaultMaxPerPage,
			wantOffset: 0,
		}, {
			name:       "when per page is set",
			perPage:    3,
			wantLimit:  3,
			wantOffset: 0,
		}, {
			name:       "when per page > max per page",
			perPage:    100,
			perPageMax: 3,
			wantLimit:  3,
			wantOffset: 0,
		}, {
			name:       "when page is set",
			perPage:    10,
			page:       2,
			wantLimit:  10,
			wantOffset: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pager := Pager{}
			pager.SetPage(tt.page).SetPerPage(tt.perPage).SetPerPageMax(tt.perPageMax)

			require.Equal(t, tt.wantLimit, pager.GetLimit())
			require.Equal(t, tt.wantOffset, pager.GetOffset())
		})
	}
}

func TestNewPagePer(t *testing.T) {
	t.Run("it works", func(t *testing.T) {
		pager := NewPagePer(1, 2)

		require.EqualValues(t, 1, pager.page)
		require.EqualValues(t, 2, pager.perPage)
	})
}

func TestPager_GetPage(t *testing.T) {
	t.Run("it works", func(t *testing.T) {
		pager := New().SetPage(1)
		require.Equal(t, pager.page, pager.GetPage())
	})
}

func TestPager_GetPerPage(t *testing.T) {
	t.Run("it works", func(t *testing.T) {
		pager := New().SetPerPage(1)
		require.Equal(t, pager.perPage, pager.GetPerPage())
	})
}
