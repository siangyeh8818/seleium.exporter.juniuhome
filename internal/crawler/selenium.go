package crawler

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	tool "github.com/siangyeh8818/seleium.exporter.juniuhome/internal/tool"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
)

func CallSelium() {
	log.Println("---------------------CallSelium()---------------------")
	for {
		remaining_device := RunSelium()
		log.Println("---------------------write remaining_device tp out.csv---------------------")
		tool.WriteWithIoutil("output.csv", remaining_device)
		internal_time, _ := time.ParseDuration(os.Getenv("SELEIUM_INTERNAL_TIME"))
		//internal_time, _ := strconv.Atoi(os.Getenv("SELEIUM_INTERNAL_TIME"))
		time.Sleep(time.Duration(internal_time))
	}
}

func RunSelium() string {

	const (
		seleniumPath = `/usr/local/bin/chromedriver`
		//port         = 9515
	)

	port, err := PickUnusedPort()

	fmt.Println("port", port)

	//如果seleniumServer沒有啟動，就啟動一個seleniumServer所需要的引數，可以為空，示例請參見https://github.com/tebeka/selenium/blob/master/example_test.go
	opts := []selenium.ServiceOption{}
	//opts := []selenium.ServiceOption{
	//    selenium.StartFrameBuffer(),           // Start an X frame buffer for the browser to run in.
	//    selenium.GeckoDriver(geckoDriverPath), // Specify the path to GeckoDriver in order to use Firefox.
	//}

	//selenium.SetDebug(true)
	service, err := selenium.NewChromeDriverService(seleniumPath, port, opts...)
	if nil != err {
		fmt.Println("start a chromedriver service falid", err.Error())
		//return
	}
	//注意這裡，server關閉之後，chrome視窗也會關閉
	defer service.Stop()

	//連結本地的瀏覽器 chrome
	caps := selenium.Capabilities{
		"browserName": "chrome",
	}

	//禁止圖片載入，加快渲染速度
	imagCaps := map[string]interface{}{
		"profile.managed_default_content_settings.images": 2,
	}
	chromeCaps := chrome.Capabilities{
		Prefs: imagCaps,
		Path:  "",
		Args: []string{
			"--headless", // 設定Chrome無頭模式，在linux下執行，需要設定這個引數，否則會報錯
			"--no-sandbox",
			"--disable-gpu",
			"blink-settings=imagesEnabled=false",
			"--enable-features=OverlayScrollbar",
			"--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36", // 模擬user-agent，防反爬
		},
	}
	//以上是設定瀏覽器引數
	caps.AddChrome(chromeCaps)

	// 調起chrome瀏覽器
	w_b1, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		fmt.Println("connect to the webDriver faild", err.Error())
		//return
	}
	//關閉一個webDriver會對應關閉一個chrome視窗
	//但是不會導致seleniumServer關閉

	err = w_b1.Get("https://juniuhome.com/#/login")
	if err != nil {
		fmt.Println("get page faild", err.Error())
		//return
	}
	//driver.findElement(By.xpath("//input[@id='gh-ac']")).sendKeys("Guitar");
	//elementID = driver.findElement(By.id("exampleInputAmount")).sendKeys("Guitar")
	//.sendKeys("op.service@solartninc.com")
	webelement, _ := w_b1.FindElement(selenium.ByXPATH, "/html/body/div/div[2]/section/div/div/div[2]/div/form/div[1]/div/p/input")
	webelement.SendKeys(os.Getenv("JUNIUHOME_ACCOUNT"))
	webelement2, _ := w_b1.FindElement(selenium.ByXPATH, "/html/body/div/div[2]/section/div/div/div[2]/div/form/div[2]/div/p/input")
	webelement2.SendKeys(os.Getenv("JUNIUHOME_PASSWORD"))
	time.Sleep(3 * time.Second)
	webelement3, _ := w_b1.FindElement(selenium.ByXPATH, "/html/body/div/div[2]/section/div/div/div[2]/div/form/div[3]/button")
	log.Println(webelement3.Text())
	webelement3.Click()

	time.Sleep(3 * time.Second)
	err = w_b1.Get("https://juniuhome.com/#/superSignApps")
	if err != nil {
		fmt.Println("get page faild", err.Error())
		//return
	}
	webelement4, _ := w_b1.FindElement(selenium.ByXPATH, "/html/body/div[1]/div[2]/div[1]/div[2]/div/div[2]/h6")
	remaining_device, _ := webelement4.Text()
	log.Println("-------remaining_device----------")
	log.Println(remaining_device)
	/*
		f_amountstotalmoney, _ := strconv.ParseFloat(amountstotalmoney, 64)

		log.Println(f_amountstotalmoney)

			webelement5, _ := w_b1.FindElement(selenium.ByClassName, "user-quick-link")

			test, _ := webelement5.GetAttribute("href")
			log.Printf("user-quick-link : %s\n", test)
	*/
	defer w_b1.Quit()
	defer w_b1.Close()
	return remaining_device
}

func PickUnusedPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0, err
	}
	return port, nil
}
