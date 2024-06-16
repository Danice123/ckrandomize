package cmd

import (
	"os"

	"github.com/Danice123/ckrandomize/data"
	"github.com/spf13/cobra"
)

const BEGINNING_OF_JOHTO_GRASS_WILDMON_TABLE = 0x2A5F4
const BEGINNING_OF_JOHTO_WATER_WILDMON_TABLE = 0x2B128
const BEGINNING_OF_KANTO_GRASS_WILDMON_TABLE = 0x2B27F
const BEGINNING_OF_KANTO_WATER_WILDMON_TABLE = 0x2B802

func init() {
	rootCmd.AddCommand(readCmd)
}

var readCmd = &cobra.Command{
	Use:   "read [target] [output]",
	Short: "READ",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		f, err := os.ReadFile(args[0])
		if err != nil {
			return err
		}
		err = os.WriteFile(args[1], f, 0777)
		if err != nil {
			return err
		}

		rom, err := os.OpenFile(args[1], os.O_RDWR, 0777)
		if err != nil {
			return err
		}

		_, err = ReadGrassWM(rom, BEGINNING_OF_JOHTO_GRASS_WILDMON_TABLE)
		if err != nil {
			return err
		}

		_, err = ReadWaterWM(rom, BEGINNING_OF_JOHTO_WATER_WILDMON_TABLE)
		if err != nil {
			return err
		}

		_, err = ReadGrassWM(rom, BEGINNING_OF_KANTO_GRASS_WILDMON_TABLE)
		if err != nil {
			return err
		}

		_, err = ReadWaterWM(rom, BEGINNING_OF_KANTO_WATER_WILDMON_TABLE)
		if err != nil {
			return err
		}

		// _, err = rom.Seek(BEGINNING_OF_JOHTO_GRASS_WILDMON_TABLE, 0)
		// if err != nil {
		// 	return err
		// }
		// for _, wm := range johtoGrass {
		// 	_, err = rom.Write(wm[:])
		// 	if err != nil {
		// 		return err
		// 	}
		// }

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
