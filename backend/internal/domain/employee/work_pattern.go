package employee

import "fmt"

type WorkPattern struct {
	WorkDays int
	RestDays int
}

func NewWorkPattern(workDays, restDays int) (WorkPattern, error) {
	if workDays < 1 {
		return WorkPattern{}, fmt.Errorf("work days must be at least 1")
	}
	if restDays < 1 {
		return WorkPattern{}, fmt.Errorf("rest days must be at least 1")
	}
	return WorkPattern{WorkDays: workDays, RestDays: restDays}, nil
}

func DefaultWorkPattern() WorkPattern {
	return WorkPattern{WorkDays: 4, RestDays: 1}
}

func (wp WorkPattern) IsRestDay(diaDelAnio int) bool {
	cycleLength := wp.WorkDays + wp.RestDays
	positionInCycle := (diaDelAnio - 1) % cycleLength
	return positionInCycle >= wp.WorkDays
}
