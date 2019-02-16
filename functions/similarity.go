package functions

import (
	"github.com/server-may-cry/highloadcup_18/structures"
)

func CalculateSimularity(meLikes, toLikes []structures.Like) float64 {
	var total float64

	meAverageLikes := calculateAverageLikesTime(meLikes)
	toAverageLikes := calculateAverageLikesTime(toLikes)

	for id, ts := range meAverageLikes {
		tsMatch, ok := toAverageLikes[id]
		if !ok {
			continue
		}
		timeDiff := ts - tsMatch
		if timeDiff < 0 {
			timeDiff = -timeDiff
		}

		total += 1.0 / float64(timeDiff)
	}

	return total
}

func calculateAverageLikesTime(likes []structures.Like) map[int]int {
	result := make(map[int]int, len(likes))

	for _, like := range likes {
		val, ok := result[like.ID]
		if !ok {
			result[like.ID] = like.Ts
		} else {
			result[like.ID] = (val + like.Ts) / 2
		}
	}

	return result
}
