package event

//bus events handled by the playback package
const (
	QueueEnqueue     Type = "queue:enqueue"
	QueueUpdate      Type = "queue:update"
	QueueDeleteAnnon Type = "queue:deleteAnnon"
	QueueChange      Type = "queue:change"
	GetQueue         Type = "queue:get"
	QueueContent     Type = "queue:content"
)
