package model

//go:generate stringer -type Feature
type Feature int

const (
	FeatureGetMedia Feature = iota
)
