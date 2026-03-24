package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/ask", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req AskRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.Question == "" && len(req.Messages) == 0 {
			http.Error(w, "缺少问题或消息", http.StatusBadRequest)
			return
		}
		sid := KBServiceID
		if req.ServiceID != nil && *req.ServiceID != "" {
			sid = *req.ServiceID
		}
		if sid == "" {
			http.Error(w, "缺少知识服务ID(KB_SERVICE_ID)", http.StatusBadRequest)
			return
		}
		if KBApiKey == "" {
			http.Error(w, "缺少API Key(KB_API_KEY)", http.StatusBadRequest)
			return
		}
		var messages []Message
		if len(req.Messages) > 0 {
			messages = req.Messages
		} else {
			content, _ := json.Marshal(req.Question)
			messages = []Message{{Role: "user", Content: content}}
		}
		ksReq := KnowledgeServiceRequest{
			ServiceResourceID: sid,
			Messages:          messages,
			QueryParam:        req.QueryParam,
			Stream:            req.Stream,
		}
		if !ksReq.Stream {
			res, err := CallKnowledgeServiceChat(ksReq)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadGateway)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(res)
			return
		}
		u := "https://" + KBDomain + "/api/knowledge/service/chat"
		b, _ := json.Marshal(ksReq)
		httpReq, _ := http.NewRequest("POST", u, bytes.NewReader(b))
		httpReq.Header.Set("Content-Type", "application/json")
		httpReq.Header.Set("Authorization", "Bearer "+KBApiKey)
		httpReq.Header.Set("Host", KBDomain)
		client := &http.Client{Timeout: 0}
		resp, err := client.Do(httpReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			http.Error(w, string(body), resp.StatusCode)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if f, ok := w.(http.Flusher); ok {
			buf := make([]byte, 4096)
			for {
				n, rerr := resp.Body.Read(buf)
				if n > 0 {
					w.Write(buf[:n])
					f.Flush()
				}
				if rerr != nil {
					if rerr == io.EOF {
						break
					}
					return
				}
			}
		} else {
			io.Copy(w, resp.Body)
		}
	})

	mux.HandleFunc("/api/compose", func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, fmt.Sprintf("compose panic: %v", rec), http.StatusInternalServerError)
			}
		}()
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req ComposeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sid := KBSearchServiceID
		if req.ServiceID != nil && *req.ServiceID != "" {
			sid = *req.ServiceID
		}
		if sid == "" {
			http.Error(w, "缺少知识服务ID(KB_SERVICE_ID或请求体service_id)", http.StatusBadRequest)
			return
		}
		if KBApiKey == "" {
			http.Error(w, "缺少API Key(KB_API_KEY)", http.StatusBadRequest)
			return
		}
		var messages []Message
		if len(req.Messages) > 0 {
			messages = req.Messages
		} else {
			content, _ := json.Marshal(BuildComposeQuery(req))
			messages = []Message{{Role: "user", Content: content}}
		}
		ksReq := KnowledgeServiceRequest{
			ServiceResourceID: sid,
			Messages:          messages,
			QueryParam:        req.QueryParam,
			Stream:            false,
		}
		searchRes, err := CallKnowledgeServiceChat(ksReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		var refText strings.Builder
		var items int
		if data, ok := searchRes["data"].(map[string]interface{}); ok {
			if list, ok := data["result_list"].([]interface{}); ok {
				for _, it := range list {
					if items >= 30 {
						break
					}
					m, _ := it.(map[string]interface{})
					txt := ""
					if s, ok := m["content"].(string); ok && s != "" {
						txt = s
					} else if s, ok := m["md_content"].(string); ok && s != "" {
						txt = s
					} else if s, ok := m["origin_text"].(string); ok && s != "" {
						txt = s
					}
					if txt == "" {
						continue
					}
					if len(txt) > 800 {
						txt = txt[:800]
					}
					refText.WriteString("- ")
					refText.WriteString(txt)
					refText.WriteString("\n")
					items++
				}
			}
		}
		if ARKApiKey == "" {
			http.Error(w, "缺少ARK_API_KEY", http.StatusBadRequest)
			return
		}
		if ArkModelID == "" {
			http.Error(w, "缺少ARK模型ID(ARK_MODEL_ID)", http.StatusBadRequest)
			return
		}
		var prompt strings.Builder
		prompt.WriteString("你是智能出题与组卷助手，根据“组卷要求”和“参考材料”生成结构化试卷。\n")
		prompt.WriteString("输出严格为JSON，字段：title,string; sections,array; sections[].type,string; sections[].count,int; sections[].items,array; items[].id,string; items[].stem,string; items[].options,array[string],可选; items[].answer,string或array; items[].explanation,string可选。\n")
		prompt.WriteString("题型命名示例：single_choice,multiple_choice,fill_blank,short_answer。\n")
		prompt.WriteString("组卷要求：")
		prompt.WriteString(BuildComposeQuery(req))
		prompt.WriteString("\n参考材料（节选）：\n")
		prompt.WriteString(refText.String())
		promptText := prompt.String()
		llmRes, status, err := CallArkResponses(promptText)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		respObj := map[string]interface{}{
			"status": status,
			"search": searchRes,
			"llm":    llmRes,
		}
		b, _ := json.Marshal(respObj)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(b)))
		w.Write(b)
	})

	mux.HandleFunc("/api/upload", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		var req DocAddRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if KBAK == "" || KBSK == "" {
			http.Error(w, "缺少AK/SK(KB_AK, KB_SK)", http.StatusBadRequest)
			return
		}
		if req.ResourceID == "" && req.CollectionName == "" {
			if KBID != "" {
				req.ResourceID = KBID
			}
		}
		if req.Project == "" {
			req.Project = "default"
		}
		if req.ServiceResourceID == "" && KBServiceID != "" {
			req.ServiceResourceID = KBServiceID
		}
		if req.AddType == "" {
			http.Error(w, "缺少add_type", http.StatusBadRequest)
			return
		}
		res, status, err := SignAndPostDocAdd(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(res)
	})

	// 文件上传：浏览器选文件 → Go接收 → 上传TOS → 调知识库 add_doc
	mux.HandleFunc("/api/upload_file", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// 检查配置
		if KBAK == "" || KBSK == "" {
			http.Error(w, "缺少AK/SK配置(KB_AK, KB_SK)", http.StatusBadRequest)
			return
		}
		if TOSBucket == "" {
			http.Error(w, "缺少TOS_BUCKET配置", http.StatusBadRequest)
			return
		}

		// 解析 multipart 表单，最大 200MB
		if err := r.ParseMultipartForm(200 << 20); err != nil {
			http.Error(w, "文件解析失败: "+err.Error(), http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "缺少文件: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 读取表单参数
		docName := r.FormValue("doc_name")
		if docName == "" {
			docName = header.Filename
		}
		docType := r.FormValue("doc_type")
		if docType == "" {
			// 从文件扩展名推断
			ext := strings.ToLower(filepath.Ext(header.Filename))
			switch ext {
			case ".pdf":
				docType = "pdf"
			case ".docx":
				docType = "docx"
			case ".doc":
				docType = "doc"
			case ".txt":
				docType = "txt"
			case ".md":
				docType = "md"
			case ".pptx":
				docType = "pptx"
			default:
				docType = "pdf"
			}
		}
		description := r.FormValue("description")

		// 生成 TOS 对象路径：uploads/2026-03/原文件名
		now := time.Now()
		objectKey := fmt.Sprintf("uploads/%s/%d_%s", now.Format("2006-01"), now.UnixMilli(), header.Filename)

		// 第一步：上传到 TOS
		tosPath, err := UploadToTOS(r.Context(), objectKey, file, header.Size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}

		// 第二步：调知识库 add_doc API，用 TOS 路径导入
		payload := DocAddRequest{
			Project:     "default",
			AddType:     "tos",
			TOSPath:     tosPath,
			DocName:     docName,
			DocType:     docType,
			Description: description,
		}
		if KBID != "" {
			payload.ResourceID = KBID
		}
		if KBServiceID != "" {
			payload.ServiceResourceID = KBServiceID
		}

		addRes, status, err := SignAndPostDocAdd(payload)
		if err != nil {
			http.Error(w, "知识库导入失败: "+err.Error(), http.StatusBadGateway)
			return
		}

		result := map[string]interface{}{
			"tos_path":  tosPath,
			"file_name": header.Filename,
			"file_size": header.Size,
			"doc_name":  docName,
			"doc_type":  docType,
			"kb_result": addRes,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(result)
	})

	staticDir := filepath.Join("app", "static")
	mux.Handle("/", http.FileServer(http.Dir(staticDir)))
}
