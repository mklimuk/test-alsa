package annon

import (
	"github.com/mklimuk/husar/db"

	log "github.com/Sirupsen/logrus"

	r "github.com/dancannon/gorethink"
)

/*
Catalog is an interface giving access to the templates catalog
*/
type Catalog interface {
	Get(ID string) (*CatalogTemplate, error)
	Save(template *CatalogTemplate) (*CatalogTemplate, error)
	SaveAll(tps *CatalogTemplates) (*CatalogTemplates, error)
	GetAll() (map[string][]*CatalogTemplate, error)
	Delete(ID string) error
}

/*
CatalogTemplates Wrapper type over a slice of CatalogTemplate
*/
type CatalogTemplates []CatalogTemplate

/*
CatalogTemplate represents a template that can be used to generate a custom announcement
*/
type CatalogTemplate struct {
	ID           string             `gorethink:"id,omitempty" json:"id" yaml:"id"`
	Type         string             `gorethink:"type" json:"type" yaml:"type"`
	Lang         string             `gorethink:"lang" json:"lang" yaml:"lang"`
	Translations *map[string]string `gorethink:"translations" json:"translations" yaml:"translations"`
	Title        string             `gorethink:"title" json:"title" yaml:"title"`
	Description  string             `gorethink:"description" json:"description" yaml:"description"`
	Categories   []string           `gorethink:"categories" json:"categories" yaml:"categories"`
	Templates    Templates          `gorethink:"templates" json:"templates" yaml:"templates"`
}

/*
NewCatalog A catalog constructor
*/
func NewCatalog(session *r.Session, db db.Database, table db.Table, templates []string) Catalog {
	clog := log.WithFields(log.Fields{"logger": "annongen.tpl", "method": "NewCatalog"})

	clog.Info("Creating templates catalog.")
	c := &cat{session: session, db: db, table: table}
	// parse templates

	clog.WithField("templates", templates).Info("Parsing templates from file.")

	tps := Parse(templates)
	// save templates to database
	if _, err := c.SaveAll(&tps); err != nil {
		clog.WithField("templates", templates).WithError(err).
			Error("Could not parse templates catalog.")

	}
	return Catalog(c)
}

type cat struct {
	session *r.Session
	db      db.Database
	table   db.Table
}

func (c *cat) Get(ID string) (*CatalogTemplate, error) {
	resp, err := r.DB(c.db).Table(c.table).Get(ID).Run(c.session)
	if err != nil {
		log.WithFields(log.Fields{"logger": "annongen.tpl", "method": "Get", "template": ID}).
			WithError(err).Error("Error getting query results.")
		return nil, err
	}
	if resp.IsNil() {
		return nil, err
	}
	a := new(CatalogTemplate)
	err = resp.One(a)
	if err != nil {
		log.WithFields(log.Fields{"logger": "annongen.tpl", "method": "Get", "template": ID}).
			WithError(err).Error("Error scanning query results.")
		return nil, err
	}
	return a, err
}

func (c *cat) GetAll() (catalog map[string][]*CatalogTemplate, err error) {
	resp, err := r.DB(c.db).Table(c.table).Run(c.session)
	if err != nil {
		log.WithFields(log.Fields{"logger": "annongen.tpl", "method": "GetAll"}).
			WithError(err).Error("Error getting query results.")
		return nil, err
	}
	var tpl []CatalogTemplate
	err = resp.All(&tpl)
	if err != nil {
		log.WithFields(log.Fields{"logger": "annongen.tpl", "method": "GetAll"}).
			WithError(err).Error("Error scanning query results.")
		return nil, err
	}
	catalog = make(map[string][]*CatalogTemplate)
	for _, t := range tpl {
		ins := new(CatalogTemplate)
		*ins = t
		for _, c := range t.Categories {
			if v, ok := catalog[c]; ok {
				v = append(v, ins)
				catalog[c] = v
			} else {
				catalog[c] = []*CatalogTemplate{ins}
			}
		}
	}
	return catalog, err
}

func (c *cat) Save(template *CatalogTemplate) (*CatalogTemplate, error) {
	resp, err := r.DB(c.db).Table(c.table).Insert(template, r.InsertOpts{
		Conflict: "replace",
	}).RunWrite(c.session)
	if err != nil {
		log.WithFields(log.Fields{"logger": "annongen.tpl", "method": "Save", "template": template.ID}).
			WithError(err).Error("Error running write query.")
		return nil, err
	}
	if resp.Inserted > 0 {
		(*template).ID = resp.GeneratedKeys[0]
	}
	return template, err
}

func (c *cat) SaveAll(tps *CatalogTemplates) (*CatalogTemplates, error) {
	res, err := r.DB(c.db).Table(c.table).Insert(*tps, r.InsertOpts{
		Conflict: "replace",
	}).RunWrite(c.session)
	if err != nil {
		log.WithFields(log.Fields{"logger": "annongen.tpl", "method": "SaveAll"}).
			WithError(err).Error("Error running write query.")
		return tps, err
	}
	i := 0
	for _, t := range *tps {
		if t.ID == "" {
			t.ID = res.GeneratedKeys[i]
			i++
		}
	}
	return tps, err
}

func (c *cat) Delete(ID string) error {
	_, err := r.DB(c.db).Table(c.table).Get(ID).Delete().RunWrite(c.session)
	if err != nil {
		log.WithFields(log.Fields{
			"logger": "annongen.catalog", "method": "Delete", "template": ID}).
			WithError(err).Error("Error occured while deleting template.")
	}
	return err
}
