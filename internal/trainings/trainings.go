package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

const (
	TypeRunning = "Бег"
	TypeWalking = "Ходьба"
)

const dataFieldsCount = 3

var (
	ErrInvalidFormat       = errors.New("Invalid data format")
	ErrInvalidSteps        = errors.New("Invalid steps count")
	ErrInvalidDuration     = errors.New("Invalid duration")
	ErrUnknownTrainingType = errors.New("Unknown training type")
)

type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	personaldata.Personal
}

func (t *Training) Parse(datastring string) (err error) {
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

	duration, err := time.ParseDuration(parts[2])
	if err != nil {
		return fmt.Errorf("%w: %w", ErrInvalidDuration, err)
	}
	if duration <= 0 {
		return fmt.Errorf("%w: %v", ErrInvalidDuration, duration)
	}

	t.Steps = steps
	t.TrainingType = parts[1]
	t.Duration = duration

	return nil
}

func (t Training) ActionInfo() (string, error) {
	distance := spentenergy.Distance(t.Steps, t.Height)
	meanSpeed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	var (
		calories float64
		err      error
	)

	switch t.TrainingType {
	case TypeRunning:
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case TypeWalking:
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", ErrUnknownTrainingType
	}
	if err != nil {
		return "", fmt.Errorf("Calorie calculation: %w", err)
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, t.Duration.Hours(), distance, meanSpeed, calories,
	), nil
}
