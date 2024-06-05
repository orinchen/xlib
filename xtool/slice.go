package xtool

func SliceConvert[ST, TT any](source []ST, mapFunc func(ST) (TT, error)) ([]TT, error) {
	var tags []TT
	for _, item := range source {
		if tag, err := mapFunc(item); err != nil {
			return nil, err
		} else {
			tags = append(tags, tag)
		}
	}
	return tags, nil
}

func HasDuplicated[ST any](source []ST, equalFunc func(ST, ST) bool) bool {
	for i := 0; i < len(source); i++ {
		for j := i + 1; j < len(source); j++ {
			if equalFunc(source[i], source[j]) {
				return true
			}
		}
	}
	return false
}

func SameOfAll[ST any](source []ST, equalFunc func(ST, ST) bool) bool {
	for i := 0; i < len(source); i++ {
		for j := i + 1; j < len(source); j++ {
			if !equalFunc(source[i], source[j]) {
				return false
			}
		}
	}
	return true
}
