/**
 * @author Jose Nidhin
 */
package main

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func CallbackHandler(redisClient RedisClient, logger *slog.Logger) http.HandlerFunc {
	fn := func(res http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		id := chi.URLParam(req, "id")
		logger.Info("New callback request", slog.String("id", id))

		// a quick hack to demo the concept dont do this in
		// productions systems
		body, err := io.ReadAll(req.Body)
		if err != nil {
			logger.Error("Error reading body",
				slog.Any("error", err))
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		msg := Message{
			RequestId: id,
			Data:      body,
		}

		msgBytes, err := msg.Marshal()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = redisClient.Publish(ctx, ChannelName, msgBytes).Err()
		if err != nil {
			logger.Error("Error publishing to channel",
				slog.Any("error", err),
				slog.String("channelName", ChannelName))
		}

		res.WriteHeader(http.StatusOK)
	}

	return http.HandlerFunc(fn)
}
