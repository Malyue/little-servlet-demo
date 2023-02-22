package chat

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"
)

func ParseMessage(message []byte) (sendTime time.Time, receiver int, content string, err error) {
	params := make(map[string]interface{})
	err = json.Unmarshal(message, &params)
	if err != nil {
		return
	}

	receiver, _ = strconv.Atoi(fmt.Sprint(params["receiver"]))
	content = fmt.Sprint(params["content"])
	timeString := params["sendTime"].(float64)
	timeStamp := int64(timeString * math.Pow10(0))
	sendTime = time.Unix(timeStamp, 0)
	return
}
