package tables

import "time"

type Genre struct {
	GenreID           int64      `json:"genreId"`
	GenreName         string     `json:"genreName"`
	GenreTl           string     `json:"genreTl"`
	Remark            *string    `json:"remark,omitempty"`
	ActiveDatetime    time.Time  `json:"activeDatetime"`
	NonActiveDatetime *time.Time `json:"nonActiveDatetime,omitempty"`
	CreateDatetime    time.Time  `json:"createDatetime"`
	UpdateDatetime    *time.Time `json:"updateDatetime,omitempty"`
}
