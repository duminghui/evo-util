package main

import (
	"evo-util/data"
	"flag"
	"fmt"
	"strings"
)

//go:generate go generate ./data

func getEvoPathStr(evoPath []*data.EVONode) string {
	var pathStrList []string
	for _, evoNode := range evoPath {
		entity := evoNode.Entity
		pathStr := fmt.Sprintf("%v(%v)@%v", entity.Name, entity.CName, evoNode.Ord)
		if entity.EvoLock != "无" {
			pathStr = fmt.Sprintf("※%v", pathStr)
		}
		pathStrList = append(pathStrList, pathStr)
	}
	return strings.Join(pathStrList, "=>")
}

func printEntityInfo(entity *data.Entity) {
	fmt.Println("-------------------")
	fmt.Println(entity.Phase)
	fmt.Printf("EvoLock:%v\n", entity.EvoLock)
	fmt.Println("Evo Conditions:")
	avoConditionsTitle := fmt.Sprintf("%-30s%-8s%-8s%-8s%-8s%-8s%-8s%-10s%-8s", "パラメータ", "体重", "育成ミス", "ご機嫌", "しつけ", "戦闘勝利", "技数", "デコード レベル", "必要数")
	fmt.Println(avoConditionsTitle)
	fmt.Println(entity.Evo)
	fmt.Println("-------------------")
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
	idx := 0
	for _, evoPath := range data.EVOPathList {
		evoPathLen := len(evoPath)
		for i, evoNode := range evoPath {
			if i == evoPathLen-1 {
				break
			}
			if evoNode.Entity.Key == entity.Key {
				idx++
				fmt.Println(fmt.Sprintf("%3d", idx), getEvoPathStr(evoPath))
			}
		}
	}
	printEntityInfo(entity)
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

		if *t == "up" {
			printUpEVOInfo(entity)
		} else {
			printDownEvoInfo(entity)
		}
	}
}
