// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

//go:generate make mdatagen

// Package ambientweatherreceiver implements a receiver that hosts an HTTP server
// to ingest metrics pushed by Ambient Weather Wi-Fi consoles configured in
// "custom server" mode (AMBWeather upload format).
package ambientweatherreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/ambientweatherreceiver"
