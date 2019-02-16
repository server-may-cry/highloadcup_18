package listhelper_test

import (
	"testing"

	"github.com/server-may-cry/highloadcup_18/listhelper"
)

func BenchmarkIntersection(b *testing.B) {
	var list1 = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	var list2 = []int{123, 145, 1, 3, 46, 7, 89, 90}
	for i := 0; i < b.N; i++ {
		listhelper.Intersection(list1, list2)
	}
}

func BenchmarkRemoveDuplicates(b *testing.B) {
	var list3 = []int{1, 2, 3, 556, 6, 7, 74, 53, 32, 423, 432, 5, 346, 43, 2, 4356, 45, 6, 1, 2, 3, 556, 6, 7, 74, 53, 32, 423, 432, 5, 346, 43, 2, 4356, 45, 6, 1, 2, 3, 556, 6, 7, 74, 53, 32, 423, 432, 5, 346, 43, 2, 4356, 45, 6, 432}
	for i := 0; i < b.N; i++ {
		listhelper.RemoveDuplicates(list3)
	}
}

func TestIntersection(t *testing.T) {
	var listIntersectionTestData = []struct {
		a []int
		b []int
		r []int
	}{
		{
			[]int{1, 2, 3, 8},
			[]int{2, 3, 4, 9},
			[]int{2, 3},
		},
		{
			[]int{1, 2, 3},
			[]int{2},
			[]int{2},
		},
	}
	for _, v := range listIntersectionTestData {
		intersected := listhelper.Intersection(v.a, v.b)
		if len(intersected) != len(v.r) {
			t.Errorf("got %v, want %v", intersected, v.r)
		}
		for i, v := range v.r {
			if v != intersected[i] {
				t.Errorf("got %d expected %d", intersected[i], v)
			}
		}
	}
}
