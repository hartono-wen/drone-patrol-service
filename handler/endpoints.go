package handler

import (
	"encoding/json"
	"errors"
	"math"
	"net/http"

	"github.com/google/uuid"
	"github.com/hartono-wen/sawitpro-technical-interview-software-architect/generated"
	"github.com/hartono-wen/sawitpro-technical-interview-software-architect/repository"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// PostEstate is an HTTP handler that creates a new estate.
//
// It expects a JSON request body with the following fields:
//   - Length: the length of the estate
//   - Width: the width of the estate
//
// If the request is valid, it creates a new estate in the repository and returns
// a JSON response with the ID of the new estate.
// If the request is invalid or there is an error creating the estate, it returns
// an appropriate HTTP error response.
func (s *Server) PostEstate(ctx echo.Context) error {
	var req generated.PostEstateJSONRequestBody
	err := json.NewDecoder(ctx.Request().Body).Decode(&req)
	if err != nil {
		log.Print("err decoding request: ", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := ctx.Validate(req); err != nil {
		log.Print("err validating request: ", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	createEstateInput := &repository.CreateEstateInput{
		Id:     uuid.New().String(),
		Length: uint16(req.Length),
		Width:  uint16(req.Width),
	}

	output, err := s.Repository.CreateEstate(ctx.Request().Context(), createEstateInput)
	if err != nil {
		log.Print("err when creating estate: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Something happens in our end. Let us check."})
	}
	var resp generated.EstateResponse
	resp.Id, err = uuid.Parse(output.Id)
	if err != nil {
		log.Print("err when parsing estate UUID: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Something happens in our end. Let us check."})
	}
	return ctx.JSON(http.StatusOK, resp)
}

// PostEstateEstateIdTree creates a new tree for the specified estate.
// It validates the request body, checks if the estate exists, ensures the requested
// coordinates are within the estate's boundaries, and checks if a tree already exists
// at the specified coordinates. If all checks pass, it creates a new tree and returns
// the tree's ID in the response.
func (s *Server) PostEstateEstateIdTree(ctx echo.Context, estateId openapi_types.UUID) error {
	var req generated.PostEstateEstateIdTreeJSONRequestBody
	err := json.NewDecoder(ctx.Request().Body).Decode(&req)
	if err != nil {
		log.Print("err decoding request: ", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := ctx.Validate(req); err != nil {
		log.Print("err validating request: ", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	getEstateByEstateId := &repository.GetEstateByEstateIdInput{
		Id: estateId.String(),
	}
	estate, err := s.Repository.GetEstateByEstateId(ctx.Request().Context(), getEstateByEstateId)

	if err != nil {
		log.Error("err getting estate by estate id: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Something happens in our end. Let us check."})
	}

	if estate == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Estate not found"})
	}

	if (req.X > int(estate.Estate.Length) || req.X < 0) || (req.Y > int(estate.Estate.Width) || req.Y < 0) {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	isTreeExistInput := &repository.IsTreeExistInput{
		EstateId: estateId.String(),
		X:        req.X,
		Y:        req.Y,
	}
	isTreeExistOutput, err := s.Repository.IsTreeExist(ctx.Request().Context(), isTreeExistInput)
	if err != nil {
		log.Error("err checking whether tree is exist or not: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Something happens in our end. Let us check."})
	}

	if isTreeExistOutput.IsExist {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	createTreeInput := &repository.CreateTreeInput{
		Id:       uuid.New().String(),
		EstateId: estateId.String(),
		X:        req.X,
		Y:        req.Y,
		Height:   req.Height,
	}

	output, err := s.Repository.CreateTree(ctx.Request().Context(), createTreeInput)
	if err != nil {
		log.Error("err creating tree: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Something happens in our end. Let us check."})
	}
	var resp generated.TreeResponse
	resp.Id, err = uuid.Parse(output.Id)
	if err != nil {
		log.Print("err when parsing tree UUID: ", err)
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	return ctx.JSON(http.StatusOK, resp)
}

// GetEstateEstateIdStats retrieves the statistics for an estate based on the provided estate ID.
// It returns the count, maximum, minimum, and median values for the estate.
func (s *Server) GetEstateEstateIdStats(ctx echo.Context, estateId openapi_types.UUID) error {
	getEstateByEstateId := &repository.GetEstateByEstateIdInput{
		Id: estateId.String(),
	}
	estate, err := s.Repository.GetEstateByEstateId(ctx.Request().Context(), getEstateByEstateId)

	if err != nil {
		log.Error("err getting estate by estate id: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Something happens in our end. Let us check."})
	}

	if estate == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Estate not found"})
	}

	getEstateStatsByEstateIdInput := &repository.GetEstateStatsByEstateIdInput{
		EstateId: estateId.String(),
	}
	output, err := s.Repository.GetEstateStatsByEstateId(ctx.Request().Context(), getEstateStatsByEstateIdInput)
	if err != nil {
		log.Error("err getting estate stats by estate id: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Something happens in our end. Let us check."})
	}

	var resp generated.EstateStatsResponse
	resp.Count = output.Count
	resp.Max = output.Max
	resp.Min = output.Min
	resp.Median = output.Median

	return ctx.JSON(http.StatusOK, resp)
}

// GetEstateEstateIdDronePlan retrieves the estate and trees for the given estate ID,
// calculates the total distance the drone needs to travel to cover the entire estate,
// and returns the drone plan response with the total distance.
func (s *Server) GetEstateEstateIdDronePlan(ctx echo.Context, estateId openapi_types.UUID, params generated.GetEstateEstateIdDronePlanParams) error {
	getEstateEstateIdDronePlanInput := &repository.GetEstateTreesByEstateIdInput{
		EstateId: estateId.String(),
	}
	output, err := s.Repository.GetEstateTreesByEstateId(ctx.Request().Context(), getEstateEstateIdDronePlanInput)
	if err != nil {
		log.Print("err getting estate trees by estate id: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Invalid request"})
	}

	if output == nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"error": "Estate not found"})
	}

	calculateDroneDistanceInput := &repository.CalculateDroneDistanceInput{
		Estate: output.Estate,
		Trees:  output.Trees,
	}

	calculateDroneDistanceOutput, err := s.CalculateDroneDistance(calculateDroneDistanceInput, params.MaxDistance)
	if err != nil {
		log.Print("err calculating drone distance: ", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": "Something happens in our end. Let us check."})
	}

	var resp generated.DronePlanResponse
	resp.Distance = calculateDroneDistanceOutput.TotalDistance
	if params.MaxDistance != nil {
		resp.Distance = *params.MaxDistance
		resp.Rest = &struct {
			X *int `json:"x,omitempty"`
			Y *int `json:"y,omitempty"`
		}{
			X: &calculateDroneDistanceOutput.LastAchievableXCoordinate,
			Y: &calculateDroneDistanceOutput.LastAchievableYCoordinate,
		}

	}

	return ctx.JSON(http.StatusOK, resp)
}

// CalculateDroneDistance calculates the total distance the drone needs to travel to cover the entire estate, taking into account the estate dimensions and the heights of the trees.
// The function takes an input struct containing the estate details and the trees, and an optional maximum distance parameter.
// It returns a struct containing the total distance, the total horizontal distance, the total vertical distance, and the last achievable coordinates for the drone.
// If the maximum distance is provided and the calculated total distance exceeds it, the function will return the last achievable coordinates instead of the full distance.
func (s *Server) CalculateDroneDistance(input *repository.CalculateDroneDistanceInput, maxDistance *int) (*repository.CalculateDroneDistanceOutput, error) {
	// Validate that input must not be nil. If nil, return error.
	if input == nil {
		return nil, errors.New("err CalculateDroneDistance: invalid input -- nothing to calculate drone distance")
	}

	if maxDistance == nil {
		log.Info("maxDistance is nil, calculating the total drone distance.")
	} else {
		log.Info("maxDistance is NOT nil, calculating the max distance that the drone can travel.")
	}

	calculateDroneDistanceOutput := &repository.CalculateDroneDistanceOutput{}

	//totalHorizontalDistance := (input.Estate.Length*input.Estate.Width - 1) * s.Config.ScaleFactor
	totalHorizontalDistance := 0

	// Create estate and populate estate with 1 because 1 is the minimum height for the drone flying.
	plantationGridArray := make([][]int, input.Estate.Width)
	for i := range plantationGridArray {
		plantationGridArray[i] = make([]int, input.Estate.Length)
		for j := range plantationGridArray[i] {
			plantationGridArray[i][j] = 1 // Populate with 1
		}
	}

	// //Debugging purpose
	// for _, row := range plantationGridArray {
	// 	log.Print(row)
	// }

	// Populate the estate with the trees. Set also the height for the drone to patrol the tree.
	for _, t := range input.Trees {
		plantationGridArray[t.Y-1][t.X-1] = t.Height + 1
	}

	totalVerticalDistance := 0
	var currentHeight, previousHeight int
	var i, j int
	// Iterate the Y axis of the estate (hence using input.Estate.Width - not input.Estate.Length)
	for i = 0; i < input.Estate.Width; i++ {

		// Determine if need to go east to west or west to east.
		// The logic is to determine if the current row is even or odd, if even then go east to west, if odd then go west to east.
		if i%2 == 0 {

			// Now iterate the X axis of the estate (hence using input.Estate.Length).
			// The direction of the iteration is east to west (because the row is even).
			for j = 0; j < input.Estate.Length; j++ {
				//log.Printf("i: %d, j: %d\n", i, j)
				if j == 0 {
					if i == 0 {
						// Since this is the very first grid, no previous height which makes sense.
						currentHeight = plantationGridArray[i][j]
					} else {
						currentHeight = plantationGridArray[i][j]
						previousHeight = plantationGridArray[i-1][j]
					}
				} else {
					currentHeight = plantationGridArray[i][j]
					previousHeight = plantationGridArray[i][j-1]
				}

				// Calculate the difference of the height / vertical distance that the drone needs to travel.
				increment := int(math.Abs(float64(currentHeight - previousHeight)))

				//log.Printf("current coordinate: (%d, %d), current height: %d, previous coordinate: (%d, %d), previous height: %d, current row: %d\n", j+1, i+1, currentHeight, j, i+1, previousHeight, i+1)
				//log.Printf("currentHeight: %d, previousHeight: %d, increment: %d\n", currentHeight, previousHeight, increment)

				// Add the difference of the height to the total vertical distance.
				totalVerticalDistance += increment

				if !(i == 0 && j == 0) {
					totalHorizontalDistance += s.Config.ScaleFactor
				}

				if maxDistance != nil && *maxDistance < (totalHorizontalDistance+totalVerticalDistance+currentHeight) {
					return calculateDroneDistanceOutput, nil
				}

				calculateDroneDistanceOutput.LastAchievableXCoordinate = j + 1
				calculateDroneDistanceOutput.LastAchievableYCoordinate = i + 1

				if i == input.Estate.Width-1 && j == input.Estate.Length-1 {
					totalVerticalDistance += plantationGridArray[i][j]
				}

				//log.Printf("totalVerticalDistance: %d", totalVerticalDistance)
			}
		} else {
			// Since this is the odd row, the direction of the iteration is west to east.
			// Hence the iteration starts from input.Estate.Length - 1 and not 0.
			for j = input.Estate.Length - 1; j >= 0; j-- {
				//log.Printf("i: %d, j: %d\n", i, j)

				// Below condition determines if the current estate grid is the first one of the iteration.
				// If it is, then we need to determine the previous height from *below* row instead.
				// Previous row is used instead of previous column because we need to iterate from *south* to *north*
				// since the iteration has reached the end of the grid in that X (horizontal) axis
				if j == input.Estate.Length-1 {
					currentHeight = plantationGridArray[i][j]
					previousHeight = plantationGridArray[i-1][j] // use previous row instead of column

				} else {
					currentHeight = plantationGridArray[i][j]
					previousHeight = plantationGridArray[i][j+1] // use next column instead of row
				}

				// Calculate the difference of the height / vertical distance that the drone needs to travel.
				increment := int(math.Abs(float64(currentHeight - previousHeight)))
				//log.Printf("current coordinate: (%d, %d), current height: %d, previous coordinate: (%d, %d), previous height: %d, current row: %d\n", j+1, i+1, currentHeight, j, i+1, previousHeight, i+1)
				//log.Printf("currentHeight: %d, previousHeight: %d, increment: %d\n", currentHeight, previousHeight, increment)

				totalVerticalDistance += increment
				totalHorizontalDistance += s.Config.ScaleFactor

				if maxDistance != nil && *maxDistance < (totalHorizontalDistance+totalVerticalDistance+currentHeight) {
					return calculateDroneDistanceOutput, nil

				}
				calculateDroneDistanceOutput.LastAchievableXCoordinate = j + 1
				calculateDroneDistanceOutput.LastAchievableYCoordinate = i + 1

				// If reaching the last grid, don't forget to add the vertical distance of the last grid so that the drone can land.
				if i == input.Estate.Width-1 && j == 0 {
					totalVerticalDistance += plantationGridArray[i][j]
				}
			}
			//log.Printf("totalVerticalDistance: %d", totalVerticalDistance)
		}
	}
	calculateDroneDistanceOutput.TotalDistance = totalVerticalDistance + totalHorizontalDistance
	calculateDroneDistanceOutput.TotalHorizontalDistance = totalHorizontalDistance
	calculateDroneDistanceOutput.TotalVerticalDistance = totalVerticalDistance

	return calculateDroneDistanceOutput, nil
}
