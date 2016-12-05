package train

import (
	"time"

	log "github.com/Sirupsen/logrus"
	r "github.com/dancannon/gorethink"
)

// Entity represents different entities we can attach settings to
type Entity string

// UpdateMode tells us if the announcement should be managed automatically or not
type UpdateMode string

// Predefined SettingsType values
const (
	TrainInstance Entity = "train"
	TrainOrder    Entity = "trainOrder"
	TrainNumber   Entity = "trainNumber"

	Auto   UpdateMode = "AUTO"
	Manual UpdateMode = "MANUAL"

	Category string = "category"
)

// Settings represents additional options that can be set for a given train
type Settings struct {
	ID         string                       `gorethink:"id" json:"id"`
	Type       Entity                       `gorethink:"entity" json:"entity"`
	Audio      *map[AnnonOption]string      `gorethink:"audio" json:"audio"`
	Lang       *map[string]LanguageSettings `gorethink:"lang,omitempty" json:"lang,omitempty"`
	Overrides  *map[string]string           `gorethink:"overrides,omitempty" json:"overrides,omitempty"`
	Display    *map[string]string           `gorethink:"display" json:"display"`
	Mode       UpdateMode                   `gorethink:"updateMode" json:"updateMode"`
	FirstEvent *time.Time                   `gorethink:"firstEvent" json:"firstEvent"`
	Arrival    *TimetableEvent              `gorethink:"arrival" json:"arrival"`
	Departure  *TimetableEvent              `gorethink:"departure" json:"departure"`
	Delay      int                          `gorethink:"delay" json:"delay"`
}

// SettingsChanges represents changefeeds data from database
type SettingsChanges struct {
	NewVal *Settings `gorethink:"new_val,omitempty"`
	OldVal *Settings `gorethink:"old_val,omitempty"`
}

type LanguageSettings struct {
	Enabled bool                    `gorethink:"enabled" json:"enabled"`
	I18n    *map[AnnonOption]string `gorethink:"i18n,omitempty" json:"i18n,omitempty"`
}

/*
Store methods related to Settings
*/

func (store *rethinkStore) SaveSettings(st *Settings) (*Settings, error) {
	resp, err := r.DB(store.db).Table(store.settingsTable).Insert(st, r.InsertOpts{
		Conflict: "replace",
	}).RunWrite(store.session)
	if err != nil {
		log.WithFields(log.Fields{
			"logger":   "train.store.settings",
			"settings": st,
			"error":    err,
		}).Errorln("Error occured while saving settings object.")
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{
			"logger":   "train.store.settings",
			"realtime": st,
			"result":   resp,
		}).Debugln("Database save result.")
	}
	return st, err
}

func (store *rethinkStore) SaveAllSettings(sts *[]Settings) error {
	res, err := r.DB(store.db).Table(store.settingsTable).Insert(*sts, r.InsertOpts{
		Conflict: "replace",
	}).RunWrite(store.session)
	if err != nil {
		log.WithFields(log.Fields{"logger": "train.store.settings", "method": "SaveAllSettings"}).
			WithError(err).Error("Error occured while saving settings objects.")
		return err
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "train.store.settings", "method": "SaveAllSettings", "settings": sts, "result": res}).
			Debug("Database save result.")
	}
	return err
}

func (store *rethinkStore) GetSettings(ID string) (*Settings, error) {
	resp, err := r.DB(store.db).Table(store.settingsTable).Get(ID).Run(store.session)
	if err != nil {
		log.WithFields(log.Fields{"logger": "train.store.settings", "method": "GetSettings", "settings": ID}).
			WithError(err).Error("Error while retrieving settings.")
		return nil, err
	}
	if resp.IsNil() {
		return nil, err
	}
	s := new(Settings)
	err = resp.One(s)
	if err != nil {
		log.WithFields(log.Fields{"logger": "train.store.settings", "method": "GetSettings", "settings": ID}).
			WithError(err).Error("Error while parsing settings result.")
		return nil, err
	}
	return s, err
}

func (store *rethinkStore) GetAllSettings(IDs ...string) (*[]Settings, error) {
	var res *r.Cursor
	var err error
	if len(IDs) == 0 {
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "train.store.settings", "method": "GetAllSettings"}).
				Debug("Getting all settings.")
		}
		res, err = r.DB(store.db).Table(store.settingsTable).Run(store.session)
	} else {
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "train.store.settings", "method": "GetAllSettings", "settings": IDs}).
				Debug("Getting settings for IDs.")
		}
		res, err = r.DB(store.db).Table(store.settingsTable).GetAll(r.Args(IDs)).Run(store.session)
	}
	if err != nil {
		log.WithFields(log.Fields{"logger": "train.store.settings", "method": "GetAllSettings", "settings": IDs}).
			WithError(err).Error("Error loading settings for train.")
		return nil, err
	}
	s := new([]Settings)
	err = res.All(s)
	if err != nil {
		log.WithFields(log.Fields{"logger": "train.store.settings", "method": "GetAllSettings", "settings": IDs}).
			WithError(err).Error("Error parsing settings.")
		return nil, err
	}
	return s, err
}

func (store *rethinkStore) NextSettingsChange() (*Settings, *Settings, error) {
	if store.settingsChange == nil {
		err := store.initSettingsChangeCursor()
		if err != nil {
			return nil, nil, err
		}
	}
	var next bool
	changes := new(SettingsChanges)
	if next = store.settingsChange.Next(changes); !next {
		err := store.settingsChange.Err()
		log.WithFields(log.Fields{"logger": "train.store.settings", "method": "NextSettingsChange"}).
			WithError(err).Error("Error while reading from settings table changes cursor.")
		return nil, nil, err
	}
	return changes.NewVal, changes.OldVal, nil
}

func (store *rethinkStore) initSettingsChangeCursor() error {
	if log.GetLevel() > log.InfoLevel {
		log.WithFields(log.Fields{"logger": "train.store.settings", "method": "initSettingsChangeCursor", "windowSize": store.windowSize}).
			Info("Initializing change cursor for train settings changes.")
	}
	res, err := r.DB(store.db).Table(store.settingsTable).
		Changes(r.ChangesOpts{IncludeInitial: false}).Run(store.session)
	if err != nil {
		log.WithFields(log.Fields{"logger": "train.store.settings", "method": "initSettingsChangeCursor", "windowSize": store.windowSize}).
			WithError(err).Error("Error while initializing train settings changes cursor.")
		return err
	}
	store.settingsChange = res
	return nil
}
