/*
 * Copyright 2018 The OpenWallet Authors
 * This file is part of the OpenWallet library.
 *
 * The OpenWallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The OpenWallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package risevision

import (
	"github.com/blocktree/go-owcdrivers/addressEncoder"
	"github.com/blocktree/risevision-adapter/risevision_addrdec"
)

type AddressDecoder struct {
	wm *WalletManager //钱包管理者
}

//NewAddressDecoder 地址解析器
func NewAddressDecoder(wm *WalletManager) *AddressDecoder {
	decoder := AddressDecoder{}
	decoder.wm = wm
	return &decoder
}

//PrivateKeyToWIF 私钥转WIF
func (decoder *AddressDecoder) PrivateKeyToWIF(priv []byte, isTestnet bool) (string, error) {
	wif := addressEncoder.AddressEncode(priv, addressEncoder.BTC_mainnetPrivateWIFCompressed)
	return wif, nil
}

//PublicKeyToAddress 公钥转地址
func (decoder *AddressDecoder) PublicKeyToAddress(pub []byte, isTestnet bool) (string, error) {
	address, err := risevision_addrdec.Default.AddressEncode(pub)
	if err != nil {
		return "", err
	}
	return address, nil
}

//RedeemScriptToAddress 多重签名赎回脚本转地址
func (decoder *AddressDecoder) RedeemScriptToAddress(pubs [][]byte, required uint64, isTestnet bool) (string, error) {
	return "", nil
}

//WIFToPrivateKey WIF转私钥
func (decoder *AddressDecoder) WIFToPrivateKey(wif string, isTestnet bool) ([]byte, error) {
	priv, err := addressEncoder.AddressDecode(wif, addressEncoder.BTC_mainnetPrivateWIFCompressed)
	if err != nil {
		return nil, err
	}
	return priv, err
}
