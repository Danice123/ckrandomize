package data

import (
	"crypto/sha1"
	"encoding/binary"
	"math/rand"
	"slices"
)

type Randomizer struct {
	occurance  map[string]int
	seededRand *rand.Rand
}

type Tables interface {
	GetTable(int) []Wildmon
	LenTables() int
	ModifyTable(pokemon int, table int, index int) Tables
}

func NewRandomizer(seed string) (*Randomizer, error) {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(seed))
	if err != nil {
		return nil, err
	}
	var seedInt uint64 = binary.BigEndian.Uint64(hasher.Sum(nil))
	return &Randomizer{
		occurance:  make(map[string]int),
		seededRand: rand.New(rand.NewSource(int64(seedInt))),
	}, nil

}

func (ths *Randomizer) Randomize(area Tables) Tables {
	dupMap := map[string]string{}
	for i := 0; i < area.LenTables(); i++ {
		t := area.GetTable(i)
		for j, wm := range t {
			base := GetBaseMon(wm)
			if base == "NOT_FOUND" {
				panic("NOT FOUND " + wm.Name)
			}
			for poolName, pl := range RandomizerPools {
				if slices.Contains(pl, base) {
					if poolName == "NORANDOM" {
						break
					}
					var newMon string
					if _, ok := dupMap[base]; ok {
						newMon = dupMap[base]
					} else {
						newMon = pl[ths.seededRand.Intn(len(pl))]
						dupMap[base] = newMon
					}
					area = area.ModifyTable(GetMonId(EvolveMon(newMon, wm.Level)), i, j)
				}
			}
		}
	}
	return area
}

func (ths *Randomizer) RollNewMon(pool []string, original string) string {
	var total int
	for _, n := range pool {
		total += ths.occurance[n]
	}
	var average int
	if total > 0 {
		average = total / len(pool)
	}

	var newMon string
	for {
		newMon = pool[ths.seededRand.Intn(len(pool))]
		if ths.occurance[newMon] > average {
			continue
		}
		if newMon != original {
			break
		}
	}
	return newMon
}

func GetBaseMon(wm Wildmon) string {
	for base, evos := range EvoMap {
		if wm.Name == base {
			return base
		}
		for evo := range evos {
			if wm.Name == evo {
				return base
			}
		}
	}
	return "NOT_FOUND"
}

func EvolveMon(name string, level int) string {
	var minEvo int
	evo := name
	if EvoMap[name] != nil && len(EvoMap[name]) > 0 {
		for n, l := range EvoMap[name] {
			if l <= level {
				if l > minEvo {
					minEvo = l
					evo = n
				}
				if l == minEvo && rand.Intn(100) > 50 {
					evo = n
				}
			}
		}
	}
	return evo
}

func GetMonId(name string) int {
	for i, n := range PokeConsts {
		if name == n {
			return i
		}
	}
	panic("CANNOT FIND ID FOR NAME: " + name)
}
