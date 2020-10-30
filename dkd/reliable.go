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

/**
 *  Reliable Message signed by an asymmetric key
 *  ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
 *  This class is used to sign the SecureMessage
 *  It contains a 'signature' field which signed with sender's private key
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
 *      },
 *      //-- signature
 *      signature: "..."   // base64_encode()
 *  }
 */
type ReliableMessage struct {
	SecureMessage

	_signature []byte
}

func CreateReliableMessage(dictionary *map[string]interface{}) *ReliableMessage {
	return new(ReliableMessage).LoadReliableMessage(dictionary)
}

func (msg *ReliableMessage)LoadReliableMessage(dictionary *map[string]interface{}) *ReliableMessage {
	if msg.LoadMessage(dictionary) != nil {
		// lazy load
		msg._signature = nil
	}
	return msg
}

func (msg *ReliableMessage) GetDelegate() *ReliableMessageDelegate {
	delegate := msg.GetEnvelope().GetDelegate()
	handler := (*delegate).(ReliableMessageDelegate)
	return &handler
}

func (msg *ReliableMessage) SetDelegate(delegate *ReliableMessageDelegate) {
	handler := (MessageDelegate)(*delegate)
	msg.GetEnvelope().SetDelegate(&handler)
}

func (msg *ReliableMessage) GetSignature() []byte {
	if msg._signature == nil {
		base64 := msg.Get("signature")
		handler := *msg.GetDelegate()
		msg._signature = handler.DecodeSignature(base64.(string), msg)
	}
	return msg._signature
}

/**
 *  Sender's Meta
 *  ~~~~~~~~~~~~~
 *  Extends for the first message package of 'Handshake' protocol.
 *
 * @param meta - Meta object or dictionary
 */
func (msg *ReliableMessage) GetMeta() map[string]interface{} {
	value := msg.Get("meta")
	if value == nil {
		return nil
	}
	return value.(map[string]interface{})
}

func (msg *ReliableMessage) SetMeta(meta map[string]interface{}) {
	msg.Set("meta", meta)
}

/**
 *  Sender's Profile
 *  ~~~~~~~~~~~~~~~~
 *  Extends for the first message package of 'Handshake' protocol.
 *
 * @param profile - Profile object or dictionary
 */
func (msg *ReliableMessage) GetProfile() map[string]interface{} {
	value := msg.Get("profile")
	if value == nil {
		return nil
	}
	return value.(map[string]interface{})
}

func (msg *ReliableMessage) SetProfile(profile map[string]interface{}) {
	msg.Set("profile", profile)
}

/*
 *  Verify the Reliable Message to Secure Message
 *
 *    +----------+      +----------+
 *    | sender   |      | sender   |
 *    | receiver |      | receiver |
 *    | time     |  ->  | time     |
 *    |          |      |          |
 *    | data     |      | data     |  1. verify(data, signature, sender.PK)
 *    | key/keys |      | key/keys |
 *    | signature|      +----------+
 *    +----------+
 */

/**
 *  Verify 'data' and 'signature' field with sender's public key
 *
 * @return SecureMessage object
 */
func (msg *ReliableMessage) Verify() *SecureMessage {
	data := msg.GetData()
	if data == nil {
		panic("failed to decode content data")
	}
	signature := msg.GetSignature()
	if signature == nil {
		panic("failed to decode message signature")
	}
	sender := msg.GetSender()
	// 1. verify data signature with sender's public key
	handler := *msg.GetDelegate()
	if handler.VerifyDataSignature(data, signature, sender, msg) {
		// 2. pack message
		info := msg.CopyMap()
		delete(info, "signature")
		return CreateSecureMessage(&info)
	} else {
		//panic("message signature not match")
		return nil
	}
}
