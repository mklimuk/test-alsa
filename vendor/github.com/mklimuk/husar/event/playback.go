package event

// event bus events handled by the playback package
const (
	PlaybackStart      Type = "playback:start"
	PlaybackEnd        Type = "playback:end"
	PlaybackToggleMute Type = "playback:toggleMute"
	PlaybackMuteAll    Type = "playback:muteAll"
	PlaybackUnmuteAll  Type = "playback:unmuteAll"
	SetVolume          Type = "playback:setVolume"
	SetVolumeAck       Type = "playback:setVolume:ack"
	PlaybackTrigger    Type = "playback:trigger"
	PlaybackTriggerAck Type = "playback:trigger:ack"
)
