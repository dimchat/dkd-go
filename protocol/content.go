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
package protocol

import (
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
)

// ContentType
type MessageType = string

type SerialNumberType = uint64

/**
 *  Message Content
 *  <p>
 *      This class is for creating message content
 *  </p>
 *
 *  <blockquote><pre>
 *  data format: {
 *      "type"    : i2s(0),         // message type
 *      "sn"      : 0,              // serial number
 *
 *      "time"    : 123,            // message time
 *      "group"   : "{GroupID}",    // for group message
 *
 *      //-- message info
 *      "text"    : "text",         // for text message
 *      "command" : "Command Name"  // for system command
 *      //...
 *  }
 *  </pre></blockquote>
 */
type Content interface {
	Mapper

	Type() MessageType    // content type
	SN() SerialNumberType // serial number as message id

	Time() Time // message time

	// Group ID/string for group message
	//    if field 'group' exists, it means this is a group message
	Group() ID
	SetGroup(group ID)
}

/**
 *  Content Factory
 *  ~~~~~~~~~~~~~~~
 */
type ContentFactory interface {

	/**
	 *  Parse map object to content
	 *
	 * @param content - content info
	 * @return Content
	 */
	ParseContent(content StringKeyMap) Content
}

//
//  Factory method
//

func ParseContent(content interface{}) Content {
	helper := GetContentHelper()
	return helper.ParseContent(content)
}

func GetContentFactory(msgType MessageType) ContentFactory {
	helper := GetContentHelper()
	return helper.GetContentFactory(msgType)
}

func SetContentFactory(msgType MessageType, factory ContentFactory) {
	helper := GetContentHelper()
	helper.SetContentFactory(msgType, factory)
}

//
//  Conveniences
//

func ContentConvert(array interface{}) []Content {
	values := FetchList(array)
	contents := make([]Content, 0, len(values))
	var msg Content
	for _, item := range values {
		msg = ParseContent(item)
		if msg == nil {
			continue
		}
		contents = append(contents, msg)
	}
	return contents
}

func ContentRevert(contents []Content) []StringKeyMap {
	array := make([]StringKeyMap, len(contents))
	for idx, msg := range contents {
		array[idx] = msg.Map()
	}
	return array
}
