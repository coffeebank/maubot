// maubot - A plugin-based Matrix bot system written in Go.
// Copyright (C) 2018 Tulir Asokan
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package matrix

import (
	"encoding/json"

	"maubot.xyz"
	"maunium.net/go/gomatrix"
)

type Event struct {
	*gomatrix.Event
	Client *Client
}

func roundtripContent(rawContent map[string]interface{}) (content maubot.Content) {
	if len(rawContent) == 0 {
		content.Raw = rawContent
		return
	}
	data, _ := json.Marshal(&rawContent)
	json.Unmarshal(data, &content)
	content.Raw = rawContent
	return
}

func (evt *Event) Interface() *maubot.Event {
	var stateKey string
	if evt.StateKey != nil {
		stateKey = *evt.StateKey
	}
	mbEvent := &maubot.Event{
		EventFuncs: evt,
		StateKey:   stateKey,
		Sender:     evt.Sender,
		Type:       maubot.EventType(evt.Type),
		Timestamp:  evt.Timestamp,
		ID:         evt.ID,
		RoomID:     evt.RoomID,
		Content:    roundtripContent(evt.Content),
		Redacts:    evt.Redacts,
		Unsigned: maubot.Unsigned{
			PrevContent:   roundtripContent(evt.Unsigned.PrevContent),
			PrevSender:    evt.Unsigned.PrevSender,
			ReplacesState: evt.Unsigned.ReplacesState,
			Age:           evt.Unsigned.Age,
		},
	}
	RemoveReplyFallback(mbEvent)
	return mbEvent
}

func (evt *Event) Reply(text string) (string, error) {
	return evt.SendRawEvent(maubot.EventMessage,
		SetReply(
			RenderMarkdown(text),
			evt.Event))
}

func (evt *Event) ReplyContent(content maubot.Content) (string, error) {
	return evt.SendRawEvent(maubot.EventMessage, SetReply(content, evt.Event))
}

func (evt *Event) SendMessage(text string) (string, error) {
	return evt.SendRawEvent(maubot.EventMessage, RenderMarkdown(text))
}

func (evt *Event) SendContent(content maubot.Content) (string, error) {
	return evt.SendRawEvent(maubot.EventMessage, content)
}

func (evt *Event) SendRawEvent(evtType maubot.EventType, content interface{}) (string, error) {
	resp, err := evt.Client.SendMessageEvent(evt.RoomID, string(evtType), content)
	if err != nil {
		return "", err
	}
	return resp.EventID, nil
}