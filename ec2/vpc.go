package ec2

// RouteTable describes a route table which contains a set of rules, called routes
// that are used to determine where network traffic is directed.
//
// See http://goo.gl/bI9hkg for more details.
type RouteTable struct {
	Id              string                  `xml:"routeTableId"`
	VpcId           string                  `xml:"vpcId"`
	Routes          []Route                 `xml:"routeSet>item"`
	Associations    []RouteTableAssociation `xml:"associationSet>item"`
	PropagatingVgws []PropagatingVgw        `xml:"propagatingVgwSet>item"`
	Tags            []Tag                   `xml:"tagSet>item"`
}

// Route describes a route in a route table.
//
// See http://goo.gl/hE5Kxe for more details.
type Route struct {
	DestinationCidrBlock   string `xml:"destinationCidrBlock"`   // The CIDR block used for the destination match.
	GatewayId              string `xml:"gatewayId"`              // The ID of a gateway attached to your VPC.
	InstanceId             string `xml:"instanceId"`             // The ID of a NAT instance in your VPC.
	InstanceOwnerId        string `xml:"instanceOwnerId"`        // The AWS account ID of the owner of the instance.
	NetworkInterfaceId     string `xml:"networkInterfaceId"`     // The ID of the network interface.
	State                  string `xml:"state"`                  // The state of the route. Valid values: active | blackhole
	Origin                 string `xml:"origin"`                 // Describes how the route was created. Valid values: Valid values: CreateRouteTable | CreateRoute | EnableVgwRoutePropagation
	VpcPeeringConnectionId string `xml:"vpcPeeringConnectionId"` // The ID of the VPC peering connection.
}

// RouteTableAssociation describes an association between a route table and a subnet.
//
// See http://goo.gl/BZB8o8 for more details.
type RouteTableAssociation struct {
	Id           string `xml:"routeTableAssociationId"` // The ID of the association between a route table and a subnet.
	RouteTableId string `xml:"routeTableId"`            // The ID of the route table.
	SubnetId     string `xml:"subnetId"`                // The ID of the subnet.
	Main         bool   `xml:"main"`                    // Indicates whether this is the main route table.
}

// PropagatingVgw describes a virtual private gateway propagating route.
//
// See http://goo.gl/myGQtG for more details.
type PropagatingVgw struct {
	GatewayId string `xml:"gatewayID"`
}

// DescribeRouteTablesResp represents a response from a DescribeRouteTables call
//
// See http://goo.gl/T3tVsg for more details.
type DescribeRouteTablesResp struct {
	RequestId   string       `xml:"requestId"`
	RouteTables []RouteTable `xml:"routeTableSet>item"`
}

// DescribeRouteTables describes one or more of your route tables
//
// See http://goo.gl/S0RVos for more details.
func (ec2 *EC2) DescribeRouteTables(routeTableIds []string, filter *Filter) (resp *DescribeRouteTablesResp, err error) {
	params := makeParams("DescribeRouteTables")
	addParamsList(params, "RouteTableId", routeTableIds)
	filter.addParams(params)
	resp = &DescribeRouteTablesResp{}
	err = ec2.query(params, resp)
	if err != nil {
		return nil, err
	}
	return
}
