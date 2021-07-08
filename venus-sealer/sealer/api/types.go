package api

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/specs-actors/actors/builtin/miner"
	proof5 "github.com/filecoin-project/specs-actors/v5/actors/runtime/proof"
	"github.com/ipfs/go-cid"
)

type AllocateSectorSpec struct {
	AllowedMiners     []abi.ActorID
	AllowedProofTypes []abi.RegisteredSealProof
}

type AllocatedSector struct {
	ID        abi.SectorID
	ProofType abi.RegisteredSealProof
}

type PieceInfo struct {
	Size abi.PaddedPieceSize
	Cid  cid.Cid
}

type DealInfo struct {
	ID    abi.DealID
	Piece abi.PieceInfo
}

type Deals []DealInfo

type AcquireDealsSpec struct {
	MaxDeals *uint
}

type Ticket struct {
	Ticket abi.Randomness
	Epoch  abi.ChainEpoch
}

type SubmitResult uint64

const (
	SubmitUnknown SubmitResult = iota
	SubmitAccepted
	SubmitDuplicateSubmit
	SubmitMismatchedSubmission
	SubmitRejected
	Submitted
)

type OnChainState uint64

const (
	OnChainStateUnknown OnChainState = iota
	OnChainStatePending
	OnChainStatePacked
	OnChainStateLanded
	OnChainStateNotFound
	OnChainStateFailed
)

type PreCommitOnChainInfo struct {
	CommR  [32]byte
	CommD  [32]byte
	Ticket Ticket
	Deals  []abi.DealID
}

type ProofOnChainInfo struct {
	Proof []byte
}

type SubmitPreCommitResp struct {
	Res  SubmitResult
	Desc *string
}

type PollPreCommitStateResp struct {
	State OnChainState
	Desc  *string
}

type WaitSeedResp struct {
	ShouldWait bool
	Delay      int
	Seed       *Seed
}

type Seed struct {
	Seed  abi.Randomness
	Epoch abi.ChainEpoch
}

type SubmitProofResp struct {
	Res  SubmitResult
	Desc *string
}

type PollProofStateResp struct {
	State OnChainState
	Desc  *string
}

type MinerInfo struct {
	ID                  abi.ActorID
	Owner               address.Address
	Worker              address.Address
	NewWorker           address.Address
	ControlAddresses    []address.Address
	WorkerChangeEpoch   abi.ChainEpoch
	WindowPoStProofType abi.RegisteredPoStProof
	SealProofType       abi.RegisteredSealProof
}

type TipSetToken []byte

type AggregateInput struct {
	Spt   abi.RegisteredSealProof
	Info  proof5.AggregateSealVerifyInfo
	Proof []byte
}

type PreCommitEntry struct {
	Deposit abi.TokenAmount
	Pci     *miner.SectorPreCommitInfo
}
