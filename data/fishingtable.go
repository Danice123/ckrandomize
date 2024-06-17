package data

import "fmt"

type FishingTables struct {
	Slot       int
	Rod        int
	Tables     []FishingTable
	Timegroups []FishingTimegroup
}

func (ths FishingTables) GetTable(i int) []Wildmon {
	var wm []Wildmon
	for _, entry := range ths.Tables[ths.Slot].GetTable(ths.Rod) {
		if entry.isTimegroup {
			wm = append(wm, ths.Timegroups[entry.timegroup].GetTable()[i])
		} else {
			wm = append(wm, entry.wm)
		}
	}
	return wm
}

func (ths FishingTables) LenTables() int {
	return 2
}

func (ths FishingTables) ModifyTable(pokemon int, table int, index int) Tables {
	entry := ths.Tables[ths.Slot].GetTable(ths.Rod)[index]
	if entry.isTimegroup {
		ths.Timegroups[entry.timegroup] = ths.Timegroups[entry.timegroup].ModifyTable(pokemon, table)
	} else {
		ths.Tables[ths.Slot] = ths.Tables[ths.Slot].ModifyTable(pokemon, ths.Rod, index)
	}
	return ths
}

const FISHING_TABLE_ENTRY_LENGTH = 3
const OLD_ROD_LENGTH = FISHING_TABLE_ENTRY_LENGTH * 3
const GOOD_ROD_LENGTH = FISHING_TABLE_ENTRY_LENGTH * 4
const SUPER_ROD_LENGTH = FISHING_TABLE_ENTRY_LENGTH * 4

type FishingTable [OLD_ROD_LENGTH + GOOD_ROD_LENGTH + SUPER_ROD_LENGTH]byte

func (ths FishingTable) GetTable(table int) []FishingTableEntry {
	switch table {
	case 0: // Old Rod
		return ths.readTable(0, 3)
	case 1: // Good Rod
		return ths.readTable(OLD_ROD_LENGTH, 4)
	case 2: // Super Rod
		return ths.readTable(OLD_ROD_LENGTH+GOOD_ROD_LENGTH, 4)
	default:
		panic("index out of range")
	}
}

func (ths FishingTable) ModifyTable(pokemon int, table int, index int) FishingTable {
	switch table {
	case 0: // Old Rod
		ths[FISHING_TABLE_ENTRY_LENGTH*index+1] = byte(pokemon)
	case 1: // Good Rod
		ths[OLD_ROD_LENGTH+FISHING_TABLE_ENTRY_LENGTH*index+1] = byte(pokemon)
	case 2: // Super Rod
		ths[OLD_ROD_LENGTH+GOOD_ROD_LENGTH+FISHING_TABLE_ENTRY_LENGTH*index+1] = byte(pokemon)
	default:
		panic("index out of range")
	}
	return ths
}

func (ths FishingTable) readTable(offset int, length int) []FishingTableEntry {
	t := make([]FishingTableEntry, length)
	for i := 0; i < len(t); i++ {
		pNum := int(ths[offset+FISHING_TABLE_ENTRY_LENGTH*i+1])
		if pNum == 0 {
			t[i] = FishingTableEntry{
				isTimegroup: true,
				timegroup:   int(ths[offset+FISHING_TABLE_ENTRY_LENGTH*i+2]),
			}
		} else {
			t[i] = FishingTableEntry{
				wm: Wildmon{
					Level: int(ths[offset+FISHING_TABLE_ENTRY_LENGTH*i+2]),
					Name:  PokeConsts[int(ths[offset+FISHING_TABLE_ENTRY_LENGTH*i+1])],
				},
			}
			if t[i].wm.Name == "" {
				panic(fmt.Sprintf("MISSING NAME FOR %d", int(ths[offset+FISHING_TABLE_ENTRY_LENGTH*i+1])))
			}
		}
	}
	return t
}

type FishingTableEntry struct {
	isTimegroup bool
	timegroup   int
	wm          Wildmon
}

const FISHING_TIMEGROUP_LENGTH = 4

type FishingTimegroup [FISHING_TIMEGROUP_LENGTH]byte

func (ths FishingTimegroup) GetTable() []Wildmon {
	return []Wildmon{
		{
			Name:  PokeConsts[int(ths[0])],
			Level: int(ths[1]),
		},
		{
			Name:  PokeConsts[int(ths[2])],
			Level: int(ths[3]),
		},
	}
}

func (ths FishingTimegroup) ModifyTable(pokemon int, table int) FishingTimegroup {
	switch table {
	case 0: // Day
		ths[0] = byte(pokemon)
	case 1: // Night
		ths[2] = byte(pokemon)
	default:
		panic("index out of range")
	}
	return ths
}
