package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"

	"sipecut/models"
)

// SendWhatsAppNotification templates and sends a WhatsApp message
func SendWhatsAppNotification(db *gorm.DB, user *models.User, templateKey string, placeholders map[string]string) error {
	// 1. Fetch template from AppSettings
	var setting models.AppSetting
	message := ""
	if err := db.Where("key = ?", templateKey).First(&setting).Error; err == nil {
		message = setting.Value
	} else {
		// Default fallbacks if settings not found
		switch templateKey {
		case "whatsapp_template_leave_new":
			message = "Halo Admin, ada pengajuan cuti baru dari {{NAMA}} (NIP: {{NIP}})."
		case "whatsapp_template_leave_update":
			message = "Halo {{NAMA}}, status pengajuan cuti Anda telah diupdate ke: {{STATUS}}."
		case "whatsapp_template_change_new":
			message = "Halo Admin, ada pengajuan perubahan data baru dari {{NAMA}}."
		case "whatsapp_template_change_update":
			message = "Halo {{NAMA}}, pengajuan perubahan data Anda telah diupdate ke: {{STATUS}}."
		default:
			message = "Notifikasi SIPECUT: Status pengajuan Anda telah diperbarui."
		}
	}

	// 2. Replace placeholders
	for key, val := range placeholders {
		message = strings.ReplaceAll(message, "{{"+key+"}}", val)
	}

	// 3. Check endpoint config
	apiKey := os.Getenv("WHATSAPP_GATEWAY_API_KEY")
	gatewayURL := os.Getenv("WHATSAPP_GATEWAY_URL")
	if gatewayURL == "" {
		gatewayURL = "https://api.fonnte.com/send" // Default to Fonnte
	}

	status := "Simulated"
	noHP := user.NoHP
	if noHP == "" {
		noHP = "08XXXXXXXXXX"
	}

	var err error
	if apiKey != "" && user.NoHP != "" {
		// Send HTTP request to WhatsApp gateway (e.g. Fonnte)
		status, err = sendGatewayRequest(gatewayURL, apiKey, noHP, message)
		if err != nil {
			status = "Failed"
			fmt.Printf("Error sending WhatsApp message: %v\n", err)
		}
	} else {
		// Simulate
		fmt.Printf("\n--- [SIMULATOR WHATSAPP NOTIFICATION] ---\nTo: %s (%s)\nMessage: %s\n-----------------------------------------\n\n", user.NIP, noHP, message)
	}

	// 4. Save notification log
	logEntry := models.NotificationLog{
		UserID:    user.ID,
		Channel:   "WhatsApp",
		Message:   message,
		Status:    status,
		SentAt:    time.Now(),
	}
	db.Create(&logEntry)

	return err
}

func sendGatewayRequest(url, apiKey, target, message string) (string, error) {
	// Fonnte request format:
	// POST multipart/form-data or application/json
	payload := map[string]string{
		"target":  target,
		"message": message,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "Failed", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "Failed", err
	}

	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "Failed", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "Sent", nil
	}

	body, _ := io.ReadAll(resp.Body)
	return "Failed", fmt.Errorf("gateway returned status %d: %s", resp.StatusCode, string(body))
}
