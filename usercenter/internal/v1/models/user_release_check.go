package models

type UserReleaseCheck struct{}

func NewUserReleaseCheck() *UserReleaseCheck {
	return &UserReleaseCheck{}
}

func (uc *UserReleaseCheck) CheckCanReq(uid string) bool {
	//if env.ReleaseMode == "release" {
	//	phone := queryUserPhone(uid)
	//	if phone != "" && !tools.InArray(phone, releaseUserPhones) {
	//		return false
	//	}
	//}
	return true
}