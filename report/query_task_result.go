package report

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dacker-soul/getui/publics"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 通用参数结构体：含 taskId 和可选 actions
type TaskResultParam struct {
	TaskId  string   `json:"taskId"`
	Actions []string `json:"actions,omitempty"` // 可为空
}

type PushDailyStatsParam struct {
	Date             string `json:"date"`
	NeedGetuiByBrand bool   `json:"needGetuiByBrand,omitempty"` // 可为空
}

type PushReportResponse struct {
	Msg  string                            `json:"msg"`
	Code int                               `json:"code"`
	Data map[string]PushTaskStatisticsData `json:"data"` // key 为 taskId，例如 RASA_0723_bd5732498608131111fe20881f8bb689
}

type PushTaskStatisticsData struct {
	Total PushChannelStats `json:"total"`
	GT    PushChannelStats `json:"gt"`
	APN   PushChannelStats `json:"apn,omitempty"` // iOS APN 推送数据，可能没有
}

type PushChannelStats struct {
	MsgNum     int `json:"msg_num,omitempty"` // 可下发数（仅 total 中有）
	TargetNum  int `json:"target_num"`        // 实际下发数
	ReceiveNum int `json:"receive_num"`       // 接收数
	DisplayNum int `json:"display_num"`       // 展示数
	ClickNum   int `json:"click_num"`         // 点击数
}

// 响应
type PushDailyStatsResp struct {
	Code int                               `json:"code"`
	Msg  string                            `json:"msg"`
	Data map[string]PushTaskStatisticsData `json:"data"`
}

func GetTaskResult(
	ctx context.Context,
	config publics.GeTuiConfig,
	token string,
	param *TaskResultParam,
) (*PushReportResponse, error) {

	if param.TaskId == "" {
		return nil, errors.New("taskId 不能为空")
	}

	// 基础 URL
	baseUrl := fmt.Sprintf("%s%s/report/push/task/%s", publics.ApiUrl, config.AppId, param.TaskId)

	// 如果有自定义事件，拼接 query 参数
	if len(param.Actions) > 0 {
		actionStr := strings.Join(param.Actions, ",")
		baseUrl += "?actionIdList=" + url.QueryEscape(actionStr)
	}

	fmt.Println("请求地址:", baseUrl)

	// 构造请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("构造请求失败: %w", err)
	}
	req.Header.Set("token", token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(bodyBytes))
	}

	var result PushReportResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

// 请求函数
func GetPushDailyStats(
	ctx context.Context,
	config publics.GeTuiConfig,
	token string,
	param PushDailyStatsParam,
) (*PushDailyStatsResp, error) {
	if param.Date == "" {
		return nil, errors.New("date不能为空")
	}

	baseUrl := fmt.Sprintf("%s%s/report/push/date/%s", publics.ApiUrl, config.AppId, param.Date)

	// 添加query参数
	if param.NeedGetuiByBrand {
		baseUrl += "?needGetuiByBrand=true"
	}

	fmt.Println("请求地址:", baseUrl)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("构造请求失败: %w", err)
	}

	req.Header.Set("token", token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(bodyBytes))
	}

	var result PushDailyStatsResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}

type UniPushBalanceResp struct {
	Code int                          `json:"code"`
	Msg  string                       `json:"msg"`
	Data map[string]PushCountChannels `json:"data"`
}

// 推送限制各通道数据
type PushCountChannels map[string]PushCountLimit

// 推送限制详情结构
type PushCountLimit struct {
	PushNum   *int64 `json:"push_num,omitempty"`   // vv才有，表示请求量
	TotalNum  string `json:"total_num"`            // 总推送量，注意类型示例中有string和number
	RemainNum *int64 `json:"remain_num,omitempty"` // 剩余推送量，部分厂商有
	Limit     bool   `json:"limit"`                // 是否被限量
}

func GetUniPushBalance(ctx context.Context, config publics.GeTuiConfig, token string) (*UniPushBalanceResp, error) {
	url := fmt.Sprintf("%s%s/report/push/count", publics.ApiUrl, config.AppId)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("构造请求失败: %w", err)
	}

	req.Header.Set("token", token)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("请求失败，状态码: %d，响应: %s", resp.StatusCode, string(body))
	}

	var result UniPushBalanceResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return &result, nil
}
