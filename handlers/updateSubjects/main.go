package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kinghunter58/jwe"

	qs "gitlab.com/zapochvam-ei-sq/plannerix-backend/models/QS"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/errors"
	"gitlab.com/zapochvam-ei-sq/plannerix-backend/models/profile"
)

type Request struct {
	Token    string    `json:"token"`
	Subjects []Subject `json:"subjects"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Subject struct {
	id   string `json:"id"`
	name string `json:"name"`
}

func handler(ctx context.Context, req interface{}) (qs.Response, error) {
	body := Request{}
	err := qs.GetBody(req, &body)

	if err != nil {
		return qs.NewError(errors.LambdaError.Error(), -1)
	}

	key, err := jwe.GetPrivateKeyFromEnv("RSAPRIVATEKEY")

	if err != nil {
		return qs.NewError("Internal Server Error", 6)
	}

	p := profile.Payload{}
	err = jwe.ParseEncryptedToken(body.Token, key, &p)
	if err != nil {
		return qs.NewError("Internal Server Error", 6) // Fix output
	}

	fmt.Println(p)
	res := Response{
		Success: true,
		Message: "schedule updated successfully",
	}
	return qs.NewResponse(200, res)
}

func main() {
	lambda.Start(handler)
}
