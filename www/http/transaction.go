package http

import (
	"encoding/hex"
	"net/http"

	"github.com/gorilla/mux"
	pactus "github.com/pactus-project/pactus/www/grpc/gen/go"
)

func (s *Server) GetTransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := hex.DecodeString(vars["id"])
	if err != nil {
		s.writeError(w, err)
		return
	}

	res, err := s.transaction.GetTransaction(s.ctx,
		&pactus.GetTransactionRequest{
			Id:        id,
			Verbosity: pactus.TransactionVerbosity_TRANSACTION_DATA,
		},
	)
	if err != nil {
		s.writeError(w, err)
		return
	}

	tm := newTableMaker()
	txToTable(res.Transaction, tm)
	s.writeHTML(w, tm.html())
}

func txToTable(trx *pactus.TransactionInfo, tm *tableMaker) {
	if trx == nil {
		return
	}
	tm.addRowTxID("ID", trx.Id)
	tm.addRowBytes("Data", trx.Data)
	tm.addRowInt("Version", int(trx.Version))
	tm.addRowBytes("Stamp", trx.Stamp)
	tm.addRowInt("Sequence", int(trx.Sequence))
	tm.addRowInt("Fee", int(trx.Fee))
	tm.addRowString("Memo", trx.Memo)
	switch trx.Type {
	case pactus.PayloadType_SEND_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Send).Send
		tm.addRowString("Payload type", "Send")
		tm.addRowAccAddress("Sender", pld.Sender)
		tm.addRowAccAddress("Receiver", pld.Receiver)
		tm.addRowAmount("Amount", pld.Amount)

	case pactus.PayloadType_BOND_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Bond).Bond
		tm.addRowString("Payload type", "Bond")
		tm.addRowAccAddress("Sender", pld.Sender)
		tm.addRowValAddress("Receiver", pld.Receiver)
		tm.addRowAmount("Stake", pld.Stake)

	case pactus.PayloadType_SORTITION_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Sortition).Sortition
		tm.addRowString("Payload type", "Sortition")
		tm.addRowValAddress("Address", pld.Address)
		tm.addRowBytes("Proof", pld.Proof)

	case pactus.PayloadType_UNBOND_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Unbond).Unbond
		tm.addRowString("Payload type", "Unbond")
		tm.addRowValAddress("Validator", pld.Validator)

	case pactus.PayloadType_WITHDRAW_PAYLOAD:
		pld := trx.Payload.(*pactus.TransactionInfo_Withdraw).Withdraw
		tm.addRowString("Payload type", "Withdraw")
		tm.addRowValAddress("Sender", pld.From)
		tm.addRowAccAddress("Receiver", pld.To)
		tm.addRowAmount("Amount", pld.Amount)
	}
	if trx.PublicKey != "" {
		tm.addRowString("PublicKey", trx.PublicKey)
	}
	if trx.Signature != nil {
		tm.addRowBytes("Signature", trx.Signature)
	}
}
