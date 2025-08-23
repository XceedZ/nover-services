package tables

type Mission struct {
	MissionID   int64  `json:"missionId" db:"mission_id"`
	Title       string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
}

type MissionTier struct {
	TierID       int64 `json:"tierId" db:"tier_id"`
	Threshold    int   `json:"threshold" db:"threshold"`
	RewardAmount int   `json:"rewardAmount" db:"reward_amount"`
	TierOrder    int   `json:"tierOrder" db:"tier_order"`
}

type MissionStatus struct {
	MissionInfo      Mission       `json:"missionInfo"`
	Tiers            []MissionTier `json:"tiers"`
	CurrentProgress  int           `json:"currentProgress"`
	LastClaimedTier  int64         `json:"lastClaimedTier"`
}