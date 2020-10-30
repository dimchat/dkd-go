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
package types

func CopyMap(origin map[string]interface{}) map[string]interface{} {
	clone := make(map[string]interface{})
	for key, value := range origin {
		clone[key] = value
	}
	return clone
}

type Dictionary struct {
	dictionary *map[string]interface{}
}

func (dict *Dictionary) LoadDictionary(dictionary *map[string]interface{}) *Dictionary {
	dict.dictionary = dictionary
	return dict
}

func (dict *Dictionary) InitDictionary() *Dictionary {
	dictionary := make(map[string]interface{})
	return dict.LoadDictionary(&dictionary)
}

func (dict *Dictionary) CopyMap() map[string]interface{} {
	return CopyMap(dict.GetMap())
}

func (dict *Dictionary) GetMap() map[string]interface{} {
	return *dict.dictionary
}

func (dict *Dictionary) Get(key string) interface{} {
	return (*dict.dictionary)[key]
}

func (dict *Dictionary) Set(key string, value interface{}) {
	if value == nil {
		delete(*dict.dictionary, key)
	} else {
		(*dict.dictionary)[key] = value
	}
}
