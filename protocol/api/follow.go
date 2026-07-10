package api

// FollowCountReturn models a condenser_api.get_follow_count response.
// See conveyor/src/user-search/client.ts FollowCountReturn.
type FollowCountReturn struct {
	Account        string `json:"account"`
	FollowerCount  int    `json:"follower_count"`
	FollowingCount int    `json:"following_count"`
}

// FollowReturn models a single entry in a condenser_api.get_followers /
// get_following response. The `what` array contains role strings such as
// ["blog"] or ["ignore"]. See conveyor/src/user-search/client.ts FollowReturn.
type FollowReturn struct {
	Follower  string   `json:"follower"`
	Following string   `json:"following"`
	What      []string `json:"what"`
}
