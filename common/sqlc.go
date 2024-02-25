package common

import (
	"errors"
	. "github.com/knaka/go-utils"
	"github.com/knaka/magefiles-shared/shdir"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

type Sqlc mg.Namespace

// Gen generates sqlc queries code.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (s Sqlc) Gen() (err error) {
	err = sqlcGen()
	if err != nil {
		return
	}
	err = Gogen("./db/.")
	if err != nil {
		return
	}
	err = sqlcGen()
	return
}

func Gogen(target string) error {
	V0(Gobin("go-generate-fast"))
	return RunWith("", nil, "go-generate-fast", target)
}

func sqlcGen() error {
	paths := V(filepath.Glob("db/schema*.sql"))
	paths = append(paths,
		filepath.Join("db", "migrations"),
		filepath.Join("db", "queries"),
		filepath.Join("db", "sqlc.yaml"),
	)
	sourceNewestStamp := V(target.NewestModTime(
		paths...,
	))
	destPath := filepath.Join("db", "sqlcgen")
	destPathBak := destPath + ".bak"
	destNewestStamp := V(target.NewestModTime(destPath))
	if !destNewestStamp.IsZero() && destNewestStamp.Compare(sourceNewestStamp) > 0 {
		return nil
	}
	V0(os.Rename(destPath, destPathBak))
	err := shdir.RunWith("db", nil, "sqlc", "generate")
	if err == nil {
		V0(sh.Rm(destPathBak))
	} else {
		V0(os.Rename(destPathBak, destPath))
	}
	return err
}

// Vet runs queries through a set of lint rules.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Sqlc) Vet() error {
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return errors.New("no DB_URL")
	}
	u, err := url.ParseRequestURI(dbUrl)
	if err != nil {
		return err
	}
	values := u.Query()
	values.Set("options", "-c enable_seqscan=off")
	u.RawQuery = values.Encode()
	err = os.Setenv("DB_URL", u.String())
	if err != nil {
		log.Panicf("panic 8b901e4 (%v)", err)
	}
	// Linting queries — sqlc 1.23.0 documentation https://docs.sqlc.dev/en/stable/howto/vet.html
	return RunWith("db", map[string]string{
		//"SQLCDEBUG": "dumpvetenv=1,dumpexplain=1", // Environment variables — sqlc 1.23.0 documentation https://docs.sqlc.dev/en/stable/reference/environment-variables.html#sqlcdebug
	}, "sqlc", "vet")
}

func init() {
	AddGenFn(Sqlc.Gen)
}
