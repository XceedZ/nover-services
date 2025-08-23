package tables

// DailyCheckinReward merepresentasikan data dari tabel daily_checkin_rewards.
type DailyCheckinReward struct {
	DayNumber    int `json:"dayNumber" db:"day_number"`
	RewardAmount int `json:"rewardAmount" db:"reward_amount"`
}

// CheckinStatusResponse adalah struct untuk response API status check-in.
type CheckinStatusResponse struct {
	// ✨ NAMA FIELD DIPERBARUI
	TotalCheckinsThisMonth int                  `json:"totalCheckinsThisMonth"`
	TodayCheckedIn         bool                 `json:"todayCheckedIn"`
	Rewards                []DailyCheckinReward `json:"rewards"`
	CheckedInDates         []string             `json:"checkedInDates"` // Format "YYYY-MM-DD"
}