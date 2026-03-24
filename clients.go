package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
	"github.com/volcengine/volc-sdk-golang/base"
)

func CallKnowledgeServiceChat(req KnowledgeServiceRequest) (map[string]interface{}, error) {
	u := "https://" + KBDomain + "/api/knowledge/service/chat"
	b, _ := json.Marshal(req)
	httpReq, _ := http.NewRequest("POST", u, bytes.NewReader(b))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+KBApiKey)
	httpReq.Header.Set("Host", KBDomain)
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var out map[string]interface{}
	if err := json.Unmarshal(body, &out); err != nil {
		return map[string]interface{}{"raw": string(body)}, nil
	}
	return out, nil
}

func SignAndPostDocAdd(payload DocAddRequest) (map[string]interface{}, int, error) {
	u := "https://" + KBDomain + "/api/knowledge/doc/add"
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest(strings.ToUpper("POST"), u, bytes.NewReader(b))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	if AccountID != "" {
		req.Header.Set("V-Account-Id", AccountID)
	}
	req.Host = KBDomain
	cred := base.Credentials{
		AccessKeyID:     KBAK,
		SecretAccessKey: KBSK,
		Service:         "air",
		Region:          "cn-beijing",
	}
	req = cred.Sign(req)
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var out map[string]interface{}
	if err := json.Unmarshal(body, &out); err != nil {
		return map[string]interface{}{"raw": string(body)}, resp.StatusCode, nil
	}
	return out, resp.StatusCode, nil
}

func CallArkResponses(prompt string) (map[string]interface{}, int, error) {
	u := "https://ark.cn-beijing.volces.com/api/v3/responses"
	body := map[string]interface{}{
		"model": ArkModelID,
		"input": []map[string]interface{}{
			{
				"role": "user",
				"content": []map[string]interface{}{
					{"type": "input_text", "text": prompt},
				},
			},
		},
	}
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", u, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ARKApiKey)
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	resBody, _ := io.ReadAll(resp.Body)
	var out map[string]interface{}
	if err := json.Unmarshal(resBody, &out); err != nil {
		return map[string]interface{}{"raw": string(resBody)}, resp.StatusCode, nil
	}
	return out, resp.StatusCode, nil
}

// UploadToTOS 将文件上传到 TOS，返回 TOS 路径（bucket/key 格式）
func UploadToTOS(ctx context.Context, objectKey string, data io.Reader, size int64) (string, error) {
	tosClient, err := tos.NewClientV2(TOSEndpoint,
		tos.WithRegion(TOSRegion),
		tos.WithCredentials(tos.NewStaticCredentials(KBAK, KBSK)),
	)
	if err != nil {
		return "", fmt.Errorf("TOS客户端初始化失败: %w", err)
	}

	_, err = tosClient.PutObjectV2(ctx, &tos.PutObjectV2Input{
		PutObjectBasicInput: tos.PutObjectBasicInput{
			Bucket:        TOSBucket,
			Key:           objectKey,
			ContentLength: size,
		},
		Content: data,
	})
	if err != nil {
		return "", fmt.Errorf("TOS上传失败: %w", err)
	}

	// 返回 TOS 路径，格式为 bucket/key，知识库 add_doc API 需要这个格式
	return TOSBucket + "/" + objectKey, nil
}
