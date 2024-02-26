package xtool

func SliceConvert[ST, TT any](source []ST, mapFunc func(ST) TT) []TT {
	var tags []TT
	for _, item := range source {
		tags = append(tags, mapFunc(item))
	}
	return tags
}
