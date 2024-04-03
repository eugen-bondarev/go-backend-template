package model

import "github.com/eugen-bondarev/go-slice-helpers/parallel"

type ModelMapper[TModel any, T any] interface {
	FromModel(TModel) T
	ToModel(T) TModel
}

func ManyFromModel[TModel any, T any](m ModelMapper[TModel, T], items []TModel) []T {
	return parallel.Map(items, m.FromModel)
}

func ManyToModel[TModel any, T any](m ModelMapper[TModel, T], items []T) []TModel {
	return parallel.Map(items, m.ToModel)
}
