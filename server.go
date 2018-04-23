package main

import (
	"TFRP/azurerm"
	"TFRP/datadog"
	"TFRP/kubernetes"
	"TFRP/pkg/core/consts"
	"TFRP/pkg/core/engines"
	"TFRP/pkg/core/storage"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/gorilla/mux"
	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/spf13/pflag"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	addr       = pflag.String("insecure-address", ":8080", "The <host>:<port> for insecure (HTTP) serving")
	secureAddr = pflag.String("secure-address", ":443", "The <host>:<port> for secure (HTTPS) serving")
)

func main() {
	// router := mux.NewRouter()
	// router.NotFoundHandler = http.HandlerFunc(NotFound)
	// router.HandleFunc("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroup}/providers/Microsoft.Terraform-OSS/provider/{provider}", putProvider).Methods("PUT")
	// router.HandleFunc("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroup}/providers/Microsoft.Terraform-OSS/resource/{resource}", getResource).Methods("GET")
	// router.HandleFunc("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroup}/providers/Microsoft.Terraform-OSS/resource/{resource}", putResource).Methods("PUT")
	// router.HandleFunc("/subscriptions/{subscriptionId}/resourcegroups/{resourceGroup}/providers/Microsoft.Terraform-OSS/resource/{resource}", deleteResource).Methods("DELETE")
	// //log.Fatal(http.ListenAndServeTLS(":443", "fullchain.pem", "privkey.pem", router))
	// log.Fatal(http.ListenAndServe(":8080", router))

	pflag.Parse()

	initRoutes()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initRoutes() {
	webService := new(restful.WebService)
	webService.
		Path(consts.SubscriptionsURLPrefix).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	addProvidersOperationRoutes(webService)
	addResourcesOperationRoutes(webService)

	restful.Add(webService)
}

func addProvidersOperationRoutes(webService *restful.WebService) {
	webService.Route(webService.
		PUT(consts.ProviderRegistrationOperationRoute).
		To(putProviderRegistrationController).
		Doc("Create/update a provider registration").
		Operation(consts.PutProviderRegistrationControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathProviderRegistrationParameter, "Name of provider registration").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))
}

func addResourcesOperationRoutes(webService *restful.WebService) {
	webService.Route(webService.
		GET(consts.ResourceOperationRoute).
		To(getResourceController).
		Doc("Get a resource").
		Operation(consts.GetResourceControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceNameParameter, "Name of resource").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))

	webService.Route(webService.
		PUT(consts.ResourceOperationRoute).
		To(putResourceController).
		Doc("Create/update a resource").
		Operation(consts.PutResourceControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceNameParameter, "Name of resource").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))

	webService.Route(webService.
		DELETE(consts.ResourceOperationRoute).
		To(deleteResourceController).
		Doc("Delete a resource").
		Operation(consts.DeleteResourceControllerName).
		Param(webService.PathParameter(consts.PathSubscriptionIDParameter, "Identifier of customer subscription").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceGroupNameParameter, "Name of resource group").DataType("string")).
		Param(webService.PathParameter(consts.PathResourceNameParameter, "Name of resource").DataType("string")).
		Param(webService.QueryParameter(consts.RequestAPIVersionParameterName, "API Version").DataType("string")))
}

type Error struct {
	Error ErrorDetails `json:"error,omitempty"`
}

type ErrorDetails struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type Package struct {
	Id           bson.ObjectId `bson:"_id,omitempty"`
	ResourceId   string
	ProviderName string
	Config       string
}

type ResourcePackage struct {
	Id           bson.ObjectId `bson:"_id,omitempty"`
	ResourceID   string
	StateID      string
	Config       string
	ResourceType string
}

type Provider struct {
	Location   string
	Properties PoviderProperties
}

type PoviderProperties struct {
	ProviderName string
	Settings     PoviderSettings
}

type PoviderSettings struct {
	Config string
}

type Resource struct {
	Location   string
	Properties ResourceProperties
}

type ResourceProperties struct {
	ProviderID   string
	ResourceName string
	Settings     interface{}
}

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider
var datadogTestAccProvider *schema.Provider
var kubernetesTestAccProvider *schema.Provider
var database string
var password string
var dialInfo *mgo.DialInfo

func init() {
	testAccProvider = azurerm.Provider().(*schema.Provider)
	datadogTestAccProvider = datadog.Provider().(*schema.Provider)
	kubernetesTestAccProvider = kubernetes.Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"azurerm":    testAccProvider,
		"datadog":    datadogTestAccProvider,
		"kubernetes": kubernetesTestAccProvider,
	}
	database = "tfrp001"
	password = "TXWxRsEbZBrBUCJaq3Zu2NqdfafLJcdbKu8rJ6dwKBnjRzfSIwJ8vh23gxRof7GNhOgfeZjfqKL1M7fMWiWQEw=="

	// DialInfo holds options for establishing a session with a MongoDB cluster.
	dialInfo = &mgo.DialInfo{
		Addrs:    []string{fmt.Sprintf("%s.documents.azure.com:10255", database)}, // Get HOST + PORT
		Timeout:  60 * time.Second,
		Database: database, // It can be anything
		Username: database, // Username
		Password: password, // PASSWORD
		DialServer: func(addr *mgo.ServerAddr) (net.Conn, error) {
			return tls.Dial("tcp", addr.String(), &tls.Config{})
		},
	}
}

func getKubernetesTemplateInJson(configFile []byte, resource Resource, resourceID string, resourceSpec []byte) string {
	return fmt.Sprintf(`
		{
			"provider": {
				"kubernetes": {
					"inline_config": %q
				}
			},
			"resource": {
				"%s": {
					"%s": %s
				}
			}
		}
`, configFile, resource.Properties.ResourceName, resourceID, string(resourceSpec))
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(Error{Error: ErrorDetails{Code: "RequestUriInvalid", Message: "Invalid request URI."}})
}

func putProvider(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	fullyQualifiedResourceID := "/subscriptions/" + params["subscriptionId"] + "/resourceGroups/" + params["resourceGroup"] + "/providers/Microsoft.Terraform-OSS/provider/" + params["provider"]

	provider := Provider{}
	defer req.Body.Close()
	json.NewDecoder(req.Body).Decode(&provider)

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}

	defer session.Close()

	// SetSafe changes the session safety mode.
	// If the safe parameter is nil, the session is put in unsafe mode, and writes become fire-and-forget,
	// without error checking. The unsafe mode is faster since operations won't hold on waiting for a confirmation.
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode.
	session.SetSafe(&mgo.Safe{})

	// get collection
	collection := session.DB(database).C("provider")

	// insert Document in collection
	err = collection.Insert(&Package{
		ResourceId:   fullyQualifiedResourceID,
		ProviderName: provider.Properties.ProviderName,
		Config:       provider.Properties.Settings.Config,
	})

	if err != nil {
		log.Fatal("Problem inserting data: ", err)
		return
	}

	// Get Document from collection
	result := Package{}
	err = collection.Find(bson.M{"resourceid": fullyQualifiedResourceID}).One(&result)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	responseBody, _ := json.Marshal(provider)
	w.Write(responseBody)
}

// Put resources
func putResource(w http.ResponseWriter, req *http.Request) {
	resource := Resource{}
	defer req.Body.Close()
	json.NewDecoder(req.Body).Decode(&resource)

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}

	defer session.Close()

	// SetSafe changes the session safety mode.
	// If the safe parameter is nil, the session is put in unsafe mode, and writes become fire-and-forget,
	// without error checking. The unsafe mode is faster since operations won't hold on waiting for a confirmation.
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode.
	session.SetSafe(&mgo.Safe{})

	// get collection
	collection := session.DB(database).C("provider")

	// Get Document from collection
	result := Package{}
	err = collection.Find(bson.M{"resourceid": resource.Properties.ProviderID}).One(&result)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	// Prepare config file
	inlineconfig := result.Config
	decoded, _ := base64.StdEncoding.DecodeString(inlineconfig)
	resourceSpec, _ := json.Marshal(resource.Properties.Settings)

	params := mux.Vars(req)
	configFile := getKubernetesTemplateInJson(decoded, resource, params["resource"], resourceSpec)
	fmt.Printf("%s", configFile)

	provider := testAccProviders["kubernetes"]

	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Printf("%s", err)
	}

	// Init provider
	for _, v := range cfg.ProviderConfigs {
		provider.Configure(terraform.NewResourceConfig(v.RawConfig))
	}

	info := &terraform.InstanceInfo{
		Type: resource.Properties.ResourceName,
	}

	for _, v := range cfg.Resources {
		state := new(terraform.InstanceState)
		state.Init()
		diff, err := provider.Diff(info, state, terraform.NewResourceConfig(v.RawConfig))
		if err != nil {
			fmt.Printf("%s", err)
		}

		// Call apply to create resource
		resultState, _ := provider.Apply(info, state, diff)
		fmt.Printf("%s", resultState.ID)

		// Storage operations
		collection := session.DB(database).C("resource")

		fullyQualifiedResourceID := "/subscriptions/" + params["subscriptionId"] + "/resourceGroups/" + params["resourceGroup"] + "/providers/Microsoft.Terraform-OSS/resource/" + params["resource"]

		// insert Document in collection
		err = collection.Insert(&ResourcePackage{
			ResourceID:   fullyQualifiedResourceID,
			StateID:      resultState.ID,
			Config:       configFile,
			ResourceType: resource.Properties.ResourceName,
		})

		if err != nil {
			log.Fatal("Problem inserting data: ", err)
			return
		}
	}

	responseBody, _ := json.Marshal(resource)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

// Delete resources
func deleteResource(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	fullyQualifiedResourceID := "/subscriptions/" + params["subscriptionId"] + "/resourceGroups/" + params["resourceGroup"] + "/providers/Microsoft.Terraform-OSS/resource/" + params["resource"]

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}

	defer session.Close()

	// SetSafe changes the session safety mode.
	// If the safe parameter is nil, the session is put in unsafe mode, and writes become fire-and-forget,
	// without error checking. The unsafe mode is faster since operations won't hold on waiting for a confirmation.
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode.
	session.SetSafe(&mgo.Safe{})

	// get collection
	collection := session.DB(database).C("resource")

	// Get Document from collection
	result := ResourcePackage{}
	err = collection.Find(bson.M{"resourceid": fullyQualifiedResourceID}).One(&result)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	fmt.Printf("%s", result.Config)

	provider := testAccProviders["kubernetes"]

	cfg, err := config.Load(result.Config)
	if err != nil {
		fmt.Printf("%s", err)
	}

	// Init provider
	for _, v := range cfg.ProviderConfigs {
		provider.Configure(terraform.NewResourceConfig(v.RawConfig))
	}

	info := &terraform.InstanceInfo{
		Type: result.ResourceType,
	}

	state := new(terraform.InstanceState)
	state.ID = result.StateID

	diff := new(terraform.InstanceDiff)
	diff.Destroy = true

	// Call apply to delete resource
	resultState, _ := provider.Apply(info, state, diff)

	responseBody, _ := json.Marshal(resultState)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

// Get resources
func getResource(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	fullyQualifiedResourceID := "/subscriptions/" + params["subscriptionId"] + "/resourceGroups/" + params["resourceGroup"] + "/providers/Microsoft.Terraform-OSS/resource/" + params["resource"]

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	session, err := mgo.DialWithInfo(dialInfo)

	if err != nil {
		fmt.Printf("Can't connect to mongo, go error %v\n", err)
		os.Exit(1)
	}

	defer session.Close()

	// SetSafe changes the session safety mode.
	// If the safe parameter is nil, the session is put in unsafe mode, and writes become fire-and-forget,
	// without error checking. The unsafe mode is faster since operations won't hold on waiting for a confirmation.
	// http://godoc.org/labix.org/v2/mgo#Session.SetMode.
	session.SetSafe(&mgo.Safe{})

	// get collection
	collection := session.DB(database).C("resource")

	// Get Document from collection
	result := ResourcePackage{}
	err = collection.Find(bson.M{"resourceid": fullyQualifiedResourceID}).One(&result)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	fmt.Printf("%s", result.Config)

	provider := testAccProviders["kubernetes"]

	cfg, err := config.Load(result.Config)
	if err != nil {
		fmt.Printf("%s", err)
	}

	// Init provider
	for _, v := range cfg.ProviderConfigs {
		provider.Configure(terraform.NewResourceConfig(v.RawConfig))
	}

	info := &terraform.InstanceInfo{
		Type: result.ResourceType,
	}

	state := new(terraform.InstanceState)
	state.Init()
	state.ID = result.StateID

	// Call refresh
	resultState, err := provider.Refresh(info, state)
	if err != nil {
		fmt.Printf("%s", err)
	}

	responseBody, _ := json.Marshal(resultState)
	w.Header().Set("Content-Type", "application/json")
	w.Write(responseBody)
}

func putProviderRegistrationController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedProviderRegistrationID(request)

	provider := Provider{}
	rawBody, err := ioutil.ReadAll(request.Request.Body)
	err = json.Unmarshal(rawBody, &provider)

	// insert Document in collection
	err = storage.GetProviderRegistrationDataProvider().Insert(&Package{
		ResourceId:   fullyQualifiedResourceID,
		ProviderName: provider.Properties.ProviderName,
		Config:       provider.Properties.Settings.Config,
	})

	if err != nil {
		log.Fatal("Problem inserting data: ", err)
		return
	}

	// Get Document from collection
	result := Package{}
	err = storage.GetProviderRegistrationDataProvider().Find(bson.M{"resourceid": fullyQualifiedResourceID}, &result)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	responseBody, _ := json.Marshal(provider)
	response.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	response.Write(responseBody)
}

// Get resources
func getResourceController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedResourceID(request)

	// Get Document from collection
	result := ResourcePackage{}
	err := storage.GetResourceDataProvider().Find(bson.M{"resourceid": fullyQualifiedResourceID}, &result)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	fmt.Printf("%s", result.Config)

	provider := testAccProviders["kubernetes"]

	cfg, err := config.Load(result.Config)
	if err != nil {
		fmt.Printf("%s", err)
	}

	// Init provider
	for _, v := range cfg.ProviderConfigs {
		provider.Configure(terraform.NewResourceConfig(v.RawConfig))
	}

	info := &terraform.InstanceInfo{
		Type: result.ResourceType,
	}

	state := new(terraform.InstanceState)
	state.Init()
	state.ID = result.StateID

	// Call refresh
	resultState, err := provider.Refresh(info, state)
	if err != nil {
		fmt.Printf("%s", err)
	}

	responseBody, _ := json.Marshal(resultState)
	response.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	response.Write(responseBody)
}

// Put resources
func putResourceController(request *restful.Request, response *restful.Response) {
	resource := Resource{}

	rawBody, err := ioutil.ReadAll(request.Request.Body)
	err = json.Unmarshal(rawBody, &resource)

	// Get Document from collection
	result := Package{}
	err = storage.GetProviderRegistrationDataProvider().Find(bson.M{"resourceid": resource.Properties.ProviderID}, &result)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	// Prepare config file
	inlineconfig := result.Config
	decoded, _ := base64.StdEncoding.DecodeString(inlineconfig)
	resourceSpec, _ := json.Marshal(resource.Properties.Settings)

	configFile := getKubernetesTemplateInJson(decoded, resource, engines.GetResourceName(request), resourceSpec)
	fmt.Printf("%s", configFile)

	provider := testAccProviders["kubernetes"]

	cfg, err := config.Load(configFile)
	if err != nil {
		fmt.Printf("%s", err)
	}

	// Init provider
	for _, v := range cfg.ProviderConfigs {
		provider.Configure(terraform.NewResourceConfig(v.RawConfig))
	}

	info := &terraform.InstanceInfo{
		Type: resource.Properties.ResourceName,
	}

	for _, v := range cfg.Resources {
		state := new(terraform.InstanceState)
		state.Init()
		diff, err := provider.Diff(info, state, terraform.NewResourceConfig(v.RawConfig))
		if err != nil {
			fmt.Printf("%s", err)
		}

		// Call apply to create resource
		resultState, _ := provider.Apply(info, state, diff)
		fmt.Printf("%s", resultState.ID)

		fullyQualifiedResourceID := engines.GetFullyQualifiedResourceID(request)

		// insert Document in collection
		err = storage.GetResourceDataProvider().Insert(&ResourcePackage{
			ResourceID:   fullyQualifiedResourceID,
			StateID:      resultState.ID,
			Config:       configFile,
			ResourceType: resource.Properties.ResourceName,
		})

		if err != nil {
			log.Fatal("Problem inserting data: ", err)
			return
		}
	}

	responseBody, _ := json.Marshal(resource)
	response.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	response.Write(responseBody)
}

// Delete resources
func deleteResourceController(request *restful.Request, response *restful.Response) {
	fullyQualifiedResourceID := engines.GetFullyQualifiedResourceID(request)

	// Get Document from collection
	result := ResourcePackage{}
	err := storage.GetResourceDataProvider().Find(bson.M{"resourceid": fullyQualifiedResourceID}, &result)
	if err != nil {
		log.Fatal("Error finding record: ", err)
		return
	}

	fmt.Printf("%s", result.Config)

	provider := testAccProviders["kubernetes"]

	cfg, err := config.Load(result.Config)
	if err != nil {
		fmt.Printf("%s", err)
	}

	// Init provider
	for _, v := range cfg.ProviderConfigs {
		provider.Configure(terraform.NewResourceConfig(v.RawConfig))
	}

	info := &terraform.InstanceInfo{
		Type: result.ResourceType,
	}

	state := new(terraform.InstanceState)
	state.ID = result.StateID

	diff := new(terraform.InstanceDiff)
	diff.Destroy = true

	// Call apply to delete resource
	resultState, _ := provider.Apply(info, state, diff)

	responseBody, _ := json.Marshal(resultState)
	response.Header().Set(restful.HEADER_ContentType, restful.MIME_JSON)
	response.Write(responseBody)
}
