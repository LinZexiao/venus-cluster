package mock

import (
	"context"
	"sync"

	"github.com/filecoin-project/go-state-types/abi"
	venusMessager "github.com/filecoin-project/venus-messager/api/client"

	"github.com/dtynn/venus-cluster/venus-sealer/sealer/api"
)

var _ api.CommitmentManager = (*commitMgr)(nil)

func NewCommitManager() api.CommitmentManager {
	cmgr := &commitMgr{}

	cmgr.pres.commits = map[abi.SectorID]api.PreCommitOnChainInfo{}
	cmgr.proofs.proofs = map[abi.SectorID]api.ProofOnChainInfo{}
	return cmgr
}

func NewMessagerClient() venusMessager.IMessager {
	return nil
}

type commitMgr struct {
	pres struct {
		sync.RWMutex
		commits map[abi.SectorID]api.PreCommitOnChainInfo
	}

	proofs struct {
		sync.RWMutex
		proofs map[abi.SectorID]api.ProofOnChainInfo
	}
}

func (c *commitMgr) SubmitPreCommit(ctx context.Context, sid abi.SectorID, info api.PreCommitOnChainInfo) (api.SubmitPreCommitResp, error) {
	c.pres.Lock()
	defer c.pres.Unlock()

	if _, ok := c.pres.commits[sid]; ok {
		return api.SubmitPreCommitResp{
			Res:  api.SubmitDuplicateSubmit,
			Desc: nil,
		}, nil
	}

	c.pres.commits[sid] = info

	return api.SubmitPreCommitResp{
		Res:  api.SubmitAccepted,
		Desc: nil,
	}, nil
}

func (c *commitMgr) PreCommitState(ctx context.Context, sid abi.SectorID) (api.PollPreCommitStateResp, error) {
	c.pres.RLock()
	defer c.pres.RUnlock()

	if _, ok := c.pres.commits[sid]; ok {
		return api.PollPreCommitStateResp{
			State: api.OnChainStateLanded,
			Desc:  nil,
		}, nil
	}

	return api.PollPreCommitStateResp{
		State: api.OnChainStateNotFound,
		Desc:  nil,
	}, nil
}

func (c *commitMgr) SubmitProof(ctx context.Context, sid abi.SectorID, info api.ProofOnChainInfo) (api.SubmitProofResp, error) {
	c.proofs.Lock()
	defer c.proofs.Unlock()

	if _, ok := c.proofs.proofs[sid]; ok {
		return api.SubmitProofResp{
			Res:  api.SubmitDuplicateSubmit,
			Desc: nil,
		}, nil
	}

	c.proofs.proofs[sid] = info

	return api.SubmitProofResp{
		Res:  api.SubmitAccepted,
		Desc: nil,
	}, nil
}

func (c *commitMgr) ProofState(ctx context.Context, sid abi.SectorID) (api.PollProofStateResp, error) {
	c.proofs.RLock()
	defer c.proofs.RUnlock()

	if _, ok := c.proofs.proofs[sid]; ok {
		return api.PollProofStateResp{
			State: api.OnChainStateLanded,
			Desc:  nil,
		}, nil
	}

	return api.PollProofStateResp{
		State: api.OnChainStateNotFound,
		Desc:  nil,
	}, nil
}
