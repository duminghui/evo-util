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

func printEVOInfo(entity *data.Entity) {
	var idx = 0
	for _, evoPath := range data.EVOPathList {
		if evoPath[len(evoPath)-1].Key == entity.Key {
			idx++
			var pathStrList []string
			for _, node := range evoPath {
				pathStr := fmt.Sprintf("%v(%v)", node.Name, node.CName)
				if node.EvoLock != "无" {
					pathStr = fmt.Sprintf("※%v", pathStr)
				}
				pathStrList = append(pathStrList, pathStr)
			}
			fmt.Println(idx, strings.Join(pathStrList, "=>"))
		}
	}
	fmt.Println("-------------------")
	fmt.Println(entity.Phase)
	fmt.Printf("EvoLock:%v\n", entity.EvoLock)
	fmt.Println("Evo Conditions:")
	fmt.Println(entity.Evo)
	fmt.Println("-------------------")
}

func main() {
	name := flag.String("name", "", "名称")
	flag.Parse()
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
