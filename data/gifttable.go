package data

import "fmt"

type GiftTable struct {
	Area string
	Data []byte
}

func (ths GiftTable) Map() string {
	return ths.Area
}

func (ths GiftTable) Type() string {
	return "gift"
}

func (ths GiftTable) GetTable(_ int) []Wildmon {
	t := make([]Wildmon, len(ths.Data))
	for i := 0; i < len(t); i++ {
		t[i] = Wildmon{
			Name: PokeConsts[int(ths.Data[i])],
		}
		if t[i].Name == "" {
			panic(fmt.Sprintf("MISSING NAME FOR %d", int(ths.Data[i])))
		}
	}
	return t
}

func (ths GiftTable) LenTables() int {
	return 1
}

func (ths GiftTable) ModifyTable(pokemon int, _ int, index int) Tables {
	ths.Data[index] = byte(pokemon)
	return ths
}
