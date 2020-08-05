package ifua

import (
	"context"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"tapera.integrasi/api/util/httpx"

	mic "tapera.integrasi/grpc/client/mitraintegrasi/v1"
)

type (
	InvsFund struct {
		SaCode     string `json:"sa_code"`
		SID        string `json:"sid"`
		InvstrName string `json:"investor_name"`
		ClientCode string `json:"client_code"`
	}
)

func (c *Controller) sendDataIfuaReq(w http.ResponseWriter, r *http.Request) {
	ext := httpx.New(w, r)
	ctx := context.Background()
	var param []InvsFund
	var createData mic.ListInvsFundRegData
	if err := ext.BindJSON(&param); err != nil {
		ext.JSONerr(http.StatusBadRequest, "Invalid request payload")
		return
	}

	for _, v := range param {
		createData.InvsFundRegData = append(createData.InvsFundRegData, &mic.InvsFundRegData{
			SaCode:     v.SaCode,
			SID:        v.SID,
			InvstrName: v.InvstrName,
			ClientCode: v.ClientCode,
		})
	}

	conn, err := grpc.Dial("10.172.24.63:80", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	client := mic.NewIfuaRegToKSEIClient(conn)

	stream, err := client.SendInvsFundReg(ctx)
	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}

	if err := stream.Send(&createData); err != nil {
		log.Fatalf("can not send %v", err)
	}
	if _, err := stream.CloseAndRecv(); err != nil {
		log.Println(err)
	}

	ext.JSON(http.StatusOK, "done")
}
