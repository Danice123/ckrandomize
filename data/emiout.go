package data

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/Danice123/emidocgen/package/ckp"
)

type EmiEncounterFile struct {
	Encounters []ckp.Encounter `json:"encounters"`
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
