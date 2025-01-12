package utils

import (
	"fmt"
	"strings"
)

func ValidatePair(pair string) error {
	if strings.Contains(pair, "/") {
		parts := strings.Split(pair, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid pair format: %s", pair)
		}
	} else {
		return fmt.Errorf("invalid pair format: %s", pair)
	}

	return nil
}

func ConvertToPointerSlice[T any](inputSlice []T) []*T {
	pointerSlice := make([]*T, len(inputSlice))
	for i := range inputSlice {
		pointerSlice[i] = &inputSlice[i]
	}
	return pointerSlice
}
