package table

import (
	"gopkg.in/mgo.v2/bson"
)

// 只存笔记基本信息
// 内容不存放
type Notes struct {
	NoteId bson.ObjectId `bson:"_id,omitempty"` // 必须要设置bson:"_id" 不然mgo不会认为是主键
	Title  string        `title`                // 标题
	Desc   string        `desc`                 // 描述, 非html
}
