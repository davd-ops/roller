package testutils

import (
	"path/filepath"

	"errors"

	initconfig "github.com/dymensionxyz/roller/cmd/config/init"
	"github.com/dymensionxyz/roller/cmd/consts"
)

func getRollappKeysDir(root string) string {
	return filepath.Join(root, consts.ConfigDirName.Rollapp, innerKeysDirName)
}

func VerifyRollappKeys(root string) error {
	rollappKeysDir := getRollappKeysDir(root)
	sequencerKeyInfoPath := filepath.Join(rollappKeysDir, consts.KeyNames.RollappSequencer+".info")
	if err := verifyFileExists(sequencerKeyInfoPath); err != nil {
		return err
	}
	relayerKeyInfoPath := filepath.Join(rollappKeysDir, consts.KeyNames.HubSequencer+".info")
	if err := verifyFileExists(relayerKeyInfoPath); err != nil {
		return err
	}
	for i := 0; i < 2; i++ {
		err := verifyAndRemoveFilePattern(addressPattern, rollappKeysDir)
		if err != nil {
			return err
		}
	}
	nodeKeyPath := getNodeKeyPath(root)
	if err := verifyFileExists(nodeKeyPath); err != nil {
		return err
	}
	privValKeyPath := getPrivValKeyPath(root)
	if err := verifyFileExists(privValKeyPath); err != nil {
		return err
	}
	return nil
}

func getNodeKeyPath(root string) string {
	return filepath.Join(initconfig.RollappConfigDir(root), "node_key.json")
}

func getPrivValKeyPath(root string) string {
	return filepath.Join(initconfig.RollappConfigDir(root), "priv_validator_key.json")
}

func SanitizeGenesis(genesisPath string) error {
	params := []initconfig.PathValue{
		{
			Path:  "genesis_time",
			Value: "PLACEHOLDER_TIMESTAMP",
		},
		{
			Path:  "app_state.auth.accounts.0.base_account.address",
			Value: "PLACEHOLDER_ADDRESS",
		},
		{
			Path:  "app_state.bank.balances.0.address",
			Value: "PLACEHOLDER_ADDRESS",
		},
	}

	err := initconfig.UpdateJSONParams(genesisPath, params)
	if err != nil {
		return err
	}
	return nil
}

func VerifyRollerConfig(rollappConfig initconfig.InitConfig) error {
	existingConfig, err := initconfig.LoadConfigFromTOML(rollappConfig.Home)
	if err != nil {
		return err
	}
	if rollappConfig == existingConfig {
		return nil
	}
	return errors.New("roller config does not match")
}