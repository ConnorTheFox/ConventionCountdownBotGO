package main

import (
	"conBot/helper"

	"github.com/globalsign/mgo/bson"
	tgAPI "gopkg.in/tucnak/telebot.v2"
)

//Handlers From Telegram

func findSubOrUnsubKeyboard(chatID int64) [][]tgAPI.InlineButton {
	var keyboardToSend [][]tgAPI.InlineButton
	if db.ItemExists("users", bson.M{"chatId": chatID}) == true {
		keyboardToSend = keyboards["mainUnsub"]
	} else {
		keyboardToSend = keyboards["mainSub"]
	}
	return keyboardToSend
}

func getChatUser(chat *tgAPI.Chat, user *tgAPI.User) (*tgAPI.ChatMember, error) {
	chatMember, err := bot.ChatMemberOf(chat, user)
	if err != nil {
		return nil, err
	}
	return chatMember, nil
}

func checkForAdmin(chatMember *tgAPI.ChatMember) bool {
	if chatMember.Role == "creator" || chatMember.Role == "administrator" {
		return true
	}
	return false
}

func handleStart(msg *tgAPI.Message) {
	if msg.Chat.Type == "channel" || msg.Chat.Type == "privatechannel" {
		return
	}
	bot.Send(msg.Sender, config.WelcomeMsg, &tgAPI.ReplyMarkup{
		InlineKeyboard: findSubOrUnsubKeyboard(msg.Chat.ID),
	})
}

func handleGroupAdd(msg *tgAPI.Message) {
	bot.Send(msg.Chat, config.GroupAddMsg, &tgAPI.ReplyMarkup{
		InlineKeyboard: findSubOrUnsubKeyboard(msg.Chat.ID),
	})
}

//Handalers For Keybaord
func handleSubBtn(c *tgAPI.Callback) {
	if c.Message.FromGroup() == true {
		chatUser, err := getChatUser(c.Message.Chat, c.Sender)
		if err != nil {
			handleBtnClick("An Error Occured", keyboards["back"], c)
			handleErr(err)
			return
		}
		isAdmin := checkForAdmin(chatUser)
		if isAdmin == false {
			handleBtnClick(config.GroupNotAdminMsg, keyboards["back"], c)
			return
		}
	}
	status := handleSub(c.Message)
	if status == true {
		handleBtnClick(config.SubMsg, keyboards["back"], c)
	} else {
		handleBtnClick(config.AlreadySubMsg, keyboards["back"], c)
	}
}

func handleUnsubBtn(c *tgAPI.Callback) {
	if c.Message.FromGroup() == true {
		chatUser, err := getChatUser(c.Message.Chat, c.Sender)
		if err != nil {
			handleBtnClick("An Error Occured", keyboards["back"], c)
			handleErr(err)
			return
		}
		isAdmin := checkForAdmin(chatUser)
		if isAdmin == false {
			handleBtnClick(config.GroupNotAdminMsg, keyboards["back"], c)
			return
		}
	}
	status := handleUnsub(c.Message)
	if status == true {
		handleBtnClick(config.UnsubMsg, keyboards["back"], c)
	} else {
		handleBtnClick(config.NotSubMsg, keyboards["back"], c)
	}
}

func handleCommandBtn(c *tgAPI.Callback) {
	handleBtnClick(config.CmdMsg, keyboards["cmd"], c)
}

func handleHomeBtn(c *tgAPI.Callback) {
	handleBtnClick(config.WelcomeMsg, findSubOrUnsubKeyboard(c.Message.Chat.ID), c)
}

func handleInfoBtn(c *tgAPI.Callback) {
	handleBtnClick(config.InfoMsg, keyboards["back"], c)
}

func handleDaysBtn(c *tgAPI.Callback) {
	dayStr := helper.GetDays(config.Date) + " Days Until " + config.Con + " !"
	handleBtnClick(dayStr, keyboards["back"], c)
}

func handleSub(msg *tgAPI.Message) bool {
	if db.ItemExists("users", bson.M{"chatId": msg.Chat.ID}) == true {
		return false
	}
	itemToInsert := helper.User{
		ChatID: msg.Chat.ID,
		Name:   msg.Chat.Username,
		Group:  msg.FromGroup(),
	}
	db.Insert("users", itemToInsert)
	return true
}

func handleUnsub(msg *tgAPI.Message) bool {
	if db.ItemExists("users", bson.M{"chatId": msg.Chat.ID}) == false {
		return false
	}
	db.RemoveOne("users", bson.M{"chatId": msg.Chat.ID})
	return true
}
