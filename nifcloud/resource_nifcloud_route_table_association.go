package nifcloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/awserr"
	"github.com/shztki/nifcloud-sdk-go/service/computing"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceNifcloudRouteTableAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudRouteTableAssociationCreate,
		Read:   resourceNifcloudRouteTableAssociationRead,
		Update: resourceNifcloudRouteTableAssociationUpdate,
		Delete: resourceNifcloudRouteTableAssociationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceNifcloudRouteTableAssociationImport,
		},

		Schema: map[string]*schema.Schema{
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceNifcloudRouteTableAssociationCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	log.Printf(
		"[INFO] Creating route table association: %s => %s",
		d.Get("router_id").(string),
		d.Get("route_table_id").(string))

	associationOpts := computing.AssociateRouteTableInput{
		RouteTableId: nifcloud.String(d.Get("route_table_id").(string)),
		RouterId:     nifcloud.String(d.Get("router_id").(string)),
		Agreement:    nifcloud.Bool(false),
	}

	//var associationID string
	var resp *computing.AssociateRouteTableOutput
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		var err error
		resp, err = conn.AssociateRouteTable(&associationOpts)
		if err != nil {
			if awsErr, ok := err.(awserr.Error); ok {
				if awsErr.Code() == "Client.InvalidParameterNotFound.RouterId" {
					return resource.RetryableError(awsErr)
				}
			}
			return resource.NonRetryableError(err)
		}
		//associationID = *resp.AssociationId
		return nil
	})
	if isResourceTimeoutError(err) {
		resp, err = conn.AssociateRouteTable(&associationOpts)
	}
	if err != nil {
		return fmt.Errorf("Error creating route table association: %s", err)
	}

	// Set the ID and return
	//d.SetId(associationID)
	//log.Printf("[INFO] Association ID: %s", d.Id())

	return resourceNifcloudRouteTableAssociationRead(d, meta)
}

func resourceNifcloudRouteTableAssociationRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	// Get the routing table that this association belongs to
	rtRaw, _, err := resourceNifcloudRouteTableAssociationStateRefreshFunc(
		conn, d.Get("route_table_id").(string))()
	if err != nil {
		return err
	}
	if rtRaw == nil {
		return nil
	}
	rt := rtRaw.(*computing.RouteTableSetItem)
	log.Printf("[INFO] route table set: %v", rt)
	d.Set("route_table_id", rt.RouteTableId)

	// Inspect that the association exists
	found := false
	for _, a := range rt.AssociationSet {
		log.Printf("[INFO] Association Set: %v", a)
		if *a.RouterId == d.Get("router_id").(string) {
			found = true
			d.Set("router_id", *a.RouterId)
			d.SetId(*a.RouteTableAssociationId)
			break
		}
	}

	if !found {
		// It seems it doesn't exist anymore, so clear the ID
		d.SetId("")
	}

	log.Printf("[INFO] Association ID: %s", d.Id())
	return nil
}

func resourceNifcloudRouteTableAssociationUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	log.Printf(
		"[INFO] Creating route table association: %s => %s",
		d.Get("router_id").(string),
		d.Get("route_table_id").(string))

	req := &computing.ReplaceRouteTableAssociationInput{
		AssociationId: nifcloud.String(d.Id()),
		RouteTableId:  nifcloud.String(d.Get("route_table_id").(string)),
		Agreement:     nifcloud.Bool(false),
	}
	resp, err := conn.ReplaceRouteTableAssociation(req)

	if err != nil {
//		ec2err, ok := err.(awserr.Error)
//		if ok && ec2err.Code() == "Client.InvalidParameterNotFound.AssociationId" {
//			// Not found, so just create a new one
//			return resourceNifcloudRouteTableAssociationCreate(d, meta)
//		}
		return err
	}

	// Update the ID
	d.SetId(*resp.NewAssociationId)
	log.Printf("[INFO] Association ID: %s", d.Id())

	return nil
}

func resourceNifcloudRouteTableAssociationDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	log.Printf("[INFO] Deleting route table association: %s", d.Id())
	_, err := conn.DisassociateRouteTable(&computing.DisassociateRouteTableInput{
		AssociationId: nifcloud.String(d.Id()),
		Agreement:     nifcloud.Bool(false),
	})
	if err != nil {
		ec2err, ok := err.(awserr.Error)
		if ok && ec2err.Code() == "Client.InvalidParameterNotFound.AssociationId" {
			return nil
		}

		return fmt.Errorf("Error deleting route table association: %s", err)
	}

	return nil
}

func resourceNifcloudRouteTableAssociationImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return []*schema.ResourceData{}, fmt.Errorf("Wrong format for import: %s. Use 'router ID/route table ID'", d.Id())
	}

	routerID := parts[0]
	routeTableID := parts[1]

	log.Printf("[DEBUG] Importing route table association, router: %s, route table: %s", routerID, routeTableID)

	conn := meta.(*NifcloudClient).computingconn

	routerFilter := &computing.RequestFilterStruct{
		Name:         nifcloud.String("association.router-id"),
		RequestValue: []*string{nifcloud.String(routerID)},
	}
	routeTableFilter := &computing.RequestFilterStruct{
		Name:         nifcloud.String("association.route-table-id"),
		RequestValue: []*string{nifcloud.String(routeTableID)},
	}
	input := &computing.DescribeRouteTablesInput{
		Filter: []*computing.RequestFilterStruct{routerFilter,routeTableFilter},
	}

	output, err := conn.DescribeRouteTables(input)
	if err != nil || len(output.RouteTableSet) == 0 {
		return nil, fmt.Errorf("Error finding route table: %v", err)
	}

	rt := output.RouteTableSet[0]

	var associationID string
	for _, a := range rt.AssociationSet {
		if nifcloud.StringValue(a.RouterId) == routerID {
			associationID = nifcloud.StringValue(a.RouteTableAssociationId)
			break
		}
	}
	if associationID == "" {
		return nil, fmt.Errorf("Error finding route table, ID: %v", *rt.RouteTableId)
	}

	d.SetId(associationID)
	d.Set("router_id", routerID)
	d.Set("route_table_id", routeTableID)

	return []*schema.ResourceData{d}, nil
}

// resourceNifcloudRouteTableAssociationStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a RouteTableAssociation.
func resourceNifcloudRouteTableAssociationStateRefreshFunc(conn *computing.Computing, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var resp *computing.DescribeRouteTablesOutput
		var err error
		err = resource.Retry(5*time.Minute, func() *resource.RetryError {
			resp, err = conn.DescribeRouteTables(&computing.DescribeRouteTablesInput{
				RouteTableId: []*string{nifcloud.String(id)},
			})
			if err != nil {
				if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "Client.InvalidParameterNotFound.RouteTableId" {
					resp = nil
				} else {
					log.Printf("Error on RouteTableStateRefresh: %s", err)
					return nil
				}
			} else if resp.RouteTableSet[0].AssociationSet == nil {
				return resource.RetryableError(fmt.Errorf("not finding route table association (%s) still deleting", id))
			}
			return nil
		})

		if resp == nil {
			// Sometimes AWS just has consistency issues and doesn't see
			// our instance yet. Return an empty state.
			return nil, "", nil
		}

		rt := resp.RouteTableSet[0]
		return rt, "ready", nil
	}
}