package soap

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"
)

// against "unused imports"
var _ time.Time
var _ xml.Name

type PDPServiceFaultDetailType string

const (
	PDPServiceFaultDetailTypeInfo PDPServiceFaultDetailType = "Info"

	PDPServiceFaultDetailTypeWarning PDPServiceFaultDetailType = "Warning"

	PDPServiceFaultDetailTypeError PDPServiceFaultDetailType = "Error"
)

type PobierzRozkladPlanowyNastepneZamkniecieDlaWszystkichStacji struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PobierzRozkladPlanowyNastepneZamkniecieDlaWszystkichStacji"`
}

type PobierzRozkladPlanowyNastepneZamkniecieDlaWszystkichStacjiResponse struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PobierzRozkladPlanowyNastepneZamkniecieDlaWszystkichStacjiResponse"`

	PobierzRozkladPlanowyNastepneZamkniecieDlaWszystkichStacjiResult struct {
		RozkladPlanowy    *RozkladPlanowyTyp           `xml:"RozkladPlanowy,omitempty"`
		Utrudnienia       *ArrayOfUtrudnienieTyp       `xml:"Utrudnienia,omitempty"`
		Stacje            *ArrayOfStacjaTyp            `xml:"Stacje,omitempty"`
		Tlumaczenia       *ArrayOfTlumaczenieTyp       `xml:"Tlumaczenia,omitempty"`
		Uslugi            *ArrayOfUslugaTyp            `xml:"Uslugi,omitempty"`
		RodzajeWagonow    *ArrayOfRodzajWagonuTyp      `xml:"RodzajeWagonow,omitempty"`
		RodzajePowiazan   *ArrayOfRodzajPowiazaniaTyp  `xml:"RodzajePowiazan,omitempty"`
		Przewoznicy       *ArrayOfPrzewoznikTyp        `xml:"Przewoznicy,omitempty"`
		KategorieHandlowe *ArrayOfKategoriaHandlowaTyp `xml:"KategorieHandlowe,omitempty"`
		WersjeJezykowe    *ArrayOfWersjaJezykowaTyp    `xml:"WersjeJezykowe,omitempty"`
	} `xml:"PobierzRozkladPlanowyNastepneZamkniecieDlaWszystkichStacjiResult,omitempty"`
}

type PobierzRozkladPlanowyTygodniowyDlaWszystkichStacji struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PobierzRozkladPlanowyTygodniowyDlaWszystkichStacji"`
}

type PobierzRozkladPlanowyTygodniowyDlaWszystkichStacjiResponse struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PobierzRozkladPlanowyTygodniowyDlaWszystkichStacjiResponse"`

	PobierzRozkladPlanowyTygodniowyDlaWszystkichStacjiResult struct {
		RozkladPlanowy    *RozkladPlanowyTyp           `xml:"RozkladPlanowy,omitempty"`
		Utrudnienia       *ArrayOfUtrudnienieTyp       `xml:"Utrudnienia,omitempty"`
		Stacje            *ArrayOfStacjaTyp            `xml:"Stacje,omitempty"`
		Tlumaczenia       *ArrayOfTlumaczenieTyp       `xml:"Tlumaczenia,omitempty"`
		Uslugi            *ArrayOfUslugaTyp            `xml:"Uslugi,omitempty"`
		RodzajeWagonow    *ArrayOfRodzajWagonuTyp      `xml:"RodzajeWagonow,omitempty"`
		RodzajePowiazan   *ArrayOfRodzajPowiazaniaTyp  `xml:"RodzajePowiazan,omitempty"`
		Przewoznicy       *ArrayOfPrzewoznikTyp        `xml:"Przewoznicy,omitempty"`
		KategorieHandlowe *ArrayOfKategoriaHandlowaTyp `xml:"KategorieHandlowe,omitempty"`
		WersjeJezykowe    *ArrayOfWersjaJezykowaTyp    `xml:"WersjeJezykowe,omitempty"`
	} `xml:"PobierzRozkladPlanowyTygodniowyDlaWszystkichStacjiResult,omitempty"`
}

type PobierzRozkladPlanowyTygodniowyDlaStacji struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PobierzRozkladPlanowyTygodniowyDlaStacji"`

	IdStacji int32 `xml:"idStacji,omitempty"`
}

type PobierzRozkladPlanowyTygodniowyDlaStacjiResponse struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PobierzRozkladPlanowyTygodniowyDlaStacjiResponse"`

	PobierzRozkladPlanowyTygodniowyDlaStacjiResult struct {
		RozkladPlanowy    *RozkladPlanowyTyp           `xml:"RozkladPlanowy,omitempty"`
		Utrudnienia       *ArrayOfUtrudnienieTyp       `xml:"Utrudnienia,omitempty"`
		Stacje            *ArrayOfStacjaTyp            `xml:"Stacje,omitempty"`
		Tlumaczenia       *ArrayOfTlumaczenieTyp       `xml:"Tlumaczenia,omitempty"`
		Uslugi            *ArrayOfUslugaTyp            `xml:"Uslugi,omitempty"`
		RodzajeWagonow    *ArrayOfRodzajWagonuTyp      `xml:"RodzajeWagonow,omitempty"`
		RodzajePowiazan   *ArrayOfRodzajPowiazaniaTyp  `xml:"RodzajePowiazan,omitempty"`
		Przewoznicy       *ArrayOfPrzewoznikTyp        `xml:"Przewoznicy,omitempty"`
		KategorieHandlowe *ArrayOfKategoriaHandlowaTyp `xml:"KategorieHandlowe,omitempty"`
		WersjeJezykowe    *ArrayOfWersjaJezykowaTyp    `xml:"WersjeJezykowe,omitempty"`
	} `xml:"PobierzRozkladPlanowyTygodniowyDlaStacjiResult,omitempty"`
}

type PobierzRozkladPlanowyTygodniowyDlaListyStacji struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PobierzRozkladPlanowyTygodniowyDlaListyStacji"`

	ListaStacji *ArrayOfInt `xml:"listaStacji,omitempty"`
}

type PobierzRozkladPlanowyTygodniowyDlaListyStacjiResponse struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PobierzRozkladPlanowyTygodniowyDlaListyStacjiResponse"`

	PobierzRozkladPlanowyTygodniowyDlaListyStacjiResult struct {
		RozkladPlanowy    *RozkladPlanowyTyp           `xml:"RozkladPlanowy,omitempty"`
		Utrudnienia       *ArrayOfUtrudnienieTyp       `xml:"Utrudnienia,omitempty"`
		Stacje            *ArrayOfStacjaTyp            `xml:"Stacje,omitempty"`
		Tlumaczenia       *ArrayOfTlumaczenieTyp       `xml:"Tlumaczenia,omitempty"`
		Uslugi            *ArrayOfUslugaTyp            `xml:"Uslugi,omitempty"`
		RodzajeWagonow    *ArrayOfRodzajWagonuTyp      `xml:"RodzajeWagonow,omitempty"`
		RodzajePowiazan   *ArrayOfRodzajPowiazaniaTyp  `xml:"RodzajePowiazan,omitempty"`
		Przewoznicy       *ArrayOfPrzewoznikTyp        `xml:"Przewoznicy,omitempty"`
		KategorieHandlowe *ArrayOfKategoriaHandlowaTyp `xml:"KategorieHandlowe,omitempty"`
		WersjeJezykowe    *ArrayOfWersjaJezykowaTyp    `xml:"WersjeJezykowe,omitempty"`
	} `xml:"PobierzRozkladPlanowyTygodniowyDlaListyStacjiResult,omitempty"`
}

type PobierzRozkladRzeczywistyDlaWszystkichStacji struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PobierzRozkladRzeczywistyDlaWszystkichStacji"`
}

type PobierzRozkladRzeczywistyDlaWszystkichStacjiResponse struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PobierzRozkladRzeczywistyDlaWszystkichStacjiResponse"`

	PobierzRozkladRzeczywistyDlaWszystkichStacjiResult struct {
		RozkladRzeczywisty *RozkladRzeczywistyTyp `xml:"RozkladRzeczywisty,omitempty"`
		Stacje             *ArrayOfStacjaTyp      `xml:"Stacje,omitempty"`
		Tlumaczenia        *ArrayOfTlumaczenieTyp `xml:"Tlumaczenia,omitempty"`
	} `xml:"PobierzRozkladRzeczywistyDlaWszystkichStacjiResult,omitempty"`
}

type PobierzRozkladRzeczywistyDlaStacji struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PobierzRozkladRzeczywistyDlaStacji"`

	IdStacji int32 `xml:"idStacji,omitempty"`
}

type PobierzRozkladRzeczywistyDlaStacjiResponse struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PobierzRozkladRzeczywistyDlaStacjiResponse"`

	PobierzRozkladRzeczywistyDlaStacjiResult struct {
		RozkladRzeczywisty *RozkladRzeczywistyTyp `xml:"RozkladRzeczywisty,omitempty"`
		Stacje             *ArrayOfStacjaTyp      `xml:"Stacje,omitempty"`
		Tlumaczenia        *ArrayOfTlumaczenieTyp `xml:"Tlumaczenia,omitempty"`
	} `xml:"PobierzRozkladRzeczywistyDlaStacjiResult,omitempty"`
}

type PobierzRozkladRzeczywistyDlaListyStacji struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PobierzRozkladRzeczywistyDlaListyStacji"`

	ListaStacji *ArrayOfInt `xml:"listaStacji,omitempty"`
}

type PobierzRozkladRzeczywistyDlaListyStacjiResponse struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PobierzRozkladRzeczywistyDlaListyStacjiResponse"`

	PobierzRozkladRzeczywistyDlaListyStacjiResult struct {
		RozkladRzeczywisty *RozkladRzeczywistyTyp `xml:"RozkladRzeczywisty,omitempty"`
		Stacje             *ArrayOfStacjaTyp      `xml:"Stacje,omitempty"`
		Tlumaczenia        *ArrayOfTlumaczenieTyp `xml:"Tlumaczenia,omitempty"`
	} `xml:"PobierzRozkladRzeczywistyDlaListyStacjiResult,omitempty"`
}

type RozkladPlanowyTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 RozkladPlanowyTyp"`

	CzasGeneracji time.Time        `xml:"CzasGeneracji,omitempty"`
	ZakresOd      time.Time        `xml:"ZakresOd,omitempty"`
	ZakresDo      time.Time        `xml:"ZakresDo,omitempty"`
	Trasy         *ArrayOfTrasaTyp `xml:"Trasy,omitempty"`
}

type ArrayOfTrasaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfTrasaTyp"`

	Trasa []*TrasaTyp `xml:"Trasa,omitempty"`
}

type TrasaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 TrasaTyp"`

	RozkladID                    int16                    `xml:"RozkladID,omitempty"`
	ZamowienieSKRJID             int32                    `xml:"ZamowienieSKRJID,omitempty"`
	NumerKrajowy                 string                   `xml:"NumerKrajowy,omitempty"`
	NumerMiedzynarodowyWjazdowy  string                   `xml:"NumerMiedzynarodowyWjazdowy,omitempty"`
	NumerMiedzynarodowyWyjazdowy string                   `xml:"NumerMiedzynarodowyWyjazdowy,omitempty"`
	Nazwa                        string                   `xml:"Nazwa,omitempty"`
	KategorieHandlowe            string                   `xml:"KategorieHandlowe,omitempty"`
	Przewoznik                   string                   `xml:"Przewoznik,omitempty"`
	RelacjaPoczatkowaID          int32                    `xml:"RelacjaPoczatkowaID,omitempty"`
	RelacjaKoncowaID             int32                    `xml:"RelacjaKoncowaID,omitempty"`
	Powiazana                    bool                     `xml:"Powiazana,omitempty"`
	KalendarzKursowania          *ArrayOfDate             `xml:"KalendarzKursowania,omitempty"`
	StacjePlanowe                *ArrayOfStacjaPlanowaTyp `xml:"StacjePlanowe,omitempty"`
}

type ArrayOfDate struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfDate"`

	Data []time.Time `xml:"Data,omitempty"`
}

type ArrayOfStacjaPlanowaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfStacjaPlanowaTyp"`

	StacjaPlanowa []*StacjaPlanowaTyp `xml:"StacjaPlanowa,omitempty"`
}

type StacjaPlanowaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 StacjaPlanowaTyp"`

	StacjaID           int32                               `xml:"StacjaID,omitempty"`
	NumerNaTrasie      int16                               `xml:"NumerNaTrasie,omitempty"`
	KategoriaWjazdowa  string                              `xml:"KategoriaWjazdowa,omitempty"`
	NumerWjazdowy      string                              `xml:"NumerWjazdowy,omitempty"`
	PeronWjazdowy      string                              `xml:"PeronWjazdowy,omitempty"`
	TorWjazdowy        string                              `xml:"TorWjazdowy,omitempty"`
	DzienPrzyjazdu     int16                               `xml:"DzienPrzyjazdu,omitempty"`
	CzasPrzyjazdu      time.Time                           `xml:"CzasPrzyjazdu,omitempty"`
	DzienOdjazdu       int16                               `xml:"DzienOdjazdu,omitempty"`
	CzasOdjazdu        time.Time                           `xml:"CzasOdjazdu,omitempty"`
	KategoriaWyjazdowa string                              `xml:"KategoriaWyjazdowa,omitempty"`
	NumerWyjazdowy     string                              `xml:"NumerWyjazdowy,omitempty"`
	PeronWyjazdowy     string                              `xml:"PeronWyjazdowy,omitempty"`
	TorWyjazdowy       string                              `xml:"TorWyjazdowy,omitempty"`
	Sklady             *ArrayOfSkladTyp                    `xml:"Sklady,omitempty"`
	Powiazania         *ArrayOfPowiazanieTyp               `xml:"Powiazania,omitempty"`
	Utrudnienia        *ArrayOfUtrudnienieStacjaPlanowaTyp `xml:"Utrudnienia,omitempty"`
}

type ArrayOfSkladTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfSkladTyp"`

	Sklad []*SkladTyp `xml:"Sklad,omitempty"`
}

type SkladTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 SkladTyp"`

	ID                            int64            `xml:"ID,omitempty"`
	StacjaPoczatkowaNumerNaTrasie int16            `xml:"StacjaPoczatkowaNumerNaTrasie,omitempty"`
	Dlugosc                       int16            `xml:"Dlugosc,omitempty"`
	KalendarzKursowania           *ArrayOfDate     `xml:"KalendarzKursowania,omitempty"`
	Uslugi                        *ArrayOfString   `xml:"Uslugi,omitempty"`
	Wagony                        *ArrayOfWagonTyp `xml:"Wagony,omitempty"`
}

type ArrayOfString struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfString"`

	Kod []string `xml:"Kod,omitempty"`
}

type ArrayOfWagonTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfWagonTyp"`

	Wagon []*WagonTyp `xml:"Wagon,omitempty"`
}

type WagonTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 WagonTyp"`

	ID                  int64          `xml:"ID,omitempty"`
	NumerKolejny        int16          `xml:"NumerKolejny,omitempty"`
	NumerHandlowy       string         `xml:"NumerHandlowy,omitempty"`
	RelacjaPoczatkowaID int32          `xml:"RelacjaPoczatkowaID,omitempty"`
	RelacjaKoncowaID    int32          `xml:"RelacjaKoncowaID,omitempty"`
	RodzajWagonuKod     string         `xml:"RodzajWagonuKod,omitempty"`
	Uslugi              *ArrayOfString `xml:"Uslugi,omitempty"`
}

type ArrayOfPowiazanieTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfPowiazanieTyp"`

	Powiazanie []*PowiazanieTyp `xml:"Powiazanie,omitempty"`
}

type PowiazanieTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PowiazanieTyp"`

	ID                  int64        `xml:"ID,omitempty"`
	ZamowienieSKRJID    int32        `xml:"ZamowienieSKRJID,omitempty"`
	StacjaID            int32        `xml:"StacjaID,omitempty"`
	StacjaNumerNaTrasie int16        `xml:"StacjaNumerNaTrasie,omitempty"`
	RodzajPowiazaniaKod string       `xml:"RodzajPowiazaniaKod,omitempty"`
	KalendarzPowiazania *ArrayOfDate `xml:"KalendarzPowiazania,omitempty"`
}

type ArrayOfUtrudnienieStacjaPlanowaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfUtrudnienieStacjaPlanowaTyp"`

	UtrudnienieStacjaPlanowa []*UtrudnienieStacjaPlanowaTyp `xml:"UtrudnienieStacjaPlanowa,omitempty"`
}

type UtrudnienieStacjaPlanowaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 UtrudnienieStacjaPlanowaTyp"`

	UtrudnienieID       int64        `xml:"UtrudnienieID,omitempty"`
	KalendarzKursowania *ArrayOfDate `xml:"KalendarzKursowania,omitempty"`
}

type ArrayOfUtrudnienieTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfUtrudnienieTyp"`

	Utrudnienie []*UtrudnienieTyp `xml:"Utrudnienie,omitempty"`
}

type UtrudnienieTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 UtrudnienieTyp"`

	ID            int64  `xml:"ID,omitempty"`
	Tytul         string `xml:"Tytul,omitempty"`
	TytulGUID     *Guid  `xml:"TytulGUID,omitempty"`
	Komunikat     string `xml:"Komunikat,omitempty"`
	KomunikatGUID *Guid  `xml:"KomunikatGUID,omitempty"`
}

type ArrayOfStacjaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfStacjaTyp"`

	Stacja []*StacjaTyp `xml:"Stacja,omitempty"`
}

type StacjaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 StacjaTyp"`

	ID                int32  `xml:"ID,omitempty"`
	Nazwa             string `xml:"Nazwa,omitempty"`
	NazwaGUID         *Guid  `xml:"NazwaGUID,omitempty"`
	NazwaSkrocona     string `xml:"NazwaSkrocona,omitempty"`
	NazwaSkroconaGUID *Guid  `xml:"NazwaSkroconaGUID,omitempty"`
}

type ArrayOfTlumaczenieTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfTlumaczenieTyp"`

	Tlumaczenie []*TlumaczenieTyp `xml:"Tlumaczenie,omitempty"`
}

type TlumaczenieTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 TlumaczenieTyp"`

	ID        int64  `xml:"ID,omitempty"`
	Kod       string `xml:"Kod,omitempty"`
	TrescGUID *Guid  `xml:"TrescGUID,omitempty"`
	Tresc     string `xml:"Tresc,omitempty"`
}

type ArrayOfUslugaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfUslugaTyp"`

	Usluga []*UslugaTyp `xml:"Usluga,omitempty"`
}

type UslugaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 UslugaTyp"`

	Kod      string `xml:"Kod,omitempty"`
	Opis     string `xml:"Opis,omitempty"`
	OpisGUID *Guid  `xml:"OpisGUID,omitempty"`
}

type ArrayOfRodzajWagonuTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfRodzajWagonuTyp"`

	RodzajWagonu []*RodzajWagonuTyp `xml:"RodzajWagonu,omitempty"`
}

type RodzajWagonuTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 RodzajWagonuTyp"`

	Kod       string `xml:"Kod,omitempty"`
	Nazwa     string `xml:"Nazwa,omitempty"`
	NazwaGUID *Guid  `xml:"NazwaGUID,omitempty"`
}

type ArrayOfRodzajPowiazaniaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfRodzajPowiazaniaTyp"`

	RodzajPowiazania []*RodzajPowiazaniaTyp `xml:"RodzajPowiazania,omitempty"`
}

type RodzajPowiazaniaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 RodzajPowiazaniaTyp"`

	Kod       string `xml:"Kod,omitempty"`
	Nazwa     string `xml:"Nazwa,omitempty"`
	NazwaGUID *Guid  `xml:"NazwaGUID,omitempty"`
}

type ArrayOfPrzewoznikTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfPrzewoznikTyp"`

	Przewoznik []*PrzewoznikTyp `xml:"Przewoznik,omitempty"`
}

type PrzewoznikTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PrzewoznikTyp"`

	Skrot string `xml:"Skrot,omitempty"`
	Nazwa string `xml:"Nazwa,omitempty"`
}

type ArrayOfKategoriaHandlowaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfKategoriaHandlowaTyp"`

	KategoriaHandlowa []*KategoriaHandlowaTyp `xml:"KategoriaHandlowa,omitempty"`
}

type KategoriaHandlowaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 KategoriaHandlowaTyp"`

	Symbol          string `xml:"Symbol,omitempty"`
	PrzewoznikSkrot string `xml:"PrzewoznikSkrot,omitempty"`
	Nazwa           string `xml:"Nazwa,omitempty"`
	NazwaGUID       *Guid  `xml:"NazwaGUID,omitempty"`
}

type ArrayOfWersjaJezykowaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfWersjaJezykowaTyp"`

	WersjaJezykowa []*WersjaJezykowaTyp `xml:"WersjaJezykowa,omitempty"`
}

type WersjaJezykowaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 WersjaJezykowaTyp"`

	Kod       string `xml:"Kod,omitempty"`
	Nazwa     string `xml:"Nazwa,omitempty"`
	NazwaGUID *Guid  `xml:"NazwaGUID,omitempty"`
}

type PDPServiceFault struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PDPServiceFault"`

	Details *ArrayOfPDPServiceFaultDetail `xml:"Details,omitempty"`
}

type ArrayOfPDPServiceFaultDetail struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfPDPServiceFaultDetail"`

	PDPServiceFaultDetail []*PDPServiceFaultDetail `xml:"PDPServiceFaultDetail,omitempty"`
}

type PDPServiceFaultDetail struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 PDPServiceFaultDetail"`

	Code    string                     `xml:"Code,omitempty"`
	Type    *PDPServiceFaultDetailType `xml:"Type,omitempty"`
	Message string                     `xml:"Message,omitempty"`
}

type ArrayOfInt struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfInt"`

	Int []int32 `xml:"int,omitempty"`
}

type RozkladRzeczywistyTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 RozkladRzeczywistyTyp"`

	CzasGeneracji time.Time                `xml:"CzasGeneracji,omitempty"`
	Trasy         *ArrayOfTrasaWykonanaTyp `xml:"Trasy,omitempty"`
}

type ArrayOfTrasaWykonanaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfTrasaWykonanaTyp"`

	Trasa []*TrasaWykonanaTyp `xml:"Trasa,omitempty"`
}

type TrasaWykonanaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 TrasaWykonanaTyp"`

	RozkladID        int16                     `xml:"RozkladID,omitempty"`
	ZamowienieSKRJID int32                     `xml:"ZamowienieSKRJID,omitempty"`
	DataKursowania   time.Time                 `xml:"DataKursowania,omitempty"`
	StacjeWykonane   *ArrayOfStacjaWykonanaTyp `xml:"StacjeWykonane,omitempty"`
}

type ArrayOfStacjaWykonanaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 ArrayOfStacjaWykonanaTyp"`

	StacjaWykonana []*StacjaWykonanaTyp `xml:"StacjaWykonana,omitempty"`
}

type StacjaWykonanaTyp struct {
	XMLName xml.Name `xml:"http://sdip.plk-sa.pl/v1.1 StacjaWykonanaTyp"`

	StacjaID      int32     `xml:"StacjaID,omitempty"`
	NumerNaTrasie int16     `xml:"NumerNaTrasie,omitempty"`
	Przyjazd      time.Time `xml:"Przyjazd,omitempty"`
	Odjazd        time.Time `xml:"Odjazd,omitempty"`
	Zatwierdzenie bool      `xml:"Zatwierdzenie,omitempty"`
	Odwolanie     bool      `xml:"Odwolanie,omitempty"`
}

type Guid string

const ()

type IRozkladJazdy struct {
	client *SOAPClient
}

func NewIRozkladJazdy(url string, tls *tls.Config, auth *BasicAuth) *IRozkladJazdy {
	client := NewSOAPClient(url, tls, auth)
	return &IRozkladJazdy{
		client: client,
	}
}

// Error can be either of the following types:
//
//   - PDPServiceFaultFault

func (service *IRozkladJazdy) PobierzRozkladPlanowyNastepneZamkniecieDlaWszystkichStacji(request *PobierzRozkladPlanowyNastepneZamkniecieDlaWszystkichStacji) (*PobierzRozkladPlanowyNastepneZamkniecieDlaWszystkichStacjiResponse, error) {
	response := new(PobierzRozkladPlanowyNastepneZamkniecieDlaWszystkichStacjiResponse)
	err := service.client.Call("http://sdip.plk-sa.pl/v1.1/IRozkladJazdy/PobierzRozkladPlanowyNastepneZamkniecieDlaWszystkichStacji", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - PDPServiceFaultFault

func (service *IRozkladJazdy) PobierzRozkladPlanowyTygodniowyDlaWszystkichStacji(request *PobierzRozkladPlanowyTygodniowyDlaWszystkichStacji) (*PobierzRozkladPlanowyTygodniowyDlaWszystkichStacjiResponse, error) {
	response := new(PobierzRozkladPlanowyTygodniowyDlaWszystkichStacjiResponse)
	err := service.client.Call("http://sdip.plk-sa.pl/v1.1/IRozkladJazdy/PobierzRozkladPlanowyTygodniowyDlaWszystkichStacji", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - PDPServiceFaultFault

func (service *IRozkladJazdy) PobierzRozkladPlanowyTygodniowyDlaStacji(request *PobierzRozkladPlanowyTygodniowyDlaStacji) (*PobierzRozkladPlanowyTygodniowyDlaStacjiResponse, error) {
	response := new(PobierzRozkladPlanowyTygodniowyDlaStacjiResponse)
	err := service.client.Call("http://sdip.plk-sa.pl/v1.1/IRozkladJazdy/PobierzRozkladPlanowyTygodniowyDlaStacji", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - PDPServiceFaultFault

func (service *IRozkladJazdy) PobierzRozkladPlanowyTygodniowyDlaListyStacji(request *PobierzRozkladPlanowyTygodniowyDlaListyStacji) (*PobierzRozkladPlanowyTygodniowyDlaListyStacjiResponse, error) {
	response := new(PobierzRozkladPlanowyTygodniowyDlaListyStacjiResponse)
	err := service.client.Call("http://sdip.plk-sa.pl/v1.1/IRozkladJazdy/PobierzRozkladPlanowyTygodniowyDlaListyStacji", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - PDPServiceFaultFault

func (service *IRozkladJazdy) PobierzRozkladRzeczywistyDlaWszystkichStacji(request *PobierzRozkladRzeczywistyDlaWszystkichStacji) (*PobierzRozkladRzeczywistyDlaWszystkichStacjiResponse, error) {
	response := new(PobierzRozkladRzeczywistyDlaWszystkichStacjiResponse)
	err := service.client.Call("http://sdip.plk-sa.pl/v1.1/IRozkladJazdy/PobierzRozkladRzeczywistyDlaWszystkichStacji", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - PDPServiceFaultFault

func (service *IRozkladJazdy) PobierzRozkladRzeczywistyDlaStacji(request *PobierzRozkladRzeczywistyDlaStacji) (*PobierzRozkladRzeczywistyDlaStacjiResponse, error) {
	response := new(PobierzRozkladRzeczywistyDlaStacjiResponse)
	err := service.client.Call("http://sdip.plk-sa.pl/v1.1/IRozkladJazdy/PobierzRozkladRzeczywistyDlaStacji", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

// Error can be either of the following types:
//
//   - PDPServiceFaultFault

func (service *IRozkladJazdy) PobierzRozkladRzeczywistyDlaListyStacji(request *PobierzRozkladRzeczywistyDlaListyStacji) (*PobierzRozkladRzeczywistyDlaListyStacjiResponse, error) {
	response := new(PobierzRozkladRzeczywistyDlaListyStacjiResponse)
	err := service.client.Call("http://sdip.plk-sa.pl/v1.1/IRozkladJazdy/PobierzRozkladRzeczywistyDlaListyStacji", request, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

var timeout = time.Duration(30 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`

	Body SOAPBody
}

type SOAPHeader struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Header"`

	Header interface{}
}

type SOAPBody struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`

	Fault   *SOAPFault  `xml:",omitempty"`
	Content interface{} `xml:",omitempty"`
}

type SOAPFault struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Fault"`

	Code   string `xml:"faultcode,omitempty"`
	String string `xml:"faultstring,omitempty"`
	Actor  string `xml:"faultactor,omitempty"`
	Detail string `xml:"detail,omitempty"`
}

type BasicAuth struct {
	Login    string
	Password string
}

type SOAPClient struct {
	url  string
	tls  *tls.Config
	auth *BasicAuth
}

func (b *SOAPBody) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if b.Content == nil {
		return xml.UnmarshalError("Content must be a pointer to a struct")
	}

	var (
		token    xml.Token
		err      error
		consumed bool
	)

Loop:
	for {
		if token, err = d.Token(); err != nil {
			return err
		}

		if token == nil {
			break
		}

		switch se := token.(type) {
		case xml.StartElement:
			if consumed {
				return xml.UnmarshalError("Found multiple elements inside SOAP body; not wrapped-document/literal WS-I compliant")
			} else if se.Name.Space == "http://schemas.xmlsoap.org/soap/envelope/" && se.Name.Local == "Fault" {
				b.Fault = &SOAPFault{}
				b.Content = nil

				err = d.DecodeElement(b.Fault, &se)
				if err != nil {
					return err
				}

				consumed = true
			} else {
				if err = d.DecodeElement(b.Content, &se); err != nil {
					return err
				}

				consumed = true
			}
		case xml.EndElement:
			break Loop
		}
	}

	return nil
}

func (f *SOAPFault) Error() string {
	return f.String
}

func NewSOAPClient(url string, tls *tls.Config, auth *BasicAuth) *SOAPClient {
	return &SOAPClient{
		url:  url,
		tls:  tls,
		auth: auth,
	}
}

func (s *SOAPClient) Call(soapAction string, request, response interface{}) error {
	envelope := SOAPEnvelope{
	//Header:        SoapHeader{},
	}

	envelope.Body.Content = request
	buffer := new(bytes.Buffer)

	encoder := xml.NewEncoder(buffer)
	//encoder.Indent("  ", "    ")

	if err := encoder.Encode(envelope); err != nil {
		return err
	}

	if err := encoder.Flush(); err != nil {
		return err
	}

	log.Println(buffer.String())

	req, err := http.NewRequest("POST", s.url, buffer)
	if err != nil {
		return err
	}
	if s.auth != nil {
		req.SetBasicAuth(s.auth.Login, s.auth.Password)
	}

	req.Header.Add("Content-Type", "text/xml; charset=\"utf-8\"")
	if soapAction != "" {
		req.Header.Add("SOAPAction", soapAction)
	}

	req.Header.Set("User-Agent", "gowsdl/0.1")
	req.Close = true

	tr := &http.Transport{
		TLSClientConfig: s.tls,
		Dial:            dialTimeout,
	}

	client := &http.Client{Transport: tr}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	rawbody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if len(rawbody) == 0 {
		log.Println("empty response")
		return nil
	}

	log.Println(string(rawbody))
	respEnvelope := new(SOAPEnvelope)
	respEnvelope.Body = SOAPBody{Content: response}
	err = xml.Unmarshal(rawbody, respEnvelope)
	if err != nil {
		return err
	}

	fault := respEnvelope.Body.Fault
	if fault != nil {
		return fault
	}

	return nil
}
