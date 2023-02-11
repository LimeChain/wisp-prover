package handlers

import (
	"github.com/LimeChain/crc-prover/pkg/app/configs"
	"github.com/LimeChain/crc-prover/pkg/app/rest"
	"github.com/LimeChain/crc-prover/pkg/log"
	"github.com/LimeChain/crc-prover/pkg/proof"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"net/http"
)

// ZKHandler is handler for zkp operations
type ZKHandler struct {
	ProverConfig configs.ProverConfig
}

// GenerateReq is request for proof generation
type GenerateReq struct {
	Circuit string         `json:"circuit"`
	Inputs  proof.ZKInputs `json:"inputs"`
}

// VerifyReq is request for proof verification
type VerifyReq struct {
	Circuit string          `json:"circuit"`
	ZKP     proof.FullProof `json:"zkp"`
}

// VerifyResp is response for proof verification
type VerifyResp struct {
	Valid bool `json:"valid"`
}

// NewZKHandler creates new instance of handler
func NewZKHandler(proverConfig configs.ProverConfig) *ZKHandler {
	return &ZKHandler{
		proverConfig,
	}
}

// GenerateProof is a handler for proof generation
// POST /api/v1/proof/generate
func (h *ZKHandler) GenerateProof(w http.ResponseWriter, r *http.Request) {
	var req GenerateReq
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		rest.ErrorJSON(w, r, http.StatusBadRequest, err, "invalid request", 0)
		return
	}
	log.WithContext(r.Context()).Debugw("Proof generation request", "inputs", req, "type", req.Inputs)

	prover, err := h.parseRequestedCircuit(req.Circuit)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusBadRequest, err, req.Circuit, 0)
		return
	}
	generatedProof, err := prover.GenerateProof(r.Context(), req.Inputs)
	if err != nil {
		rest.ErrorJSON(w, r, http.StatusInternalServerError, errors.New("Internal Server Error"), "", 0)
		return
	}

	render.JSON(w, r, generatedProof)
}

func (h *ZKHandler) parseRequestedCircuit(circuit string) (*proof.Proof, error) {
	_, exists := h.ProverConfig.SupportedCircuits[circuit]
	if !exists {
		return nil, errors.New("unsupported circuit provided")
	}

	return proof.NewProof(h.ProverConfig, circuit, circuit), nil
}

// VerifyProof is a handler for zkp verification
// POST /api/v1/proof/verify
func (h *ZKHandler) VerifyProof(w http.ResponseWriter, r *http.Request) {

	valid := false

	var req VerifyReq
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		rest.ErrorJSON(w, r, http.StatusBadRequest, err, "can't bind request", 0)
		return
	}

	log.WithContext(r.Context()).Debugw("Proof verification request", "inputs", req)
	//err := proof.VerifyZkProof(r.Context(), circuitPath, &req.ZKP)
	//if err == nil {
	//	valid = true
	//}

	render.JSON(w, r, VerifyResp{Valid: valid})
}
