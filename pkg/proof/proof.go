package proof

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/LimeChain/crc-prover/pkg/app/configs"
	"github.com/iden3/go-rapidsnark/prover"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"time"

	"github.com/LimeChain/crc-prover/pkg/log"
	"github.com/iden3/go-rapidsnark/types"
	"github.com/iden3/go-rapidsnark/verifier"
	"github.com/pkg/errors"
)

const (
	proofDataFormat   string = "2006-01-02-15-04-05.000"
	inputJsonFileName string = "input.json"
	witnessFileName   string = "witness.wtns"
	zkeyFileName      string = "final.zkey"
	vkeyFileName      string = "vkey.json"
)

// Proof is a struct that is able to compute a witness and proof for a given circuit
type Proof struct {
	config            configs.ProverConfig
	circuitBinaryPath string // Path to the circuit binary to be used for witness generation
	zkeyPath          string // Path to the zkey to be used for proof generation
	vkeyPath          string // Path to the verification key to be used for verifying the proof
	dataPath          string // Path to the folder where input.json will be saved, witness and proof will be generated
}

// NewProof assumes that the `binary`, `zkey` and `vkey` are placed in the `./{baseCircuitPaths}/{circuit}/` directory
func NewProof(config configs.ProverConfig, circuitName, binaryName string) *Proof {
	return &Proof{
		config:            config,
		circuitBinaryPath: config.CircuitsBasePath + "/" + circuitName + "/" + binaryName,
		zkeyPath:          config.CircuitsBasePath + "/" + circuitName + "/" + zkeyFileName,
		vkeyPath:          config.CircuitsBasePath + "/" + circuitName + "/" + vkeyFileName,
		dataPath:          config.ProofsBasePath + "/" + circuitName + "/" + time.Now().Format(proofDataFormat),
	}
}

// ZKInputs are inputs for proof generation
type ZKInputs map[string]interface{}

// ZKProof is structure that represents SnarkJS library result of proof generation
type ZKProof struct {
	A        []string   `json:"pi_a"`
	B        [][]string `json:"pi_b"`
	C        []string   `json:"pi_c"`
	Protocol string     `json:"protocol"`
}

// FullProof is ZKP proof with public signals
type FullProof struct {
	Proof      *ZKProof `json:"proof"`
	PubSignals []string `json:"pub_signals"`
}

// GenerateProof executes snarkjs groth16prove function and returns proof only if it's valid
func (p *Proof) GenerateProof(ctx context.Context, inputs ZKInputs) (*types.ZKProof, error) {
	if inputs == nil {
		return nil, fmt.Errorf("no inputs provided")
	}
	// Create proof data directory
	err := os.MkdirAll(p.dataPath, 0775)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create proofs folder on file system")
	}

	// Create input.json proofs folder
	inputJsonPath := p.dataPath + "/" + inputJsonFileName
	inputJson, err := os.Create(inputJsonPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create input.json on file system")
	}
	defer inputJson.Close()
	encoder := json.NewEncoder(inputJson)
	err = encoder.Encode(inputs)
	if err != nil {
		return nil, errors.Wrap(err, "failed to populate input.json file with input")
	}

	// Generate witness
	witnessStartTime := time.Now()
	witnessPath := p.dataPath + "/" + witnessFileName
	cmd := exec.Command(p.circuitBinaryPath, inputJsonPath, witnessPath)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate witness")
	}
	witness, err := os.ReadFile(witnessPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read newly generated witness file")
	}
	log.WithContext(ctx).Infow("Successfully generated witness", "duration", time.Now().Sub(witnessStartTime).String())

	// Generate Proof
	proofStartTime := time.Now()
	zkeyBytes, err := os.ReadFile(p.zkeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read zkey file")
	}
	proof, err := prover.Groth16Prover(zkeyBytes, witness)
	if err != nil {
		log.WithContext(ctx).Errorw("failed to generate proof", "proof", proof, "error", err)
		return nil, errors.Wrap(err, "failed to generate proof")
	}
	log.WithContext(ctx).Infow("Generated proof", "duration", time.Now().Sub(proofStartTime).String())

	// Verify generated Proof
	vkeyBytes, err := os.ReadFile(p.vkeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read verification_key file")
	}

	err = verifier.VerifyGroth16(*proof, vkeyBytes)
	if err != nil {
		log.WithContext(ctx).Errorw("failed to verify proof", "proof", proof, "error", err)
		return nil, errors.Wrap(err, "failed to verify proof")
	}
	log.WithContext(ctx).Info("Successfully verified proof")

	// Delete witness & input.json
	err = os.RemoveAll(filepath.Dir(inputJsonPath))
	if err != nil {
		log.WithContext(ctx).Errorw("failed to delete input.json folder", "path", filepath.Dir(inputJsonPath), "error", err)
	}

	return proof, nil
}

// VerifyZkProof executes snarkjs verify function and returns if proof is valid
func VerifyZkProof(ctx context.Context, circuitPath string, zkp *FullProof) error {

	if path.Clean(circuitPath) != circuitPath {
		return fmt.Errorf("illegal circuitPath")
	}

	vkeyBytes, err := os.ReadFile(circuitPath + "/verification_key.json")
	if err != nil {
		return errors.Wrap(err, "failed to read verification_key file")
	}

	proof := types.ZKProof{
		Proof: &types.ProofData{
			A: zkp.Proof.A,
			B: zkp.Proof.B,
			C: zkp.Proof.C,
			//Protocol: "groth16",
		},
		PubSignals: zkp.PubSignals,
	}
	err = verifier.VerifyGroth16(proof, vkeyBytes)
	if err != nil {
		log.WithContext(ctx).Errorw("failed to verify proof", "proof", zkp, "error", err)
		return errors.Wrap(err, "failed to verify proof")
	}

	return nil
}
