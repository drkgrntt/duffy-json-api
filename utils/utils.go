package utils

func ForEach[
	Type interface{},
](
	slice []Type,
	callback func(item Type, index int, slice []Type),
) {
	for index, item := range slice {
		callback(item, index, slice)
	}
}

func Map[
	InputType interface{},
	OutputType interface{},
](
	slice []InputType,
	callback func(item InputType, index int, slice []InputType) OutputType,
) []OutputType {
	result := make([]OutputType, len(slice))
	for index, item := range slice {
		newItem := callback(item, index, slice)
		result[index] = newItem
	}
	return result
}

func Filter[
	Type interface{},
](
	slice []Type,
	callback func(item Type, index int, slice []Type) bool,
) []Type {
	result := make([]Type, 0)
	for index, item := range slice {
		shouldInclude := callback(item, index, slice)
		if shouldInclude {
			result = append(result, item)
		}
	}
	return result
}

func Some[
	Type interface{},
](
	slice []Type,
	callback func(item Type, index int, slice []Type) bool,
) bool {
	for index, item := range slice {
		matchesCondition := callback(item, index, slice)
		if matchesCondition {
			return true
		}
	}
	return false
}

func Reduce[
	InputType interface{},
	OutputType interface{},
](
	slice []InputType,
	callback func(accumulator OutputType, item InputType, index int, slice []InputType) OutputType,
	accumulator OutputType,
) OutputType {
	for index, item := range slice {
		accumulator = callback(accumulator, item, index, slice)
	}
	return accumulator
}
