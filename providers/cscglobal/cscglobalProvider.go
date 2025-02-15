package cscglobal

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/StackExchange/dnscontrol/v4/providers"
)

/*

CSC Global Registrar:

Info required in `creds.json`:
   - api-key             Api Key
   - user-token          User Token
   - notification_emails (optional) Comma separated list of email addresses to send notifications to
*/

type providerClient struct {
	key          string
	token        string
	notifyEmails []string
}

var features = providers.DocumentationNotes{
	providers.CanGetZones:            providers.Can(),
	providers.CanUseCAA:              providers.Can(),
	providers.CanUseSRV:              providers.Can(),
	providers.DocOfficiallySupported: providers.Can(),
}

// Set cscDebug to true if you want to see the JSON of important API requests and responses.
var cscDebug = false

func newReg(conf map[string]string) (providers.Registrar, error) {
	return newProvider(conf)
}

func newDsp(conf map[string]string, meta json.RawMessage) (providers.DNSServiceProvider, error) {
	return newProvider(conf)
}

func newProvider(m map[string]string) (*providerClient, error) {
	api := &providerClient{}

	api.key, api.token = m["api-key"], m["user-token"]
	if api.key == "" || api.token == "" {
		return nil, fmt.Errorf("missing CSC Global api-key and/or user-token")
	}

	if m["notification_emails"] != "" {
		api.notifyEmails = strings.Split(m["notification_emails"], ",")
	}

	return api, nil
}

func init() {
	providers.RegisterRegistrarType("CSCGLOBAL", newReg)

	fns := providers.DspFuncs{
		Initializer:   newDsp,
		RecordAuditor: AuditRecords,
	}
	providers.RegisterDomainServiceProviderType("CSCGLOBAL", fns, features)
}
