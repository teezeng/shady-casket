package spentenergy

import (
	"errors"
	"fmt"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе.
)

var (
	ErrInvalidSteps    = errors.New("Steps count must be greater than zero")
	ErrInvalidWeight   = errors.New("Weight must be greater than zero")
	ErrInvalidHeight   = errors.New("Height must be greater than zero")
	ErrInvalidDuration = errors.New("Duration must be greater than zero")
)

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := runningCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	return calories * walkingCaloriesCoefficient, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	return runningCalories(steps, weight, height, duration)
}

func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps <= 0 || height <= 0 || duration <= 0 {
		return 0
	}

	return Distance(steps, height) / duration.Hours()
}

func Distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient

	return float64(steps) * stepLength / mInKm
}

func runningCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validate(steps, weight, height, duration); err != nil {
		return 0, err
	}

	meanSpeed := MeanSpeed(steps, height, duration)

	return weight * meanSpeed * duration.Minutes() / minInH, nil
}

func validate(steps int, weight, height float64, duration time.Duration) error {
	switch {
	case steps <= 0:
		return fmt.Errorf("%w: %d", ErrInvalidSteps, steps)
	case weight <= 0:
		return fmt.Errorf("%w: %v", ErrInvalidWeight, weight)
	case height <= 0:
		return fmt.Errorf("%w: %v", ErrInvalidHeight, height)
	case duration <= 0:
		return fmt.Errorf("%w: %v", ErrInvalidDuration, duration)
	}

	return nil
}
