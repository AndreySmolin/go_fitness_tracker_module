package spentcalories

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// parseTraining возвращает количесво шагов, вид и продолжительность активности
func parseTraining(data string) (int, string, time.Duration, error) {
	sliceData := strings.Split(data, ",")
	if len(sliceData) != 3 {
		return 0, "", 0, fmt.Errorf("length of slice %s is not equal to 3", sliceData)
	}
	steps, err := strconv.Atoi(sliceData[0])
	if err != nil {
		return 0, "", 0, fmt.Errorf("number %s conversion error: %w", sliceData[0], err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("the number of steps %d must not be negative", steps)
	}
	duration, er := time.ParseDuration(sliceData[2])
	if er != nil {
		return 0, "", 0, fmt.Errorf("error converting string %s: %w", sliceData[2], er)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("duration cannot be less than or equal to 0")
	}
	return steps, sliceData[1], duration, nil
}

// distance расчитывает  дистанцию
func distance(steps int, height float64) float64 {
	return float64(steps) * height * stepLengthCoefficient / mInKm
}

// meanSpeed расчитывает среднюю скорость
func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	return distance(steps, height) / duration.Hours()
}

// TrainingInfo возвращает  строку с информацией о тренировке
func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		return "", fmt.Errorf("parseTraining function error: %w", err)
	}
	switch activity {
	case "Ходьба":
		distance := distance(steps, height)
		speed := meanSpeed(steps, height, duration)
		calories, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("calorie counting error: %w", err)
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, duration.Hours(), distance, speed, calories), nil
	case "Бег":
		distanceRun := distance(steps, height)
		speedRun := meanSpeed(steps, height, duration)
		caloriesRun, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", fmt.Errorf("calorie counting error: %w", err)
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, duration.Hours(), distanceRun, speedRun, caloriesRun), nil
	}
	return "", fmt.Errorf("неизвестный тип тренировки")
}

// RunningSpentCalories возвращает  количество калорий, потраченных при беге
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 {
		return 0, fmt.Errorf("error: weight %.2f cannot be negative or equal to zero", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("error: height %.2f cannot be negative or equal to zero", weight)
	}
	if steps <= 0 {
		return 0, fmt.Errorf("error: steps %d cannot be negative or equal to zero", steps)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("error: duration cannot be negative or equal to zero")
	}
	speed := meanSpeed(steps, height, duration)
	return (weight * duration.Minutes() * speed / 60), nil
}

// WalkingSpentCalories возвращает количество калорий, потраченных при ходьбе
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if weight <= 0 {
		return 0, fmt.Errorf("error: weight %.2f cannot be negative or equal to zero", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("error: height %.2f cannot be negative or equal to zero", weight)
	}
	if steps <= 0 {
		return 0, fmt.Errorf("error: steps %d cannot be negative or equal to zero", steps)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("error: duration cannot be negative or equal to zero")
	}
	speed := meanSpeed(steps, height, duration)
	return (weight * duration.Minutes() * speed / 60) * walkingCaloriesCoefficient, nil
}
