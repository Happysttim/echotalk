package handler

const (
	LimitSmall  = 10
	LimitMedium = 30
	LimitLarge  = 30
)

const (
	TypeID        = "_id"
	TypeUpdatedAt = "updated_at"
	TypeRate      = "rate_up"
)

var Limits = []int{LimitSmall, LimitMedium, LimitLarge}
var SurveyTypes = []string{TypeID, TypeUpdatedAt}
var AnswerTypes = []string{TypeID, TypeRate}
