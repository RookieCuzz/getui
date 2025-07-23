package report

import (
	"context"
	"github.com/dacker-soul/getui/auth"
	"github.com/dacker-soul/getui/publics"
	"testing"
)

var (
	Conf = publics.GeTuiConfig{
		AppId:        "XHy5KG2B6v6bfeU9inrOV4",
		AppSecret:    "q8g91dL0jp7jS7auDRIjX9",
		AppKey:       "3RNEaOqHTz83XL6lOAVYn7",
		MasterSecret: "u6WNgxgfNh7d6lvSfhlg12",
	}
	Cid   = "" // clientId
	Ctx   = context.Background()
	Token = ""
	Alias = "test_user"
)

func TestGetTaskResultWithAction(t *testing.T) {
	token, err := auth.GetToken(Ctx, Conf)
	if err != nil {
		t.Error(err)
	}
	var param = TaskResultParam{
		TaskId:  "RASA_0723_bd5732498608131111fe20881f8bb689",
		Actions: []string{"12312312", "1312312"},
	}
	resp, err := GetTaskResult(context.Background(), Conf, token.Data.Token, &param)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)
}

func TestGetPushDailyStats(t *testing.T) {
	token, err := auth.GetToken(Ctx, Conf)
	if err != nil {
		t.Error(err)
	}
	var param = PushDailyStatsParam{
		//Date: time.Now().Format("2006-01-02"),
		Date: "2025-07-23",
	}
	resp, err := GetPushDailyStats(context.Background(), Conf, token.Data.Token, param)
	if err != nil {
		t.Error(err)
	}
	t.Log(resp)

}

func TestGetUniPushBalance(t *testing.T) {

	token, err := auth.GetToken(Ctx, Conf)
	if err != nil {
		t.Error(err)
	}
	balance, err := GetUniPushBalance(context.Background(), Conf, token.Data.Token)
	if err != nil {
		t.Error(err)
	}
	t.Log(balance)
}
