package entity

type Follower struct {
	FollowingId string `json:"following_id"`
	FollowerId  string `json:"follower_id"`
	UnFollowed  bool   `json:"followed"`
}
