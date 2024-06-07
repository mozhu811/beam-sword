package lark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	BatchRequestUrl      = "https://open.feishu.cn/open-apis/bitable/v1/apps/%s/tables/%s/records/batch_create"
	TenantAccessTokenUrl = "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"
)

func CreateRecord(records string) error {
	//if !gjson.Valid(records) {
	//	return fmt.Errorf("records is not valid JSON: %s", records)
	//}
	appToken := os.Getenv("BIT_TABLE_APP_TOKEN")
	tableId := os.Getenv("TABLE_ID")
	url := fmt.Sprintf(BatchRequestUrl, appToken, tableId)
	r := gjson.Parse(records)
	var reqBody BatchCreateReq
	r.ForEach(func(key, value gjson.Result) bool {
		var record Record
		date := value.Get("date").String()
		rDate, _ := time.Parse("2006/01/02", date)
		fields := Fields{
			Event:  value.Get("event").String(),
			Type:   value.Get("type").String(),
			Amount: value.Get("amount").Num,
			Tag:    value.Get("tag").String(),
			Date:   rDate.UnixMilli(),
		}
		record.Fields = fields
		reqBody.Records = append(reqBody.Records, record)
		return true
	})

	bodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal body json fail: %w", err)
	}
	request, err := http.NewRequest("POST", url,
		bytes.NewBuffer(bodyJson))

	if err != nil {
		return err
	}
	token, err := requestTenantAccessToken()
	if err != nil {
		return fmt.Errorf("request lark token fail: %w", err)
	}
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	c := http.Client{}
	response, err := c.Do(request)
	if err != nil {
		return fmt.Errorf("request lark api fail: %w", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("read response body fail: %w", err)
	}

	fmt.Println(string(body))

	return nil
}

func requestTenantAccessToken() (string, error) {
	appID := os.Getenv("LARK_APP_ID")
	appSecret := os.Getenv("LARK_APP_SECRET")

	reqBody := make(map[string]string)
	reqBody["app_id"] = appID
	reqBody["app_secret"] = appSecret

	reqBodyJson, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal body json fail: %w", err)
	}
	request, err := http.NewRequest("POST", TenantAccessTokenUrl, bytes.NewBuffer(reqBodyJson))
	if err != nil {
		return "", fmt.Errorf("create request body fail: %w", err)
	}

	c := http.Client{}
	response, err := c.Do(request)
	if err != nil {
		return "", fmt.Errorf("request tenant access token fail: %w", err)
	}
	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("read response body fail: %w", err)
	}

	fmt.Println("access token api response body:", string(respBody))
	respBodyJson := gjson.ParseBytes(respBody)
	token := respBodyJson.Get("tenant_access_token").String()
	if token == "" {
		return "", fmt.Errorf("get tenant access token fail")
	}

	return token, nil
}
