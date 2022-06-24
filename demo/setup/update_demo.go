package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	math_rand "math/rand"
	"os"
	"regexp"
	"strings"

	"github.com/run-x/cloudgrep/pkg/config"
	"github.com/run-x/cloudgrep/pkg/datastore"
	"github.com/run-x/cloudgrep/pkg/model"
	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

const configInput = "./demo/setup/demo-setup.yaml"
const configOutput = "./demo/demo.yaml"
const targetDB = "demo/demo.db"

//generate the demo data from the original dump
func main() {

	//clean up previous target db
	e := os.Remove(targetDB)
	if e != nil {
		log.Fatal(e)
	}

	input := readConfig(configInput)
	output := readConfig(configOutput)
	ctx := input.ctx

	//copy DB over
	allResources := input.getResources(nil)
	output.logger.Sugar().Infof("Starting with original DB: %v resources", len(allResources))
	if err := output.ds.WriteResources(output.ctx, allResources); err != nil {
		log.Fatal(e)
	}

	//copy the tags from the ec2.instance to their ec2.NetworkInterface
	for _, instance := range input.getResources(map[string]string{"type": "ec2.Instance"}) {
		eni := getEni(instance)
		if eni != "" {
			eniResource := input.getResources(map[string]string{"id": eni})[0]
			eniResource.Tags = instance.Tags
			if err := output.ds.WriteResources(ctx, model.Resources{eniResource}); err != nil {
				log.Fatal(e)
			}
		}
	}

	// add a few more teams, the orignal teams are consumer & marketplace
	newTeams := map[string](string){
		"order-management": "consumer",
		"billing":          "consumer",
		"data-management":  "marketplace"}
	for team, original := range newTeams {
		toCopy := output.getResources(map[string]string{"team": original})
		fnResource := func(r model.Resource) model.Resource {
			r.Tags = r.Tags.Delete("team").Add("team", team)
			return r
		}
		//add a few random resources to make it more realistic (so they don't all have the same number)
		extra := toCopy[0:math_rand.Intn(len(toCopy)/2)]
		output.copyResources(append(toCopy, extra...), fnResource)
	}

	//assign 95% of the ec2.NetworkInterface, ec2.SecurityGroup, ec2.Subnet without tags to infra team
	//leave some aside for default vpc and such
	for _, _type := range []string{"ec2.NetworkInterface", "ec2.SecurityGroup", "ec2.Subnet"} {
		toUpdate := output.getResources(map[string]string{"type": _type, "team": "(missing)"})
		fnResource := func(r model.Resource) model.Resource {
			r.Tags = model.Tags{model.Tag{Key: "team", Value: "infra"}}
			return r
		}
		output.updateResources(percent(toUpdate, 95), fnResource)
	}

	// set a wrong value -> prod -> production for 1 resource
	output.updateTag(
		map[string]string{
			"env":    "prod",
			"type":   "eks.Cluster",
			"market": "Europe",
			"team":   "marketplace",
		},
		"env", "production",
	)

	//set the "managed-by" tag
	// identify the cloudformation stuff
	output.updateTag(
		map[string]string{
			"aws:cloudformation:logical-id": "(not null)",
		},
		"managed-by", "cloudformation",
	)
	// for the rest - 80% is terraform - leave some unassigned for demo
	toUpdate := output.getResources(map[string]string{"managed-by": "(missing)"})
	fnResource := func(r model.Resource) model.Resource {
		r.Tags = model.Tags{model.Tag{Key: "managed-by", Value: "terraform"}}
		return r
	}
	output.updateResources(percent(toUpdate, 80), fnResource)

	// fix a wrong tag (in the original data)
	output.updateTag(
		map[string]string{
			"market": "North America/",
		},
		"market", "North America",
	)

	//create some EC2 intances that looks like manually created (without any tag)
	ec2Instances := output.getResources(map[string]string{"type": "ec2.Instance"})[0:2]
	fnResource = func(r model.Resource) model.Resource {
		r.Tags = model.Tags{}
		return r
	}
	output.copyResources(ec2Instances, fnResource)

	// remove the tag "team" to one of the RDS
	rdsInstance := output.getResources(map[string]string{"type": "rds.DBInstance", "team": "(not null)"})[0:1]
	fnResource = func(r model.Resource) model.Resource {
		team := r.Tags.Find("team").Value
		//remove the tag, add another one to help with the demo
		r.Tags = r.Tags.Delete("team").Add("description", fmt.Sprintf("database for %v data", team))
		return r
	}
	output.updateResources(rdsInstance, fnResource)

	output.logger.Sugar().Infof("Demo data is ready")
}

type demo struct {
	ctx    context.Context
	ds     datastore.Datastore
	logger *zap.Logger
}

func readConfig(path string) *demo {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	var cfg config.Config
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// connect to the database
	ctx := context.Background()
	logger, _ := zap.NewProduction()
	ds, err := datastore.NewDatastore(ctx, cfg, logger)
	if err != nil {
		log.Fatal(err)
	}
	return &demo{
		ctx:    ctx,
		logger: logger,
		ds:     ds,
	}
}

func percent(resources model.Resources, percent int) model.Resources {
	newSize := len(resources) * percent / 100
	return resources[0:newSize]
}

func (d *demo) getResources(filter map[string]string) model.Resources {
	filterObj := map[string]interface{}{
		"filter": filter,
		"limit":  2000,
	}
	query, err := json.Marshal(filterObj)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := d.ds.GetResources(d.ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Resources
}

func (d *demo) updateTag(filter map[string]string, key string, val string) {
	resources := d.getResources(filter)
	for i, r := range resources {
		d.logger.Sugar().Infof("Updating resource %v, setting tag %v=%v", r.Id, key, val)
		t := r.Tags.Find(key)
		if t != nil {
			r.Tags = r.Tags.Delete(key)
			if val != "(missing)" {
				r.Tags = r.Tags.Add(key, val)
			}
		} else {
			//need to create the tag
			r.Tags = r.Tags.Add(key, val)
		}
		resources[i] = r
	}
	err := d.ds.WriteResources(d.ctx, resources)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *demo) updateResources(resources model.Resources, fnResource func(r model.Resource) model.Resource) {
	for i, r := range resources {
		updateR := fnResource(*r)
		resources[i] = &updateR
	}
	d.logger.Sugar().Infof("Updating %v resources", len(resources))
	err := d.ds.WriteResources(d.ctx, resources)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *demo) copyResources(resources model.Resources, fnResource func(r model.Resource) model.Resource) {
	var newResources model.Resources
	for _, r := range resources {
		rNew := *r
		rNew.Id = genId(r.Id)
		rNew = fnResource(rNew)
		newResources = append(newResources, &rNew)
	}
	d.logger.Sugar().Infof("Adding %v resources", len(resources))
	err := d.ds.WriteResources(d.ctx, newResources)
	if err != nil {
		log.Fatal(err)
	}
}

func genId(id string) string {
	var hexaStr string
	if strings.HasPrefix(id, "arn:") {
		//arn:aws:elasticloadbalancing:eu-west-3:309944644246:loadbalancer/app/WebPr-Appli-V4D1N28N99UF/71469eb0acf9280b
		//update the hexa after the /
		hexaStr = right(id, "/")
	} else if strings.Contains(id, "-") {
		hexaStr = right(id, "-")
	} else {
		//replace with a random string
		return randSeq(len(id))
	}
	//update the hexa part
	//ex: i-002c2c600b59bd24c will update 002c2c600b59bd24c
	return strings.Replace(id, hexaStr, randomHex(len(hexaStr)), 1)
}

func right(s string, separator string) string {
	parts := strings.Split(s, separator)
	return parts[len(parts)-1]
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[math_rand.Intn(len(letters))]
	}
	return string(b)
}

// randToken generates a random hex value.
func randomHex(n int) string {
	bytes := make([]byte, n/2)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(bytes)
}

func getEni(r *model.Resource) string {
	var re = regexp.MustCompile(`"(NetworkInterfaceId)":"((\\"|[^"])*)"`)
	matches := re.FindStringSubmatch(string(r.RawData))
	if len(matches) > 2 {
		eni := matches[2]
		return eni
	}
	return ""
}
