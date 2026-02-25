/* license: https://mit-license.org
 *
 *  Dao-Ke-Dao: Universal Message Module
 *
 *                                Written in 2021 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2021 Albert Moky
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 * ==============================================================================
 */
package protocol

import . "github.com/dimchat/mkm-go/types"

// InstantMessage represents a message type designed for immediate delivery,
// extending the base Message with a dedicated content payload.
//
// It includes standard envelope fields along with the actual message content.
//
//	data format: {
//		//-- envelope
//		"sender"   : "moki@xxx",
//		"receiver" : "hulk@yyy",
//		"time"     : 123,
//
//		//-- content
//		"content"  : {...}
//	}
type InstantMessage interface {
	Message

	Content() Content
	/*/
	// only for rebuild content
	SetContent(body Content)
	/*/
}

/**
 *  Message Factory
 */

type InstantMessageFactory interface {

	// GenerateSerialNumber generates SN for message content
	//
	// Parameters:
	//   - msgType: content type
	//   - now: message time
	// Returns: SN (serial number as msg id)
	GenerateSerialNumber(msgType MessageType, now Time) SerialNumberType

	// CreateInstantMessage creates instant message with envelope & content
	//
	// Parameters:
	//   - head: message envelope
	//   - body: message content
	// Returns: InstantMessage
	CreateInstantMessage(head Envelope, body Content) InstantMessage

	// ParseInstantMessage parses map object to message
	//
	// Parameters:
	//   - msg: message info
	// Returns: InstantMessage
	ParseInstantMessage(msg StringKeyMap) InstantMessage
}

//
//  Factory methods
//

func CreateInstantMessage(head Envelope, body Content) InstantMessage {
	helper := GetInstantMessageHelper()
	return helper.CreateInstantMessage(head, body)
}

func ParseInstantMessage(msg any) InstantMessage {
	helper := GetInstantMessageHelper()
	return helper.ParseInstantMessage(msg)
}

func GenerateSerialNumber(msgType MessageType, now Time) SerialNumberType {
	helper := GetInstantMessageHelper()
	return helper.GenerateSerialNumber(msgType, now)
}

func GetInstantMessageFactory() InstantMessageFactory {
	helper := GetInstantMessageHelper()
	return helper.GetInstantMessageFactory()
}

func SetInstantMessageFactory(factory InstantMessageFactory) {
	helper := GetInstantMessageHelper()
	helper.SetInstantMessageFactory(factory)
}

//
//  Conveniences
//

func InstantMessageConvert(array any) []InstantMessage {
	values := FetchList(array)
	messages := make([]InstantMessage, 0, len(values))
	var msg InstantMessage
	for _, item := range values {
		msg = ParseInstantMessage(item)
		if msg == nil {
			continue
		}
		messages = append(messages, msg)
	}
	return messages
}

func InstantMessageRevert(messages []InstantMessage) []StringKeyMap {
	array := make([]StringKeyMap, len(messages))
	for idx, msg := range messages {
		array[idx] = msg.Map()
	}
	return array
}
