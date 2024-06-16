package data

const GRASS_WILDMON_HEADER_LENGTH = 5
const GRASS_WILDMON_TIME_TABLE_LENGTH = 14
const MORNING = 0
const DAY = 1
const NIGHT = 2

type GrassWildmonTable [GRASS_WILDMON_HEADER_LENGTH + GRASS_WILDMON_TIME_TABLE_LENGTH*3]byte

func (ths GrassWildmonTable) Map() string {
	return MapConsts[int(ths[0])][int(ths[1])]
}

func (ths GrassWildmonTable) MorningTable() []Wildmon {
	return readWildmonTable(ths[:], 7, GRASS_WILDMON_HEADER_LENGTH)
}

func (ths GrassWildmonTable) DayTable() []Wildmon {
	return readWildmonTable(ths[:], 7, GRASS_WILDMON_HEADER_LENGTH+GRASS_WILDMON_TIME_TABLE_LENGTH)
}

func (ths GrassWildmonTable) NightTable() []Wildmon {
	return readWildmonTable(ths[:], 7, GRASS_WILDMON_HEADER_LENGTH+GRASS_WILDMON_TIME_TABLE_LENGTH*2)
}

func (ths GrassWildmonTable) ModifyTable(pokemon int, table int, index int) GrassWildmonTable {
	ths[GRASS_WILDMON_HEADER_LENGTH+GRASS_WILDMON_TIME_TABLE_LENGTH*table+2*index+1] = byte(pokemon)
	return ths
}

const WATER_WILDMON_HEADER_LENGTH = 3
const WATER_WILDMON_TABLE_LENGTH = 6

type WaterWildmonTable [WATER_WILDMON_HEADER_LENGTH + WATER_WILDMON_TABLE_LENGTH]byte

func (ths WaterWildmonTable) Map() string {
	return MapConsts[int(ths[0])][int(ths[1])]
}

func (ths WaterWildmonTable) Table() []Wildmon {
	return readWildmonTable(ths[:], 3, WATER_WILDMON_HEADER_LENGTH)
}

func (ths WaterWildmonTable) ModifyTable(pokemon int, index int) WaterWildmonTable {
	ths[WATER_WILDMON_HEADER_LENGTH+2*index+1] = byte(pokemon)
	return ths
}

type Wildmon struct {
	Name  string
	Level int
}

func readWildmonTable(table []byte, lenWildmons int, offset int) []Wildmon {
	t := make([]Wildmon, lenWildmons)
	for i := 0; i < len(t); i++ {
		t[i] = Wildmon{
			Level: int(table[offset+i*2]),
			Name:  PokeConsts[int(table[offset+i*2+1])],
		}
	}
	return t
}
