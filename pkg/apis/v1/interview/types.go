package interview

type (
	createInterviewReq struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Date     int64  `json:"date"`
		Location string `json:"location"`
		Position string `json:"position"`

		IntervieweeUid string `json:"interviewee_uid"`
	}

	deleteInterviewReq struct {
		ID int64 `json:"id"`
	}
)
