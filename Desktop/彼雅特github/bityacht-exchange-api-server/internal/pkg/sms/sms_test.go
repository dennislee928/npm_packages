package sms

import (
	"context"
	"fmt"
	"testing"
	"time"

	"bityacht-exchange-api-server/configs"

	"github.com/spf13/viper"
)

func TestSend(t *testing.T) {
	viper.AddConfigPath("../../../configs")
	configs.Init()
	configs.Config.SMS.Enable = true
	Init()

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := Send(ctx, Message{
		Phone:        "0912345678",
		Message:      "test",
		ReceiverName: "TestName", // Optional
		CallbackURL:  "",         // Optional, Callback URL to receive the delivery receipt of the message.
		DeliveryTime: nil,        // Optional
		ValidateTime: nil,        // Optional
	}); err != nil {
		t.Error(err)
	}

	// Test only one message in batch
	if err := SendBatch(ctx, []Message{{
		Phone:        "0912345678",
		Message:      "test-batch(1)",
		ReceiverName: "TestName", // Optional
		CallbackURL:  "",         // Optional, Callback URL to receive the delivery receipt of the message.
		DeliveryTime: nil,        // Optional
		ValidateTime: nil,        // Optional
	}}); err != nil {
		t.Error(err)
	}

	messages := make([]Message, 10)
	for i := 0; i < len(messages); i++ {
		messages[i] = Message{
			Phone:        "0912345678",
			Message:      fmt.Sprintf("test-batch(%d/%d)", i, len(messages)),
			ReceiverName: "TestName", // Optional
			CallbackURL:  "",         // Optional, Callback URL to receive the delivery receipt of the message.
			DeliveryTime: nil,        // Optional
			ValidateTime: nil,        // Optional
		}
	}
	if err := SendBatch(ctx, messages); err != nil {
		t.Error(err)
	}
}
