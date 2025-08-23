// File: pkg/dao/mission_dao.go
package dao

import (
	"context"
	"fmt"
	"noversystem/pkg/tables"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MissionDao struct {
	DB *pgxpool.Pool
}

func NewMissionDao(db *pgxpool.Pool) *MissionDao {
	return &MissionDao{DB: db}
}

// GetActiveMissionsWithUserProgress mengambil semua misi aktif dan progres user hari ini.
func (d *MissionDao) GetActiveMissionsWithUserProgress(ctx context.Context, userID int64) ([]tables.MissionStatus, error) {
	today := time.Now().Format("2006-01-02")
	
	// Struct sementara untuk menampung hasil query flat
	type flatMissionData struct {
		tables.Mission
		tables.MissionTier
		CurrentProgress *int   `db:"current_value"`
		LastClaimedTier *int64 `db:"last_claimed_tier_id"`
	}

	var flatResults []flatMissionData
	
	query := `
		SELECT
			m.mission_id, m.title, m.description,
			mt.tier_id, mt.threshold, mt.reward_amount, mt.tier_order,
			ump.current_value,
			ump.last_claimed_tier_id
		FROM
			missions m
		JOIN
			mission_tiers mt ON m.mission_id = mt.mission_id
		LEFT JOIN
			user_mission_progress ump ON m.mission_id = ump.mission_id 
			AND ump.user_id = $1 AND ump.progress_date = $2
		WHERE
			m.is_active = TRUE
		ORDER BY
			m.mission_id, mt.tier_order ASC
	`
	err := pgxscan.Select(ctx, d.DB, &flatResults, query, userID, today)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data misi: %w", err)
	}

	// Proses hasil flat menjadi struktur nested yang rapi
	missionMap := make(map[int64]*tables.MissionStatus)
	var finalMissions []tables.MissionStatus

	for _, row := range flatResults {
		// Jika misi ini belum ada di map, buat entri baru
		if _, exists := missionMap[row.MissionID]; !exists {
			missionMap[row.MissionID] = &tables.MissionStatus{
				MissionInfo: tables.Mission{
					MissionID:   row.MissionID,
					Title:       row.Title,
					Description: row.Description,
				},
				Tiers:           make([]tables.MissionTier, 0),
				CurrentProgress: 0,
				LastClaimedTier: 0,
			}
			// Karena query-nya diurutkan, kita bisa langsung tambahkan ke slice final
			finalMissions = append(finalMissions, *missionMap[row.MissionID])
		}

		// Tambahkan tier ke misi yang sesuai
		missionStatus := missionMap[row.MissionID]
		missionStatus.Tiers = append(missionStatus.Tiers, row.MissionTier)

		// Update progress (hanya akan di-update sekali per misi)
		if row.CurrentProgress != nil {
			missionStatus.CurrentProgress = *row.CurrentProgress
		}
		if row.LastClaimedTier != nil {
			missionStatus.LastClaimedTier = *row.LastClaimedTier
		}
	}

	// Karena kita menggunakan pointer, perubahan di map akan tercermin di slice
	// Kita perlu update slice final dengan data dari map setelah loop selesai
	for i := range finalMissions {
		finalMissions[i] = *missionMap[finalMissions[i].MissionInfo.MissionID]
	}

	return finalMissions, nil
}