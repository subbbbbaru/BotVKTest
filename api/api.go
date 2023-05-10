package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/subbbbbaru/botvktest/utils"
)

const (
	URL_VK = "https://api.vk.com/method/"
)

type Params map[string]interface{}

// VK struct.
type MYVK struct {
	accessToken string
	MethodURL   string
	Version     string
}

func NewVK(token string, version string) *MYVK {
	var vk MYVK

	vk.accessToken = token
	vk.Version = version

	vk.MethodURL = URL_VK

	return &vk
}
func (myVK *MYVK) GetLongpollServer(groupId string) (response utils.GroupsLongPollServer, err error) {
	params := url.Values{}
	params.Add("access_token", myVK.accessToken)
	params.Add("group_id", groupId)
	params.Add("v", myVK.Version)

	u, _ := url.ParseRequestURI(myVK.MethodURL)
	u.Path += "groups.getLongPollServer"
	u.RawQuery = params.Encode()
	reqUrl := fmt.Sprintf("%v", u)
	resp, err := http.Get(reqUrl)

	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	resp_body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(resp_body, &response); err != nil {
		return response, err
	}

	if response.Error.ErrorCode != 0 {
		return response, fmt.Errorf("error in getlongpollserver: %s", response.Error.ErrorMsg)
	}
	return
}

func getRandomId() (random_id string) {
	rand.Seed(time.Now().UnixNano())
	random_id = strconv.Itoa(rand.Intn(1000000))
	return
}

func buildQueryParams(sliceParams ...Params) url.Values {
	query := url.Values{}
	for _, params := range sliceParams {
		for key, value := range params {
			query.Set(key, FmtValue(value))
		}
	}
	return query
}

func FmtValue(value interface{}) string {
	if value == nil {
		return ""
	}

	switch f := value.(type) {
	case bool:
		return fmtBool(f)
	case utils.Attachments:
		return f.ToAttachment()
	case utils.JSONObject:
		return f.ToJSON()
	case string:
		return f
	}

	return "error"
}

func fmtBool(value bool) string {
	if value {
		return "1"
	}
	return "0"
}

func (myVK *MYVK) MessageSend(userId string, message string, params ...Params) error {
	queryParams := buildQueryParams(params...)
	queryParams.Set("user_id", userId)
	queryParams.Set("message", message)
	queryParams.Set("access_token", myVK.accessToken)
	queryParams.Set("v", myVK.Version)
	queryParams.Set("random_id", getRandomId())

	u, _ := url.ParseRequestURI(myVK.MethodURL)
	u.Path += "messages.send"
	u.RawQuery = queryParams.Encode()
	reqUrl := fmt.Sprintf("%v", u)
	resp, err := http.Get(reqUrl)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil

}
func uploadURL(responseBody []byte, param string) (string, error) {
	var result map[string]any
	err := json.Unmarshal(responseBody, &result)
	if err != nil {
		return "", fmt.Errorf("error unmarshal json: %s", err.Error())
	}

	if param == "response" {
		resp := result["response"].(map[string]any)
		return fmt.Sprintf("%s", resp["upload_url"]), nil
	}
	if param == "file" {
		return fmt.Sprintf("%s", result[param]), nil
	}

	return "", fmt.Errorf("%s", "Error param")
}

func (myVK *MYVK) GetMessagesUploadServer(peer_id string, fileName string) (utils.Attachment, error) {
	queryParams := url.Values{}
	queryParams.Set("access_token", myVK.accessToken)
	queryParams.Set("peer_id", peer_id)
	queryParams.Set("v", myVK.Version)

	u, _ := url.ParseRequestURI(myVK.MethodURL)
	u.Path += "docs.getMessagesUploadServer"
	u.RawQuery = queryParams.Encode()
	reqUrl := fmt.Sprintf("%v", u)
	resp, err := http.Get(reqUrl)
	if err != nil {
		return utils.Attachment{}, err
	}
	defer resp.Body.Close()

	resp_body, _ := ioutil.ReadAll(resp.Body)
	upload_url, err := uploadURL(resp_body, "response")

	if err != nil {
		return utils.Attachment{}, err
	}

	fileUrlSave, err := uploadFile(fileName, upload_url)

	if err != nil {
		return utils.Attachment{}, err
	}
	attachment, err := myVK.docSave(fileName, fileUrlSave)
	if err != nil {
		return utils.Attachment{}, err
	}
	return attachment, nil
}

func uploadFile(fileName string, uploadUrl string) (string, error) {
	bodyBuff := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuff)
	fileWriter, err := bodyWriter.CreateFormFile("file", fileName)
	if err != nil {
		panic(err)
	}
	fh, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fh.Close()

	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		panic(err)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(uploadUrl, contentType, bodyBuff)
	if err != nil {
		panic(err)
	}
	resp_body2, _ := ioutil.ReadAll(resp.Body)
	fileUrlSave, err := uploadURL(resp_body2, "file")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return fileUrlSave, err
}

func (myVK *MYVK) docSave(filePath string, fileUrlSave string) (utils.Attachment, error) {

	params := url.Values{}
	fileName := strings.FieldsFunc(filePath, func(r rune) bool {
		if r == '/' {
			return true
		}
		return false
	})

	params.Add("access_token", myVK.accessToken)
	params.Add("file", fileUrlSave)
	params.Add("title", fileName[len(fileName)-1])
	params.Add("v", myVK.Version)
	u, _ := url.ParseRequestURI(myVK.MethodURL)
	u.Path += "docs.save"
	u.RawQuery = params.Encode()

	urlStr3 := fmt.Sprintf("%v", u)

	resp, err := http.Get(urlStr3)
	if err != nil {
		return utils.Attachment{}, err
	}
	resp_body3, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return utils.Attachment{}, err
	}
	attachment, _ := utils.UnmarshalWelcome(resp_body3)
	if err != nil {
		return utils.Attachment{}, err
	}

	defer resp.Body.Close()

	return attachment, nil

}
