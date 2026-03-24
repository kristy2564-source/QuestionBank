package main

import "encoding/json"

// ========================
// 前端请求类型
// ========================

type AskRequest struct {
	Question   string           `json:"question"`
	KbID       *string          `json:"kb_id,omitempty"`
	Messages   []Message        `json:"messages,omitempty"`
	QueryParam *KnowledgeFilter `json:"query_param,omitempty"`
	ServiceID  *string          `json:"service_id,omitempty"`
	Stream     bool             `json:"stream,omitempty"`
}

type Message struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

type KnowledgeFilter struct {
	DocFilter json.RawMessage `json:"doc_filter"`
}

type KnowledgeServiceRequest struct {
	ServiceResourceID string           `json:"service_resource_id"`
	Messages          []Message        `json:"messages"`
	QueryParam        *KnowledgeFilter `json:"query_param,omitempty"`
	Stream            bool             `json:"stream"`
}

// ========================
// 知识服务响应类型（来自火山方舟官方SDK）
// ========================

type ServiceChatResponse struct {
	Code    int64                              `json:"code"`
	Message string                             `json:"message,omitempty"`
	Data    *CollectionServiceChatResponseData `json:"data,omitempty"`
}

type CollectionServiceChatResponseData struct {
	CollectionSearchKnowledgeResponseData
	*CollectionChatCompletionResponseData
}

type CollectionSearchKnowledgeResponseData struct {
	Count        int32                           `json:"count"`
	RewriteQuery string                          `json:"rewrite_query,omitempty"`
	TokenUsage   *TotalTokenUsage                `json:"token_usage,omitempty"`
	ResultList   []*CollectionSearchResponseItem `json:"result_list,omitempty"`
}

type TotalTokenUsage struct {
	EmbeddingUsage *ModelTokenUsage `json:"embedding_token_usage,omitempty"`
	RerankUsage    *int64           `json:"rerank_token_usage,omitempty"`
	LLMUsage       *ModelTokenUsage `json:"llm_token_usage,omitempty"`
	RewriteUsage   *ModelTokenUsage `json:"rewrite_token_usage,omitempty"`
}

type ModelTokenUsage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

type CollectionSearchResponseItem struct {
	Id                  string                              `json:"id"`
	Content             string                              `json:"content"`
	MdContent           string                              `json:"md_content,omitempty"`
	Score               float64                             `json:"score"`
	PointId             string                              `json:"point_id"`
	OriginText          string                              `json:"origin_text,omitempty"`
	OriginalQuestion    string                              `json:"original_question,omitempty"`
	ChunkTitle          string                              `json:"chunk_title,omitempty"`
	ChunkId             int                                 `json:"chunk_id"`
	ProcessTime         int64                               `json:"process_time"`
	RerankScore         float64                             `json:"rerank_score,omitempty"`
	DocInfo             CollectionSearchResponseItemDocInfo `json:"doc_info,omitempty"`
	RecallPosition      int32                               `json:"recall_position"`
	RerankPosition      int32                               `json:"rerank_position,omitempty"`
	ChunkType           string                              `json:"chunk_type,omitempty"`
	ChunkSource         string                              `json:"chunk_source,omitempty"`
	UpdateTime          int64                               `json:"update_time"`
	ChunkAttachmentList []ChunkAttachment                   `json:"chunk_attachment,omitempty"`
	TableChunkFields    []PointTableChunkField              `json:"table_chunk_fields,omitempty"`
	OriginalCoordinate  *ChunkPositions                     `json:"original_coordinate,omitempty"`
}

type CollectionSearchResponseItemDocInfo struct {
	DocId      string `json:"doc_id"`
	DocName    string `json:"doc_name"`
	CreateTime int64  `json:"create_time"`
	DocType    string `json:"doc_type"`
	DocMeta    string `json:"doc_meta,omitempty"`
	Source     string `json:"source"`
	Title      string `json:"title,omitempty"`
}

type ChunkAttachment struct {
	UUID    string `json:"uuid,omitempty"`
	Caption string `json:"caption"`
	Type    string `json:"type"`
	Link    string `json:"link,omitempty"`
}

type PointTableChunkField struct {
	FieldName  string      `json:"field_name"`
	FieldValue interface{} `json:"field_value"`
}

type ChunkPositions struct {
	PageNo []int       `json:"page_no"`
	BBox   [][]float64 `json:"bbox"`
}

type CollectionChatCompletionResponseData struct {
	GenerateAnswer   string  `json:"generated_answer"`
	ReasoningContent string  `json:"reasoning_content,omitempty"`
	Prompt           *string `json:"prompt,omitempty"`
	End              bool    `json:"end,omitempty"`
}

// ========================
// 文档上传类型
// ========================

type UploadItem struct {
	Question string   `json:"question"`
	Answer   string   `json:"answer"`
	Tags     []string `json:"tags,omitempty"`
}

type UploadRequest struct {
	Items []UploadItem `json:"items"`
	KbID  *string      `json:"kb_id,omitempty"`
}

type DocAddRequest struct {
	CollectionName    string          `json:"collection_name,omitempty"`
	Project           string          `json:"project,omitempty"`
	ResourceID        string          `json:"resource_id,omitempty"`
	ServiceResourceID string          `json:"service_resource_id,omitempty"`
	AddType           string          `json:"add_type,omitempty"`
	DocID             string          `json:"doc_id,omitempty"`
	DocName           string          `json:"doc_name,omitempty"`
	DocType           string          `json:"doc_type,omitempty"`
	Description       string          `json:"description,omitempty"`
	LarkFile          json.RawMessage `json:"lark_file,omitempty"`
	TOSPath           string          `json:"tos_path,omitempty"`
	URL               string          `json:"url,omitempty"`
	Meta              json.RawMessage `json:"meta,omitempty"`
}

type LocalUploadRequest struct {
	Path              string `json:"path"`
	ResourceID        string `json:"resource_id"`
	CollectionName    string `json:"collection_name"`
	ServiceResourceID string `json:"service_resource_id,omitempty"`
	Project           string `json:"project,omitempty"`
	DocName           string `json:"doc_name,omitempty"`
	DocType           string `json:"doc_type,omitempty"`
	Description       string `json:"description,omitempty"`
	DryRun            bool   `json:"dry_run,omitempty"`
}

// ========================
// 组卷请求类型
// ========================

type ComposeRequest struct {
	Title      string           `json:"title,omitempty"`
	Subject    string           `json:"subject,omitempty"`
	Grade      string           `json:"grade,omitempty"`
	Difficulty string           `json:"difficulty,omitempty"`
	Counts     map[string]int   `json:"counts,omitempty"`
	Tags       []string         `json:"tags,omitempty"`
	Messages   []Message        `json:"messages,omitempty"`
	QueryParam *KnowledgeFilter `json:"query_param,omitempty"`
	ServiceID  *string          `json:"service_id,omitempty"`
	Format     string           `json:"format,omitempty"`
}
