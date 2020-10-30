/* license: https://mit-license.org
 *
 *  Dao-Ke-Dao: Universal Message Module
 *
 *                                Written in 2020 by Moky <albert.moky@gmail.com>
 *
 * ==============================================================================
 * The MIT License (MIT)
 *
 * Copyright (c) 2020 Albert Moky
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
package dkd

import (
	"dkd-go/protocol"
	"dkd-go/types"
	"math/rand"
	"time"
)

/**
 *  Message Content
 *  ~~~~~~~~~~~~~~~
 *  This class is for creating message content
 *
 *  data format: {
 *      'type'    : 0x00,            // message type
 *      'sn'      : 0,               // serial number
 *
 *      'group'   : 'Group ID',      // for group message
 *
 *      //-- message info
 *      'text'    : 'text',          // for text message
 *      'command' : 'Command Name',  // for system command
 *      //...
 *  }
 */
type Content struct {
	types.Dictionary

	_delegate *MessageDelegate

	// message type: text, image, ...
	_type protocol.ContentType

	// serial number: random number to identify message content
	_sn uint32

	// message time
	_time time.Time

	// extra info
	_group interface{}
}

func (content *Content)LoadContent(dictionary *map[string]interface{}) *Content {
	if content.LoadDictionary(dictionary) != nil {
		// lazy load
		content._type = 0
		content._sn = 0
		content._time = time.Unix(0, 0)
	}
	return content
}

func (content *Content) InitContent(t protocol.ContentType) *Content {
	if content.InitDictionary() != nil {
		// message time
		now := time.Now()
		stamp := now.Unix()
		// serial number
		rand.Seed(stamp)
		sn := rand.Uint32()
		// initialize
		content._type = t
		content._sn = sn
		content._time = now

		content.Set("type", uint8(t))
		content.Set("sn", sn)
		content.Set("time", stamp)
	}
	return content
}

func (content *Content) GetDelegate() *MessageDelegate {
	return content._delegate
}

func (content *Content) SetDelegate(delegate *MessageDelegate) {
	content._delegate = delegate
}

func (content *Content) GetType() protocol.ContentType {
	if content._type == 0 {
		t := content.Get("type")
		content._type = protocol.ContentType(t.(uint8))
	}
	return content._type
}

func (content *Content) GetSerialNumber() uint32 {
	if content._sn == 0 {
		sn := content.Get("sn")
		content._sn = sn.(uint32)
	}
	return content._sn
}

func (content *Content) GetTime() time.Time {
	if content._time.IsZero() {
		timestamp := content.Get("time")
		if timestamp != nil {
			content._time = time.Unix(timestamp.(int64), 0)
		}
	}
	return content._time
}

// Group ID/string for group message
//    if field 'group' exists, it means this is a group message
func (content *Content) GetGroup() interface{} {
	if content._group == nil {
		group := content.Get("group")
		if group != nil {
			handler := *content.GetDelegate()
			content._group = handler.GetID(group)
		}
	}
	return content._group
}

func (content *Content) SetGroup(group interface{})  {
	content.Set("group", group.(string))
	content._group = group
}
