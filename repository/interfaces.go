// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	CreateEstate(ctx context.Context, input *CreateEstateInput) (output *CreateEstateOutput, err error)
	GetEstateByEstateId(ctx context.Context, input *GetEstateByEstateIdInput) (output *GetEstateByEstateIdOutput, err error)
	IsTreeExist(ctx context.Context, input *IsTreeExistInput) (output *IsTreeExistOutput, err error)
	CreateTree(ctx context.Context, input *CreateTreeInput) (output *CreateTreeOutput, err error)
	GetEstateStatsByEstateId(ctx context.Context, input *GetEstateStatsByEstateIdInput) (output *GetEstateStatsByEstateIdOutput, err error)
	GetEstateTreesByEstateId(ctx context.Context, input *GetEstateTreesByEstateIdInput) (output *GetEstateTreesByEstateIdOutput, err error)
}
