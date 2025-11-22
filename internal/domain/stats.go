package domain

type UserReviewStat struct {
	UserID       string `json:"user_id"`
	ReviewsCount int    `json:"reviews_count"`
}
type Stats struct {
	TotalPR        int              `json:"total_pr"`
	OpenPR         int              `json:"open_pr"`
	MergedPR       int              `json:"merged_pr"`
	ReviewsPerUser []UserReviewStat `json:"reviews_per_user"`
}
