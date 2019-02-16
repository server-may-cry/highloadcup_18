package db

import (
	"sort"

	"github.com/server-may-cry/highloadcup_18/listhelper"
)

type Grouped struct {
	Count int
	Keys  map[string]string
}

type toGroup struct {
	Type  string
	found map[string][]int
}

func (s *Storage) Group(f Filter, keys []string, limit int, ascending bool) []Grouped {
	filtered := s.FilterNoDataFetch(f)
	groups := make([]toGroup, 0, len(keys))
	for _, k := range keys {
		switch k {
		case "sex":
			group := listhelper.GroupIntersection(filtered, s.sexIndex)
			groups = append(groups, toGroup{
				Type:  k,
				found: group,
			})
		case "status":
			group := listhelper.GroupIntersection(filtered, s.statusIndex)
			groups = append(groups, toGroup{
				Type:  k,
				found: group,
			})
		case "interests":
			group := listhelper.GroupIntersection(filtered, s.interestsIndex)
			groups = append(groups, toGroup{
				Type:  k,
				found: group,
			})
		case "country":
			group := listhelper.GroupIntersection(filtered, s.countryIndex)
			groups = append(groups, toGroup{
				Type:  k,
				found: group,
			})
		case "city":
			group := listhelper.GroupIntersection(filtered, s.cityIndex)
			groups = append(groups, toGroup{
				Type:  k,
				found: group,
			})
		}
	}
	result := countGroupBy(groups)
	result = sortGroupResult(result, ascending, keys)
	l := len(result)
	if limit > l {
		limit = l
	}
	return result[:limit]
}

func sortGroupResult(result []Grouped, ascending bool, keys []string) []Grouped {
	sort.SliceStable(result, func(i, j int) bool {
		a := result[i]
		b := result[j]

		if a.Count == b.Count {
			for _, k := range keys {
				aa, aaOk := a.Keys[k]
				bb, bbOk := b.Keys[k]
				if !bbOk && !aaOk {
					continue
				}
				if !bbOk {
					return ascending
				}
				if aa != bb {
					if ascending {
						return aa > bb
					}
					return aa < bb
				}
			}
		}
		if ascending {
			return a.Count > b.Count
		}
		return a.Count < b.Count
	})
	return result
}

type groupedChunk struct {
	ids  []int
	path map[string]string
}

func countGroupBy(groups []toGroup) []Grouped {
	if len(groups) == 0 {
		return []Grouped{}
	}
	first := groups[0]
	var chunks []groupedChunk
	for i, ids := range first.found {
		path := map[string]string{}
		if i != "" {
			path[first.Type] = i
		}
		chunk := groupedChunk{
			ids:  ids,
			path: path,
		}
		if len(chunk.ids) > 0 {
			chunks = append(chunks, chunk)
		}
	}
	for _, toGroup := range groups[1:] {
		var newChunks []groupedChunk
		for i, ids := range toGroup.found {
			for _, chunk := range chunks {
				newPath := make(map[string]string, len(chunk.path))
				for k, v := range chunk.path {
					newPath[k] = v
				}
				cp := groupedChunk{
					ids:  listhelper.Intersection(chunk.ids, ids),
					path: newPath,
				}
				if i != "" {
					cp.path[toGroup.Type] = i
				}
				if len(cp.ids) > 0 {
					newChunks = append(newChunks, cp)
				}
			}
		}
		chunks = newChunks
	}
	return chunksToGroupeds(chunks)
}

func chunksToGroupeds(gc []groupedChunk) []Grouped {
	g := make([]Grouped, 0, len(gc))
	for _, v := range gc {
		g = append(g, Grouped{
			Count: len(v.ids),
			Keys:  v.path,
		})
	}
	return g
}
