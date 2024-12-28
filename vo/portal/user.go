package portalvo

type UserInfoVO struct {
	ID          int64   `json:"id"`
	Username    string  `json:"username"`
	Nickname    *string `json:"nickname,omitempty"`
	Capacity    int64   `json:"capacity"`
	UseCapacity int64   `json:"useCapacity"`
}
