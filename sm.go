package ishumei

import (
	"context"
	"net/http"
)

type SmService service

// ===================================数美文本审核开始======================================

type TextAuditingJobOptions struct {
	AccessKey    string `json:"accessKey"`
	AppID        string `json:"appId"`
	EventID      string `json:"eventId"`
	Type         string `json:"type"`
	BusinessType string `json:"businessType,omitempty"`
	Data         *TextData  `json:"data"`
}

type TextData struct{
	Text     string `json:"text"`
	TokenID  string `json:"tokenId"`
	Lang     string `json:"lang,omitempty"`
	Nickname string `json:"nickname,omitempty"`
	IP       string `json:"ip,omitempty"`
	DeviceID string `json:"deviceId,omitempty"`
	Extra   *TextDataExtra  `json:"extra,omitempty"`
}

type TextDataExtra struct{
	ReceiveTokenID  string       `json:"receiveTokenId,omitempty"`
	Topic           string       `json:"topic,omitempty"`
	AtID            string       `json:"atId,omitempty"`
	Room            string       `json:"room,omitempty"`
	Level           int          `json:"level,omitempty"`
	Role            string       `json:"role,omitempty"`
	Sex             int          `json:"sex,omitempty"`
	IsTokenSeparate int          `json:"isTokenSeparate,omitempty"`
	PassThrough     interface{} `json:"passThrough,omitempty"`
}


type GetTextAuditingJobResult struct {
	Code               int                   `json:"code"`
	Message            string                `json:"message"`
	RequestID          string                `json:"requestId"`
	RiskLevel          string                `json:"riskLevel"`
	RiskLabel1         string                `json:"riskLabel1"`
	RiskLabel2         string                `json:"riskLabel2"`
	RiskLabel3         string                `json:"riskLabel3"`
	RiskDescription    string                `json:"riskDescription"`
	RiskDetail         *GetTextRiskDetail           `json:"riskDetail"`
	TokenLabels        *GetTextTokenLabels          `json:"tokenLabels"`
	AuxInfo            *GetTextAuxInfo              `json:"auxInfo"`
	AllLabels          []*GetTextAllLabels          `json:"allLabels"`
	BusinessLabels     []*GetTextBusinessLabels     `json:"businessLabels"`
	TokenProfileLabels []*GetTextTokenProfileLabels `json:"tokenProfileLabels,omitempty"`
	TokenRiskLabels    []*GetTextTokenRiskLabels    `json:"tokenRiskLabels,omitempty"`
}
type GetTextWords struct {
	Word     string `json:"word,omitempty"`
	Position []int  `json:"position,omitempty"`
}
type GetTextMatchedLists struct {
	Name  string   `json:"name,omitempty"`
	Words []*GetTextWords `json:"words,omitempty"`
}
type GetTextRiskSegments struct {
	Segment  string `json:"segment,omitempty"`
	Position []int  `json:"position,omitempty"`
}
type GetTextRiskDetail struct {
	MatchedLists []*GetTextMatchedLists `json:"matchedLists,omitempty"`
	RiskSegments []*GetTextRiskSegments `json:"riskSegments,omitempty"`
}
type GetTextUGCAccountRisk struct {
	SexyRiskTokenid float64 `json:"sexy_risk_tokenid,omitempty"`
}
type GetTextTokenLabels struct {
	UGCAccountRisk *GetTextUGCAccountRisk `json:"UGC_account_risk,omitempty"`
}
type GetTextContactResult struct {
	ContactString string `json:"contactString,omitempty"`
	ContactType   int    `json:"contactType,omitempty"`
}
type GetTextAuxInfo struct {
	FilteredText  string           `json:"filteredText,omitempty"`
	PassThrough   interface{}     `json:"passThrough,omitempty"`
	ContactResult []*GetTextContactResult `json:"contactResult,omitempty"`
}
type GetTextAllLabels struct {
	RiskLabel1      string      `json:"riskLabel1"`
	RiskLabel2      string      `json:"riskLabel2"`
	RiskLabel3      string      `json:"riskLabel3"`
	RiskDescription string      `json:"riskDescription"`
	Probability     float64     `json:"probability"`
	RiskDetail      *GetTextRiskDetail `json:"riskDetail"`
	RiskLevel       string      `json:"riskLevel"`
}
type GetTextBusinessDetail struct {}
type GetTextBusinessLabels struct {
	BusinessLabel1      string          `json:"businessLabel1"`
	BusinessLabel2      string          `json:"businessLabel2"`
	BusinessLabel3      string          `json:"businessLabel3"`
	BusinessDescription string          `json:"businessDescription"`
	Probability         float64         `json:"probability"`
	BusinessDetail      *GetTextBusinessDetail `json:"businessDetail"`
}
type GetTextTokenProfileLabels struct {
	Label1      string `json:"label1,omitempty"`
	Label2      string `json:"label2,omitempty"`
	Label3      string `json:"label3,omitempty"`
	Description string `json:"description,omitempty"`
	Timestamp   int64  `json:"timestamp,omitempty"`
}
type GetTextTokenRiskLabels struct {
	Label1      string `json:"label1,omitempty"`
	Label2      string `json:"label2,omitempty"`
	Label3      string `json:"label3,omitempty"`
	Description string `json:"description,omitempty"`
	Timestamp   int64  `json:"timestamp,omitempty"`
}

// PutTextAuditingJobResult is the result of PutTextAuditingJob
type PutTextAuditingJobResult GetTextAuditingJobResult

//TextAuditingJob 文本审核-创建任务
func (s *SmService) TextAuditingJob(ctx context.Context, opt *TextAuditingJobOptions) (*PutTextAuditingJobResult, *Response, error) {
	var res PutTextAuditingJobResult
	sendOpt := sendOptions{
		baseURL: s.client.BaseURL.SMURL,
		uri:     "/text/v4",
		method:  http.MethodPost,
		body:    opt,
		result:  &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}

// ===================================数美文本审核结束======================================


// ===================================数美图片审核开始======================================

type ImageRecognitionOptions struct {
	AccessKey string `json:"accessKey"`
	AppID string `json:"appId"`
	EventID string `json:"eventId"`
	Type string `json:"type"`
	Data *ImageData `json:"data"`
}

type ImageData struct {
	TokenID string `json:"tokenId"`
	Img string `json:"img"`
	Room string `json:"room,omitempty"`
	IP string `json:"ip,omitempty"`
	DeviceID string `json:"deviceId,omitempty"`
	MaxFrame int `json:"maxFrame,omitempty"`
	Interval int `json:"interval,omitempty"`
	Extra *ImageExtra `json:"extra,omitempty"`
}

type ImageExtra struct {
	IsIgnoreTLS bool `json:"isIgnoreTls,omitempty"`
}

type GetImgAuditingJobResult struct {
	Code int `json:"code"`
	Message string `json:"message"`
	RequestID string `json:"requestId"`
	RiskLevel string `json:"riskLevel"`
	RiskLabel1 string `json:"riskLabel1"`
	RiskLabel2 string `json:"riskLabel2"`
	RiskLabel3 string `json:"riskLabel3"`
	RiskDescription string `json:"riskDescription"`
	RiskDetail *GetImgRiskDetail `json:"riskDetail"`
	AuxInfo *GetImgAuxInfo `json:"auxInfo"`
	AllLabels []*GetImgAllLabels `json:"allLabels"`
	BusinessLabels []*GetImgBusinessLabels `json:"businessLabels,omitempty"`
}
type GetImgFaces struct {
	Name string `json:"name,omitempty"`
	Location []int `json:"location,omitempty"`
	Probability float64 `json:"probability,omitempty"`
}
type GetImgObjects struct {
	Name string `json:"name,omitempty"`
	Location []int `json:"location,omitempty"`
	Probability float64 `json:"probability,omitempty"`
}
type GetImgWords struct {
	Word string `json:"word"`
	Position []int `json:"position"`
}
type GetImgMatchedLists struct {
	Name string `json:"name,omitempty"`
	Words []*GetImgWords `json:"words,omitempty"`
}
type GetImgRiskSegments struct {
	Segment string `json:"segment,omitempty"`
	Position []int `json:"position,omitempty"`
}
type GetImgOcrText struct {
	Text string `json:"text,omitempty"`
	MatchedLists []*GetImgMatchedLists `json:"matchedLists,omitempty"`
	RiskSegments *GetImgRiskSegments `json:"riskSegments,omitempty"`
}
type GetImgRiskDetail struct {
	Faces []*GetImgFaces `json:"faces,omitempty"`
	Objects []*GetImgObjects `json:"objects,omitempty"`
	OcrText *GetImgOcrText `json:"ocrText,omitempty"`
}
type GetImgTypeVersion struct {
	POLITICS string `json:"POLITICS,omitempty"`
	VIOLENCE string `json:"VIOLENCE,omitempty"`
	BAN string `json:"BAN,omitempty"`
	PORN string `json:"PORN,omitempty"`
	MINOR string `json:"MINOR,omitempty"`
	AD string `json:"AD,omitempty"`
	SPAM string `json:"SPAM,omitempty"`
	LOGO string `json:"LOGO,omitempty"`
	STAR string `json:"STAR,omitempty"`
	OCR string `json:"OCR,omitempty"`
	IMGTEXT string `json:"IMGTEXT,omitempty"`
	SCREEN string `json:"SCREEN,omitempty"`
	SCENCE string `json:"SCENCE,omitempty"`
	QR string `json:"QR,omitempty"`
	FACE string `json:"FACE,omitempty"`
	QUALITY string `json:"QUALITY,omitempty"`
	PORTRAIT string `json:"PORTRAIT,omitempty"`
	ANIMAL string `json:"ANIMAL,omitempty"`
}

type GetImgAuxInfo struct {
	Segments int `json:"segments"`
	TypeVersion *GetImgTypeVersion `json:"typeVersion"`
	ErrorCode int `json:"errorCode,omitempty"`
	PassThrough interface{} `json:"passThrough,omitempty"`
}
type GetImgAllLabels struct {
	RiskLabel1 string `json:"riskLabel1"`
	RiskLabel2 string `json:"riskLabel2"`
	RiskLabel3 string `json:"riskLabel3"`
	RiskDescription string `json:"riskDescription"`
	Probability float64 `json:"probability,omitempty"`
	RiskDetail *GetImgRiskDetail `json:"riskDetail"`
}
type GetImgBusinessDetail struct {
	Name string `json:"name,omitempty"`
	Probability float64 `json:"probability,omitempty"`
	FaceRatio float64 `json:"face_ratio,omitempty"`
	FaceNum int `json:"face_num,omitempty"`
	Location []int `json:"location,omitempty"`
	PersonNum int `json:"person_num,omitempty"`
	PersonRatio float64 `json:"person_ratio,omitempty"`
}
type GetImgBusinessLabels struct {
	BusinessLabel1 string `json:"businessLabel1,omitempty"`
	BusinessLabel2 string `json:"businessLabel2,omitempty"`
	BusinessLabel3 string `json:"businessLabel3,omitempty"`
	BusinessDescription string `json:"businessDescription,omitempty"`
	BusinessDetail *GetImgBusinessDetail `json:"businessDetail,omitempty"`
	Probability float64 `json:"probability,omitempty"`
	ConfidenceLevel float64 `json:"confidenceLevel,omitempty"`
}

type PutImgAuditingJobResult GetImgAuditingJobResult

//ImageSingleAuditing 单张图片同步
func (s *SmService) ImageSingleAuditing(ctx context.Context, opt *ImageRecognitionOptions) (*PutImgAuditingJobResult, *Response, error) {
	var res PutImgAuditingJobResult
	sendOpt := sendOptions{
		baseURL:  s.client.BaseURL.SMURL,
		uri:      "/image/v4",
		method:   http.MethodPost,
		body:    opt,
		result:   &res,
	}
	resp, err := s.client.send(ctx, &sendOpt)
	return &res, resp, err
}
// ===================================数美图片审核结束======================================
