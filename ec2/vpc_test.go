package ec2_test

import (
	"time"

	"github.com/czos/goamz/aws"
	"github.com/czos/goamz/ec2"
	"github.com/motain/gocheck"
)

func (s *S) TestCreateRouteTable(c *gocheck.C) {
	testServer.Response(200, nil, CreateRouteTableExample)

	resp, err := s.ec2.CreateRouteTable("vpc-11ad4878")

	req := testServer.WaitRequest()
	c.Assert(req.Form["Action"], gocheck.DeepEquals, []string{"CreateRouteTable"})
	c.Assert(req.Form["VpcId"], gocheck.DeepEquals, []string{"vpc-11ad4878"})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "59abcd43-35bd-4eac-99ed-be587EXAMPLE")
	c.Assert(resp.RouteTable.Id, gocheck.Equals, "rtb-f9ad4890")
	c.Assert(resp.RouteTable.VpcId, gocheck.Equals, "vpc-11ad4878")
	c.Assert(resp.RouteTable.Routes, gocheck.HasLen, 1)
	c.Assert(resp.RouteTable.Routes[0], gocheck.DeepEquals, ec2.Route{
		DestinationCidrBlock: "10.0.0.0/22",
		GatewayId:            "local",
		State:                "active",
	})
	c.Assert(resp.RouteTable.Associations, gocheck.HasLen, 0)
	c.Assert(resp.RouteTable.Tags, gocheck.HasLen, 0)
}

func (s *S) TestDescribeRouteTables(c *gocheck.C) {
	testServer.Response(200, nil, DescribeRouteTablesExample)

	filter := ec2.NewFilter()
	filter.Add("key1", "value1")
	filter.Add("key2", "value2", "value3")

	resp, err := s.ec2.DescribeRouteTables([]string{"rt1", "rt2"}, nil)

	req := testServer.WaitRequest()
	c.Assert(req.Form["Action"], gocheck.DeepEquals, []string{"DescribeRouteTables"})
	c.Assert(req.Form["RouteTableId.1"], gocheck.DeepEquals, []string{"rt1"})
	c.Assert(req.Form["RouteTableId.2"], gocheck.DeepEquals, []string{"rt2"})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "6f570b0b-9c18-4b07-bdec-73740dcf861aEXAMPLE")
	c.Assert(resp.RouteTables, gocheck.HasLen, 2)

	rt1 := resp.RouteTables[0]
	c.Assert(rt1.Id, gocheck.Equals, "rtb-13ad487a")
	c.Assert(rt1.VpcId, gocheck.Equals, "vpc-11ad4878")
	c.Assert(rt1.Routes, gocheck.DeepEquals, []ec2.Route{
		{DestinationCidrBlock: "10.0.0.0/22", GatewayId: "local", State: "active", Origin: "CreateRouteTable"},
	})
	c.Assert(rt1.Associations, gocheck.DeepEquals, []ec2.RouteTableAssociation{
		{Id: "rtbassoc-12ad487b", RouteTableId: "rtb-13ad487a", Main: true},
	})

	rt2 := resp.RouteTables[1]
	c.Assert(rt2.Id, gocheck.Equals, "rtb-f9ad4890")
	c.Assert(rt2.VpcId, gocheck.Equals, "vpc-11ad4878")
	c.Assert(rt2.Routes, gocheck.DeepEquals, []ec2.Route{
		{DestinationCidrBlock: "10.0.0.0/22", GatewayId: "local", State: "active", Origin: "CreateRouteTable"},
		{DestinationCidrBlock: "0.0.0.0/0", GatewayId: "igw-eaad4883", State: "active"},
	})
	c.Assert(rt2.Associations, gocheck.DeepEquals, []ec2.RouteTableAssociation{
		{Id: "rtbassoc-faad4893", RouteTableId: "rtb-f9ad4890", SubnetId: "subnet-15ad487c"},
	})
}

func (s *S) TestAssociateRouteTable(c *gocheck.C) {
	testServer.Response(200, nil, AssociateRouteTableExample)

	resp, err := s.ec2.AssociateRouteTable("rtb-e4ad488d", "subnet-15ad487c")

	req := testServer.WaitRequest()
	c.Assert(req.Form["Action"], gocheck.DeepEquals, []string{"AssociateRouteTable"})
	c.Assert(req.Form["RouteTableId"], gocheck.DeepEquals, []string{"rtb-e4ad488d"})
	c.Assert(req.Form["SubnetId"], gocheck.DeepEquals, []string{"subnet-15ad487c"})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "59dbff89-35bd-4eac-99ed-be587EXAMPLE")
	c.Assert(resp.AssociationId, gocheck.Equals, "rtbassoc-f8ad4891")
}

func (s *S) TestDisassociateRouteTable(c *gocheck.C) {
	testServer.Response(200, nil, DisassociateRouteTableExample)

	resp, err := s.ec2.DisassociateRouteTable("rtbassoc-f8ad4891")

	req := testServer.WaitRequest()
	c.Assert(req.Form["Action"], gocheck.DeepEquals, []string{"DisassociateRouteTable"})
	c.Assert(req.Form["AssociationId"], gocheck.DeepEquals, []string{"rtbassoc-f8ad4891"})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "59dbff89-35bd-4eac-99ed-be587EXAMPLE")
	c.Assert(resp.Return, gocheck.Equals, true)
}

func (s *S) TestReplaceRouteTableAssociation(c *gocheck.C) {
	testServer.Response(200, nil, ReplaceRouteTableAssociationExample)

	resp, err := s.ec2.ReplaceRouteTableAssociation("rtbassoc-f8ad4891", "rtb-f9ad4890")

	req := testServer.WaitRequest()
	c.Assert(req.Form["Action"], gocheck.DeepEquals, []string{"ReplaceRouteTableAssociation"})
	c.Assert(req.Form["RouteTableId"], gocheck.DeepEquals, []string{"rtb-f9ad4890"})
	c.Assert(req.Form["AssociationId"], gocheck.DeepEquals, []string{"rtbassoc-f8ad4891"})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "59dbff89-35bd-4eac-88ed-be587EXAMPLE")
	c.Assert(resp.NewAssociationId, gocheck.Equals, "rtbassoc-faad2958")
}

func (s *S) TestDeleteRouteTable(c *gocheck.C) {
	testServer.Response(200, nil, DeleteRouteTableExample)

	resp, err := s.ec2.DeleteRouteTable("rtb-f9ad4890")

	req := testServer.WaitRequest()
	c.Assert(req.Form["Action"], gocheck.DeepEquals, []string{"DeleteRouteTable"})
	c.Assert(req.Form["RouteTableId"], gocheck.DeepEquals, []string{"rtb-f9ad4890"})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "49dbff89-35bd-4eac-99ed-be587EXAMPLE")
	c.Assert(resp.Return, gocheck.Equals, true)
}

// VPC tests with example responses
func (s *S) TestCreateVPCExample(c *gocheck.C) {
	testServer.Response(200, nil, CreateVpcExample)

	resp, err := s.ec2.CreateVPC("10.0.0.0/16", "default")
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], gocheck.DeepEquals, []string{"CreateVpc"})
	c.Assert(req.Form["CidrBlock"], gocheck.DeepEquals, []string{"10.0.0.0/16"})
	c.Assert(req.Form["InstanceTenancy"], gocheck.DeepEquals, []string{"default"})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
	vpc := resp.VPC
	c.Assert(vpc.Id, gocheck.Equals, "vpc-1a2b3c4d")
	c.Assert(vpc.State, gocheck.Equals, "pending")
	c.Assert(vpc.CIDRBlock, gocheck.Equals, "10.0.0.0/16")
	c.Assert(vpc.DHCPOptionsId, gocheck.Equals, "dopt-1a2b3c4d2")
	c.Assert(vpc.Tags, gocheck.HasLen, 0)
	c.Assert(vpc.IsDefault, gocheck.Equals, false)
	c.Assert(vpc.InstanceTenancy, gocheck.Equals, "default")
}

func (s *S) TestDeleteVPCExample(c *gocheck.C) {
	testServer.Response(200, nil, DeleteVpcExample)

	resp, err := s.ec2.DeleteVPC("vpc-id")
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], gocheck.DeepEquals, []string{"DeleteVpc"})
	c.Assert(req.Form["VpcId"], gocheck.DeepEquals, []string{"vpc-id"})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
}

func (s *S) TestVPCsExample(c *gocheck.C) {
	testServer.Response(200, nil, DescribeVpcsExample)

	resp, err := s.ec2.VPCs([]string{"vpc-1a2b3c4d"}, nil)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], gocheck.DeepEquals, []string{"DescribeVpcs"})
	c.Assert(req.Form["VpcId.1"], gocheck.DeepEquals, []string{"vpc-1a2b3c4d"})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
	c.Assert(resp.VPCs, gocheck.HasLen, 1)
	vpc := resp.VPCs[0]
	c.Assert(vpc.Id, gocheck.Equals, "vpc-1a2b3c4d")
	c.Assert(vpc.State, gocheck.Equals, "available")
	c.Assert(vpc.CIDRBlock, gocheck.Equals, "10.0.0.0/23")
	c.Assert(vpc.DHCPOptionsId, gocheck.Equals, "dopt-7a8b9c2d")
	c.Assert(vpc.Tags, gocheck.HasLen, 0)
	c.Assert(vpc.IsDefault, gocheck.Equals, false)
	c.Assert(vpc.InstanceTenancy, gocheck.Equals, "default")
}

// VPC tests to run against either a local test server or live on EC2.
func (s *ServerTests) TestVPCs(c *gocheck.C) {
	resp1, err := s.ec2.CreateVPC("10.0.0.0/16", "")
	c.Assert(err, gocheck.IsNil)
	assertVPC(c, resp1.VPC, "", "10.0.0.0/16")
	id1 := resp1.VPC.Id

	resp2, err := s.ec2.CreateVPC("10.1.0.0/16", "default")
	c.Assert(err, gocheck.IsNil)
	assertVPC(c, resp2.VPC, "", "10.1.0.0/16")
	id2 := resp2.VPC.Id

	// We only check for the VPCs we just created, because the user
	// might have others in his account (when testing against the EC2
	// servers). In some cases it takes a short while until both VPCs
	// are created, so we need to retry a few times to make sure.
	var list *ec2.VPCsResp
	done := false
	testAttempt := aws.AttemptStrategy{
		Total: 2 * time.Minute,
		Delay: 5 * time.Second,
	}
	for a := testAttempt.Start(); a.Next(); {
		c.Logf("waiting for %v to be created", []string{id1, id2})
		list, err = s.ec2.VPCs(nil, nil)
		if err != nil {
			c.Logf("retrying; VPCs returned: %v", err)
			continue
		}
		found := 0
		for _, vpc := range list.VPCs {
			c.Logf("found VPC %v", vpc)
			switch vpc.Id {
			case id1:
				assertVPC(c, vpc, id1, resp1.VPC.CIDRBlock)
				found++
			case id2:
				assertVPC(c, vpc, id2, resp2.VPC.CIDRBlock)
				found++
			}
			if found == 2 {
				done = true
				break
			}
		}
		if done {
			c.Logf("all VPCs were created")
			break
		}
	}
	if !done {
		c.Fatalf("timeout while waiting for VPCs %v", []string{id1, id2})
	}

	list, err = s.ec2.VPCs([]string{id1}, nil)
	c.Assert(err, gocheck.IsNil)
	c.Assert(list.VPCs, gocheck.HasLen, 1)
	assertVPC(c, list.VPCs[0], id1, resp1.VPC.CIDRBlock)

	f := ec2.NewFilter()
	f.Add("cidr", resp2.VPC.CIDRBlock)
	list, err = s.ec2.VPCs(nil, f)
	c.Assert(err, gocheck.IsNil)
	c.Assert(list.VPCs, gocheck.HasLen, 1)
	assertVPC(c, list.VPCs[0], id2, resp2.VPC.CIDRBlock)

	_, err = s.ec2.DeleteVPC(id1)
	c.Assert(err, gocheck.IsNil)
	_, err = s.ec2.DeleteVPC(id2)
	c.Assert(err, gocheck.IsNil)
}

// deleteVPCs ensures the given VPCs are deleted, by retrying until a
// timeout or all VPC cannot be found anymore.  This should be used to
// make sure tests leave no VPCs around.
func (s *ServerTests) deleteVPCs(c *gocheck.C, ids []string) {
	testAttempt := aws.AttemptStrategy{
		Total: 2 * time.Minute,
		Delay: 5 * time.Second,
	}
	for a := testAttempt.Start(); a.Next(); {
		deleted := 0
		c.Logf("deleting VPCs %v", ids)
		for _, id := range ids {
			_, err := s.ec2.DeleteVPC(id)
			if err == nil || err.Error() == "InvalidVpcID.NotFound" {
				c.Logf("VPC %s deleted", id)
				deleted++
				continue
			}
			if err != nil {
				c.Logf("retrying; DeleteVPC returned: %v", err)
			}
		}
		if deleted == len(ids) {
			c.Logf("all VPCs deleted")
			return
		}
	}
	c.Fatalf("timeout while waiting %v VPCs to get deleted!", ids)
}

func assertVPC(c *gocheck.C, obtained ec2.VPC, expectId, expectCidr string) {
	if expectId != "" {
		c.Assert(obtained.Id, gocheck.Equals, expectId)
	} else {
		c.Assert(obtained.Id, gocheck.Matches, `^vpc-[0-9a-f]+$`)
	}
	c.Assert(obtained.State, gocheck.Matches, "(available|pending)")
	if expectCidr != "" {
		c.Assert(obtained.CIDRBlock, gocheck.Equals, expectCidr)
	} else {
		c.Assert(obtained.CIDRBlock, gocheck.Matches, `^\d+\.\d+\.\d+\.\d+/\d+$`)
	}
	c.Assert(obtained.DHCPOptionsId, gocheck.Matches, `^dopt-[0-9a-f]+$`)
	c.Assert(obtained.IsDefault, gocheck.Equals, false)
	c.Assert(obtained.Tags, gocheck.HasLen, 0)
	c.Assert(obtained.InstanceTenancy, gocheck.Matches, "(default|dedicated)")
}
