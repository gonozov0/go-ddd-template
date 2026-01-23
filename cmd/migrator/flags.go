package main

import (
	"flag"
	"strconv"
	"strings"

	"fmt"
)

type DBFlag string

const (
	Postgres DBFlag = "postgres"
	YDB      DBFlag = "ydb"
)

type Flags struct {
	DB               DBFlag
	MigrationVersion *uint
}

func parseFlags() (Flags, error) {
	var (
		versionFlag string
		dbFlag      string
	)

	flag.StringVar(&versionFlag, "version", "", "The version of the migration to apply")
	flag.StringVar(&dbFlag, "db", "", "Database type (postgres|ydb)")
	flag.Parse()

	db, err := parseDBFlag(dbFlag)
	if err != nil {
		return Flags{}, fmt.Errorf("failed to parse db flag: %w", err)
	}

	migrationVersion, err := parseMigrationVersion(versionFlag)
	if err != nil {
		return Flags{}, fmt.Errorf("failed to parse migration's version: %w", err)
	}

	return Flags{
		DB:               db,
		MigrationVersion: migrationVersion,
	}, nil
}

func parseDBFlag(value string) (DBFlag, error) {
	if value == "" {
		return "", fmt.Errorf("database type must be specified with -db flag (postgres|ydb)")
	}

	switch strings.ToLower(value) {
	case string(Postgres):
		return Postgres, nil
	case string(YDB):
		return YDB, nil
	default:
		return "", fmt.Errorf("unknown db type: %s", value)
	}
}

func parseMigrationVersion(value string) (*uint, error) {
	if value == "" {
		return nil, nil
	}

	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("failed parse uint: %w", err)
	}

	version := uint(v)

	return &version, nil
}
