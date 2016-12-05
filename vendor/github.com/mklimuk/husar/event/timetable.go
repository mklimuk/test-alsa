package event

// events handled by the timetable package
const (
	TimetableStations   Type = "timetable:stations"
	TimetableContent    Type = "timetable:timetable"
	GetTimetable        Type = "timetable:get"
	TrainUpdateEvent    Type = "timetable:event:update"
	TrainUpdateMode     Type = "timetable:mode:update"
	TrainUpdate         Type = "timetable:train:update"
	AudioSettingsUpdate Type = "timetable:train:updateAudioSettings"
)
