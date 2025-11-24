package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	stepLength = 0.65 // Длина одного шага в метрах
	mInKm      = 1000 // Количество метров в одном километре
)

// parsePackage возвращает количесво шагов и продолжительность прогулки
func parsePackage(data string) (int, time.Duration, error) {
	sliceData := strings.Split(data, ",")
	if len(sliceData) != 2 {
		return 0, 0, fmt.Errorf("length of slice %s is not equal to 2", sliceData)
	}
	steps, err := strconv.Atoi(sliceData[0])
	if err != nil {
		return 0, 0, fmt.Errorf("number %s conversion error: %w", sliceData[0], err)
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("the number of steps %d must not be negative", steps)
	}
	duration, er := time.ParseDuration(sliceData[1])
	if er != nil {
		return 0, 0, fmt.Errorf("error converting string %s: %w", sliceData[1], er)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("duration cannot be less than or equal to 0")
	}
	return steps, duration, nil
}

// DayActionInfo возвращает количесво шагов,дистанцию и сожженные калории
func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Printf("Function error parsePackage: %v", err)
		return ""
	}
	distance := float64(steps) * stepLength / mInKm
	calories, er := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if er != nil {
		log.Printf("Function error WalkingSpentCalories: %v", er)
		return ""
	}
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distance, calories)
}
