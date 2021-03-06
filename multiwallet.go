package multiwallet

import (
	"github.com/OpenBazaar/multiwallet/bitcoin"
	"github.com/OpenBazaar/multiwallet/config"
	"github.com/OpenBazaar/wallet-interface"
	"github.com/op/go-logging"
	"github.com/tyler-smith/go-bip39"
	"time"
)

var log = logging.MustGetLogger("multiwallet")

type MultiWallet map[wallet.CoinType]wallet.Wallet

func NewMultiWallet(cfg *config.Config) (MultiWallet, error) {
	log.SetBackend(logging.AddModuleLevel(cfg.Logger))

	if cfg.Mnemonic == "" {
		ent, err := bip39.NewEntropy(128)
		if err != nil {
			return nil, err
		}
		mnemonic, err := bip39.NewMnemonic(ent)
		if err != nil {
			return nil, err
		}
		cfg.Mnemonic = mnemonic
		cfg.CreationDate = time.Now()
	}

	multiwallet := make(MultiWallet) // TODO: change to wallet interface when BitcoinWallet conforms
	var err error
	for _, coin := range cfg.Coins {
		var w *bitcoin.BitcoinWallet // TODO: change to wallet interface when BitcoinWallet conforms
		switch coin.CoinType {
		case wallet.Bitcoin:
			w, err = bitcoin.NewBitcoinWallet(coin, cfg.Mnemonic, cfg.Params, cfg.Proxy)
			if err != nil {
				return nil, err
			}
			multiwallet[coin.CoinType] = w
		}
	}
	return multiwallet, nil
}

func (w *MultiWallet) Start() {
	for _, wallet := range *w {
		wallet.Start()
	}
}
