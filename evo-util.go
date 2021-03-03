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
	pathStrTmpl := "%v(%v)%v%v%v"
	for _, evoNode := range evoPath {
		entity := evoNode.Entity
		lockFlag := ""
		if entity.EvoLock != "" {
			lockFlag = "🔒"
		}
		itemFlag := ""
		if entity.Evo == "" && entity.EvoItem != "" {
			itemFlag = "\U0001F9EA"
		}
		ordStr := ""
		if evoNode.Ord != 0 {
			ordStr = fmt.Sprintf("@%v", evoNode.Ord)
		}
		pathStr := fmt.Sprintf(pathStrTmpl, entity.Name, entity.CName, lockFlag, itemFlag, ordStr)
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
	mainStrTmplTmpl := "%%%vv %%%vv %%-%vs %%-%vs %%-%vs %%-%vs %%-%vs %%-%vs %%-%vs %%-%vs %%-%vs %%-%vs %%-%vs %%-%vs"
	titleList := []string{"※", "#", "", "パラメータ", "体重", "育成ミス", "ご機嫌", "しつけ", "戦闘勝利", "技数", "[デコード レベル]", "必要数", "EvoItem", "EvoLock"}
	var spaceMaxList = [14]int{}
	for idx, title := range titleList {
		spaceMaxList[idx] = getStrSpace(title)
	}
	var infoListList [][]string
	for _, ordInfo := range ordInfoList {
		entity, _ := data.EntityMap[ordInfo.EKey]
		var infoList []string

		idx := 0
		cFlag := ""
		if entity.Key == cKey {
			cFlag = "※"
		}
		spaceMaxList[idx] = 1
		infoList = append(infoList, cFlag)

		idx = 1
		ordStr := fmt.Sprintf("%v", ordInfo.Ord)
		spaceMaxList[idx] = int(math.Max(float64(spaceMaxList[idx]), float64(len(ordStr))))
		infoList = append(infoList, ordStr)

		idx = 2
		name := fmt.Sprintf("%v(%v)", entity.Name, entity.CName)
		spaceMaxList[idx] = int(math.Max(float64(spaceMaxList[idx]), float64(getStrSpace(name))))
		infoList = append(infoList, name)

		idx = 3
		var evoCondList []string
		if entity.Evo != "" {
			evoCondList = strings.Split(entity.Evo, ",")
		}
		evoCondStrListLen := len(evoCondList)

		tmpIdx := idx
		for i, evoCond := range evoCondList {
			idx = tmpIdx + i
			spaceMaxList[idx] = int(math.Max(float64(spaceMaxList[idx]), float64(getStrSpace(evoCond))))
			infoList = append(infoList, evoCond)
		}
		idx += evoCondStrListLen % 9
		tmpIdx = idx
		fillLen := 9 - evoCondStrListLen
		for i := 0; i < fillLen; i++ {
			idx = tmpIdx + i
			spaceMaxList[idx] = int(math.Max(float64(spaceMaxList[idx]), float64(1)))
			infoList = append(infoList, "-")
		}

		idx = 12
		evoItem := entity.EvoItem
		if evoItem == "" {
			evoItem = "-"
		}
		spaceMaxList[idx] = int(math.Max(float64(spaceMaxList[idx]), float64(getStrSpace(evoItem))))
		infoList = append(infoList, evoItem)

		idx = 13
		lockFlag := "-"
		if entity.EvoLock != "" {
			lockFlag = "🔒"
		}
		spaceMaxList[idx] = int(math.Max(float64(spaceMaxList[idx]), float64(getStrSpace(lockFlag))))
		infoList = append(infoList, lockFlag)
		infoListList = append(infoListList, infoList)
	}

	var titleSpaceIList []interface{}
	var titleIList []interface{}
	for i, title := range titleList {
		titleSpaceIList = append(titleSpaceIList, getFillLen(spaceMaxList[i], title))
		titleIList = append(titleIList, title)
	}
	titleTmpl := fmt.Sprintf(mainStrTmplTmpl, titleSpaceIList...)
	fmt.Printf(titleTmpl, titleIList...)
	fmt.Println()

	for _, infoList := range infoListList {
		var infoSpaceIList []interface{}
		var infoIList []interface{}
		for i, info := range infoList {
			infoSpaceIList = append(infoSpaceIList, getFillLen(spaceMaxList[i], info))
			infoIList = append(infoIList, info)
		}
		infoTmpl := fmt.Sprintf(mainStrTmplTmpl, infoSpaceIList...)
		fmt.Printf(infoTmpl, infoIList...)
		fmt.Println()
	}
}

func printEvoCompareInfo(entity *data.Entity) {
	for _, p := range entity.P {
		pEntity, _ := data.EntityMap[p.EKey]
		fmt.Printf("==========%v(%v)==========\n", pEntity.Name, pEntity.CName)
		printEntitiesInfo(entity.Key, pEntity.N)
	}
	if len(entity.P) == 0 {
		printEntitiesInfo(entity.Key, []data.OrdInfo{{entity.Key, 1}})
	}
}

func getStrSpace(str string) int {
	nameLen := utf8.RuneCountInString(str)
	for _, c := range str {
		if c >= utf8.RuneSelf && c != '※' {
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
				lenSubEvoPath := len(subEvoPath)
				if len(subEvoPath[lenSubEvoPath-1].Entity.N) == 0 {
					//if len(subEvoPath) > 1 {
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

func printLockInfo() {
	idx := 0
	maxNameSpace := 0
	var entityList []*data.Entity
	for _, entity := range data.AllList {
		if entity.EvoLock != "" {
			idx++
			name := fmt.Sprintf("%v(%v)", entity.Name, entity.CName)
			maxNameSpace = int(math.Max(float64(maxNameSpace), float64(getStrSpace(name))))
			entityList = append(entityList, entity)
		}
	}
	msgTmplTmpl := "%%2v %%-%vs %%v"
	for i, entity := range entityList {
		name := fmt.Sprintf("%v(%v)", entity.Name, entity.CName)
		fillLen := getFillLen(maxNameSpace, name)
		msgTmpl := fmt.Sprintf(msgTmplTmpl, fillLen)
		fmt.Printf(msgTmpl, i+1, name, entity.EvoLock)
		fmt.Println()
	}
}

func printPhaseList(phase string) {
	var ordInfoList []data.OrdInfo
	idx := 0
	for _, entity := range data.AllList {
		if entity.Phase == phase {
			idx++
			ordInfoList = append(ordInfoList, data.OrdInfo{EKey: entity.Key, Ord: idx})
		}
	}
	printEntitiesInfo("", ordInfoList)
}

func main() {
	_t := flag.String("t", "", "方式:one,lock,phase,up,down,comp")
	_name := flag.String("n", "", "名称")
	flag.Parse()
	t := *_t
	name := *_name
	if t != "" {
		fmt.Println("成長期に2才、成熟期に3才、完全体に6才、究極体に11才程で進化する。")
	}
	if t == "one" {
		printOnlyOnePath()
	} else if t == "lock" {
		printLockInfo()
	} else if t == "phase" {
		phaseList := []string{"幼1", "幼2", "长", "熟", "完", "究"}
		phase := ""
		switch name {
		case "幼1":
			phase = "幼生期1"
		case "幼2":
			phase = "幼生期2"
		case "长":
			phase = "成长期"
		case "熟":
			phase = "成熟期"
		case "完":
			phase = "完全体"
		case "究":
			phase = "究极体"
		}
		if phase == "" {
			fmt.Println(phaseList)
		} else {
			printPhaseList(phase)
		}
	} else {
		if name == "" {
			fmt.Println("无名称输入")
			return
		}
		fmt.Printf("搜索: %v\n", name)

		entity, err := data.FindEntityByName(name)
		if err != nil {
			fmt.Printf("查找'%v'出错: %v", name, err)
			return
		}

		switch t {
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
