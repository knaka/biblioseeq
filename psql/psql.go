package psql

import (
	. "github.com/knaka/go-utils"
	common "github.com/knaka/magefiles-common"
	"github.com/magefile/mage/mg"
	"os"
	"strconv"
)

func setVerbose() (err error) {
	defer Catch(&err)
	V0(os.Setenv(mg.VerboseEnv, strconv.FormatBool(true)))
	return nil
}

// Psql executes psql(1) command on the database defined in .env* file
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Psql() (err error) {
	defer Catch(&err)
	V0(setVerbose())
	for _, dbUrl := range []string{
		os.Getenv("DB_URL"),
		os.Getenv("ADMIN_DB_URL"),
	} {
		err = common.RunWith("", nil, "psql", dbUrl)
		if err == nil {
			return nil
		}
	}
	return err
}
