package services

import (
	"context"
	"goshell/internal/adapters/db/pgsql"
	"goshell/internal/entities"
	"log"
)

//type CommandService struct
type CommandService struct {
}

//func NewCommandService() *CommandService
func NewCommandService() *CommandService {
	return &CommandService{}
}

type ifCommandStorage interface {
	CommExec(ctx context.Context, id int) error
	CommSave(ctx context.Context, bs []byte) (int, error)
	GetList(ctx context.Context) (entities.Command_count, error)
	GetOne(ctx context.Context, i int) (entities.Command_count, error)
}

func (esv *CommandService) PostExec(ctx context.Context, bs []byte) error {
	var est ifCommandStorage
	est = pgsql.NewCommandStorage()

	id, err := est.CommSave(ctx, bs)
	if err != nil {
		log.Println(err.Error(), "CommSave error")
		return err
	}

	err = est.CommExec(ctx, id)
	if err != nil {
		log.Println(err.Error(), "CommExec error")
		return err
	}

	return nil
}

// func CommGetList(ctx context.Context) (entities.Command_count, error) - возвращает список команд
func (esv *CommandService) ExecOne(ctx context.Context, id int) error {
	var est ifCommandStorage
	est = pgsql.NewCommandStorage()

	err := est.CommExec(ctx, id)
	if err != nil {
		log.Println(err.Error(), "commands_list")
		return err
	}

	return nil
}

// func CommGetList(ctx context.Context) (entities.Command_count, error) - возвращает список команд
func (esv *CommandService) Exec(ctx context.Context, ids []int) error {
	var est ifCommandStorage
	est = pgsql.NewCommandStorage()

	for _, id := range ids {
		go est.CommExec(ctx, id)
	}

	return nil
}

// func CommGetList(ctx context.Context) (entities.Command_count, error) - возвращает список команд
func (esv *CommandService) GetList(ctx context.Context) (entities.Command_count, error) {
	var est ifCommandStorage
	est = pgsql.NewCommandStorage()

	out_arr_count, err := est.GetList(ctx)
	if err != nil {
		log.Println(err.Error(), "commands_list")
		return entities.Command_count{Values: []entities.Command{}, Count: 0}, err
	}

	return out_arr_count, nil
}

// func CommGetOne(ctx context.Context, id int) (entities.Command_count, error) - вывод команды по id
func (esv *CommandService) GetOne(ctx context.Context, id int) (entities.Command_count, error) {
	var est ifCommandStorage
	est = pgsql.NewCommandStorage()

	out_arr_count, err := est.GetOne(ctx, id)
	if err != nil {
		log.Println(err.Error(), "commands_one")
		return entities.Command_count{Values: []entities.Command{}, Count: 0}, err
	}

	return out_arr_count, nil
}
