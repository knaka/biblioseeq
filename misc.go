package common

import (
	"fmt"
	. "github.com/knaka/go-utils"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/tidwall/gjson"
	"net/url"
	"os"
)

var DirsToCleanUp []string

// Clean cleans up generated files.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Clean() {
	fmt.Println("Cleaning...")
	for _, dir := range DirsToCleanUp {
		Ensure0(os.RemoveAll(dir))
	}
}

//goland:noinspection GoUnusedExportedType, GoUnnecessarilyExportedIdentifiers
type Env mg.Namespace

// Compose prints text in the .env format that references the Docker Compose configuration.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Env) Compose() {
	json := Ensure(sh.Output("docker", "compose", "config", "--format", "json"))
	host := "127.0.0.1"
	publishedPort := gjson.Get(json, "services.db.ports.0.published").Int()
	urlDb := url.URL{
		Scheme: "postgresql",
		Host:   fmt.Sprintf("%s:%d", host, publishedPort),
		User: url.UserPassword(
			gjson.Get(json, "services.ap.environment.DB_USER").String(),
			gjson.Get(json, "services.ap.environment.DB_PASSWORD").String(),
		),
		Path:    "/" + gjson.Get(json, "services.ap.environment.DB_DATABASE").String(),
		RawPath: "sslmode=disable",
	}
	Ensure0(fmt.Fprintf(os.Stdout, "DB_URL=%s\n", urlDb.String()))
}

// Print (host string, port int, user, password, database string) prints text in the .env format that references the CDK configuration.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func (Env) Print(host string, port int, user, password, database string) {
	urlDb := url.URL{
		Scheme: "postgresql",
		Host:   fmt.Sprintf("%s:%d", host, port),
		User: url.UserPassword(
			user,
			password,
		),
		Path:     "/" + database,
		RawQuery: "sslmode=disable",
	}
	Ensure0(fmt.Fprintf(os.Stdout, "DB_URL=%s\n", urlDb.String()))
}

// Lint analyses.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Lint() error {
	mg.Deps(
		Sqlc.Vet,
	)
	return nil
}

// Gen generates binding codes.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Gen() error {
	// Docker build does not generate dependent files.
	if os.Getenv("NO_GEN") == "" && os.Getenv("NO_GENERATE") == "" {
		mg.Deps(
			Sqlc.Gen,
			Bufgen,
			Dockerfiles,
		)
	}
	return nil
}

// Generate is an alias for Gen.
//
//goland:noinspection GoUnusedExportedFunction, GoUnnecessarilyExportedIdentifiers
func Generate() error {
	return Gen()
}
