package listhelper

func RemoveDuplicates(list []int) []int {
	xs := make([]int, len(list))
	copy(xs, list)
	found := make(map[int]struct{}, len(xs))
	j := 0
	for i, x := range xs {
		if _, ok := found[x]; !ok {
			found[x] = struct{}{}
			xs[j] = xs[i]
			j++
		}
	}
	return xs[:j]
}

func Diff(X, Y []int) []int {
	m := make(map[int]struct{}, len(Y))

	for _, y := range Y {
		m[y] = struct{}{}
	}

	ret := make([]int, 0, len(X))
	var i int
	for _, x := range X {
		if _, ok := m[x]; ok {
			continue
		}
		ret = append(ret, x)
		i++
	}

	return ret[:i]
}

// Intersection input lists must be sorted, returned list will be sorted
func Intersection(b, a []int) []int {
	aLen := len(a)
	bLen := len(b)
	shortest := aLen
	if bLen < aLen {
		shortest = bLen
	}
	if 0 == aLen || 0 == bLen {
		return a[0:0]
	}
	var l, i, j int

	r := make([]int, 0, shortest)
	for {
		if a[i] == b[j] {
			r = append(r, b[j])
			i++
			j++
			l++
		} else if a[i] < b[j] {
			i++
		} else {
			j++
		}

		if i == aLen || j == bLen {
			return r[:l]
		}
	}
}

func GroupIntersection(a []int, b map[string][]int) map[string][]int {
	group := make(map[string][]int, len(b))
	for g, ids := range b {
		intersection := Intersection(a, ids)
		group[g] = intersection
	}
	return group
}
