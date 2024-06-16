package data

import (
	"encoding/hex"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var MapConsts map[int]map[int]string
var PokeConsts map[int]string

func init() {
	m, err := readMapConstants()
	if err != nil {
		panic(err)
	}
	MapConsts = m

	m2, err := readPokemonConstants()
	if err != nil {
		panic(err)
	}
	PokeConsts = m2
}

func readPokemonConstants() (map[int]string, error) {
	d, err := os.ReadFile("ref/pokemon_constants.asm")
	if err != nil {
		return nil, err
	}
	regex := regexp.MustCompile(`const\s(\w+)\s+;\s(\w+)`)
	m := map[int]string{}
	for _, line := range strings.Split(string(d), "\n") {
		ret := regex.FindStringSubmatch(line)
		if ret != nil {
			n, err := hex.DecodeString(ret[2])
			if err != nil {
				return nil, err
			}
			m[int(n[0])] = ret[1]
		}
	}
	return m, nil
}

func readMapConstants() (map[int]map[int]string, error) {
	d, err := os.ReadFile("ref/map_constants.asm")
	if err != nil {
		return nil, err
	}
	groupRegex := regexp.MustCompile(`newgroup (\w+)\s+;\s+(\d+)`)
	mapRegex := regexp.MustCompile(`map_const (\w+),[\s\d+,]+;\s+(\d+)`)

	m := map[int]map[int]string{}
	var group int
	for _, line := range strings.Split(string(d), "\n") {
		grpRet := groupRegex.FindStringSubmatch(line)
		if grpRet != nil {
			// fmt.Printf("Found map group: %s\n", grpRet[1])
			group, err = strconv.Atoi(grpRet[2])
			if err != nil {
				return nil, err
			}
			m[group] = map[int]string{}
		}
		mapRet := mapRegex.FindStringSubmatch(line)
		if mapRet != nil {
			// fmt.Printf("Found map: %s\n", mapRet[1])
			mId, err := strconv.Atoi(mapRet[2])
			if err != nil {
				return nil, err
			}
			m[group][mId] = mapRet[1]
		}
	}
	return m, nil
}
