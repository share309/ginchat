package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

/**
{
    "cmd":"send",
    "data":{
        "to_user_id": 2,
        "message":"hello 2"
    }
}
*/

var OnlineUser = make(map[uint]*websocket.Conn)

type WSReq struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
}

type WSResp struct {
	Code uint        `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type SendS struct {
	ToUserId uint   `json:"to_user_id"`
	Message  string `json:"message"`
}

type WSSendReq struct {
	Cmd  string `json:"cmd"`
	Data SendS  `json:"data"`
}

func Chat(c *gin.Context) {

	//get middle userId
	userId := c.GetUint("userId")

	// update websocket
	var upgrader = websocket.Upgrader{} // use default options

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("upgrade:", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "websocket fail",
			"data": nil,
		})
		return
	}

	defer conn.Close()

	//online user data
	OnlineUser[userId] = conn

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		var req WSReq

		err = json.Unmarshal(message, &req)
		if err != nil {

			WSRespErr(conn, 1000, "illegal json")
			continue
		}

		switch req.Cmd {
		case "send":
			go Send(conn, message)
		default:
			WSRespErr(conn, 10001, "no function")
		}

	}
}

func Send(conn *websocket.Conn, message []byte) {

	var req WSSendReq
	err := json.Unmarshal(message, &req)
	if err != nil {
		WSRespErr(conn, 20001, "Parsing json failed ")
		return
	}

	if req.Data.ToUserId < 1 {
		WSRespErr(conn, 20002, "not ToUserId params")
		return
	}

	if req.Data.Message == "" {
		WSRespErr(conn, 20004, "not Message params")
		return
	}

	if OnlineUser[req.Data.ToUserId] == nil {
		WSRespErr(conn, 20003, "not online")
		return
	}

	WSRespSuccess(OnlineUser[req.Data.ToUserId], req.Data.Message)

	WSRespSuccess(conn, "send ok")

}

func WSRespErr(conn *websocket.Conn, code uint, msg string) {
	response := WSResp{
		Code: code,
		Msg:  msg,
		Data: nil,
	}

	responseStr, _ := json.Marshal(response)

	err := conn.WriteMessage(websocket.TextMessage, responseStr)
	if err != nil {
		log.Println("write:", err)
	}
}

func WSRespSuccess(conn *websocket.Conn, data interface{}) {
	response := WSResp{
		Code: 0,
		Msg:  "ok",
		Data: data,
	}

	responseStr, _ := json.Marshal(response)

	err := conn.WriteMessage(websocket.TextMessage, responseStr)
	if err != nil {
		log.Println("write:", err)
	}
}
