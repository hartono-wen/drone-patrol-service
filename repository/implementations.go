package repository

import (
	"context"
	"database/sql"
	"log"
)

// CreateEstate creates a new estate in the plantation management service.
// It takes a CreateEstateInput struct as input, which contains the length and width
// of the new estate. It returns a CreateEstateOutput struct, which contains the
// ID of the newly created estate.
// This function uses a transaction to ensure atomicity of the estate creation.
// If the estate already exists with the same length and width, the function
// will update the created_at timestamp of the existing estate.
func (r *Repository) CreateEstate(ctx context.Context, input *CreateEstateInput) (output *CreateEstateOutput, err error) {
	sqlStatement := `
		INSERT INTO plantation_management_service.estates (
			id
			,length
			,width
			,created_at
		)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (length, width)
		DO UPDATE SET
			created_at = now()
		RETURNING id;
   `
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("err starting transaction to create estate: ", err)
		return nil, err
	}
	defer tx.Rollback()
	output = &CreateEstateOutput{}
	err = tx.QueryRow(sqlStatement, input.Id, input.Length, input.Width).Scan(&output.Id)
	if err != nil {
		log.Println("err executiing query to create estate: ", err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.Println("err committing transaction to create estate: ", err)
		return nil, err
	}
	return output, nil
}

// GetEstateByEstateId retrieves the length and width of an estate by its ID.
// It takes a context.Context and a *GetEstateByEstateIdInput as input, and returns
// a *GetEstateByEstateIdOutput and an error.
// The function executes a SQL query to fetch the length and width of the estate
// from the `plantation_management_service.estates` table, where the estate ID
// matches the provided input.
// If an error occurs during the query execution, the function will return the
// error.
func (r *Repository) GetEstateByEstateId(ctx context.Context, input *GetEstateByEstateIdInput) (output *GetEstateByEstateIdOutput, err error) {
	sqlStatement := `
		SELECT
			estates.length
			,estates.width
		FROM
			plantation_management_service.estates
		WHERE estates.id = $1;
   `
	row := r.Db.QueryRowContext(ctx, sqlStatement, input.Id)
	output = &GetEstateByEstateIdOutput{}
	err = row.Scan(&output.Estate.Length, &output.Estate.Width)
	if err == sql.ErrNoRows {
		log.Println("err no estate is found:", err)
		return nil, nil
	} else if err != nil {
		log.Println("err executing query to select the length and the width of the estate:", err)
		return nil, err
	}
	return output, nil
}

// IsTreeExist checks if a tree with the given estate ID, x, and y coordinates exists in the plantation management service.
// The input parameter input contains the estate ID, x, and y coordinates to check for.
// The output parameter output contains a boolean indicating whether the tree exists or not.
func (r *Repository) IsTreeExist(ctx context.Context, input *IsTreeExistInput) (output *IsTreeExistOutput, err error) {
	sqlStatement := `
		SELECT EXISTS(
			SELECT 1
			FROM
				plantation_management_service.trees
			WHERE trees.estate_id = $1 AND trees.x = $2 AND trees.y = $3
		);
		
   `
	row := r.Db.QueryRowContext(ctx, sqlStatement, input.EstateId, input.X, input.Y)
	output = &IsTreeExistOutput{}
	err = row.Scan(&output.IsExist)
	if err != nil {
		log.Println("err executing query to check whether a certain tree exist or not:", err)
		return nil, err
	}
	return output, nil
}

// CreateTree creates a new tree in the plantation management service.
// The input parameter input contains the details of the new tree to be created, including its ID, estate ID, x and y coordinates, and height.
// The output parameter output contains the ID of the newly created tree.
func (r *Repository) CreateTree(ctx context.Context, input *CreateTreeInput) (output *CreateTreeOutput, err error) {
	sqlStatement := `
		INSERT INTO plantation_management_service.trees (
			id
			,estate_id
			,x
			,y
			,height
			,created_at
		)
		VALUES ($1, $2, $3, $4, $5, now())
		RETURNING id;
   `
	tx, err := r.Db.BeginTx(ctx, nil)
	if err != nil {
		log.Println("err starting transaction to create tree: ", err)
		return nil, err
	}
	defer tx.Rollback()
	output = &CreateTreeOutput{}
	err = tx.QueryRow(sqlStatement, input.Id, input.EstateId, input.X, input.Y, input.Height).Scan(&output.Id)
	if err != nil {
		log.Println("err executing query to create tree: ", err)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		log.Println("err committing transaction to create tree: ", err)
		return nil, err
	}
	return output, nil
}

// GetEstateStatsByEstateId retrieves various statistics about the trees in an estate, including the total number of trees, the maximum and minimum tree heights, and the median tree height.
// The input parameter EstateId specifies the ID of the estate to retrieve the statistics for.
// The output is a GetEstateStatsByEstateIdOutput struct containing the requested statistics.
func (r *Repository) GetEstateStatsByEstateId(ctx context.Context, input *GetEstateStatsByEstateIdInput) (output *GetEstateStatsByEstateIdOutput, err error) {
	sqlStatement := `
		SELECT
			COUNT(trees.height) AS total_trees
			,COALESCE(MAX(trees.height), 0) AS max_height
			,COALESCE(MIN(trees.height), 0) AS min_height
			,COALESCE(PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY trees.height), 0) AS median_height
		FROM
			plantation_management_service.trees
		WHERE trees.estate_id = $1;
   `
	row := r.Db.QueryRowContext(ctx, sqlStatement, input.EstateId)
	output = &GetEstateStatsByEstateIdOutput{}
	err = row.Scan(&output.Count, &output.Max, &output.Min, &output.Median)
	if err != nil {
		log.Println("err executing query to get estate stats by estate id:", err)
		return nil, err
	}
	return output, nil
}

// GetEstateTreesByEstateId retrieves the trees for a given estate, including their x, y coordinates and height.
// The input parameter EstateId specifies the ID of the estate to retrieve the trees for.
// The output is a GetEstateTreesByEstateIdOutput struct containing the requested tree data, as well as the length and width of the estate.
func (r *Repository) GetEstateTreesByEstateId(ctx context.Context, input *GetEstateTreesByEstateIdInput) (output *GetEstateTreesByEstateIdOutput, err error) {
	sqlStatement := `
		SELECT
			trees.x
			,trees.y
			,trees.height
		FROM
			plantation_management_service.trees
		WHERE trees.estate_id = $1;
   `
	rows, err := r.Db.QueryContext(ctx, sqlStatement, input.EstateId)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		log.Println("err executing query to get the trees belonging to a certain estate id:", err)
		return nil, err
	}
	defer rows.Close()

	output = &GetEstateTreesByEstateIdOutput{}
	for rows.Next() {
		var tree Tree

		err := rows.Scan(&tree.X, &tree.Y, &tree.Height)
		if err != nil {
			log.Println("err when reading the rows as result from the query:", err)
			return nil, err
		}

		output.Trees = append(output.Trees, tree)
	}

	sqlStatement = `
		SELECT
			estates.length
			,estates.width
		FROM plantation_management_service.estates
		WHERE estates.id = $1;
   `

	row := r.Db.QueryRowContext(ctx, sqlStatement, input.EstateId)
	var estate Estate
	err = row.Scan(&estate.Length, &estate.Width)
	if err != nil {
		log.Println("err executing query to get the estate length and the estate width:", err)
		return nil, err
	}
	output.Estate = estate

	return output, err
}
