package handlers

import (
	"context"

	"start-feishubot/services"

	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

func NewClearCardHandler(cardMsg CardMsg, m MessageHandler) CardHandlerFunc {
	return func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		if cardMsg.Kind == ClearCardKind {
			newCard, err, done := CommonProcessClearCache(cardMsg, m.sessionCache)
			if done {
				return newCard, err
			}
			return nil, nil
		}
		return nil, ErrNextHandler
	}
}

func CommonProcessClearCache(cardMsg CardMsg, session services.SessionServiceCacheInterface) (
	interface{}, error, bool) {
	if cardMsg.Value == "1" {
		session.Clear(cardMsg.SessionId)
		newCard, _ := newSendCard(
			withHeader("️🆑 Robot reminder", larkcard.TemplateGrey),
			withMainMd("The context information of this topic has been deleted"),
			withNote("We can start a brand new topic, keep looking for me to chat"),
		)
		//fmt.Printf("session: %v", newCard)
		return newCard, nil, true
	}
	if cardMsg.Value == "0" {
		newCard, _ := newSendCard(
			withHeader("️🆑 Robot reminder", larkcard.TemplateGreen),
			withMainMd("Still retain the context information of this topic"),
			withNote("We can continue to explore this topic and look forward to chatting with you.If you have other questions or topics you want to discuss, please tell me"),
		)
		return newCard, nil, true
	}
	return nil, nil, false
}
