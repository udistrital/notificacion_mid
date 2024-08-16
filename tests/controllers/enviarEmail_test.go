package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/udistrital/notificacion_mid/models"
)

func TestPostSendTemplatedEmail(t *testing.T) {
	toStringPointers := func(strs []string) []*string {
		ptrs := make([]*string, len(strs))
		for i, s := range strs {
			ptrs[i] = &s
		}
		return ptrs
	}

	input := models.InputTemplatedEmail{
		Source:   "notificacion_sisgplan@udistrital.edu.co",
		Template: "SISGPLAN_PLANTILLA_F",
		Destinations: []models.DestinationTemplate{
			{
				Destination: &ses.Destination{
					ToAddresses: toStringPointers([]string{"correo@udistrital.edu.co"}),
				},
				ReplacementTemplateData: map[string]interface{}{
					"NOMBRE_UNIDAD": "VICERRECTORIA ACADEMICA",
					"NOMBRE_PLAN":   "Metas pruebas fabian ",
					"VIGENCIA":      "2023",
					"TRIMESTRE":     "T1",
				},
			},
		},
		DefaultTemplateData: map[string]interface{}{
			"NOMBRE_UNIDAD": "XXXX",
			"NOMBRE_PLAN":   "XXXX",
			"VIGENCIA":      "XXXX",
		},
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Error en el marshal: %v", err)
	}

	url := "http://localhost:8080/v1/email/enviar_templated_email/"
	if response, err := http.Post(url, "application/json", bytes.NewBuffer(body)); err == nil {
		t.Log(response.StatusCode)
		if response.StatusCode != 200 {
			t.Error("Error TestPostSendTemplatedEmail Se esperaba 200 y se obtuvo", response.StatusCode)
			t.Fail()
		} else {
			t.Log("TestPostSendTemplatedEmail Finalizado Correctamente (OK)")
		}
	} else {
		t.Error("Error TestPostSendTemplatedEmail:", err.Error())
		t.Fail()
	}
}
