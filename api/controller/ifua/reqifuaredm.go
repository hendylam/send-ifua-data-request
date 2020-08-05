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
	RedmPay struct {
		InvsFundAC  string `json:"investor_fund_ac"`
		BankBicCd   string `json:"bank_bic_code"`
		BankBiMemCd string `json:"bank_bic_member_code"`
		BankName    string `json:"bank_name"`
		BankCntry   string `json:"bank_country"`
		BankBranch  string `json:"bank_branch"`
		AcCcy       string `json:"ac_ccy"`
		BankAcNo    string `json:"bank_ac_no"`
		BankAcName  string `json:"bank_ac_name"`
	}
)

func (c *Controller) sendDataRedmReq(w http.ResponseWriter, r *http.Request) {
	ext := httpx.New(w, r)
	ctx := context.Background()
	var param []RedmPay
	var createData mic.ListRedmPayData
	if err := ext.BindJSON(&param); err != nil {
		ext.JSONerr(http.StatusBadRequest, "Invalid request payload")
		return
	}

	for _, v := range param {
		createData.RedmPayRegData = append(createData.RedmPayRegData, &mic.RedmPayRegData{
			InvsFundAC:  v.InvsFundAC,
			BankBicCd:   v.BankBicCd,
			BankBiMemCd: v.BankBiMemCd,
			BankName:    v.BankName,
			BankCntry:   v.BankCntry,
			BankBranch:  v.BankBranch,
			AcCcy:       v.AcCcy,
			BankAcNo:    v.BankAcNo,
			BankAcName:  v.BankAcName,
		})
	}

	conn, err := grpc.Dial("10.172.24.63:80", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	client := mic.NewIfuaRegToKSEIClient(conn)

	stream, err := client.SendRedmPayReg(ctx)
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
