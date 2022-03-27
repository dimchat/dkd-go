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
type Content interface {
	IContent
	Map
}
type IContent interface {

	Type() uint8      // message type
	SN() uint32       // serial number as message id

	Time() time.Time  // message time

	// Group ID/string for group message
	//    if field 'group' exists, it means this is a group message
	Group() ID
	SetGroup(group ID)
}

func ContentGetType(content map[string]interface{}) uint8 {
	msgType, ok := content["type"].(uint8)
	if ok {
		return msgType
	} else {
		return 0
	}
}

func ContentGetSN(content map[string]interface{}) uint32 {
	sn, ok := content["sn"].(uint32)
	if ok {
		return sn
	} else {
		return 0
	}
}

func ContentGetTime(content map[string]interface{}) time.Time {
	timestamp, ok := content["time"].(int64)
	if ok {
		return time.Unix(timestamp, 0)
	} else {
		return time.Time{}
	}
}

func ContentGetGroup(content map[string]interface{}) ID {
	return IDParse(content["group"])
}

func ContentSetGroup(content map[string]interface{}, group ID) {
	if ValueIsNil(group) {
		delete(content, "group")
	} else {
		content["group"] = group.String()
	}
}

/**
 *  Content Factory
 *  ~~~~~~~~~~~~~~~
 */
type ContentFactory interface {
	IContentFactory
}
type IContentFactory interface {

	/**
	 *  Parse map object to content
	 *
	 * @param content - content info
	 * @return Content
	 */
	ParseContent(content map[string]interface{}) Content
}

var contentFactories = make(map[uint8]ContentFactory)

func ContentSetFactory(msgType uint8, factory ContentFactory) {
	contentFactories[msgType] = factory
}

func ContentGetFactory(msgType uint8) ContentFactory {
	return contentFactories[msgType]
}

//
//  Factory method
//
func ContentParse(content interface{}) Content {
	if ValueIsNil(content) {
		return nil
	}
	value, ok := content.(Content)
	if ok {
		return value
	}
	// get content info
	var info map[string]interface{}
	wrapper, ok := content.(Map)
	if ok {
		info = wrapper.GetMap(false)
	} else {
		info, ok = content.(map[string]interface{})
		if !ok {
			panic(content)
			return nil
		}
	}
	// get content factory by type
	msgType := ContentGetType(info)
	factory := ContentGetFactory(msgType)
	if factory == nil {
		factory = ContentGetFactory(0)  // unknown
	}
	return factory.ParseContent(info)
}
