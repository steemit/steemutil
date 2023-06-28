package operation

import (
	"github.com/steemit/steemutil/encoder"
	"github.com/steemit/steemutil/util"
)

// FC_REFLECT( steemit::chain::vote_operation,
//             (voter)
//             (author)
//             (permlink)
//             (weight) )

type VoteOperation struct {
	Voter    string     `json:"voter"`
	Author   string     `json:"author"`
	Permlink string     `json:"permlink"`
	Weight   util.Int16 `json:"weight"`
}

func (op *VoteOperation) Type() OpType {
	return TypeVote
}

func (op *VoteOperation) Data() interface{} {
	return op
}

func (op *VoteOperation) MarshalTransaction(encoderObj *encoder.Encoder) (err error) {
	encoderObj.EncodeUVarint(uint64(TypeVote.Code()))
	encoderObj.Encode(op.Voter)
	encoderObj.Encode(op.Author)
	encoderObj.Encode(op.Permlink)
	encoderObj.Encode(op.Weight)
	return
}
