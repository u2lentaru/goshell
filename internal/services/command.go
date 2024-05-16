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

// func CommExec(ctx context.Context, id int) - выполняет скрипт и сохраняет результат
func CommExec(ctx context.Context, id int) error {
	var est ifCommandStorage
	est = pgsql.NewCommandStorage()

	err := est.CommExec(ctx, id)

	if err != nil {
		log.Println("Failed execute command add!")
		return err
	}

	return nil
}

// func CommSave(ctx context.Context, bs []byte) (int, error) - сохраняет скрипт в базу
func CommSave(ctx context.Context, bs []byte) (int, error) {
	var est ifCommandStorage
	est = pgsql.NewCommandStorage()

	cid, err := est.CommSave(ctx, bs)
	if err != nil {
		log.Println(err.Error(), "commands_list")
		return 0, err
	}

	return cid, nil
}

// func CommGetList(ctx context.Context) (entities.Command_count, error) - возвращает список команд
func CommGetList(ctx context.Context) (entities.Command_count, error) {
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
