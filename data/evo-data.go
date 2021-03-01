package data

import (
	"errors"
)

//go:generate go run ./gen-data-2.go

type Entity struct {
	Key     string
	Name    string
	CName   string
	Phase   string
	EvoLock string
	Evo     string
	P       []OrdInfo
	N       []OrdInfo
}

type OrdInfo struct {
	EKey string
	Ord  int
}

type EVONode struct {
	Entity *Entity
	Ord    int
}

var NotExistEntityError = errors.New("entity not find")

func FindEntityByName(name string) (*Entity, error) {
	for _, entity := range AllList {
		if entity.Name == name || entity.CName == name {
			return entity, nil
		}
	}
	return &Entity{}, NotExistEntityError
}
