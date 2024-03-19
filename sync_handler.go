/**
 * @author Jose Nidhin
 */
package main

import (
	"log/slog"
	"net/http"
)

func SyncHandler(redisClient RedisClient, logger *slog.Logger) http.HandlerFunc {
	fn := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		reqId := GetCtxString(ctx, ReqIdCtxKey)
		logger.Info("New sync request", slog.String("requestId", reqId))

		sub := redisClient.Subscribe(ctx, ChannelName)
		defer sub.Close()
		ch := sub.Channel()

		for chMsg := range ch {
			logger.Info("Message received",
				slog.String("channelName", chMsg.Channel),
				slog.String("msgPayload", chMsg.Payload))

			msg, err := NewMessageFromJSON([]byte(chMsg.Payload))
			if err != nil {
				logger.Error("Error parsing message from channel", slog.Any("error", err))
				continue
			}

			if msg.RequestId == reqId {
				res.Header().Set("Content-Type", "application/json; charset=utf-8")
				res.Write(msg.Data)
				return
			}
		}

		res.WriteHeader(http.StatusBadRequest)
	}

	return http.HandlerFunc(fn)
}
