package common

func WrapLike(v string) string {
	return "%" + v + "%"
}

func WrapLikeLeft(v string) string {
	return "%" + v
}

func WrapLikeRight(v string) string {
	return "%" + v
}
