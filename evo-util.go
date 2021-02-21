package main

import (
	"evo-util/data"
	"flag"
	"fmt"
	"math"
	"strings"
	"unicode/utf8"
)

//go:generate go generate ./data

func getEvoPathStr(evoPath []*data.EVONode) string {
	var pathStrList []string
	for _, evoNode := range evoPath {
		entity := evoNode.Entity
		pathStr := fmt.Sprintf("%v(%v)@%v", entity.Name, entity.CName, evoNode.Ord)
		if evoNode.Ord == 0 {
			pathStr = fmt.Sprintf("%v(%v)", entity.Name, entity.CName)
		}
		if entity.EvoLock != "无" {
			pathStr = fmt.Sprintf("※%v", pathStr)
		}
		pathStrList = append(pathStrList, pathStr)
	}
	return strings.Join(pathStrList, " => ")
}

func printEntityInfo(entity *data.Entity) {
	fmt.Println()
	fmt.Println("-------------------")
	fmt.Println(entity.Phase)
	fmt.Printf("EvoLock:%v\n", entity.EvoLock)
	fmt.Println("Evo Conditions:")
	ordInfoList := []data.OrdInfo{
		{
			Ord:  1,
			EKey: entity.Key,
		},
	}
	printEntitiesInfo("", ordInfoList)
	fmt.Println("-------------------")
}

func printEntitiesInfo(cKey string, ordInfoList []data.OrdInfo) {
	mainStrTmplTmpl := "%%1v%%1v %%-%vs %%-%vs %%-%vs %%-%vs %%-%vs %%-%vs %%-%vs %%-%vs %%-%vs %%-%vs"
	titleList := []string{"", "パラメータ", "体重", "育成ミス", "ご機嫌", "しつけ", "戦闘勝利", "技数", "[デコード レベル]", "必要数"}
	spaceMaxList := []int{1, getStrSpace(titleList[1]), getStrSpace(titleList[2]), getStrSpace(titleList[3]),
		getStrSpace(titleList[4]), getStrSpace(titleList[5]), getStrSpace(titleList[6]),
		getStrSpace(titleList[7]), getStrSpace(titleList[8]), getStrSpace(titleList[9])}

	for _, ordInfo := range ordInfoList {
		entity, _ := data.EntityMap[ordInfo.EKey]
		title0Str := fmt.Sprintf("%v(%v)", entity.Name, entity.CName)
		title0StrLen := getStrSpace(title0Str)
		spaceMaxList[0] = int(math.Max(float64(spaceMaxList[0]), float64(title0StrLen)))
		evoCondStrList := strings.Split(entity.Evo, ",")
		for i := 0; i < len(evoCondStrList); i++ {
			evoCondStr := getStrSpace(evoCondStrList[i])
			idx := i + 1
			spaceMaxList[idx] = int(math.Max(float64(spaceMaxList[idx]), float64(evoCondStr)))
		}
	}

	var titleSpaceIList []interface{}
	for i, v := range titleList {
		titleSpaceIList = append(titleSpaceIList, getFillLen(spaceMaxList[i], v))
	}
	//* 1 xxx(xxx) .....
	titleIList := []interface{}{"", ""}
	for _, title := range titleList {
		titleIList = append(titleIList, title)
	}

	titleTmpl := fmt.Sprintf(mainStrTmplTmpl, titleSpaceIList...)
	//fmt.Println(titleTmpl)
	fmt.Printf(titleTmpl, titleIList...)
	fmt.Println()
	for _, ordInfo := range ordInfoList {
		var evoInfoIList []interface{}
		var evoCondSpaceIList []interface{}
		entity, _ := data.EntityMap[ordInfo.EKey]
		if entity.Key == cKey {
			evoInfoIList = append(evoInfoIList, "※")
		} else {
			evoInfoIList = append(evoInfoIList, "")
		}
		evoInfoIList = append(evoInfoIList, ordInfo.Ord)
		nameStr := fmt.Sprintf("%v(%v)", entity.Name, entity.CName)
		evoInfoIList = append(evoInfoIList, nameStr)
		evoCondSpaceIList = append(evoCondSpaceIList, getFillLen(spaceMaxList[0], nameStr))

		evoCondStrList := strings.Split(entity.Evo, ",")
		if len(evoCondStrList) == 0 {
			evoInfoIList = append(evoInfoIList, "", "", "", "", "", "", "", "", "")
			evoCondSpaceIList = append(evoCondSpaceIList, 0, 0, 0, 0, 0, 0, 0, 0, 0)
		} else if len(evoCondStrList) == 1 {
			evoInfoIList = append(evoInfoIList, entity.Evo, "", "", "", "", "", "", "", "")
			evoCondSpaceIList = append(evoCondSpaceIList, 0, 0, 0, 0, 0, 0, 0, 0, 0)
		} else {
			for i, evoCond := range evoCondStrList {
				evoInfoIList = append(evoInfoIList, evoCond)
				evoCondSpaceIList = append(evoCondSpaceIList, getFillLen(spaceMaxList[i+1], evoCond))
			}
		}
		evoCondTmpl := fmt.Sprintf(mainStrTmplTmpl, evoCondSpaceIList...)
		fmt.Printf(evoCondTmpl, evoInfoIList...)
		fmt.Println()
		//fmt.Println(evoCondTmpl, evoInfoIList)
	}

}

func printEvoCompareInfo(entity *data.Entity) {
	for _, p := range entity.P {
		pEntity, _ := data.EntityMap[p.EKey]
		fmt.Printf("==========%v(%v)==========\n", pEntity.Name, pEntity.CName)
		printEntitiesInfo(entity.Key, pEntity.N)
	}
}

func getStrSpace(name string) int {
	nameLen := utf8.RuneCountInString(name)
	for _, c := range name {
		if c >= utf8.RuneSelf {
			nameLen++
		}
	}
	return nameLen
}
func getFillLen(fillSize int, str string) int {
	outAsciiLen := 0
	for _, c := range str {
		if c >= utf8.RuneSelf && c != '※' {
			outAsciiLen++
		}
	}
	return fillSize - outAsciiLen
}

func printUpEVOInfo(entity *data.Entity) {
	idx := 0
	for _, evoPath := range data.EVOPathList {
		if evoPath[len(evoPath)-1].Entity.Key == entity.Key {
			idx++
			fmt.Println(fmt.Sprintf("%3d", idx), getEvoPathStr(evoPath))
		}
	}
	printEntityInfo(entity)
}

func printDownEvoInfo(entity *data.Entity) {
	var entityDownEvoPathList [][]*data.EVONode
	pathKeyMap := make(map[string]int)
	for _, evoPath := range data.EVOPathList {
		for i, evoNode := range evoPath {
			if evoNode.Entity == entity {
				subEvoPath := evoPath[i:]
				if len(subEvoPath) > 1 {
					var pathKeyList []string
					for _, subEvoNode := range subEvoPath {
						pathKeyList = append(pathKeyList, fmt.Sprintf("%v%v", subEvoNode.Entity.Key, subEvoNode.Ord))
					}
					pathKey := strings.Join(pathKeyList, "")
					_, ok := pathKeyMap[pathKey]
					if !ok {
						entityDownEvoPathList = append(entityDownEvoPathList, subEvoPath)
						pathKeyMap[pathKey] = 1
					}
				}
				break
			}
		}
	}
	idx := 0
	for _, evoPath := range entityDownEvoPathList {
		idx++
		fmt.Println(fmt.Sprintf("%3d", idx), getEvoPathStr(evoPath))
	}
	printEntityInfo(entity)
}

func printOnlyOnePath() {
	idx := 0
	for _, evoPath := range data.EVOPathList {
		evoPathLen := len(evoPath)
		if evoPathLen > 3 {
			entityKey := evoPath[evoPathLen-1].Entity.Key
			count := 0
			for _, cEvoPath := range data.EVOPathList {
				latestEntityKey := cEvoPath[len(cEvoPath)-1].Entity.Key
				if entityKey == latestEntityKey {
					count++
					if count > 1 {
						break
					}
				}
			}
			if count == 1 {
				idx++
				fmt.Println(fmt.Sprintf("%3d", idx), getEvoPathStr(evoPath))
			}
		}
	}
}

func main() {
	t := flag.String("t", "", "方式:one,up,down,comp")
	name := flag.String("n", "", "名称")
	flag.Parse()
	if *t == "one" {
		printOnlyOnePath()
	} else {
		if *name == "" {
			fmt.Println("无名称输入")
			return
		}
		fmt.Printf("搜索: %v\n", *name)

		entity, err := data.FindEntityByName(*name)
		if err != nil {
			fmt.Printf("查找'%v'出错: %v", *name, err)
			return
		}

		switch *t {
		case "up":
			printUpEVOInfo(entity)
		case "down":
			printDownEvoInfo(entity)
		case "comp":
			printEvoCompareInfo(entity)
		default:
			fmt.Println("No Info")
		}
	}
}
