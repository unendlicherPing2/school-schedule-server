package schedule_service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type result struct {
	Success  bool `json:"success"`
	Message  any  `json:"message"`
	Messages any  `json:"messages"`
	Data     struct {
		Hint     string `json:"hint"`
		Schedule []struct {
			Date               string `json:"date"`
			Important          string `json:"important"`
			TeacherToken       string `json:"teacher_token"`
			TeacherForename    string `json:"teacher_forename"`
			TeacherSurname     string `json:"teacher_surname"`
			SubTeacherToken    string `json:"sub_teacher_token"`
			SubTeacherSurname  string `json:"sub_teacher_surname"`
			SubTeacherForename string `json:"sub_teacher_forename"`
			Hour               string `json:"hour"`
			Block              string `json:"block"`
			SubjectToken       string `json:"subject_token"`
			SubjectLong        string `json:"subject_long"`
			SubSubjectToken    string `json:"sub_subject_token"`
			SubSubjectLong     string `json:"sub_subject_long"`
			ClassToken         string `json:"class_token"`
			ClassLong          string `json:"class_long"`
			Hints              string `json:"hints"`
			RoomToken          string `json:"room_token"`
			SubRoomToken       string `json:"sub_room_token"`
			SubRoomLong        string `json:"sub_room_long"`
			Teacher            string `json:"teacher"`
			SubTeacher         string `json:"sub_teacher"`
			HlRoom             bool   `json:"hl_room"`
		} `json:"schedule"`
		Supervision []struct {
			Hour               string `json:"hour"`
			HourLabel          string `json:"hour_label"`
			TeacherToken       string `json:"teacher_token"`
			TeacherSurname     string `json:"teacher_surname"`
			TeacherForename    string `json:"teacher_forename"`
			SubTeacherToken    string `json:"sub_teacher_token"`
			SubTeacherSurname  string `json:"sub_teacher_surname"`
			SubTeacherForename string `json:"sub_teacher_forename"`
			CorridorToken      string `json:"corridor_token"`
			CorridorLong       string `json:"corridor_long"`
			Mode               string `json:"mode"`
			Teacher            string `json:"teacher"`
			SubTeacher         string `json:"sub_teacher"`
			Del                bool   `json:"del"`
		} `json:"supervision"`
		ScheduleSum        string `json:"schedule_sum"`
		SupervisionSum     string `json:"supervision_sum"`
		ScheduleEnabled    bool   `json:"schedule_enabled"`
		SupervisionEnabled bool   `json:"supervision_enabled"`
	} `json:"data"`
}

type ScheduleEntries []struct {
	TeacherToken       string `json:"teacher_token"`
	TeacherForename    string `json:"teacher_forename"`
	TeacherSurname     string `json:"teacher_surname"`
	SubTeacherToken    string `json:"sub_teacher_token"`
	SubTeacherSurname  string `json:"sub_teacher_surname"`
	SubTeacherForename string `json:"sub_teacher_forename"`
	Hour               string `json:"hour"`
	Block              string `json:"block"`
	SubjectToken       string `json:"subject_token"`
	SubjectLong        string `json:"subject_long"`
	SubSubjectToken    string `json:"sub_subject_token"`
	SubSubjectLong     string `json:"sub_subject_long"`
	ClassToken         string `json:"class_token"`
	ClassLong          string `json:"class_long"`
	Hints              string `json:"hints"`
	RoomToken          string `json:"room_token"`
	SubRoomToken       string `json:"sub_room_token"`
	SubRoomLong        string `json:"sub_room_long"`
	Teacher            string `json:"teacher"`
	SubTeacher         string `json:"sub_teacher"`
}

type Schedule struct {
	Hint     string          `db:"hint"`
	Date     string          `db:"date"`
	Schedule ScheduleEntries `db:"schedule" json:"schedule"`
}

type ScheduleDto struct {
	Hint     string `db:"hint"`
	Date     string `db:"date"`
	Schedule string `db:"schedule"`
}

func (schedule *Schedule) fromResponse(response result, date string) {
	schedule.Hint = response.Data.Hint
	schedule.Date = date

	schedule.Schedule = make([]struct {
		TeacherToken       string `json:"teacher_token"`
		TeacherForename    string `json:"teacher_forename"`
		TeacherSurname     string `json:"teacher_surname"`
		SubTeacherToken    string `json:"sub_teacher_token"`
		SubTeacherSurname  string `json:"sub_teacher_surname"`
		SubTeacherForename string `json:"sub_teacher_forename"`
		Hour               string `json:"hour"`
		Block              string `json:"block"`
		SubjectToken       string `json:"subject_token"`
		SubjectLong        string `json:"subject_long"`
		SubSubjectToken    string `json:"sub_subject_token"`
		SubSubjectLong     string `json:"sub_subject_long"`
		ClassToken         string `json:"class_token"`
		ClassLong          string `json:"class_long"`
		Hints              string `json:"hints"`
		RoomToken          string `json:"room_token"`
		SubRoomToken       string `json:"sub_room_token"`
		SubRoomLong        string `json:"sub_room_long"`
		Teacher            string `json:"teacher"`
		SubTeacher         string `json:"sub_teacher"`
	}, len(response.Data.Schedule))

	for i, item := range response.Data.Schedule {
		schedule.Schedule[i] = struct {
			TeacherToken       string `json:"teacher_token"`
			TeacherForename    string `json:"teacher_forename"`
			TeacherSurname     string `json:"teacher_surname"`
			SubTeacherToken    string `json:"sub_teacher_token"`
			SubTeacherSurname  string `json:"sub_teacher_surname"`
			SubTeacherForename string `json:"sub_teacher_forename"`
			Hour               string `json:"hour"`
			Block              string `json:"block"`
			SubjectToken       string `json:"subject_token"`
			SubjectLong        string `json:"subject_long"`
			SubSubjectToken    string `json:"sub_subject_token"`
			SubSubjectLong     string `json:"sub_subject_long"`
			ClassToken         string `json:"class_token"`
			ClassLong          string `json:"class_long"`
			Hints              string `json:"hints"`
			RoomToken          string `json:"room_token"`
			SubRoomToken       string `json:"sub_room_token"`
			SubRoomLong        string `json:"sub_room_long"`
			Teacher            string `json:"teacher"`
			SubTeacher         string `json:"sub_teacher"`
		}{
			TeacherToken:       item.TeacherToken,
			TeacherForename:    item.TeacherForename,
			TeacherSurname:     item.TeacherSurname,
			SubTeacherToken:    item.SubTeacherToken,
			SubTeacherSurname:  item.SubTeacherSurname,
			SubTeacherForename: item.SubTeacherForename,
			Hour:               item.Hour,
			Block:              item.Block,
			SubjectToken:       item.SubjectToken,
			SubjectLong:        item.SubjectLong,
			SubSubjectToken:    item.SubSubjectToken,
			SubSubjectLong:     item.SubSubjectLong,
			ClassToken:         item.ClassToken,
			ClassLong:          item.ClassLong,
			Hints:              item.Hints,
			RoomToken:          item.RoomToken,
			SubRoomToken:       item.SubRoomToken,
			SubRoomLong:        item.SubRoomLong,
			Teacher:            item.Teacher,
			SubTeacher:         item.SubTeacher,
		}
	}
}

func GetSchedule(date string) (Schedule, error) {
	var result result
	var schedule Schedule

	url := fmt.Sprintf("https://lenne-schule.de/unterricht/vertretungsplan?format=json&date=%s", date)
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return schedule, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		return schedule, err
	}

	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Printf("error making http request: %s\n", err)
		return schedule, err
	}

	schedule.fromResponse(result, date)

	return schedule, nil
}
