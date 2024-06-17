package data

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/Danice123/emidocgen/package/ckp"
)

type EmiEncounterFile struct {
	Pools      ckp.EncounterPools `json:"encounter_pools"`
	Encounters []ckp.Encounter    `json:"encounters"`
}

func ReadEmiEncounterList() (*EmiEncounterFile, error) {
	b, err := os.ReadFile("ref/emi_encounter.json")
	if err != nil {
		return nil, err
	}
	var data EmiEncounterFile
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

type TablesForEmi interface {
	Map() string
	Type() string
	GetTable(int) []Wildmon
	LenTables() int
}

func (ths *EmiEncounterFile) Translate(area TablesForEmi) {
	translatedName := strings.ReplaceAll(strings.ToLower(area.Map()), "_", "-")
	for encounter := range ths.Encounters {
		if ths.Encounters[encounter].Area == translatedName {
			for pool := range ths.Encounters[encounter].Pools {
				if ths.Encounters[encounter].Pools[pool].Type == area.Type() {
					if area.LenTables() > 1 {
						for tod := 0; tod < area.LenTables(); tod++ {
							for i, wm := range area.GetTable(tod) {
								ths.Encounters[encounter].Pools[pool].PoolMap[translateTimeOfDay(tod)][i].Pokemon = translateName(wm.Name)
							}
						}
					} else {
						for i, wm := range area.GetTable(0) {
							ths.Encounters[encounter].Pools[pool].PoolSlice[i].Pokemon = translateName(wm.Name)
						}
					}
				}
			}
		}
	}
}

func (ths *EmiEncounterFile) TranslateHeadbutt(t TreemonTable) {
	for name := range ths.Pools.Headbutt {
		switch name {
		case "city":
			fallthrough
		case "forest":
			fallthrough
		case "kanto":
			fallthrough
		case "lake":
			fallthrough
		case "town":
			for i, wm := range t.GetTable(0) {
				ths.Pools.Headbutt[name][i].Pokemon = translateName(wm.Name)
			}
		case "route":
			for i, wm := range t.GetTable(1) {
				ths.Pools.Headbutt[name][i].Pokemon = translateName(wm.Name)
			}
		case "canyon":
			for i, wm := range t.GetTable(1) {
				ths.Pools.Headbutt[name][i].Pokemon = translateName(wm.Name)
			}
			ths.Pools.Headbutt[name][5].Pokemon = translateName(t.GetTable(1)[4].Name)
		}
	}
	for i, wm := range t.GetTable(2) {
		ths.Pools.Rock["rock"][i].Pokemon = translateName(wm.Name)
	}
}

func (ths *EmiEncounterFile) TranslateFishing(fish FishingTables) {
	for name := range ths.Pools.Fishing {
		switch name {
		case "dratini":
			fish.Slot = 4
		case "dratini-2":
			fish.Slot = 8
		case "gyarados":
			fish.Slot = 7
		case "lake":
			fish.Slot = 2
		case "ocean":
			fish.Slot = 1
		case "pond":
			fish.Slot = 3
		case "qwilfish":
			fish.Slot = 10
		case "qwilfish-noswarm":
			fish.Slot = 12
		case "qwilfish-swarm":
			fish.Slot = 5
		case "remoraid":
			fish.Slot = 11
		case "remoraid-swarm":
			fish.Slot = 6
		case "shore":
			fish.Slot = 0
		case "whirlislands":
			fish.Slot = 9
		}
		fish.Rod = 0
		for tod := 0; tod < fish.LenTables(); tod++ {
			for i, wm := range fish.GetTable(tod) {
				ths.Pools.Fishing[name].Old[translateTimeOfDay(tod+1)][i].Pokemon = translateName(wm.Name)
			}
		}
		fish.Rod = 1
		for tod := 0; tod < fish.LenTables(); tod++ {
			for i, wm := range fish.GetTable(tod) {
				ths.Pools.Fishing[name].Good[translateTimeOfDay(tod+1)][i].Pokemon = translateName(wm.Name)
			}
		}
		fish.Rod = 2
		for tod := 0; tod < fish.LenTables(); tod++ {
			for i, wm := range fish.GetTable(tod) {
				ths.Pools.Fishing[name].Super[translateTimeOfDay(tod+1)][i].Pokemon = translateName(wm.Name)
			}
		}

	}
}

func translateName(name string) string {
	n := strings.ToLower(
		strings.ReplaceAll(
			strings.ReplaceAll(name, "_", "-"),
			"--", "-"),
	)
	if n == "farfetch-d" {
		return "farfetchd"
	}
	return n
}

func translateTimeOfDay(tod int) ckp.TIME_OF_DAY {
	switch tod {
	case MORNING:
		return ckp.MORNING
	case DAY:
		return ckp.DAY
	case NIGHT:
		return ckp.NIGHT
	default:
		panic("WHAT")
	}
}
