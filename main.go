package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/4frag/padalko/tasks"
	"github.com/4frag/padalko/utils"
	"github.com/yarlson/tap"
)

func inputMatrix(ctx context.Context) error {
	size_validator := func (input int) error {
		if input == 0 {
			return errors.New("Invalid value: float must be greater than 0")
		}
		return nil
	}

	tasks.CurrentData.SizeA = utils.InputNumber(ctx, "Input count of resources", "", 0, size_validator)
	tasks.CurrentData.SizeB = utils.InputNumber(ctx, "Input count of objects", "", 0, size_validator)

	tasks.CurrentData.Model = make([][]float64, tasks.CurrentData.SizeA)
	for i := 0; i < tasks.CurrentData.SizeA; i++ {
		row, err := utils.InputMatrixRow[float64](
			ctx,
			fmt.Sprintf("Row %d:", i+1),
			tasks.CurrentData.SizeB,
			nil,
		)
		if err != nil {
			return err
		}
		tasks.CurrentData.Model[i] = row
	}
	return nil
}

func inputResources(ctx context.Context) error {
	result, err := utils.InputMatrixRow[float64](ctx, "Input resources vector", tasks.CurrentData.SizeA, nil)
	if err != nil {return err}
	tasks.CurrentData.A = result
	return nil
}

func inputPlan(ctx context.Context) error {
	result, err := utils.InputMatrixRow[float64](ctx, "Input plan vector", tasks.CurrentData.SizeB, nil)
	if err != nil {return err}
	tasks.CurrentData.B = result
	return nil
}

func PrintModelData(ctx context.Context, data tasks.ModelData) {
	tap.Message(fmt.Sprintf("ModelData: SizeA = %d, SizeB = %d", data.SizeA, data.SizeB))

	// Вывод матрицы Model
	tap.Message("Matrix (Model):")
	if len(data.Model) == 0 || len(data.Model[0]) == 0 {
		tap.Message("  <empty>")
	} else {
		for i, row := range data.Model {
			rowStr := ""
			for _, val := range row {
				rowStr += fmt.Sprintf("%10.2f ", val)
			}
			tap.Message(fmt.Sprintf("Row %d: %s", i, rowStr))
		}
	}

	// Вывод вектора ресурсов A
	tap.Message("Vector A (resources):")
	if len(data.A) == 0 {
		tap.Message("  <empty>")
	} else {
		rowStr := ""
		for _, val := range data.A {
			rowStr += fmt.Sprintf("%10.2f ", val)
		}
		tap.Message("  " + rowStr)
	}

	// Вывод вектора плана B
	tap.Message("Vector B (plan):")
	if len(data.B) == 0 {
		tap.Message("  <empty>")
	} else {
		rowStr := ""
		for _, val := range data.B {
			rowStr += fmt.Sprintf("%10.2f ", val)
		}
		tap.Message("  " + rowStr)
	}
}

const (
	ActionQuit = iota
	ActionInputMatrix
	ActionInputResources
	ActionInputPlan
	ActionCalculatePlan
	ActionCalculateWithCriteria
	ActionInfo
)

func main() {
	ctx := context.Background()
	tap.Intro(fmt.Sprintf("%sWelcome! 👋%s", tap.Green, tap.Green))

	for {
		selected := utils.CreateMenu(ctx, "Select action", []utils.MenuItem[int]{
			{
				ID: ActionInputMatrix,
				Name: "Input matrix",
				Description:  "Задать зависимость между переменными",
			},
			{
				ID: ActionInputResources,
				Name: "Input resources vector",
			},
			{
				ID: ActionInputPlan,
				Name: "Input plan vector",
			},
			{
				ID: ActionInfo,
				Name: "Print info",
			},
			{
				ID: ActionCalculatePlan,
				Name: "Calculate plan",
			},
			{
				ID: ActionCalculateWithCriteria,
				Name: "Calculate with criteria",
			},
			{
				ID: ActionQuit,
				Name: "Quit",
				Description:  "Exit program",
			},
		})

		switch selected {
		case ActionInputMatrix:
			inputMatrix(ctx)
		case ActionInputResources:
			inputResources(ctx)
		case ActionInputPlan:
			inputPlan(ctx)
		case ActionCalculatePlan:
			tasks.CurrentData.CalculatePlan()
		case ActionCalculateWithCriteria:
			criteria, err := utils.InputMatrixRow[float64](ctx, "Input criteria", tasks.CurrentData.SizeB, nil)
			if err != nil {panic(err)}
			result, err := tasks.CurrentData.SolveWithCriteria(criteria)
			result_str := make([]string, len(result))
			for i, v := range result {
				result_str[i] = strconv.FormatFloat(v, 'f', -1, 64)
			}
			tap.Message(strings.Join(result_str, " "))
		case ActionQuit:
			tap.Message("Bye 👋")
			return
		case ActionInfo:
			PrintModelData(ctx, tasks.CurrentData)
		}
	}
}
