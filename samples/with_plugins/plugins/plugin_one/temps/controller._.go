package temps

import (
	"time"

	"github.com/tuxounet/k2-sdk/bases"
	"github.com/tuxounet/k2-sdk/types"
)

type Controller struct {
	bases.BaseAppController
}

func NewController(component types.IAppComponent) types.IAppController {
	base := bases.NewBaseAppController(component, "temps", 100, nil, types.AccessPolicyPublic)
	return &Controller{
		base,
	}
}

// TimeResponse représente la réponse de l'API horloge
type TimeResponse struct {
	// Timestamp Unix en secondes
	Timestamp int64 `json:"timestamp"`
	// Date formatée en français
	DateFr string `json:"date_fr"`
	// Heure formatée en français
	HeureFr string `json:"heure_fr"`
	// Date et heure complète en français
	DateHeureFr string `json:"date_heure_fr"`
	// Année
	Annee int `json:"year"`
	// Mois (1-12)
	Mois int `json:"month"`
	// Jour du mois
	Jour int `json:"day"`
	// Heure (0-23)
	Heure int `json:"hours"`
	// Minutes (0-59)
	Minutes int `json:"minutes"`
	// Secondes (0-59)
	Secondes int `json:"seconds"`
	// Jour de la semaine en français
	JourSemaine string `json:"day_of_week"`
	// Nom du mois en français
	NomMois string `json:"month_name"`
}

var joursSemaine = []string{
	"Dimanche",
	"Lundi",
	"Mardi",
	"Mercredi",
	"Jeudi",
	"Vendredi",
	"Samedi",
}

var nomsMois = []string{
	"Janvier",
	"Février",
	"Mars",
	"Avril",
	"Mai",
	"Juin",
	"Juillet",
	"Août",
	"Septembre",
	"Octobre",
	"Novembre",
	"Décembre",
}

// GetCurrentTime retourne l'heure actuelle du serveur formatée
func (c *Controller) GetCurrentTime() TimeResponse {
	now := time.Now()

	jourSemaine := joursSemaine[now.Weekday()]
	nomMois := nomsMois[now.Month()-1]

	dateFr := now.Format("02/01/2006")
	heureFr := now.Format("15:04:05")
	dateHeureFr := jourSemaine + " " + now.Format("02") + " " + nomMois + " " + now.Format("2006") + " à " + heureFr

	return TimeResponse{
		Timestamp:   now.Unix(),
		DateFr:      dateFr,
		HeureFr:     heureFr,
		DateHeureFr: dateHeureFr,
		Annee:       now.Year(),
		Mois:        int(now.Month()),
		Jour:        now.Day(),
		Heure:       now.Hour(),
		Minutes:     now.Minute(),
		Secondes:    now.Second(),
		JourSemaine: jourSemaine,
		NomMois:     nomMois,
	}
}
