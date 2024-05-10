// This file contains types that are used in the repository layer.
package repository

type CreateEstateInput struct {
	Id            string
	Length, Width uint16
}

type CreateEstateOutput struct {
	Id string
}

type GetEstateByEstateIdInput struct {
	Id string
}

type GetEstateByEstateIdOutput struct {
	Estate Estate
}

type IsTreeExistInput struct {
	EstateId string
	X, Y     int
}

type IsTreeExistOutput struct {
	IsExist bool
}

type CreateTreeInput struct {
	Id, EstateId string
	X, Y, Height int
}

type CreateTreeOutput struct {
	Id string
}

type GetEstateTreesByEstateIdInput struct {
	EstateId string
}

type GetEstateTreesByEstateIdOutput struct {
	Trees  []Tree
	Estate Estate
}

type Tree struct {
	X, Y, Height int
}

type Estate struct {
	Length, Width int
}

type CalculateDroneDistanceInput struct {
	Trees  []Tree
	Estate Estate
}

type CalculateDroneDistanceOutput struct {
	TotalDistance             int
	TotalVerticalDistance     int
	TotalHorizontalDistance   int
	LastAchievableXCoordinate int
	LastAchievableYCoordinate int
}

type GetEstateStatsByEstateIdInput struct {
	EstateId string
}

type GetEstateStatsByEstateIdOutput struct {
	Count, Max, Min int
	Median          float32
}
