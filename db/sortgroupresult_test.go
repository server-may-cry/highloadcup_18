package db

import (
	"testing"
)

func TestSortGroupResult(t *testing.T) {
	sorted := sortGroupResult([]Grouped{
		{
			Keys: map[string]string{
				"country": "Индизия",
				"sex":     "f",
			},
		},
		{
			Keys: map[string]string{
				"sex": "f",
			},
		},
	}, false, []string{"country", "sex"})
	first := sorted[0]
	if len(first.Keys) != 1 {
		t.Errorf("wrong sorted order %+v", sorted)
	}
}

func TestSortGroupResult2(t *testing.T) {
	sorted := sortGroupResult([]Grouped{
		{
			Keys: map[string]string{
				"city": "Зеленокомск",
			},
		},
		{
			Keys: map[string]string{
				"city": "Новоква",
			},
		},
	}, true, []string{"city"})
	first := sorted[0]
	if first.Keys["city"] != "Новоква" {
		t.Errorf("wrong sorted order %+v", sorted)
	}
}

func TestSortGroupResult3(t *testing.T) {
	sorted := sortGroupResult([]Grouped{
		{
			Keys: map[string]string{
				"country": "Индизия",
				"sex":     "f",
			},
		},
		{
			Keys: map[string]string{
				"sex": "f",
			},
		},
	}, false, []string{"country", "sex"})
	first := sorted[0]
	if len(first.Keys) != 1 {
		t.Errorf("wrong sorted order %+v", sorted)
	}
}
