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

type RomData struct {
	johtoGrass []data.GrassWildmonTable
	johtoWater []data.WaterWildmonTable
	kantoGrass []data.GrassWildmonTable
	kantoWater []data.WaterWildmonTable
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

func PrintTable(t []data.Wildmon) {
	for _, wm := range t {
		fmt.Printf(" %s,%d\t|", wm.Name, wm.Level)
	}
	fmt.Println("")
}
