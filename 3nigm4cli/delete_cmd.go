//
// 3nigm4 3nigm4cli package
// Author: Guido Ronchetti <dyst0ni3@gmail.com>
// v1.0 16/06/2016
//

package main

// Golang std libs
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Internal dependencies
import (
	crypto3n "github.com/nexocrew/3nigm4/lib/crypto"
	fm "github.com/nexocrew/3nigm4/lib/filemanager"
	sc "github.com/nexocrew/3nigm4/lib/storageclient"
)

// Third party libs
import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// DeleteCmd removes remote resources starting from a reference
// file
var DeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "Removes remote resources",
	Long:    "Removes remote resources starting from a reference.",
	Example: "3nigm4cli store delete -r /tmp/resources.3rf -v",
}

func init() {
	// i/o paths
	setArgument(DeleteCmd, "referencein", &arguments.referenceInPath)
	// working queue setup
	setArgument(DeleteCmd, "workerscount", &arguments.workers)
	setArgument(DeleteCmd, "queuesize", &arguments.queue)

	// files parameters
	DeleteCmd.RunE = deleteReference
}

// deleteReference uses datastorage struct to remotely delete all chunks
// pointed by a reference file.
func deleteReference(cmd *cobra.Command, args []string) error {
	// load config file
	err := manageConfigFile()
	if err != nil {
		return err
	}

	// check for token presence
	if pss.Token == "" {
		return fmt.Errorf("you are not logged in, please call \"login\" command before invoking any other functionality")
	}

	// prepare PGP private key
	privateEntityList, err := checkAndLoadPgpPrivateKey(viper.GetString(am["privatekey"].name))
	if err != nil {
		return err
	}

	// create new store manager
	ds, err := sc.NewStorageClient(
		viper.GetString(am["storageaddress"].name),
		viper.GetInt(am["storageport"].name),
		pss.Token,
		viper.GetInt(am["workerscount"].name),
		viper.GetInt(am["queuesize"].name))
	if err != nil {
		return err
	}
	defer ds.Close()

	// get reference
	encBytes, err := ioutil.ReadFile(viper.GetString(am["referencein"].name))
	if err != nil {
		return fmt.Errorf("unable to access reference file: %s", err.Error())
	}
	// decrypt it
	refenceBytes, err := crypto3n.OpenPgpDecrypt(encBytes, privateEntityList)
	if err != nil {
		return fmt.Errorf("unable to decrypt reference file: %s", err.Error())
	}
	// unmarshal it
	var reference fm.ReferenceFile
	err = json.Unmarshal(refenceBytes, &reference)
	if err != nil {
		return fmt.Errorf("unable to decode reference file: %s", err.Error())
	}

	// delete resources from reference
	err = fm.DeleteChunks(ds, &reference)
	if err != nil {
		return err
	}

	return nil
}
