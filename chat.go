package ginchat

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

/**
{
    "cmd": "online"
}

{
    "cmd":"send",
    "data":{
        "to_user_id": 2,
        "message":"hello 2"
    }
}
*/

var OnlineUser = make(map[int]*websocket.Conn)

type WSReq struct {
	Cmd  string      `json:"cmd"`
	Data interface{} `json:"data"`
}

type WSResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type SendS struct {
	ToUserId int    `json:"to_user_id"`
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
		case "online":
			go Online(conn, message, int(userId))
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

	if OnlineUser[req.Data.ToUserId] == nil {
		WSRespErr(conn, 20003, "not online")
		return
	}

	WSRespSuccess(OnlineUser[req.Data.ToUserId], req.Data.Message)

	WSRespSuccess(conn, "send ok")

}

func Online(conn *websocket.Conn, message []byte, userId int) {

	OnlineUser[userId] = conn
	WSRespSuccess(conn, userId)
}

func WSRespErr(conn *websocket.Conn, code int, msg string) {
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
