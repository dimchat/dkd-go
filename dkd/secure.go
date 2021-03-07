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
)

/**
 *  Secure Message
 *  ~~~~~~~~~~~~~~
 *  Instant Message encrypted by a symmetric key
 *
 *  data format: {
 *      //-- envelope
 *      sender   : "moki@xxx",
 *      receiver : "hulk@yyy",
 *      time     : 123,
 *      //-- content data and key/keys
 *      data     : "...",  // base64_encode(symmetric)
 *      key      : "...",  // base64_encode(asymmetric)
 *      keys     : {
 *          "ID1": "key1", // base64_encode(asymmetric)
 *      }
 *  }
 */
type EncryptedMessage struct {
	BaseMessage
	ISecureMessage

	_data []byte
	_key []byte
	_keys map[string]string
}

func NewEncryptedMessage(dict map[string]interface{}) *EncryptedMessage {
	return new(EncryptedMessage).Init(dict)
}

func (msg *EncryptedMessage) Init(dict map[string]interface{}) *EncryptedMessage {
	if msg.BaseMessage.Init(dict) != nil {
		// lazy load
		msg._data = nil
		msg._key = nil
		msg._keys = nil
	}
	return msg
}

//-------- ISecureMessage

func (msg *EncryptedMessage) EncryptedData() []byte {
	if msg._data == nil {
		base64 := msg.Get("data")
		delegate := msg.Delegate()
		msg._data = delegate.DecodeData(base64.(string), msg)
	}
	return msg._data
}

func (msg *EncryptedMessage) EncryptedKey() []byte {
	if msg._key == nil {
		base64 := msg.Get("key")
		if base64 == nil {
			// check 'keys'
			keys := msg.EncryptedKeys()
			if keys != nil {
				receiver := msg.Receiver()
				base64 = keys[receiver.String()]
			}
		}
		if base64 != nil {
			delegate := msg.Delegate()
			msg._key = delegate.DecodeKey(base64.(string), msg)
		}
	}
	return msg._key
}

func (msg *EncryptedMessage) EncryptedKeys() map[string]string {
	if msg._keys == nil {
		keys := msg.Get("keys")
		if keys != nil {
			msg._keys = keys.(map[string]string)
		}
	}
	return msg._keys
}

/*
 *  Decrypt the Secure Message to Instant Message
 *
 *    +----------+      +----------+
 *    | sender   |      | sender   |
 *    | receiver |      | receiver |
 *    | time     |  ->  | time     |
 *    |          |      |          |  1. PW      = decrypt(key, receiver.SK)
 *    | data     |      | content  |  2. content = decrypt(data, PW)
 *    | key/keys |      +----------+
 *    +----------+
 */

/**
 *  Decrypt message, replace encrypted 'data' with 'content' field
 *
 * @return InstantMessage object
 */
func (msg *EncryptedMessage) Decrypt() InstantMessage {
	var sender = msg.Sender()
	var receiver ID
	var group = msg.Group()
	if group == nil {
		// personal message
		// not split group message
		receiver = msg.Receiver()
	} else {
		// group message
		receiver = group
	}

	// 1. decrypt 'message.key' to symmetric key
	delegate := msg.Delegate()
	// 1.1. decode encrypted key data
	key := msg.EncryptedKey()
	// 1.2. decrypt key data
	if key != nil {
		key = delegate.DecryptKey(key, sender, receiver, msg)
		if key == nil {
			panic("failed to decrypt key in msg")
		}
	}
	// 1.3. deserialize key
	//      if key is empty, means it should be reused, get it from key cache
	password := delegate.DeserializeKey(key, sender, receiver, msg)
	if password == nil {
		panic("failed to get msg key")
	}

	// 2. decrypt 'message.data' to 'message.content'
	// 2.1. decode encrypted content data
	data := msg.EncryptedData()
	if data == nil {
		panic("failed to decode content data")
	}
	// 2.2. decrypt content data
	data = delegate.DecryptContent(data, password, msg)
	if data == nil {
		panic("failed to decrypt data with key")
	}
	// 2.3. deserialize content
	content := delegate.DeserializeContent(data, password, msg)
	if content == nil {
		panic("failed to deserialize content")
	}
	// 2.4. check attachment for File/Image/Audio/Video message content
	//      if file data not download yet,
	//          decrypt file data with password;
	//      else,
	//          save password to 'message.content.password'.
	//      (do it in 'core' module)

	// 3. pack message
	info := msg.GetMap(true)
	delete(info, "key")
	delete(info, "keys")
	delete(info, "data")
	info["content"] = content.GetMap(false)
	return InstantMessageParse(info)
}

/*
 *  Sign the Secure Message to Reliable Message
 *
 *    +----------+      +----------+
 *    | sender   |      | sender   |
 *    | receiver |      | receiver |
 *    | time     |  ->  | time     |
 *    |          |      |          |
 *    | data     |      | data     |
 *    | key/keys |      | key/keys |
 *    +----------+      | signature|  1. signature = sign(data, sender.SK)
 *                      +----------+
 */

/**
 *  Sign message.data, add 'signature' field
 *
 * @return ReliableMessage object
 */
func (msg *EncryptedMessage) Sign() ReliableMessage {
	delegate := msg.Delegate()
	sender := msg.Sender()
	data := msg.EncryptedData()
	// 1. sign with sender's private key
	signature := delegate.SignData(data, sender, msg)
	// 2. encode signature
	base64 := delegate.EncodeSignature(signature, msg)
	// 3. pack message
	info := msg.GetMap(true)
	info["signature"] = base64
	return ReliableMessageParse(info)
}

/*
 *  Split/Trim group message
 *
 *  for each members, get key from 'keys' and replace 'receiver' to member ID
 */

/**
 *  Split the group message to single person messages
 *
 *  @param members - group members
 *  @return secure/reliable message(s)
 */
func (msg *EncryptedMessage) Split(members []ID) []SecureMessage {
	info := msg.GetMap(true)
	// check 'keys'
	keys := msg.EncryptedKeys()
	if keys == nil {
		keys = make(map[string]string)
	} else {
		delete(info, "keys")
	}

	// 1. move the receiver(group ID) to 'group'
	//    this will help the receiver knows the group ID
	//    when the group message separated to multi-messages;
	//    if don't want the others know your membership,
	//    DON'T do this.
	info["group"] = msg.Receiver()

	messages := make([]SecureMessage, 0, len(members))
	for _, member := range members {
		// 2. change 'receiver' to each group member
		info["receiver"] = member
		// 3. get encrypted key
		base64 := keys[member.String()]
		if base64 == "" {
			delete(info, "key")
		} else {
			info["key"] = base64
		}
		// 4. repack message
		sMsg := SecureMessageParse(CloneMap(info))
		if sMsg != nil {
			messages = append(messages, sMsg)
		}
	}
	return messages
}

/**
 *  Trim the group message for a member
 *
 * @param member - group member ID/string
 * @return SecureMessage
 */
func (msg *EncryptedMessage) Trim(member ID) SecureMessage {
	info := msg.GetMap(true)
	// check 'keys'
	keys := msg.EncryptedKeys()
	if keys != nil {
		// move key data from 'keys' to key
		base64 := keys[member.String()]
		if base64 != "" {
			info["key"] = base64
		}
		delete(info, "keys")
	}
	// check 'group'
	group := msg.Group()
	if group == nil {
		// if 'group' not exists, the 'receiver' must be a group ID here, and
		// it will not be equal to the member of course,
		// so move 'receiver' to 'group'
		info["group"] = msg.Receiver()
	}
	info["receiver"] = member
	// repack
	return SecureMessageParse(info)
}

/**
 *  General Factory
 *  ~~~~~~~~~~~~~~~
 */
type EncryptedMessageFactory struct {
	SecureMessageFactory
}

func (factory *EncryptedMessageFactory) ParseSecureMessage(msg map[string]interface{}) SecureMessage {
	if _, exists := msg["signature"]; exists {
		// this should be a reliable message
		return NewRelayMessage(msg)
	} else {
		return NewEncryptedMessage(msg)
	}
}

func BuildSecureMessageFactory() SecureMessageFactory {
	factory := SecureMessageGetFactory()
	if factory == nil {
		factory = new(EncryptedMessageFactory)
		SecureMessageSetFactory(factory)
	}
	return factory
}
