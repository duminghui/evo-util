package main

import (
	"evo-util/data"
	"flag"
	"fmt"
	"strings"
)

//go:generate go generate ./data

//type EVOPath struct {
//	entity *data.Entity
//	p      *EVOPath
//	n      *EVOPath
//}
//
//func getParentEVOPath(entity *data.Entity, nEVOPath *EVOPath) []*EVOPath {
//	var evoPathList []*EVOPath
//	if len(entity.P) == 0 {
//		evoPath := &EVOPath{
//			entity: entity,
//			n:      nEVOPath,
//		}
//		evoPathList = append(evoPathList, evoPath)
//	} else {
//		for _, pKey := range entity.P {
//			fmt.Println("=====", pKey, "=====")
//			pEntity, _ := data.EntityMap[pKey]
//			evoPath := &EVOPath{
//				entity: entity,
//				n:      nEVOPath,
//			}
//			parentEvoPathList := getParentEVOPath(pEntity, evoPath)
//			fmt.Println("=====", pKey, len(parentEvoPathList))
//			if len(parentEvoPathList) == 0 {
//				evoPathList = append(evoPathList, evoPath)
//			} else {
//				for _, parentEvoPath := range parentEvoPathList {
//					evoPath := &EVOPath{
//						entity: entity,
//						p:      parentEvoPath,
//						n:      nEVOPath,
//					}
//					evoPathList = append(evoPathList, evoPath)
//				}
//			}
//		}
//	}
//	return evoPathList
//}
//
//func printPath(evoPath *EVOPath) {
//	if evoPath.p != nil {
//		printPath(evoPath.p)
//		fmt.Print(evoPath.entity.CName)
//		if evoPath.n != nil {
//			fmt.Print("=>")
//		} else {
//			fmt.Println()
//		}
//	} else {
//		fmt.Print(evoPath.entity.CName)
//		fmt.Print("=>")
//	}
//}

func getEvoPathStr(evoPath []*data.Entity) string {
	var pathStrList []string
	for _, entity := range evoPath {
		pathStr := fmt.Sprintf("%v(%v)", entity.Name, entity.CName)
		if entity.EvoLock != "无" {
			pathStr = fmt.Sprintf("※%v", pathStr)
		}
		pathStrList = append(pathStrList, pathStr)
	}
	return strings.Join(pathStrList, "=>")
}

func printEVOInfo(entity *data.Entity) {
	idx := 0
	for _, evoPath := range data.EVOPathList {
		if evoPath[len(evoPath)-1].Key == entity.Key {
			idx++
			fmt.Println(idx, getEvoPathStr(evoPath))
		}
	}
	fmt.Println("-------------------")
	fmt.Println(entity.Phase)
	fmt.Printf("EvoLock:%v\n", entity.EvoLock)
	fmt.Println("Evo Conditions:")
	avoConditionsTitle := fmt.Sprintf("%-30s%-8s%-8s%-8s%-8s%-8s%-8s%-10s%-8s", "パラメータ", "体重", "育成ミス", "ご機嫌", "しつけ", "戦闘勝利", "技数", "デコード レベル", "必要数")
	fmt.Println(avoConditionsTitle)
	fmt.Println(entity.Evo)
	fmt.Println("-------------------")
}

func printOnlyOnePath() {
	idx := 0
	for _, evoPath := range data.EVOPathList {
		evoPathLen := len(evoPath)
		if evoPathLen > 3 {
			entity := evoPath[evoPathLen-1]
			count := 0
			for _, cEvoPath := range data.EVOPathList {
				latestEntity := cEvoPath[len(cEvoPath)-1]
				if entity == latestEntity {
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
	t := flag.String("t", "", "方式:all,one")
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
		printEVOInfo(entity)
	}
}
