package gocodecc

import "github.com/astaxie/beego/orm"

const (
	_ = iota
	MessageTypeReply
	MessageTypeMail
)

type MessageModel struct {
	Id         int `orm:pk;auto`
	Receiver   uint32
	Sender     uint32
	SenderName string `orm:size(21)`
	Type       int
	Message    string
	Url        string
	Read       int
}

var (
	messageTableName = "message"
)

func (m *MessageModel) TableName() string {
	return messageTableName
}

func init() {
	orm.RegisterModel(new(MessageModel))
}

func modelMessageNew(receiver uint32, tp int, msg string, url string, sender *WebUser) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}

	_, err = db.Exec("INSERT INTO message (receiver, sender, sender_name, type, message, url, read) VALUES (?, ?, ?, ?, ?, ?, 0)",
		receiver, sender.Uid, sender.UserName, tp, msg, url)
	return err
}

func modelMessageGetCountByReceiver(receiver uint32) (int, error) {
	db, err := getRawDB()
	if nil != err {
		return 0, err
	}

	row := db.QueryRow("SELECT COUNT(*) FROM message WHERE receiver = ? AND read = 0", receiver)
	var cnt int

	if err = row.Scan(&cnt); nil != err {
		return 0, err
	}

	return cnt, nil
}

func modelMessageGetByID(id int) (*MessageModel, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	row := db.QueryRow("SELECT type, receiver, message, url, sender, sender_name FROM message WHERE id = ?", id)
	var message MessageModel
	if err = row.Scan(&message.Type,
		&message.Receiver,
		&message.Message,
		&message.Url,
		&message.Sender,
		&message.SenderName); nil != err {
		return nil, err
	}
	message.Id = id
	return &message, nil
}

func modelMessageGetByReceiver(receiver uint32) ([]*MessageModel, error) {
	db, err := getRawDB()
	if nil != err {
		return nil, err
	}

	rows, err := db.Query("SELECT id, type, message, url, sender, sender_name FROM message WHERE receiver = ? AND read = 0", receiver)
	if nil != err {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*MessageModel, 0, 32)
	for rows.Next() {
		var msg MessageModel
		msg.Receiver = receiver

		if err = rows.Scan(&msg.Id, &msg.Type, &msg.Message, &msg.Url, &msg.Sender, &msg.SenderName); nil != err {
			return nil, err
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

func modelMessageMarkRead(uid uint32, id int) error {
	db, err := getRawDB()
	if nil != err {
		return err
	}

	_, err = db.Exec("UPDATE message SET read = 1 WHERE id = ? AND receiver = ?", id, uid)
	return err
}
