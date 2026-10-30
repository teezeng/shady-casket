package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

const dataFieldsCount = 2

var (
	ErrInvalidFormat   = errors.New("Invalid data format")
	ErrInvalidSteps    = errors.New("Invalid steps count")
	ErrInvalidDuration = errors.New("Invalid duration")
)

type DaySteps struct {
	Steps    int
	Duration time.Duration
	personaldata.Personal
}

func (ds *DaySteps) Parse(datastring string) (err error) {
	parts := strings.Split(datastring, ",")
	if len(parts) != dataFieldsCount {
		return fmt.Errorf("%w: Expected %d comma-separated values, got %d",
			ErrInvalidFormat, dataFieldsCount, len(parts))
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidSteps, err)
	}
	if steps <= 0 {
		return fmt.Errorf("%w: %d", ErrInvalidSteps, steps)
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidDuration, err)
	}
	if duration <= 0 {
		return fmt.Errorf("%w: %v", ErrInvalidDuration, duration)
	}

	ds.Steps = steps
	ds.Duration = duration

	return nil
}

func (ds DaySteps) ActionInfo() (string, error) {
	distance := spentenergy.Distance(ds.Steps, ds.Height)

	calories, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)
	if err != nil {
		return "", fmt.Errorf("Calorie calculation: %w", err)
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps, distance, calories,
	), nil
}
