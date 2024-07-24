package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	aoim := NewAOIManager(0, 250, 5, 0, 250, 5)
	fmt.Println(aoim)
}

func TestAOIManager_GetSurroundGridsByGid(t *testing.T) {
	aoim := NewAOIManager(0, 250, 5, 0, 250, 5)
	for gID, _ := range aoim.grids {
		grids := aoim.GetSurroundGridsByGid(gID)
		fmt.Println("gID:", gID, "have", len(grids), "surroundGrids")
		gIDs := make([]int, 0, len(grids))
		for _, grid := range grids {
			gIDs = append(gIDs, grid.GID)
		}
		fmt.Println("these surroundGridID are:", gIDs)
	}
}
