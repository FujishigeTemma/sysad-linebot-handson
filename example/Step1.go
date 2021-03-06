// #1 やまびこの実装
package main

// 利用したい外部のコードを読み込む
import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/v7/linebot"
)

const (
	verifyToken = "00000000000000000000000000000000"
)

// main関数は最初に呼び出されることがGo言語の仕様として決まっている
func main() {
	// LINEのAPIを利用する設定
	bot, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_ACCESS_TOKEN"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// LINEサーバからのリクエストを受け取ったときの処理
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		log.Println("Accessed")

		// リクエストを扱いやすい形に変換する
		events, err := bot.ParseRequest(req)
		switch err {
		case nil:
		// 変換に失敗したとき
		case linebot.ErrInvalidSignature:
			log.Println("ParseRequest error:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		default:
			log.Println("ParseRequest error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// LINEサーバから来たメッセージによって行う処理を変える
		for _, event := range events {
			// LINEサーバからのverify時は何もしない
			if event.ReplyToken == verifyToken {
				return
			}

			switch event.Type {
			// メッセージが来たとき
			case linebot.EventTypeMessage:
				// 返信を生成する
				replyMessage := getReplyMessage(event)
				// 生成した返信を送信する
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
					log.Print(err)
				}
			// それ以外のとき
			default:
				continue
			}
		}
	})

	// LINEサーバからのリクエストを受け取るプロセスを起動
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

const helpMessage = `使い方
テキストメッセージ:
	やまびこを返すよ！
スタンプ:
	スタンプの情報を答えるよ！
それ以外:
	それ以外にはまだ対応してないよ！ごめんね...`

// 返信を生成する
func getReplyMessage(event *linebot.Event) string {
	// 来たメッセージの種類によって行う処理を変える
	switch message := event.Message.(type) {
	// テキストメッセージが来たとき
	case *linebot.TextMessage:
		return message.Text

	// スタンプが来たとき
	case *linebot.StickerMessage:
		return fmt.Sprintf("sticker id is %v, stickerResourceType is %v", message.StickerID, message.StickerResourceType)
	// それ以外のとき
	default:
		return helpMessage
	}
}
