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
	. "github.com/dimchat/dkd-go/protocol"
	. "github.com/dimchat/mkm-go/protocol"
	. "github.com/dimchat/mkm-go/types"
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
type BaseContent struct {
	Dictionary
	IContent

	_type uint8
	_sn uint32

	_time time.Time
	_group ID
}

/* designated initializer */
func (content *BaseContent) Init(dict map[string]interface{}) *BaseContent {
	if content.Dictionary.Init(dict) != nil {
		// lazy load
		content._type = 0
		content._sn = 0
		content._time = time.Time{}
		content.setGroup(nil)
	}
	return content
}

/* designated initializer */
func (content *BaseContent) InitWithType(msgType uint8) *BaseContent {
	// message time
	now := time.Now()
	stamp := now.Unix()
	// serial number
	rand.Seed(stamp)
	sn := rand.Uint32()
	// build content info
	dict := make(map[string]interface{})
	dict["type"] = msgType
	dict["sn"] = sn
	dict["time"] = stamp
	if content.Dictionary.Init(dict) != nil {
		content._type = msgType
		content._sn = sn
		content._time = now
		content.setGroup(nil)
	}
	return content
}

//func (content *BaseContent) Release() int {
//	cnt := content.Dictionary.Release()
//	if cnt == 0 {
//		// this object is going to be destroyed,
//		// release children
//		content.setGroup(nil)
//	}
//	return cnt
//}

func (content *BaseContent) setGroup(group ID)  {
	if group != content._group {
		//ObjectRetain(group)
		//ObjectRelease(content._group)
		content._group = group
	}
}

//-------- IContent

func (content *BaseContent) Type() uint8 {
	if content._type == 0 {
		content._type = ContentGetType(content.GetMap(false))
	}
	return content._type
}

func (content *BaseContent) SN() uint32 {
	if content._sn == 0 {
		content._sn = ContentGetSN(content.GetMap(false))
	}
	return content._sn
}

func (content *BaseContent) Time() time.Time {
	if content._time.IsZero() {
		content._time = ContentGetTime(content.GetMap(false))
	}
	return content._time
}

func (content *BaseContent) Group() ID {
	if content._group == nil {
		content.setGroup(ContentGetGroup(content.GetMap(false)))
	}
	return content._group
}

func (content *BaseContent) SetGroup(group ID)  {
	ContentSetGroup(content.GetMap(false), group)
	content.setGroup(group)
}
