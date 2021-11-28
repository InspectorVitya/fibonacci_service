package app

import (
	"context"
	"errors"
	"math/big"

	"github.com/inspectorvitya/fibonacci_service/internal/storage"
)

type App struct {
	Storage storage.Storage
	lastFib int
}

var (
	ErrXMoreY            = errors.New("x more y")
	ErrLessOrGreaterZero = errors.New("x or y less zero")
)

func New(db storage.Storage) *App {
	app := &App{
		Storage: db,
	}
	return app
}

func (app *App) Init(ctx context.Context) error {
	app.lastFib = 200
	err := app.Storage.SetSliceFib(ctx, FibonacciNumbers(app.lastFib))
	if err != nil {
		return err
	}
	return nil
}

func (app *App) GetFibSlice(ctx context.Context, x, y int) ([]string, error) {
	if x <= 0 || y <= 0 {
		return nil, ErrLessOrGreaterZero
	}
	if x > y {
		return nil, ErrXMoreY
	}

	if app.lastFib > y {
		cachedFib, err := app.Storage.GetSliceFib(ctx, x, y)
		if err != nil {
			return nil, err
		}
		return cachedFib, nil
	}
	fibSlice := FibonacciNumbers(y)
	err := app.Storage.SetSliceFib(ctx, fibSlice)
	if err != nil {
		return nil, err
	}
	cachedFib, err := app.Storage.GetSliceFib(ctx, x, y)
	if err != nil {
		return nil, err
	}
	return cachedFib, nil
}

func FibonacciNumbers(n int) []string {
	test := make([]string, n)
	F1 := big.NewInt(0)
	F2 := big.NewInt(1)
	test[0] = F1.String()
	for i := 1; i < n; i++ {
		F1.Add(F1, F2)
		F1, F2 = F2, F1
		test[i] = F1.String()
	}
	return test
}
