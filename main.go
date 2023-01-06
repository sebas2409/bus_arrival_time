package main

import (
	"encoding/json"
	"fmt"
	"github.com/tebeka/selenium"
	"net/http"
	"strings"
)

const (
	seleniumPath = "./selenium-server-standalone-3.14.0.jar"
	port         = 8080
	chromeDriver = "./chromedriver"
	busStop      = "3042"
)

type Parada struct {
	Linea   string   `json:"linea"`
	Minutos []string `json:"minutos"`
	Destino string   `json:"destino"`
}

func main() {
	service, err := selenium.NewSeleniumService(seleniumPath, port, selenium.ChromeDriver(chromeDriver))
	if err != nil {
		panic(err)
	}
	defer service.Stop()

	wd, err := selenium.NewRemote(selenium.Capabilities{"browserName": "chrome"}, "http://localhost:8080/wd/hub")
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	// Access to the page for the information of the bus stops
	err = wd.Get(fmt.Sprintf("https://busmadrid.welbits.com/stop/%s", busStop))
	if err != nil {
		panic(err)
	}

	var filas []selenium.WebElement
	err = wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		filas, err = wd.FindElements(selenium.ByXPATH, "//table/tbody/tr")
		if err != nil {
			return false, err
		} else {
			return len(filas) > 0, nil
		}
	})

	var paradas []Parada
	for _, el := range filas {
		dato, _ := el.Text()
		rs := strings.Split(dato, "\n")
		paradas = append(paradas, Parada{
			Linea:   rs[0],
			Minutos: []string{rs[1], rs[2]},
			Destino: rs[3],
		})
	}

	//send information to client as push notification

	newJson, _ := json.Marshal(paradas)
	_, err = http.Post("https://ntfy.sh/information_for_sebas", "application/json", strings.NewReader(string(newJson)))
	if err != nil {
		panic(err)
	}
}
