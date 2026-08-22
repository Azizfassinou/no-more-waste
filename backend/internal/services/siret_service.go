package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type SiretInfo struct {
	IsValid        bool   `json:"is_valid"`
	CompanyName    string `json:"company_name"`
	CompanyAddress string `json:"company_address"`
	EtatAdmin      string `json:"etat_admin"`
	ErrorMessage   string `json:"error_message,omitempty"`
}

type gouvApiResponse struct {
	Results []struct {
		Siret           string `json:"siret"`
		NomComplet      string `json:"nom_complet"`
		NomRaisonSocial string `json:"nom_raison_sociale"`
		EtatAdmin       string `json:"etat_administratif"`
		Siege           struct {
			Adresse string `json:"adresse"`
		} `json:"siege"`
	} `json:"results"`
}

func ValidateSiret(siret string) (*SiretInfo, error) {
	cleaned := strings.ReplaceAll(strings.ReplaceAll(siret, " ", ""), "-", "")

	matched, _ := regexp.MatchString(`^\d{14}$`, cleaned)
	if !matched {
		return &SiretInfo{
			IsValid:      false,
			ErrorMessage: "Le numéro SIRET doit contenir exactement 14 chiffres numériques.",
		}, nil
	}

	url := fmt.Sprintf("https://recherche-entreprises.api.gouv.fr/search?q=%s", cleaned)

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return &SiretInfo{
			IsValid:        true,
			CompanyName:    "",
			CompanyAddress: "",
			ErrorMessage:   "Validation réseau indisponible (Mode hors-ligne actif).",
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &SiretInfo{
			IsValid:      false,
			ErrorMessage: fmt.Sprintf("L'API de l'État a répondu avec le statut HTTP %d.", resp.StatusCode),
		}, nil
	}

	var data gouvApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return &SiretInfo{
			IsValid:      true,
			ErrorMessage: "Erreur de lecture de la réponse de l'État.",
		}, nil
	}

	if len(data.Results) == 0 {
		return &SiretInfo{
			IsValid:      false,
			ErrorMessage: "Numéro SIRET introuvable dans le répertoire national SIRENE.",
		}, nil
	}

	res := data.Results[0]

	if res.EtatAdmin == "F" {
		return &SiretInfo{
			IsValid:      false,
			EtatAdmin:    "Fermée",
			ErrorMessage: "Cette entreprise est répertoriée comme FERMÉE / LIQUIDÉE au registre SIRENE.",
		}, nil
	}

	companyName := res.NomComplet
	if companyName == "" {
		companyName = res.NomRaisonSocial
	}

	return &SiretInfo{
		IsValid:        true,
		CompanyName:    companyName,
		CompanyAddress: res.Siege.Adresse,
		EtatAdmin:      "Actif",
	}, nil
}