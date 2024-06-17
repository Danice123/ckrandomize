package data

import "fmt"

type TreemonTable []TreemonPool

func (ths TreemonTable) GetTable(i int) []Wildmon {
	switch i {
	case 0: // Spearow headbutt
		return ths[0].GetTable()
	case 1: // Hoothoot headbutt
		return ths[6].GetTable()
	case 2: // Rocksmash
		return ths[14].GetTable()
	default:
		panic("index out of range")
	}
}

func (ths TreemonTable) LenTables() int {
	return 3
}

func (ths TreemonTable) ModifyTable(pokemon int, table int, index int) Tables {
	switch table {
	case 0: // Spearow headbutt
		ths[0] = ths[0].ModifyTable(pokemon, -1, index)
		ths[1] = ths[1].ModifyTable(pokemon, -1, index)
		ths[4] = ths[4].ModifyTable(pokemon, -1, index)
		ths[5] = ths[5].ModifyTable(pokemon, -1, index)
		ths[8] = ths[8].ModifyTable(pokemon, -1, index)
		ths[9] = ths[9].ModifyTable(pokemon, -1, index)
		ths[10] = ths[10].ModifyTable(pokemon, -1, index)
		ths[11] = ths[11].ModifyTable(pokemon, -1, index)
		ths[12] = ths[12].ModifyTable(pokemon, -1, index)
		ths[13] = ths[13].ModifyTable(pokemon, -1, index)
	case 1: // Hoothoot headbutt
		ths[6] = ths[6].ModifyTable(pokemon, -1, index)
		ths[7] = ths[7].ModifyTable(pokemon, -1, index)
		if index != 5 {
			ths[2] = ths[2].ModifyTable(pokemon, -1, index)
			ths[3] = ths[3].ModifyTable(pokemon, -1, index)
			if index == 4 {
				ths[2] = ths[2].ModifyTable(pokemon, -1, index+1)
				ths[3] = ths[3].ModifyTable(pokemon, -1, index+1)
			}
		}
	case 2: // Rocksmash
		ths[14] = ths[14].ModifyTable(pokemon, -1, index)
	default:
		panic("index out of range")
	}
	return ths
}

type TreemonPool []byte

func (ths TreemonPool) GetTable() []Wildmon {
	t := make([]Wildmon, len(ths)/3)
	for i := 0; i < len(t); i++ {
		t[i] = Wildmon{
			Level: int(ths[3*i+2]),
			Name:  PokeConsts[int(ths[3*i+1])],
		}
		if t[i].Name == "" {
			panic(fmt.Sprintf("MISSING NAME FOR %d", int(ths[3*i+1])))
		}
	}
	return t
}

func (ths TreemonPool) LenTables() int {
	return 1
}

func (ths TreemonPool) ModifyTable(pokemon int, _ int, index int) TreemonPool {
	ths[3*index+1] = byte(pokemon)
	return ths
}
