package nifcloud

import (
	"fmt"
	"log"
	"time"

	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/awserr"
	"github.com/shztki/nifcloud-sdk-go/service/computing"
	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// AWS Route resource Schema declaration
func resourceNifcloudRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudRouteCreate,
		Read:   resourceNifcloudRouteRead,
//		Update: resourceNifcloudRouteUpdate,
		Delete: resourceNifcloudRouteDelete,
		Exists: resourceNifcloudRouteExists,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"destination_cidr_block": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"network_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"route_table_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNifcloudRouteCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	createOpts := &computing.CreateRouteInput{
		RouteTableId: nifcloud.String(d.Get("route_table_id").(string)),
		DestinationCidrBlock:    nifcloud.String(d.Get("destination_cidr_block").(string)),
	}

	if v, ok := d.GetOk("ip_address"); ok {
		createOpts.IpAddress = nifcloud.String(v.(string))
	}

	if v, ok := d.GetOk("network_id"); ok {
		createOpts.NetworkId = nifcloud.String(v.(string))
	}

	if v, ok := d.GetOk("network_name"); ok {
		createOpts.NetworkName = nifcloud.String(v.(string))
	}
	log.Printf("[DEBUG] Route create config: %s", createOpts)

	// Create the route
	var err error

	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err = conn.CreateRoute(createOpts)

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
	if isResourceTimeoutError(err) {
		_, err = conn.CreateRoute(createOpts)
	}
	if err != nil {
		return fmt.Errorf("Error creating route: %s", err)
	}

	var route *computing.RouteSetItem

	if v, ok := d.GetOk("destination_cidr_block"); ok {
		err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
			route, err = resourceNifcloudRouteFindRoute(conn, d.Get("route_table_id").(string), v.(string))
			return resource.RetryableError(err)
		})
		if isResourceTimeoutError(err) {
			route, err = resourceNifcloudRouteFindRoute(conn, d.Get("route_table_id").(string), v.(string))
		}
		if err != nil {
			return fmt.Errorf("Error finding route after creating it: %s", err)
		}
	}

	d.SetId(resourceNifcloudRouteID(d, route))
	resourceNifcloudRouteSetResourceData(d, route)
	return nil
}

func resourceNifcloudRouteRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn
	routeTableId := d.Get("route_table_id").(string)

	destinationCidrBlock := d.Get("destination_cidr_block").(string)

	route, err := resourceNifcloudRouteFindRoute(conn, routeTableId, destinationCidrBlock)
	if err != nil {
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidParameterNotFound.RouteTableId" {
			log.Printf("[WARN] Route Table %q could not be found. Removing Route from state.",
				routeTableId)
			d.SetId("")
			return nil
		}
		return err
	}
	resourceNifcloudRouteSetResourceData(d, route)
	return nil
}

func resourceNifcloudRouteSetResourceData(d *schema.ResourceData, route *computing.RouteSetItem) {
	d.Set("ip_address", route.IpAddress)
	d.Set("network_id", route.NetworkId)
	d.Set("netowork_name", route.NetworkName)
}

/*
func resourceNifcloudRouteUpdate(d *schema.ResourceData, meta interface{}) error {
	return err
}
*/

func resourceNifcloudRouteDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	deleteOpts := &computing.DeleteRouteInput{
		RouteTableId: nifcloud.String(d.Get("route_table_id").(string)),
	}
	if v, ok := d.GetOk("destination_cidr_block"); ok {
		deleteOpts.DestinationCidrBlock = nifcloud.String(v.(string))
	}
	log.Printf("[DEBUG] Route delete opts: %s", deleteOpts)

	var resp *computing.DeleteRouteOutput
	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		log.Printf("[DEBUG] Trying to delete route with opts %s", deleteOpts)
		var err error
		resp, err = conn.DeleteRoute(deleteOpts)
		log.Printf("[DEBUG] Route delete result: %s", resp)

		if isNifcloudErr(err, "Client.ResourceAssociated.RouteTable", "") {
			return resource.RetryableError(err)
		}
		if err == nil {
			return nil
		}

		return resource.NonRetryableError(err)
	})
	if isResourceTimeoutError(err) {
		resp, err = conn.DeleteRoute(deleteOpts)
	}
	if err != nil {
		return fmt.Errorf("Error deleting route: %s", err)
	}
	return nil
}

func resourceNifcloudRouteExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	conn := meta.(*NifcloudClient).computingconn
	routeTableId := d.Get("route_table_id").(string)

	findOpts := &computing.DescribeRouteTablesInput{
		RouteTableId: []*string{&routeTableId},
	}

	res, err := conn.DescribeRouteTables(findOpts)
	if err != nil {
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "Client.InvalidParameterNotFound.RouteTableId" {
			log.Printf("[WARN] Route Table %q could not be found.", routeTableId)
			return false, nil
		}
		return false, fmt.Errorf("Error while checking if route exists: %s", err)
	}

	if len(res.RouteTableSet) < 1 || res.RouteTableSet[0] == nil {
		log.Printf("[WARN] Route Table %q is gone, or route does not exist.",
			routeTableId)
		return false, nil
	}

	if v, ok := d.GetOk("destination_cidr_block"); ok {
		for _, route := range (*res.RouteTableSet[0]).RouteSet {
			if route.DestinationCidrBlock != nil && *route.DestinationCidrBlock == v.(string) {
				return true, nil
			}
		}
	}

	return false, nil
}

// Helper: Create an ID for a route
func resourceNifcloudRouteID(d *schema.ResourceData, r *computing.RouteSetItem) string {

	return fmt.Sprintf("r-%s%d", d.Get("route_table_id").(string), hashcode.String(*r.DestinationCidrBlock))
}

// Helper: retrieve a route
func resourceNifcloudRouteFindRoute(conn *computing.Computing, rtbid string, cidr string) (*computing.RouteSetItem, error) {
	routeTableID := rtbid

	findOpts := &computing.DescribeRouteTablesInput{
		RouteTableId: []*string{&routeTableID},
	}

	resp, err := conn.DescribeRouteTables(findOpts)
	if err != nil {
		return nil, err
	}

	if len(resp.RouteTableSet) < 1 || resp.RouteTableSet[0] == nil {
		return nil, fmt.Errorf("Route Table %q is gone, or route does not exist.",
			routeTableID)
	}

	if cidr != "" {
		for _, route := range (*resp.RouteTableSet[0]).RouteSet {
			if route.DestinationCidrBlock != nil && *route.DestinationCidrBlock == cidr {
				return route, nil
			}
		}

		return nil, fmt.Errorf("Unable to find matching route for Route Table (%s) "+
			"and destination CIDR block (%s).", rtbid, cidr)
	}

	return nil, fmt.Errorf("When trying to find a matching route for Route Table %q "+
		"you need to specify a CIDR block", rtbid)

}