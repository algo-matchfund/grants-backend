package database

import (
	"reflect"
	"testing"
)

func makeFundAmount(count int, amount uint64) []uint64 {
	s := make([]uint64, count)
	for i := range s {
		s[i] = amount
	}
	return s
}

func TestMatches(t *testing.T) {
	funds := []*Fund{
		{
			ProjectId: "1",
			Amount:    makeFundAmount(100, 1000),
		},
		{
			ProjectId: "2",
			Amount:    makeFundAmount(250, 1000),
		},
		{
			ProjectId: "3",
			Amount:    makeFundAmount(1000, 1000),
		},
	}

	matchingPool := 1000000

	t.Run("CalculateMatches returns correct results", func(t *testing.T) {
		want := []MatchAmount{
			{
				Contributors: 100,
				Factor:       0.0009335760631097009,
				Fund:         100000,
				Match:        9242.403024786023,
				Percent:      0.9242403024786022,
				ProjectId:    "1",
			},
			{
				Contributors: 250,
				Factor:       0.0009335760631097009,
				Fund:         250000,
				Match:        58115.10992857854,
				Percent:      5.811510992857854,
				ProjectId:    "2",
			},
			{
				Contributors: 1000,
				Factor:       0.0009335760631097009,
				Fund:         1000000,
				Match:        932642.4870466354,
				Percent:      93.26424870466353,
				ProjectId:    "3",
			},
		}

		got := CalculateMatches(funds, uint64(matchingPool))

		if !reflect.DeepEqual(got, want) {
			t.Errorf("\ngot     : %+v\n\nexpected: %+v", got, want)
		}
	})

	t.Run("GetNewMatches returns correct results", func(t *testing.T) {

		want := []MatchAmount{
			{
				Contributors: 100,
				Fund:         100000,
				Match:        9232.766142596227,
				Percent:      0.9232766142596227,
				Factor:       0.0009326026406662871,
				ProjectId:    "1",
			},
			{
				Contributors: 251,
				Fund:         255000,
				Match:        59097.19583173884,
				Percent:      5.909719583173884,
				Factor:       0.0009326026406662871,
				ProjectId:    "2",
			},
			{
				Contributors: 1000,
				Fund:         1000000,
				Match:        931670.038025665,
				Percent:      93.1670038025665,
				Factor:       0.0009326026406662871,
				ProjectId:    "3",
			},
		}

		got := GetNewMatches(funds, uint64(matchingPool), "2", uint64(5000))

		if !reflect.DeepEqual(got, want) {
			t.Errorf("\ngot     : %+v\n\nexpected: %+v", got, want)
		}
	})
}
