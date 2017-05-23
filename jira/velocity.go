package jira

import (
	"encoding/json"
	"log"
	"math"
	"sort"
)

type VelocityEntry struct {
	Value float64 `json:"value"`
}

type VelocityStats struct {
	Estimated VelocityEntry `json:"estimated"`
	Completed VelocityEntry `json:"completed"`
}

type VelocityList map[string]VelocityStats

type Velocity struct {
	Stats VelocityList `json:"velocityStatEntries"`
}

func (v *Velocity) GetSortedVelocityEntries() []string {
	entries := sort.StringSlice{}
	for entry, _ := range v.Stats {
		entries = append(entries, entry)
	}
	entries.Sort()
	sort.Sort(sort.Reverse(entries[:]))
	return entries
}

func (v *Velocity) GetAverageVelocity(numberToAverage int) float64 {
	velocity := 0.0
	sortedVelocityEntries := v.GetSortedVelocityEntries()
	numVelocities := math.Min(float64(numberToAverage), float64(len(sortedVelocityEntries)))
	for i := 0; i < int(numVelocities); i++ {
		velocity += v.Stats[sortedVelocityEntries[i]].Completed.Value
	}
	velocity = velocity / float64(numVelocities)
	return velocity
}

func ParseVelocity(velocityJson []byte) Velocity {
	var velocity Velocity
	err := json.Unmarshal(velocityJson, &velocity)
	if err != nil {
		log.Fatal(err)
	}
	return velocity
}
