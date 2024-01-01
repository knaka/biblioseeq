package common

import (
	"fmt"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/tidwall/gjson"
	"net/url"
	"os"
)

var DirsToCleanUp []string

// Clean cleans up generated files.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Clean() {
	fmt.Println("Cleaning...")
	for _, dir := range DirsToCleanUp {
		_ = os.RemoveAll(dir)
	}
}

// noinspection GoUnusedExportedType, GoUnnecessarilyExportedIdentifiers
type Env mg.Namespace

// Compose prints text in the .env format that references the Docker Compose configuration.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Env) Compose() {
	json, _ := sh.Output("docker", "compose", "config", "--format", "json")
	host := "127.0.0.1"
	publishedPort := gjson.Get(json, "services.db.ports.0.published").Int()
	urlDb := url.URL{}
	urlDb.Scheme = "postgresql"
	urlDb.Host = fmt.Sprintf("%s:%d", host, publishedPort)
	urlDb.User = url.UserPassword(
		gjson.Get(json, "services.ap.environment.DB_USER").String(),
		gjson.Get(json, "services.ap.environment.DB_PASSWORD").String(),
	)
	urlDb.Path = "/" + gjson.Get(json, "services.ap.environment.DB_DATABASE").String()
	urlDb.RawQuery = "sslmode=disable"
	_, _ = fmt.Fprintf(os.Stdout, "DB_URL=%s\n", urlDb.String())
}

// Print (host string, port int, user, password, database string) prints text in the .env format that references the CDK configuration.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Env) Print(host string, port int, user, password, database string) {
	urlDb := url.URL{}
	urlDb.Scheme = "postgresql"
	urlDb.Host = fmt.Sprintf("%s:%d", host, port)
	urlDb.User = url.UserPassword(
		user,
		password,
	)
	urlDb.Path = "/" + database
	urlDb.RawQuery = "sslmode=disable"
	_, _ = fmt.Fprintf(os.Stdout, "DB_URL=%s\n", urlDb.String())
}

// Gen generates binding codes.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Gen() error {
	// Docker build does not generate dependent files.
	if os.Getenv("NO_GEN") == "" {
		mg.Deps(
			Sqlc.Gen,
			//Bufgen,
			//Dockerfiles,
		)
	}
	return nil
}

// Lint analyses.
//
// noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Lint() error {
	mg.Deps(
		Sqlc.Vet,
	)
	return nil
}
