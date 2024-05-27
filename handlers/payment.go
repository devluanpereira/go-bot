package handlers

import (
	"context"
	"fmt"

	"github.com/mercadopago/sdk-go/pkg/config"
	"github.com/mercadopago/sdk-go/pkg/payment"
	"github.com/skip2/go-qrcode"
)

func CreatePixPayment() ([]byte, string, error) {

	config, err := config.New("APP_USR-5571155767002487-032220-244efb294b0bf14a9c8772ff136b60c5-1094153639")
	if err != nil {
		return nil, "", err
	}

	client := payment.NewClient(config)

	request := payment.Request{
		TransactionAmount: 100,
		Description:       "My product",
		PaymentMethodID:   "pix",
		Payer: &payment.PayerRequest{
			Email: "luan23107@gmail.com",
			Identification: &payment.IdentificationRequest{
				Type:   "CPF",
				Number: "61511437367",
			},
		},
	}

	resource, err := client.Create(context.Background(), request)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", err
	}

	qrCode := resource.PointOfInteraction.TransactionData.QRCode
	if qrCode == "" {
		return nil, "", fmt.Errorf("QR code não gerado")
	}

	// Gerar a imagem do QR code
	qrCodeImage, err := qrcode.Encode(qrCode, qrcode.Medium, 256)
	if err != nil {
		return nil, "", fmt.Errorf("erro ao gerar a imagem do QR code: %v", err)
	}

	return qrCodeImage, qrCode, nil
}
