package nifcloud

import (
	"fmt"
	"log"

	"github.com/shztki/nifcloud-sdk-go/nifcloud"
	"github.com/shztki/nifcloud-sdk-go/nifcloud/awserr"
	"github.com/shztki/nifcloud-sdk-go/service/computing"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceNifcloudRouteTable() *schema.Resource {
	return &schema.Resource{
		Create: resourceNifcloudRouteTableCreate,
		Read:   resourceNifcloudRouteTableRead,
//		Update: resourceNifcloudRouteTableUpdate,
		Delete: resourceNifcloudRouteTableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

//		Schema: map[string]*schema.Schema{
//			"route_table_id": {
//				Type:     schema.TypeString,
//				Computed: true,
//			},
//		},
	}
}

func resourceNifcloudRouteTableCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	// Create the routing table
//	createOpts := &computing.CreateRouteTableInput{}
//	log.Printf("[DEBUG] RouteTable create config: %#v", createOpts)

	resp, err := conn.CreateRouteTable(nil) //createOpts)
	if err != nil {
		return fmt.Errorf("Error creating route table: %s", err)
	}

	// Get the ID and store it
	rt := resp.RouteTable
	d.SetId(*rt.RouteTableId)
	log.Printf("[INFO] Route Table ID: %s", d.Id())

	return resourceNifcloudRouteTableRead(d, meta)
}

func resourceNifcloudRouteTableRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	rtRaw, _, err := resourceNifcloudRouteTableStateRefreshFunc(conn, d.Id())()
	if err != nil {
		return err
	}
	if rtRaw == nil {
		d.SetId("")
		return nil
	}

//	rt := rtRaw.(*computing.RouteTable)
//	d.Set("route_table_id", rt.RouteTableId)

	return nil
}
/*
func resourceNifcloudRouteTableUpdate(d *schema.ResourceData, meta interface{}) error {

	return resourceNifcloudRouteTableRead(d, meta)
}
*/

func resourceNifcloudRouteTableDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*NifcloudClient).computingconn

	// First request the routing table since we'll have to disassociate
	// all the subnets first.
//	rtRaw, _, err := resourceNifcloudRouteTableStateRefreshFunc(conn, d.Id())()
//	if err != nil {
//		return err
//	}
//	if rtRaw == nil {
//		return nil
//	}
//	rt := rtRaw.(*computing.RouteTableSetItem)
//
//	// Do all the disassociations
//	for _, a := range rt.AssociationSet {
//		log.Printf("[INFO] Disassociating association with router: %s", *a.RouteTableAssociationId)
//		_, err := conn.DisassociateRouteTable(&computing.DisassociateRouteTableInput{
//			AssociationId: a.RouteTableAssociationId,
//			Agreement:     nifcloud.Bool(false),
//		})
//		log.Printf("[INFO] Disassociating association with router: %v", err)
////		if err != nil {
////			// First check if the association ID is not found. If this
////			// is the case, then it was already disassociated somehow,
////			// and that is okay.
////			if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "Client.InvalidParameterNotFound.AssociationId" {
////				err = nil
////			}
////		}
////		if err != nil {
////			return err
////		}
//	}
//
//	for _, a := range rt.PropagatingVgwSet {
//		log.Printf("[INFO] Disassociating association with vpn gateway: %s", *a.RouteTableAssociationId)
//		_, err := conn.NiftyDisassociateRouteTableFromVpnGateway(&computing.NiftyDisassociateRouteTableFromVpnGatewayInput{
//			AssociationId: a.RouteTableAssociationId,
//			Agreement:     nifcloud.Bool(false),
//		})
//		log.Printf("[INFO] Disassociating association with vpn gateway: %v", err)
////		if err != nil {
////			// First check if the association ID is not found. If this
////			// is the case, then it was already disassociated somehow,
////			// and that is okay.
////			if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "Client.InvalidParameterNotFound.AssociationId" {
////				err = nil
////			}
////		}
////		if err != nil {
////			return err
////		}
//	}
//
//	for _, a := range rt.ElasticLoadBalancerAssociationSet {
//		log.Printf("[INFO] Disassociating association with multi load balancer: %s", *a.RouteTableAssociationId)
//		_, err := conn.NiftyDisassociateRouteTableFromElasticLoadBalancer(&computing.NiftyDisassociateRouteTableFromElasticLoadBalancerInput{
//			AssociationId: a.RouteTableAssociationId,
//		})
//		log.Printf("[INFO] Disassociating association with multi load balancer: %v", err)
////		if err != nil {
////			// First check if the association ID is not found. If this
////			// is the case, then it was already disassociated somehow,
////			// and that is okay.
////			if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "Client.InvalidParameterNotFound.AssociationId" {
////				err = nil
////			}
////		}
////		if err != nil {
////			return err
////		}
//	}

	// Delete the route table
	log.Printf("[INFO] Deleting Route Table: %s", d.Id())
	_, err := conn.DeleteRouteTable(&computing.DeleteRouteTableInput{
		RouteTableId: nifcloud.String(d.Id()),
	})
	if err != nil {
		ec2err, ok := err.(awserr.Error)
		if ok && ec2err.Code() == "Client.InvalidParameterNotFound.RouteTableId" {
			return nil
		}

		return fmt.Errorf("Error deleting route table: %s", err)
	}

/*	// Wait for the route table to really destroy
	log.Printf(
		"[DEBUG] Waiting for route table (%s) to become destroyed",
		d.Id())

	stateConf := &resource.StateChangeConf{
		Pending: []string{"ready"},
		Target:  []string{},
		Refresh: resourceNifcloudRouteTableStateRefreshFunc(conn, d.Id()),
		Timeout: 5 * time.Minute,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf(
			"Error waiting for route table (%s) to become destroyed: %s",
			d.Id(), err)
	}
*/

	return nil
}

// resourceNifcloudRouteTableStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a RouteTable.
func resourceNifcloudRouteTableStateRefreshFunc(conn *computing.Computing, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := conn.DescribeRouteTables(&computing.DescribeRouteTablesInput{
			RouteTableId: []*string{nifcloud.String(id)},
		})
		if err != nil {
			if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "Client.InvalidParameterNotFound.RouteTableId" {
				resp = nil
			} else {
				log.Printf("Error on RouteTableStateRefresh: %s", err)
				return nil, "", err
			}
		}

		if resp == nil {
			// Sometimes AWS just has consistency issues and doesn't see
			// our instance yet. Return an empty state.
			return nil, "", nil
		}

		rt := resp.RouteTableSet[0]
		return rt, "ready", nil
	}
}