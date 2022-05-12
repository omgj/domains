package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/k0kubun/pp"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"google.golang.org/api/run/v1"
)

// RunDomainMapping requires service name already deployed and domain already verified
func RunDomainMapping(sname, domain string) {
	fmt.Printf("Domain Mapping. Service: %s\tDomain: %s\n", sname, domain)
	c, err := run.NewService(context.Background(), option.WithEndpoint(RunAdminEndpoint))
	if err != nil {
		log.Fatal(fmt.Errorf("run client init failed: %w", err))
	}
	dm := &run.DomainMapping{
		ApiVersion: "domains.cloudrun.com/v1",
		Kind:       DomainMappingKind,
		Metadata: &run.ObjectMeta{
			Name:      domain,
			Namespace: DomainMappingNamespace,
		},
		Spec: &run.DomainMappingSpec{
			CertificateMode: DomainMappingCertMode,
			RouteName:       sname,
		},
	}

	fmt.Println("mapping domain")
	yy, e := c.Namespaces.Domainmappings.Create("namespaces/domainsd", dm).Do()
	bug(e)
	pp.Print(yy.Status)
	for _, a := range yy.Status.ResourceRecords {
		fmt.Println(a)
	}
}

// RunDeploy requires a service name and an image to deploy.
func RunDeploy(sname string) {
	fmt.Println("Creating Cloud Run Service in ")
	runc, e := run.NewService(context.Background(), option.WithEndpoint(RunAdminEndpoint))
	if e != nil {
		log.Fatal(fmt.Errorf("run client init fail: %w", e))
	}

	svc := &run.Service{
		ApiVersion: RunAPIVersion,
		Kind:       RunType,
		Metadata: &run.ObjectMeta{
			Name: sname,
		},
		Spec: &run.ServiceSpec{
			Template: &run.RevisionTemplate{
				Metadata: &run.ObjectMeta{
					Name: sname + "-v1",
				},
				Spec: &run.RevisionSpec{
					Containers: []*run.Container{
						{
							Image: RunImage, // us-central1-docker.pkg.dev/domainsd/cloud-run-source-deploy/domainsd
						},
					},
				},
			},
		},
	}

	_, e = runc.Namespaces.Services.Create(ProjectNamespace, svc).Do()
	if e != nil {
		log.Fatal(fmt.Errorf("service creation fail: %w", e))
	}
	log.Printf("service create call completed")

	log.Printf("waiting for service to become ready")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	e = RunOp(ctx, runc, sname, "Ready")
	if e != nil {
		log.Fatal(e)
	}
	e = RunOp(ctx, runc, sname, "RoutesReady")
	if e != nil {
		log.Fatal(e)
	}

	log.Printf("service is ready and serving traffic!")

	// give service public access via IAM bindings.
	// we'll need to use the non-regional API endpoint with this.
	gc, err := run.NewService(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	_, err = gc.Projects.Locations.Services.SetIamPolicy(
		fmt.Sprintf("projects/domainsd/locations/us-central1/services/%s", sname),
		&run.SetIamPolicyRequest{
			Policy: &run.Policy{Bindings: []*run.Binding{{
				Members: []string{"allUsers"},
				Role:    "roles/run.invoker",
			}}},
		},
	).Do()
	if err != nil {
		log.Printf("failed setting IAM: %s", err)
	}

	// print the service URL by re-querying the service because the
	// url becomes available on the object after the Create() call
	svc, err = RunQueryService(runc, sname)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("service is deployed at: %s", svc.Status.Address.Url)
}

// RunOp waits for a service to deployed
func RunOp(ctx context.Context, c *run.APIService, sname, condition string) error {
	t := time.NewTicker(time.Second * 5)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			svc, err := RunQueryService(c, sname)
			if err != nil {
				return fmt.Errorf("failed to query service for readiness: %w", err)
			}
			for _, c := range svc.Status.Conditions {
				if c.Type == condition {
					if c.Status == "True" {
						return nil
					} else if c.Status == "False" {
						return fmt.Errorf("service could not become %q (status:%s) (reason:%s) %s",
							condition, c.Status, c.Reason, c.Message)
					}
				}
			}
		}
	}
}

// RunServiceExists will tell you if the service exists.
func RunServiceExists(c *run.APIService, region, project, name string) (bool, error) {
	_, err := c.Namespaces.Services.Get(fmt.Sprintf("namespaces/%s/services/%s", project, name)).Do()
	if err == nil {
		return true, nil
	}
	// not all errors indicate service does not exist, look for 404 status code
	v, ok := err.(*googleapi.Error)
	if !ok {
		return false, fmt.Errorf("failed to query service: %w", err)
	}
	if v.Code == http.StatusNotFound {
		return false, nil
	}
	return false, fmt.Errorf("unexpected status code=%d from get service call: %w", v.Code, err)
}

func RunQueryService(c *run.APIService, sname string) (*run.Service, error) {
	return c.Namespaces.Services.Get(fmt.Sprintf("namespaces/domainsd/services/%s", sname)).Do()
}
