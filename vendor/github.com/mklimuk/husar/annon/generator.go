package annon

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/text/language"

	"github.com/mklimuk/husar/i18n"
	"github.com/mklimuk/husar/train"

	log "github.com/Sirupsen/logrus"
)

const nightStart int = 22 * 60
const nightEnd int = 7 * 60

// Generator interface defines operations required to generate all vocal announcements for a given train
type Generator interface {
	ForTrain(t *train.Train, old *train.Train, currentAnnons *[]*Announcement, now *time.Time, cat *Catalog) (announcements []*Announcement, toDelete []*Announcement, err error)
	Generate(t *train.Train, now *time.Time, temp *CatalogTemplate, p *Producer, lang language.Tag) (*Announcement, error)
	FromTemplate(tpl *CatalogTemplate, startTime *time.Time, endTime *time.Time, dayInterval *time.Duration, nightInterval *time.Duration, now *time.Time, params TemplateParams, stationID int) (*Announcement, error)
	Preview(tpl *CatalogTemplate, params TemplateParams) (*Text, error)
}

type gen struct {
	producers  Producers
	cache      map[string]*CatalogTemplate
	priorities map[string]int
}

// NewGenerator is a constructor used tu build the default generator
func NewGenerator(p Producers, prio map[string]int) Generator {
	gen := new(gen)
	gen.producers = p
	gen.priorities = prio
	gen.cache = make(map[string]*CatalogTemplate)
	return Generator(gen)
}

// ForTrain builds a slice of announcements of different type for a train
func (g *gen) ForTrain(t *train.Train, old *train.Train, currentAnnons *[]*Announcement, now *time.Time, cat *Catalog) (announcements []*Announcement, toDelete []*Announcement, err error) {

	clog := log.WithFields(log.Fields{"logger": "generator", "method": "ForTrain"})

	// if the train is not there we quietly return
	announcements = []*Announcement{}
	if t == nil {
		return
	}
	for _, prod := range g.producers {
		var required bool
		var needsUpdate bool
		if required, needsUpdate, err = (*prod).Required(t, old, now); err != nil {
			clog.WithError(err).
				Error("Error encountered while checking if the producer is required.")
			continue
		}
		if required == true {
			if log.GetLevel() >= log.DebugLevel {
				clog = clog.WithFields(log.Fields{"train": (*t).ID, "producer": (*prod).Name()})
				clog.Debug("Producer is required for the given train.")
			}
			if needsUpdate == true {
				clog.Debug("Announcements for producer need to be updated.")
				var a *Announcement
				var tem *CatalogTemplate
				if tem, err = g.getTemplate(cat, (*prod).Name()); err != nil {
					clog.WithError(err).Error("Error getting template from the catalog.")
					continue
				}
				if tem == nil {
					clog.WithField("template", (*prod).Name()).Info("Template not found")
					continue
				}
				if a, err = g.Generate(t, now, tem, prod, i18n.DefaultLang); err != nil {
					clog.WithError(err).Error("Error while generating announcement. Aborting.")
					continue
				}
				announcements = append(announcements, a)
				// check if we need other languages
				if opt, v := train.HasI18n(t); opt == true {
					if tem.Translations == nil {
						clog.WithField("template", tem.ID).
							Warn("Translations requested but none can be found for template")
						continue
					}
					for code, l := range *v {
						if code == i18n.DefaultLang.String() {
							// this is covered before
							continue
						}
						lang := i18n.GetLang(code)
						//we add announcements for all enabled languages
						if l.Enabled {
							subID := (*tem.Translations)[code]
							if subID == "" {
								clog.WithField("template", tem.ID).
									Warn("Requested language translation not found")
								continue
							}
							if tem, err = g.getTemplate(cat, subID); err != nil {
								clog.WithError(err).
									Error("Error getting template from the catalog.")
								continue
							}
							if tem == nil {
								clog.WithField("template", subID).Warn("Template not found in the catalog")
								continue
							}
							if a, err = g.Generate(t, now, tem, prod, lang); err != nil {
								clog.WithError(err).
									Error("Error while generating announcement. Aborting.")
								continue
							}
							if a == nil {
								clog.WithFields(log.Fields{"template": (*prod).Name(), "lang": l}).
									Info("Template not found for language")
								continue
							}
							announcements = append(announcements, a)
						} else {
							//otherwise we need to delete announcements for this language if any
							found := g.byID(BuildID(t.ID, tem.Type, lang), currentAnnons)
							if found != nil {
								toDelete = append(toDelete, found)
							}
						}
					}
				}
			} else {
				// we don't need to regenerate so we just keep the old announcements
				clog.Info("Announcements update is not required. Appending current announcements to the list.")
				announcements = append(announcements, g.byType((*prod).Name(), currentAnnons)...)
			}
		} else {
			if needsUpdate {
				clog.Info("Announcement is no longer needed. Deleting.")
				// if not required but needs update we need to delete old announcements
				toDelete = append(toDelete, g.byType((*prod).Name(), currentAnnons)...)
			}
		}
	}
	return announcements, toDelete, err
}

func (g *gen) getTemplate(cat *Catalog, ID string) (temp *CatalogTemplate, err error) {
	if temp = g.cache[ID]; temp == nil {
		if temp, err = (*cat).Get(ID); err != nil {
			return nil, err
		}
		if temp == nil {
			return nil, nil
		}
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "annn.generator", "method": "getTemplate", "template": ID}).
				Debug("Parsing templates")
		}
		temp.Templates.Parse(fmt.Sprintf("%v announcement templates", ID))
		g.cache[ID] = temp
	}
	return temp, err
}

func (g *gen) Generate(t *train.Train, now *time.Time, temp *CatalogTemplate, p *Producer, lang language.Tag) (*Announcement, error) {
	clog := log.WithFields(log.Fields{"logger": "service.generator", "method": "Generate"})
	if t == nil {
		clog.Error("Received nil train for processing. This should never happen!")
		return nil, errors.New("Generator has received a nil train for processing.")
	}
	times, first, last := (*p).GetTime(t, now)
	var err error
	var human, tts, htm string
	var autoplay bool
	var params, ttsParams *TemplateParams
	if params, ttsParams, autoplay, err = (*p).BuildParams(t, lang); err != nil {
		clog.WithFields(log.Fields{"train": t.ID, "template": temp.ID}).
			WithError(err).Error("Error building template parameters.")
		return nil, err
	}
	if log.GetLevel() >= log.DebugLevel {
		clog.WithFields(log.Fields{"train": t.ID, "template": temp.ID, "params": params}).
			Debug("Generating announcement text.")
	}
	human, tts, htm, _ = ExecuteTemplates(&temp.Templates, *params, *ttsParams)

	announcement := &Announcement{
		ID:        BuildID(t.ID, temp.Type, lang),
		Time:      times,
		First:     first,
		Last:      last,
		Lang:      temp.Lang,
		Type:      temp.Type,
		Priority:  Priority(g.priorities[string(temp.Type)]),
		TrainID:   t.ID,
		StationID: t.StationID,
		Category:  Train,
		Autoplay:  autoplay,
		Text: &Text{HumanText: human,
			TtsText:  tts,
			HTMLText: htm,
		},
	}
	return announcement, err
}

func (g *gen) FromTemplate(templ *CatalogTemplate, startTime *time.Time, endTime *time.Time, dayInterval *time.Duration, nightInterval *time.Duration, now *time.Time, params TemplateParams, stationID int) (*Announcement, error) {
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "service.generator", "method": "FromTemplate", "template": templ.ID, "startTime": startTime, "endTime": endTime, "dayInterval": dayInterval, "nightInterval": nightInterval, "now": now}).
			Debug("Generating announcement from template.")
	}
	cat := Special
	var ttsParams TemplateParams
	if p := g.producers[templ.ID]; p != nil {
		ttsParams = (*p).MapParams(params, i18n.GetLang(templ.Lang))
		cat = Train
	}
	// generate announcement text based on the template
	templ.Templates.Parse(templ.Title)
	human, tts, html, err := ExecuteTemplates(&templ.Templates, params, ttsParams)
	if err != nil {
		log.WithFields(log.Fields{"logger": "service.generator", "method": "FromTemplate", "template": templ.ID}).
			WithError(err).Error("Error while executing templates")
		return nil, err
	}
	times, first, last := g.getTimes(startTime, endTime, dayInterval, nightInterval, now)
	var announcement *Announcement
	if len(times) > 0 {
		announcement = &Announcement{
			Time:      times,
			First:     *first,
			Last:      *last,
			Lang:      templ.Lang,
			Category:  cat,
			Type:      templ.Type,
			Priority:  Priority(g.priorities[string(templ.Type)]),
			StationID: stationID,
			Text: &Text{
				HumanText: human,
				HTMLText:  html,
				TtsText:   tts,
			},
			Autoplay: true,
		}
	}
	return announcement, err
}

func (g *gen) Preview(templ *CatalogTemplate, params TemplateParams) (*Text, error) {
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "annongen.generator", "method": "Preview", "template": templ.ID}).
			Debug("Generating announcement preview.")
	}
	var ttsParams TemplateParams
	if p := g.producers[templ.ID]; p != nil {
		ttsParams = (*p).MapParams(params, i18n.GetLang(templ.Lang))
	}
	// generate announcement text based on the template
	templ.Templates.Parse(templ.Title)
	human, tts, html, err := ExecuteTemplates(&templ.Templates, params, ttsParams)
	if err != nil {
		log.WithFields(log.Fields{"logger": "annongen.generator", "method": "Preview", "template": templ.ID}).
			WithError(err).Error("Error while executing templates")
		return nil, err
	}
	return &Text{HumanText: human, TtsText: tts, HTMLText: html}, nil
}

func (g *gen) getTimes(startTime *time.Time, endTime *time.Time, dayInterval *time.Duration, nightInterval *time.Duration, now *time.Time) ([]time.Time, *time.Time, *time.Time) {
	events := []time.Time{}
	var first, last *time.Time
	var current = new(time.Time)
	// this covers one shot announcements where endTime will be nil
	if startTime != nil && ((*startTime).After(*now) || (*startTime).Equal(*now)) {
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "service.generator", "method": "getTimes", "startTime": *startTime, "now": *now}).
				Debug("Adding new event to the list.")
		}
		events = append(events, *startTime)
		first = startTime
		last = startTime
	}
	if endTime != nil && dayInterval != nil {
		*current = *startTime
		if current.Equal(*endTime) {
			if log.GetLevel() >= log.DebugLevel {
				log.WithFields(log.Fields{"logger": "service.generator", "method": "getTimes", "startTime": *startTime, "endTime": *endTime, "current": *current, "now": *now}).
					Debug("Current is equal end. It means we do not need to add anything else. Exiting.")
			}
			return events, first, last
		}
		current = updateCurrent(current, dayInterval, nightInterval)
		// we increment with the interval until we reach endTime
		for current.Before(*endTime) || current.Equal(*endTime) {
			if log.GetLevel() >= log.DebugLevel {
				log.WithFields(log.Fields{
					"logger": "service.generator", "method": "getTimes", "startTime": *startTime, "endTime": *endTime, "current": *current, "now": *now}).
					Debug("Checking if current is in the future.")
			}
			if current.After(*now) {
				if log.GetLevel() >= log.DebugLevel {
					log.WithFields(log.Fields{"logger": "service.generator", "method": "getTimes", "startTime": *startTime, "endTime": *endTime, "current": *current, "now": *now}).
						Debug("Adding new event to the list.")
				}
				events = append(events, *current)
				last = current
			}
			current = updateCurrent(current, dayInterval, nightInterval)
		}
	}
	return events, first, last
}

func updateCurrent(current *time.Time, dayInterval *time.Duration, nightInterval *time.Duration) *time.Time {
	day := true
	if nightInterval != nil {
		day = isDay(current)
	}
	if day == true {
		if log.GetLevel() >= log.DebugLevel {
			log.WithFields(log.Fields{"logger": "service.generator", "method": "updateCurrent", "dayInterval": *dayInterval}).
				Debug("Adding day interval.")
		}
		res := (*current).Add(*dayInterval)
		return &res
	}
	if log.GetLevel() >= log.DebugLevel {
		log.WithFields(log.Fields{"logger": "service.generator", "method": "updateCurrent", "nightInterval": *dayInterval}).
			Debug("Adding night interval.")
	}
	res := (*current).Add(*nightInterval)
	return &res
}

func (g *gen) byType(name string, annons *[]*Announcement) (res []*Announcement) {
	res = []*Announcement{}
	for _, a := range *annons {
		if a.Type == name {
			res = append(res, a)
		}
	}
	return res
}

func (g *gen) byID(ID string, annons *[]*Announcement) *Announcement {
	for _, a := range *annons {
		if a.ID == ID {
			return a
		}
	}
	return nil
}

func (g *gen) addProducer(p *Producer) {
	g.producers[(*p).Name()] = p
}

func isDay(timestamp *time.Time) bool {
	if timestamp == nil {
		return true
	}
	minutes := timestamp.Hour()*60 + timestamp.Minute()
	return minutes > nightEnd && minutes < nightStart
}
