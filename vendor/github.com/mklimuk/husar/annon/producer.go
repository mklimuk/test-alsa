package annon

import (
	"time"

	"golang.org/x/text/language"

	"github.com/mklimuk/husar/train"
)

/*Producers is a wrapper type for a hash map of producers*/
type Producers map[string]*Producer

/*Producer generates template parameters for train announcements*/
type Producer interface {
	Required(t *train.Train, old *train.Train, reference *time.Time) (isRequired bool, needsUpdate bool, err error)
	GetTime(t *train.Train, now *time.Time) (events []time.Time, first time.Time, last time.Time)
	MapParams(in TemplateParams, lang language.Tag) TemplateParams
	BuildParams(t *train.Train, lang language.Tag) (params *TemplateParams, ttsParams *TemplateParams, autoplay bool, err error)
	Name() string
}

/*producer is a base producer implementation with a couple of useful methods*/
type producer struct {
	name string
}

// TODO change it to templateID or something similar
/*Name returns producer's type*/
func (a *producer) Name() string {
	return a.name
}

/*MapParams default implementation simply does nothing*/
func (a *producer) MapParams(params TemplateParams, lang language.Tag) TemplateParams {
	return params
}

var producers Producers = make(map[string]*Producer)

/*GetAll returns all producers registered during init*/
func GetAll() Producers {
	return producers
}
