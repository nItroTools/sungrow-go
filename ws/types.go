package ws

type Keys map[string]string

type RequestConnect struct {
	Lang    string `json:"lang"`
	Token   string `json:"token"`
	Service string `json:"service"`
}
type ResponseConnect struct {
	ResultCode int    `json:"result_code"`
	ResultMsg  string `json:"result_msg"`
	ResultData struct {
		Service     string
		Token       string
		Uid         int
		TipsDisable int `json:"tips_disable"`
	} `json:"result_data"`
}

type RequestReal struct {
	Lang       string `json:"lang"`
	Token      string `json:"token"`
	DevId      string `json:"dev_id"`
	Service    string `json:"service"`
	Time123456 int64  `json:"time123456"`
}
type ResponseReal struct {
	ResultCode int    `json:"result_code"`
	ResultMsg  string `json:"result_msg"`
	ResultData struct {
		Service string
		List    []struct {
			DataName  string `json:"data_name"`
			DataValue string `json:"data_value"`
			DataUnit  string `json:"data_unit"`
		}
		Count int
	} `json:"result_data"`
}
