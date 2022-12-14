// Copyright (c) 2021 MacEwan University. All rights reserved.
//
// This source code is licensed under the MIT-style license found in
// the LICENSE file in the root directory of this source tree.

// Package main implements a minimal working example of some the LTI library features. For simplicity, all data
// (registrations, deployments, ...) are nonpersistent and stored in the LTI library's internal nonpersistent store.
//
// On startup, the program loads all configuration data from environment variables.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/macewan-cs/lti-example/internal/env"
	lti "github.com/macewan-cs/lti-example/pkg"
	"github.com/macewan-cs/lti-example/pkg/datastore"
	"github.com/macewan-cs/lti-example/pkg/datastore/nonpersistent"
)

const keyID = "defaultKey"

// nonpersistentConfig returns a datastore.Config, which is suitable for creating LTI login handlers, LTI launch
// handlers, and after a launch, LTI connectors.
func nonpersistentConfig() datastore.Config {
	// Retrieve the registration details from environment variables.
	registration := env.RegistrationFromEnvironment()
	err := nonpersistent.DefaultStore.StoreRegistration(registration)
	if err != nil {
		log.Fatalf("registration store error: %v", err)
	}

	// Retrieve deployment details from environment variables.
	deployment := env.DeploymentFromEnvironment()
	err = nonpersistent.DefaultStore.StoreDeployment(registration.Issuer, deployment)
	if err != nil {
		log.Fatalf("deployment store error: %v", err)
	}

	// The default datastore configuration uses nonpersistent.DefaultStore.
	return lti.NewDatastoreConfig()
}

// postLaunchHandler returns an http.HandlerFunc suitable for the second argument of lti.NewLaunch.
func postLaunchHandler(datastoreConfig datastore.Config) http.HandlerFunc {
	// Retrieve the key from environment variables.
	//key := env.KeyFromEnvironment()

	return func(w http.ResponseWriter, r *http.Request) {
		// Create a connector, which is necessary to access LTI services.
		//conn, err := connector.New(datastoreConfig, lti.LaunchIDFromRequest(r), keyID)
		//if err != nil {
		//	log.Printf("cannot create connector for launch: %v", err)
		//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		//	return
		//}
		//
		//conn.SetSigningKey(key.Private)
		//
		//// Upgrade the connector to access Name and Role Provisioning Services.
		//nrps, err := conn.UpgradeNRPS()
		//if err != nil {
		//	log.Printf("cannot upgrade connector for NRPS: %v", err)
		//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		//	return
		//}
		//
		//// Get membership to demonstrate access to NRPS.
		//membership, err := nrps.GetMembership()
		//if err != nil {
		//	log.Printf("cannot get membership: %v", err)
		//	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		//	return
		//}
		lc := lti.LaunchCtxFromContext(r.Context())
		if pre, ok := lc.Token.Get("https://purl.imsglobal.org/spec/lti/claim/launch_presentation");ok {
			if lp, ok := pre.(map[string]interface{});ok {
				docTarget := lp["document_target"]
				if docTarget == "iframe" {
					fmt.Fprintf(w, `<p>This is the iframe launch!</p>
<p>Launch ID from request: %s</p>`, lc.LaunchId)
					return
				}
			}

		}

		fmt.Fprintf(w, `<p>Launch successful!</p>
<p>Launch ID from request: %s</p>`, lc.LaunchId)
//
//		fmt.Fprintf(w, `<p>Launch successful!</p>
//<p>Launch ID from request: %s</p>
//<p>Course title: %s</p>`, lti.LaunchIDFromRequest(r), membership.Context.Title)
	}
}

// logRequest logs a request made to the HTTP server.
func logRequest(r *http.Request) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.Encode(struct {
		RequestURI string `json:"requestUri"`
		Method     string `json:"method"`
		RemoteAddr string `json:"remoteAddr"`
	}{
		RequestURI: r.RequestURI,
		Method:     r.Method,
		RemoteAddr: r.RemoteAddr,
	})
}

func main() {
	var httpAddr = flag.String("addr", ":9000", "example app listen address")
	flag.Parse()

	os.Setenv("REG_ISSUER", "https://edmodoworld.com")
	os.Setenv("REG_CLIENTID", "clientid")
	os.Setenv("REG_KEYSETURI", "http://localhost:8000/certs")//get public cert to decrypt jwt token
	os.Setenv("REG_AUTHTOKENURI", "http://localhost:8000/token")
	os.Setenv("REG_AUTHLOGINURI", "http://localhost:8000/auth")
	os.Setenv("REG_TARGETLINKURI", "http://localhost:9000/launch")
	os.Setenv("DEP_DEPLOYMENTID", "1")
	key, _ := os.ReadFile("../../private.pem")
	fmt.Print(string(key))
	os.Setenv("KEY_PRIVATE", string(key))
	datastoreConfig := nonpersistentConfig()
	http.Handle("/login", lti.NewLogin(datastoreConfig))
	http.Handle("/launch", lti.NewLaunch(datastoreConfig,
		postLaunchHandler(datastoreConfig)))

	log.Printf("Listening for connections on %s...\n", *httpAddr)
	err := http.ListenAndServe(*httpAddr,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logRequest(r)
			http.DefaultServeMux.ServeHTTP(w, r)
		}),
	)
	if err != nil {
		log.Fatalf("http server error: %v", err)
	}
}
