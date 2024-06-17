package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Danice123/ckrandomize/data"
	"github.com/spf13/cobra"
)

const BEGINNING_OF_JOHTO_GRASS_WILDMON_TABLE = 0x2A5F4
const BEGINNING_OF_JOHTO_WATER_WILDMON_TABLE = 0x2B128
const BEGINNING_OF_KANTO_GRASS_WILDMON_TABLE = 0x2B27F
const BEGINNING_OF_KANTO_WATER_WILDMON_TABLE = 0x2B802
const BEGINNING_OF_TREEMON_TABLE = 0xB82FA
const END_OF_TREEMON_TABLE = 0xB840A
const BEGINNING_OF_FISHING_TABLE = 0x924E3
const LEN_FISHING_TABLE = 13
const LEN_FISHING_TIMEGROUPS = 15

const CYNDAQUIL_DEF = 0x78C8A
const CYNDAQUIL_NAME = 0x78C7F
const CYNDAQUIL_CRY = 0x78C68
const CYNDAQUIL_PIC = 0x78C66

const TOTADILE_DEF = 0x78CCC
const TOTADILE_NAME = 0x78CC1
const TOTADILE_CRY = 0x78CCA
const TOTADILE_PIC = 0x78CA8

const CHIKORITA_DEF = 0x78D08
const CHIKORITA_NAME = 0x78CFD
const CHIKORITA_CRY = 0x78CE6
const CHIKORITA_PIC = 0x78CE4

type RomData struct {
	johtoGrass []data.GrassWildmonTable
	johtoWater []data.WaterWildmonTable
	kantoGrass []data.GrassWildmonTable
	kantoWater []data.WaterWildmonTable
	treemon    data.TreemonTable
	fishing    data.FishingTables
	gifts      []data.GiftTable
}

func (data *RomData) ReadRom(filepath string) error {
	rom, err := os.OpenFile(filepath, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}

	data.johtoGrass, err = ReadGrassWM(rom, BEGINNING_OF_JOHTO_GRASS_WILDMON_TABLE)
	if err != nil {
		return err
	}
	data.johtoWater, err = ReadWaterWM(rom, BEGINNING_OF_JOHTO_WATER_WILDMON_TABLE)
	if err != nil {
		return err
	}
	data.kantoGrass, err = ReadGrassWM(rom, BEGINNING_OF_KANTO_GRASS_WILDMON_TABLE)
	if err != nil {
		return err
	}
	data.kantoWater, err = ReadWaterWM(rom, BEGINNING_OF_KANTO_WATER_WILDMON_TABLE)
	if err != nil {
		return err
	}
	data.treemon, err = ReadTreemon(rom, BEGINNING_OF_TREEMON_TABLE, END_OF_TREEMON_TABLE)
	if err != nil {
		return err
	}
	data.fishing, err = ReadFishing(rom, BEGINNING_OF_FISHING_TABLE, LEN_FISHING_TABLE, LEN_FISHING_TIMEGROUPS)
	if err != nil {
		return err
	}
	data.gifts, err = ReadGifts(rom)
	if err != nil {
		return err
	}
	return rom.Close()
}

func init() {
	rootCmd.AddCommand(readCmd)
}

var readCmd = &cobra.Command{
	Use:   "read [target] [output]",
	Short: "READ",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		rom := RomData{}
		err := rom.ReadRom(args[0])
		if err != nil {
			return err
		}

		emi, err := data.ReadEmiEncounterList()
		if err != nil {
			return err
		}

		randomizer := data.Randomizer{}
		for i := range rom.gifts {
			rom.gifts[i] = randomizer.Randomize(rom.gifts[i]).(data.GiftTable)
			emi.Translate(rom.gifts[i])
		}
		rom.treemon = randomizer.Randomize(rom.treemon).(data.TreemonTable)
		emi.TranslateHeadbutt(rom.treemon)
		for i := range rom.johtoGrass {
			rom.johtoGrass[i] = randomizer.Randomize(rom.johtoGrass[i]).(data.GrassWildmonTable)
			emi.Translate(rom.johtoGrass[i])
		}
		for i := range rom.johtoWater {
			rom.johtoWater[i] = randomizer.Randomize(rom.johtoWater[i]).(data.WaterWildmonTable)
			emi.Translate(rom.johtoWater[i])
		}
		for i := range rom.kantoGrass {
			rom.kantoGrass[i] = randomizer.Randomize(rom.kantoGrass[i]).(data.GrassWildmonTable)
			emi.Translate(rom.kantoGrass[i])
		}
		for i := range rom.kantoWater {
			rom.kantoWater[i] = randomizer.Randomize(rom.kantoWater[i]).(data.WaterWildmonTable)
			emi.Translate(rom.kantoWater[i])
		}
		for rom.fishing.Slot = 0; rom.fishing.Slot < LEN_FISHING_TABLE; rom.fishing.Slot++ {
			for rom.fishing.Rod = 0; rom.fishing.Rod < 3; rom.fishing.Rod++ {
				rom.fishing = randomizer.Randomize(rom.fishing).(data.FishingTables)
			}
		}
		emi.TranslateFishing(rom.fishing)

		emiout, err := json.MarshalIndent(emi, "", "\t")
		if err != nil {
			return err
		}
		err = os.WriteFile("emiout.json", emiout, 0777)
		if err != nil {
			return err
		}

		if len(args) > 1 {
			f, err := os.ReadFile(args[0])
			if err != nil {
				return err
			}

			err = os.WriteFile(args[1], f, 0777)
			if err != nil {
				return err
			}

			wrrom, err := os.OpenFile(args[1], os.O_RDWR, 0777)
			if err != nil {
				return err
			}

			_, err = wrrom.Seek(BEGINNING_OF_JOHTO_GRASS_WILDMON_TABLE, 0)
			if err != nil {
				return err
			}
			for _, wm := range rom.johtoGrass {
				_, err = wrrom.Write(wm[:])
				if err != nil {
					return err
				}
			}
		}
		return nil
	},
}

func ReadGrassWM(rom *os.File, offset int64) ([]data.GrassWildmonTable, error) {
	_, err := rom.Seek(offset, 0)
	if err != nil {
		return nil, err
	}
	wildmons := []data.GrassWildmonTable{}
	for {
		wm := data.GrassWildmonTable{}
		// offset, _ := rom.Seek(0, 1)
		// fmt.Printf("%x: ", offset)
		_, err = rom.Read(wm[:])
		if err != nil {
			return nil, err
		}
		if wm[0] == 0xFF {
			break
		}
		wildmons = append(wildmons, wm)
		// fmt.Printf("%s:\n", wm.Map())
	}
	return wildmons, nil
}

func ReadWaterWM(rom *os.File, offset int64) ([]data.WaterWildmonTable, error) {
	_, err := rom.Seek(offset, 0)
	if err != nil {
		return nil, err
	}
	wildmons := []data.WaterWildmonTable{}
	for {
		wm := data.WaterWildmonTable{}
		// offset, _ := rom.Seek(0, 1)
		// fmt.Printf("%x: ", offset)
		_, err = rom.Read(wm[:])
		if err != nil {
			return nil, err
		}
		if wm[0] == 0xFF {
			break
		}
		wildmons = append(wildmons, wm)
		// fmt.Printf("%s:%v\n", wm.Map(), wm.Table())
	}
	return wildmons, nil
}

func ReadTreemon(rom *os.File, offset int64, end int64) (data.TreemonTable, error) {
	_, err := rom.Seek(offset, 0)
	if err != nil {
		return nil, err
	}
	current := offset
	var t data.TreemonTable
	for {
		var pool data.TreemonPool
		for {
			next := [1]byte{}
			_, err = rom.Read(next[:])
			if err != nil {
				return nil, err
			}
			pool = append(pool, next[0])
			current++
			if next[0] == 0xFF {
				break
			}
			nextEntry := [2]byte{}
			_, err = rom.Read(nextEntry[:])
			if err != nil {
				return nil, err
			}
			pool = append(pool, nextEntry[:]...)
			current += 2
		}
		t = append(t, pool)
		if current >= end {
			break
		}
	}
	return t, nil
}

func ReadFishing(rom *os.File, offset int64, length int, tgLength int) (data.FishingTables, error) {
	_, err := rom.Seek(offset, 0)
	if err != nil {
		return data.FishingTables{}, err
	}
	tables := data.FishingTables{
		Tables:     []data.FishingTable{},
		Timegroups: []data.FishingTimegroup{},
	}
	for i := 0; i < length; i++ {
		var t data.FishingTable
		_, err := rom.Read(t[:])
		if err != nil {
			return data.FishingTables{}, err
		}
		tables.Tables = append(tables.Tables, t)
	}
	for i := 0; i < tgLength; i++ {
		var tg data.FishingTimegroup
		_, err := rom.Read(tg[:])
		if err != nil {
			return data.FishingTables{}, err
		}
		tables.Timegroups = append(tables.Timegroups, tg)
	}
	return tables, nil
}

func ReadGifts(rom *os.File) ([]data.GiftTable, error) {
	starter := data.GiftTable{Area: "new-bark-town"}
	for _, offset := range []int64{CHIKORITA_DEF, TOTADILE_DEF, CYNDAQUIL_DEF} {
		var b [1]byte
		_, err := rom.ReadAt(b[:], offset)
		if err != nil {
			return nil, err
		}
		starter.Data = append(starter.Data, b[0])
	}
	return []data.GiftTable{starter}, nil
}

func PrintTable(t []data.Wildmon) {
	for _, wm := range t {
		fmt.Printf(" %s,%d\t|", wm.Name, wm.Level)
	}
	fmt.Println("")
}
