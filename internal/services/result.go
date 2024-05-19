package services

import (
	"context"
	"goshell/internal/adapters/db/pgsql"
	"goshell/internal/entities"
	"log"
)

//type ResultService struct
type ResultService struct {
}

//func NewResultService() *ResultService
func NewResultService() *ResultService {
	return &ResultService{}
}

// type ifResultStorage interface
type ifResultStorage interface {
	GetList(ctx context.Context) (entities.Result_count, error)
}

// func (esv *ResultService) GetList(ctx context.Context) (entities.Result_count, error)
func (esv *ResultService) GetList(ctx context.Context) (entities.Result_count, error) {
	var est ifResultStorage
	est = pgsql.NewResultStorage()

	out_arr_count, err := est.GetList(ctx)
	if err != nil {
		log.Println(err.Error(), "results_list")
		return entities.Result_count{Values: []entities.Result{}, Count: 0}, err
	}

	return out_arr_count, nil
}
