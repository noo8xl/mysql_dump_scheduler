package scheduler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// @SendFileToTelegram -> send the file to the telegram
func (s *SchedulerService) sendFileToTelegram(file *os.File) error {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("chat_id", s.TelegramConfig.ChatId); err != nil {
		return fmt.Errorf("error writing chat_id: %v", err)
	}

	part, err := writer.CreateFormFile("document", filepath.Base(file.Name()))
	if err != nil {
		return fmt.Errorf("error creating form file: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("error closing writer: %v", err)
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", s.TelegramConfig.Token)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	log.Printf("Sender response is : %v", string(r))
	return nil
}

// func sendFile(dto TelegramDto) error {

// 	chatID, err := strconv.ParseInt(dto.ChatId, 10, 64)
// 	if err != nil {
// 		log.Printf("Failed to parse TELEGRAM_CHAT_ID: %v", err)
// 		return err
// 	}
// 	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
// 	if botToken == "" {
// 		log.Println("TELEGRAM_BOT_TOKEN is not set")
// 		return err
// 	}

// 	bot, err := tgbotapi.NewBotAPI(botToken)
// 	if err != nil {
// 		log.Printf("Failed to create Telegram bot: %v", err)
// 		return err
// 	}

// 	// Create a file upload message instead of a text message
// 	fileDoc := tgbotapi.NewDocument(chatID, tgbotapi.FileReader{
// 		Name:   dto.File.Name(),
// 		Reader: dto.File,
// 	})
// 	fileDoc.Caption = "Database dump file"

// 	_, err = bot.Send(fileDoc)
// 	if err != nil {
// 		log.Printf("Failed to send Telegram file: %v", err)
// 		return err
// 	}

// 	log.Printf("Telegram file sent successfully")
// 	return nil
// }
