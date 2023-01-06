package main

import (
	"github.com/tebeka/selenium"
)

const (
	seleniumPath = "./selenium-server-standalone-3.14.0.jar"
	port         = 8080
	chromeDriver = "./chromedriver.exe"
)

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

	err = wd.Get("https://busmadrid.welbits.com/stop/2001")
	if err != nil {
		panic(err)
	}

	var elements []selenium.WebElement
	err = wd.Wait(func(wd selenium.WebDriver) (bool, error) {
		elements, err = wd.FindElements(selenium.ByCSSSelector, "div[class='StopPageView__time___1Pt8j']")
		if err != nil {
			return false, err
		} else {
			return len(elements) > 0, nil
		}
	})

	for _, element := range elements {
		text, _ := element.Text()
		println(text)
	}
}
