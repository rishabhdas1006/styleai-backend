package common

const (
	ContextUserID = "user_id"
	ContextRole   = "role"
)

type Gender string

const (
	GenderMen    Gender = "men"
	GenderWomen  Gender = "women"
	GenderKids   Gender = "kids"
	GenderUnisex Gender = "unisex"
)

var validGenders = map[Gender]struct{}{
	GenderMen:    {},
	GenderWomen:  {},
	GenderKids:   {},
	GenderUnisex: {},
}

func IsValidGender(g Gender) bool {
	_, ok := validGenders[g]
	return ok
}
