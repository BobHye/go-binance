package futures

import (
	"github.com/BobHye/binance-go/log"
	"github.com/BobHye/wsc"
)

// WsHandler handle raw websocket message | 处理原始 websocket 消息
type WsHandler func(message []byte)

// ErrHandler handles errors | 处理错误
type ErrHandler func(err error)

// WsConfig webservice configuration | webservice 配置
type WsConfig struct {
	Endpoint string
}

func newWsConfig(endpoint string) *WsConfig {
	return &WsConfig{
		Endpoint: endpoint,
	}
}

var wsServe = func(cfg *WsConfig, handler WsHandler, errHandler ErrHandler) (ws *wsc.Wsc, done chan struct{}, err error) {
	done = make(chan struct{})

	ws = wsc.New(cfg.Endpoint)
	ws.OnConnected(func() {
		if log.Default.OnConnected {
			log.Default.Log("websocket connected")
		}
	})
	ws.OnConnectError(errHandler)
	ws.OnDisconnected(errHandler)
	ws.OnClose(func(code int, text string) {
		if log.Default.OnClose {
			log.Default.Log("websocket closed, code: %d, message: %s", code, text)
		}
	})
	ws.OnSentError(errHandler)
	ws.OnPingReceived(func(appData string) {
		if log.Default.OnPingReceived {
			log.Default.Log("ping received, data: %s", appData)
		}
	})
	ws.OnPongReceived(func(appData string) {
		if log.Default.OnPongReceived {
			log.Default.Log("pong received, data: %s", appData)
		}
	})
	ws.OnTextMessageReceived(handler)
	ws.OnKeepalive(func() {
		if log.Default.OnKeepalive {
			log.Default.Log("keep alive")
		}
	})

	go func() {
		ws.Connect()
		for range done {
			ws.Close()
			return
		}
	}()
	return
}
