package operation

import (
	"github.com/steemit/steemutil/encoder"
)

// FC_REFLECT( steemit::chain::comment_operation,
//             (parent_author)
//             (parent_permlink)
//             (author)
//             (permlink)
//             (title)
//             (body)
//             (json_metadata) )

// CommentOperation represents either a new post or a comment.
//
// In case Title is filled in and ParentAuthor is empty, it is a new post.
// The post category can be read from ParentPermlink.
type CommentOperation struct {
	ParentAuthor   string `json:"parent_author"`
	ParentPermlink string `json:"parent_permlink"`
	Author         string `json:"author"`
	Permlink       string `json:"permlink"`
	Title          string `json:"title"`
	Body           string `json:"body"`
	JsonMetadata   string `json:"json_metadata"`
}

func (op *CommentOperation) Type() OpType {
	return TypeComment
}

func (op *CommentOperation) Data() interface{} {
	return op
}

func (op *CommentOperation) IsStoryOperation() bool {
	return op.ParentAuthor == ""
}

func (op *CommentOperation) MarshalTransaction(encoderObj *encoder.Encoder) (err error) {
	if err = encoderObj.EncodeUVarint(uint64(op.Type().Code())); err != nil {
		return err
	}
	encoderObj.Encode(op.ParentAuthor)
	encoderObj.Encode(op.ParentPermlink)
	encoderObj.Encode(op.Author)
	encoderObj.Encode(op.Permlink)
	encoderObj.Encode(op.Title)
	encoderObj.Encode(op.Body)
	encoderObj.Encode(op.JsonMetadata)
	return
}
