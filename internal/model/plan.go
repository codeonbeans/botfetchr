package model

//go:generate stringer -type Plan
type Plan int

const (
	PlanFree Plan = iota
	PlanPro
	PlanLifetime
)
