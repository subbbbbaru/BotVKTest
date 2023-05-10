package longpoll

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/subbbbbaru/botvktest/api"
	"github.com/subbbbbaru/botvktest/utils"
)

type Longpoll struct {
	VK      *api.MYVK
	Key     string
	Ts      string
	Wait    int
	Server  string
	GroupId string
}

func NewLongpoll(vk *api.MYVK, groupId string) (*Longpoll, error) {
	lp := &Longpoll{
		VK:      vk,
		GroupId: groupId,
		Wait:    25,
	}

	err := lp.updateServer(true)

	return lp, err
}

func (longpoll *Longpoll) updateServer(updateTs bool) error {
	serverSetting, err := longpoll.VK.GetLongpollServer(longpoll.GroupId)
	if err != nil {
		return err
	}
	longpoll.Key = serverSetting.Response.Key
	longpoll.Server = serverSetting.Response.Server
	if updateTs {
		longpoll.Ts = serverSetting.Response.Ts
	}

	return nil
}
func getDirs(dirname string) map[string]string {
	mainMenuButtons := make(map[string]string)

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {

		if file.IsDir() {
			mainMenuButtons[file.Name()] = dirname + "/" + file.Name()
		}
	}
	return mainMenuButtons
}
func mainMenu() utils.Keyboard {
	startDir := "BOOK"
	mainButtons := getDirs(startDir)
	keyB := utils.Keyboard{Buttons: make([][]utils.Button, 0)}
	for dir, pathDir := range mainButtons {
		button := utils.NewButton(dir, "{\"button\":\""+pathDir+"\"}")
		row := make([]utils.Button, 0)
		row = append(row, button)
		keyB.Buttons = append(keyB.Buttons, row)
		keyB.Inline = false
	}
	return keyB //, mainButtons
}

func secondMenu(secondDir string) utils.Keyboard {
	secondButtons := getDirs(secondDir)
	keyB := utils.Keyboard{Buttons: make([][]utils.Button, 0)}
	//buttons := []utils.Button{}
	for dir, pathDir := range secondButtons {
		button := utils.NewButton(dir, "{\"button\":\""+pathDir+"\"}")
		row := make([]utils.Button, 0)
		row = append(row, button)
		keyB.Buttons = append(keyB.Buttons, row)
		keyB.Inline = false
	}
	buttonNegative := utils.NewButton("Back", nil)
	buttonNegative.Color = "negative"
	rowNegative := make([]utils.Button, 0)
	rowNegative = append(rowNegative, buttonNegative)
	keyB.Buttons = append(keyB.Buttons, rowNegative)
	return keyB //, secondButtons
}
func getFiles(dirname string) []string {
	books := make([]string, 0)

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			books = append(books, file.Name())
		}
	}

	return books
}

var payloadIOS = false
var payloadAndroid = false
var payloadWindows = false
var payloadLinux = false
var payloadC_CPP = false
var payloadGolang = false
var payloadAlgorithm = false
var payloadOOP = false

func (longpoll *Longpoll) sendWithAttach(payload *bool, msg string, vibirFile *map[int]string, fromId string) bool {
	getNumber, errnum := strconv.Atoi(msg)
	if errnum != nil {
		longpoll.VK.MessageSend(fromId, "Введи число!")
		return false
	}
	if getNumber >= len(*vibirFile)+1 || getNumber <= 0 {
		longpoll.VK.MessageSend(fromId, "Число не верно!")
		return false
	}

	if path, ok := (*vibirFile)[getNumber]; ok {
		attach, _ := longpoll.VK.GetMessagesUploadServer(fromId, path)
		typeFile := attach.Response.Type
		ownerId := fmt.Sprintf("%v", attach.Response.Doc.OwnerID)
		fileId := fmt.Sprintf("%v", attach.Response.Doc.ID)
		*payload = false
		if err := longpoll.VK.MessageSend(fromId, "", api.Params{
			"keyboard":   mainMenu(),
			"attachment": typeFile + ownerId + "_" + fileId,
		}); err != nil {
			return false
		}
	}
	*vibirFile = map[int]string{}
	//delete(*vibirFile, getNumber)
	return true
}

func negativeKey() utils.Keyboard {
	keyB := utils.Keyboard{Buttons: make([][]utils.Button, 0)}
	buttonNegative := utils.NewButton("Back", nil)
	buttonNegative.Color = "negative"
	rowNegative := make([]utils.Button, 0)
	rowNegative = append(rowNegative, buttonNegative)
	keyB.Buttons = append(keyB.Buttons, rowNegative)
	return keyB //, secondButtons
}

func (longpoll *Longpoll) LongpollHandler() error {
	vibirFile := make(map[int]string)
	for {
		updates, err := longpoll.getLongPollUpdates()
		if err != nil {
			return fmt.Errorf("error request VK: %s", err.Error())
		}
		for _, update := range updates {
			if update.Type == "message_new" {
				log.Println(update.Object.Message)
				fromId := strconv.Itoa(update.Object.Message.FromID)
				if update.Object.Message.Text == "Start" {
					if err = longpoll.VK.MessageSend(fromId, "Выбирай!", api.Params{
						"keyboard": mainMenu(),
					}); err != nil {
						return err
					}
				}
				if payloadWindows {
					if update.Object.Message.Text == "Back" && update.Object.Message.Payload == "" {
						payloadWindows = false
					} else {
						check := longpoll.sendWithAttach(&payloadWindows, update.Object.Message.Text, &vibirFile, fromId)
						if !check {
							continue
						}
					}

				} else if payloadLinux {
					if update.Object.Message.Text == "Back" && update.Object.Message.Payload == "" {
						payloadLinux = false
					} else {
						check := longpoll.sendWithAttach(&payloadLinux, update.Object.Message.Text, &vibirFile, fromId)
						if !check {
							continue
						}
					}
				} else if payloadAndroid {
					if update.Object.Message.Text == "Back" && update.Object.Message.Payload == "" {
						payloadAndroid = false
					} else {
						check := longpoll.sendWithAttach(&payloadAndroid, update.Object.Message.Text, &vibirFile, fromId)
						if !check {
							continue
						}
					}
				} else if payloadIOS {
					if update.Object.Message.Text == "Back" && update.Object.Message.Payload == "" {
						payloadIOS = false
					} else {
						check := longpoll.sendWithAttach(&payloadIOS, update.Object.Message.Text, &vibirFile, fromId)
						if !check {
							continue
						}
					}
				} else if payloadAlgorithm {
					if update.Object.Message.Text == "Back" && update.Object.Message.Payload == "" {
						payloadAlgorithm = false
					} else {
						check := longpoll.sendWithAttach(&payloadAlgorithm, update.Object.Message.Text, &vibirFile, fromId)
						if !check {
							continue
						}
					}
				} else if payloadOOP {
					if update.Object.Message.Text == "Back" && update.Object.Message.Payload == "" {
						payloadOOP = false
					} else {
						check := longpoll.sendWithAttach(&payloadOOP, update.Object.Message.Text, &vibirFile, fromId)
						if !check {
							continue
						}
					}
				} else if payloadC_CPP {
					if update.Object.Message.Text == "Back" && update.Object.Message.Payload == "" {
						payloadC_CPP = false
					} else {
						check := longpoll.sendWithAttach(&payloadC_CPP, update.Object.Message.Text, &vibirFile, fromId)
						if !check {
							continue
						}
					}
				} else if payloadGolang {
					if update.Object.Message.Text == "Back" && update.Object.Message.Payload == "" {
						payloadGolang = false
					} else {
						check := longpoll.sendWithAttach(&payloadGolang, update.Object.Message.Text, &vibirFile, fromId)
						if !check {
							continue
						}
					}
				}

				switch update.Object.Message.Payload {
				case "{\"command\":\"start\"}":
					if err = longpoll.VK.MessageSend(fromId, "Привет!\nУ меня есть немного книг для тебя!\nВыбери категорию)))", api.Params{
						"keyboard": mainMenu(),
					}); err != nil {
						return err
					}

				case "\"{\\\"button\\\":\\\"BOOK\\/Mobile\\\"}\"":
					// payloadMobile = true
					if err = longpoll.VK.MessageSend(fromId, "Выбирай!", api.Params{
						"keyboard": secondMenu("BOOK/Mobile"),
					}); err != nil {
						return err
					}
				case "\"{\\\"button\\\":\\\"BOOK\\/Mobile\\/IOS\\\"}\"":
					payloadIOS = true
					diir := "BOOK/Mobile/IOS"
					files := getFiles(diir)
					str := ""
					for idx, file := range files {
						vibirFile[idx+1] = diir + "/" + file
						str += fmt.Sprintf("%d. %s\n", idx+1, file)
					}
					if err = longpoll.VK.MessageSend(fromId, "Выбирай!\n"+str, api.Params{
						"keyboard": negativeKey(),
					}); err != nil {
						return err
					}
				case "\"{\\\"button\\\":\\\"BOOK\\/Mobile\\/Android\\\"}\"":
					payloadAndroid = true
					diir := "BOOK/Mobile/Android"
					files := getFiles(diir)
					str := ""
					for idx, file := range files {
						vibirFile[idx+1] = diir + "/" + file
						str += fmt.Sprintf("%d. %s\n", idx+1, file)
					}
					if err = longpoll.VK.MessageSend(fromId, "Выбирай!\n"+str, api.Params{
						"keyboard": negativeKey(),
					}); err != nil {
						return err
					}

				case "\"{\\\"button\\\":\\\"BOOK\\/PL\\\"}\"":
					// payloadPL = true
					if err = longpoll.VK.MessageSend(fromId, "Выбирай!", api.Params{
						"keyboard": secondMenu("BOOK/PL"),
					}); err != nil {
						return err
					}
				case "\"{\\\"button\\\":\\\"BOOK\\/PL\\/C C++\\\"}\"":
					payloadC_CPP = true
					diir := "BOOK/PL/C C++"
					files := getFiles(diir)
					str := ""
					for idx, file := range files {
						vibirFile[idx+1] = diir + "/" + file
						str += fmt.Sprintf("%d. %s\n", idx+1, file)
					}
					if err = longpoll.VK.MessageSend(fromId, "Выбирай!\n"+str, api.Params{
						"keyboard": negativeKey(),
					}); err != nil {
						return err
					}
				case "\"{\\\"button\\\":\\\"BOOK\\/PL\\/GO\\\"}\"":
					payloadGolang = true
					diir := "BOOK/PL/GO"
					files := getFiles(diir)
					str := ""
					for idx, file := range files {
						vibirFile[idx+1] = diir + "/" + file
						str += fmt.Sprintf("%d. %s\n", idx+1, file)
					}
					if err = longpoll.VK.MessageSend(fromId, "Выбирай!\n"+str, api.Params{
						"keyboard": negativeKey(),
					}); err != nil {
						return err
					}

				case "\"{\\\"button\\\":\\\"BOOK\\/Основы\\\"}\"":
					// payloadComputerScience = true
					if err = longpoll.VK.MessageSend(fromId, "Выбирай!", api.Params{
						"keyboard": secondMenu("BOOK/Основы"),
					}); err != nil {
						return err
					}
				case "\"{\\\"button\\\":\\\"BOOK\\/Основы\\/OOP\\\"}\"":
					payloadOOP = true
					diir := "BOOK/Основы/OOP"
					files := getFiles(diir)
					str := ""
					for idx, file := range files {
						vibirFile[idx+1] = diir + "/" + file
						str += fmt.Sprintf("%d. %s\n", idx+1, file)
					}
					if err = longpoll.VK.MessageSend(fromId, "Выбирай!\n"+str, api.Params{
						"keyboard": negativeKey(),
					}); err != nil {
						return err
					}
				case "\"{\\\"button\\\":\\\"BOOK\\/Основы\\/Алгоритмы\\\"}\"":
					payloadAlgorithm = true
					diir := "BOOK/Основы/Алгоритмы"
					files := getFiles(diir)
					str := ""
					for idx, file := range files {
						vibirFile[idx+1] = diir + "/" + file
						str += fmt.Sprintf("%d. %s\n", idx+1, file)
					}
					if err = longpoll.VK.MessageSend(fromId, "Выбирай!\n"+str, api.Params{
						"keyboard": negativeKey(),
					}); err != nil {
						return err
					}

				case "\"{\\\"button\\\":\\\"BOOK\\/OS\\\"}\"":
					// payloadOS = true
					if err = longpoll.VK.MessageSend(fromId, "Выбирай!", api.Params{
						"keyboard": secondMenu("BOOK/OS"),
					}); err != nil {
						return err
					}

				case "\"{\\\"button\\\":\\\"BOOK\\/OS\\/Linux\\\"}\"":
					payloadLinux = true
					diir := "BOOK/OS/Linux"
					files := getFiles(diir)
					str := ""
					for idx, file := range files {
						vibirFile[idx+1] = diir + "/" + file
						str += fmt.Sprintf("%d. %s\n", idx+1, file)
					}
					if err = longpoll.VK.MessageSend(fromId, "Выбирай!\n"+str, api.Params{
						"keyboard": negativeKey(),
					}); err != nil {
						return err
					}
				case "\"{\\\"button\\\":\\\"BOOK\\/OS\\/Windows\\\"}\"":
					payloadWindows = true
					diir := "BOOK/OS/Windows"
					files := getFiles(diir)
					str := ""
					for idx, file := range files {
						vibirFile[idx+1] = diir + "/" + file
						str += fmt.Sprintf("%d. %s\n", idx+1, file)
					}
					if err = longpoll.VK.MessageSend(fromId, "Выбирай!\n"+str, api.Params{
						"keyboard": negativeKey(),
					}); err != nil {
						return err
					}
				default:
					if err = longpoll.VK.MessageSend(fromId, "Нажми на кнопку)))\nЯ могу долго загружать книги(((\nПрости...", api.Params{
						"keyboard": mainMenu(),
					}); err != nil {
						return err
					}
				}
			}
		}

	}
}

func (longpoll *Longpoll) getLongPollUpdates() ([]utils.Update, error) {
	params := url.Values{}
	params.Add("act", "a_check")
	params.Add("key", longpoll.Key)
	params.Add("ts", longpoll.Ts)
	params.Add("wait", "25")
	params.Add("v", "5.131")
	u, _ := url.ParseRequestURI(longpoll.Server)
	u.RawQuery = params.Encode()

	urlStr2 := fmt.Sprintf("%v", u)

	resp, errHttp := http.Get(urlStr2)
	if errHttp != nil {
		return nil, errHttp
	}
	defer resp.Body.Close()
	resp_body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data struct {
		Updates []utils.Update `json:"updates"`
		TS      string         `json:"ts"`
	}

	if err := json.Unmarshal(resp_body, &data); err != nil {
		return nil, err
	}

	longpoll.Ts = data.TS

	return data.Updates, nil
}
