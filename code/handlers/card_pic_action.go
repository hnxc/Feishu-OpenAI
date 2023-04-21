package handlers

import (
	"context"

	"start-feishubot/services"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

func NewPicResolutionHandler(cardMsg CardMsg, m MessageHandler) CardHandlerFunc {
	return func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		if cardMsg.Kind == PicResolutionKind {
			CommonProcessPicResolution(cardMsg, cardAction, m.sessionCache)
			return nil, nil
		}
		return nil, ErrNextHandler
	}
}

func NewPicModeChangeHandler(cardMsg CardMsg, m MessageHandler) CardHandlerFunc {
	return func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		if cardMsg.Kind == PicModeChangeKind {
			newCard, err, done := CommonProcessPicModeChange(cardMsg, m.sessionCache)
			if done {
				return newCard, err
			}
			return nil, nil
		}
		return nil, ErrNextHandler
	}
}

func NewPicTextMoreHandler(cardMsg CardMsg, m MessageHandler) CardHandlerFunc {
	return func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		if cardMsg.Kind == PicTextMoreKind {
			go func() {
				m.CommonProcessPicMore(cardMsg)
			}()
			return nil, nil
		}
		return nil, ErrNextHandler
	}
}

func CommonProcessPicResolution(msg CardMsg,
	cardAction *larkcard.CardAction,
	cache services.SessionServiceCacheInterface) {
	option := cardAction.Action.Option
	//fmt.Println(larkcore.Prettify(msg))
	cache.SetPicResolution(msg.SessionId, services.Resolution(option))
	//send text
	replyMsg(context.Background(), "The resolution of the picture has been updated"+option,
		&msg.MsgId)
}

func (m MessageHandler) CommonProcessPicMore(msg CardMsg) {
	resolution := m.sessionCache.GetPicResolution(msg.SessionId)
	//fmt.Println("resolution: ", resolution)
	//fmt.Println("msg: ", msg)
	question := msg.Value.(string)
	bs64, _ := m.gpt.GenerateOneImage(question, resolution)
	replayImageCardByBase64(context.Background(), bs64, &msg.MsgId,
		&msg.SessionId, question)
}

func CommonProcessPicModeChange(cardMsg CardMsg,
	session services.SessionServiceCacheInterface) (
	interface{}, error, bool) {
	if cardMsg.Value == "1" {

		sessionId := cardMsg.SessionId
		session.Clear(sessionId)
		session.SetMode(sessionId,
			services.ModePicCreate)
		session.SetPicResolution(sessionId,
			services.Resolution256)

		newCard, _ :=
			newSendCard(
				withHeader("🖼️ Enter the picture creation mode", larkcard.TemplateBlue),
				withPicResolutionBtn(&sessionId),
				withNote("remind：Reply text or picture，Let AI generate related pictures。"))
		return newCard, nil, true
	}
	if cardMsg.Value == "0" {
		newCard, _ := newSendCard(
			withHeader("️🎒 Robot reminder", larkcard.TemplateGreen),
			withMainMd("Still retain the context information of this topic"),
			withNote("We can continue to explore this topic and look forward to chatting with you.If you have other questions or topics you want to discuss, please tell me"),
		)
		return newCard, nil, true
	}
	return nil, nil, false
}
