//
// goamz - Go packages to interact with the Amazon Web Services.
//
//   https://wiki.ubuntu.com/goamz
//
// Copyright (c) 2014 Canonical Ltd.
//

package ec2_test

import (
	"time"

	"github.com/czos/goamz/aws"
	"github.com/czos/goamz/ec2"
	"github.com/motain/gocheck"
)

// Subnet tests with example responses

func (s *S) TestCreateSubnetExample(c *gocheck.C) {
	testServer.Response(200, nil, CreateSubnetExample)

	resp, err := s.ec2.CreateSubnet("vpc-1a2b3c4d", "10.0.1.0/24", "us-east-1a")
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], gocheck.DeepEquals, []string{"CreateSubnet"})
	c.Assert(req.Form["VpcId"], gocheck.DeepEquals, []string{"vpc-1a2b3c4d"})
	c.Assert(req.Form["CidrBlock"], gocheck.DeepEquals, []string{"10.0.1.0/24"})
	c.Assert(req.Form["AvailabilityZone"], gocheck.DeepEquals, []string{"us-east-1a"})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
	subnet := resp.Subnet
	c.Assert(subnet.Id, gocheck.Equals, "subnet-9d4a7b6c")
	c.Assert(subnet.State, gocheck.Equals, "pending")
	c.Assert(subnet.VPCId, gocheck.Equals, "vpc-1a2b3c4d")
	c.Assert(subnet.CIDRBlock, gocheck.Equals, "10.0.1.0/24")
	c.Assert(subnet.AvailableIPCount, gocheck.Equals, 251)
	c.Assert(subnet.AvailZone, gocheck.Equals, "us-east-1a")
	c.Assert(subnet.Tags, gocheck.HasLen, 0)
}

func (s *S) TestDeleteSubnetExample(c *gocheck.C) {
	testServer.Response(200, nil, DeleteSubnetExample)

	resp, err := s.ec2.DeleteSubnet("subnet-id")
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], gocheck.DeepEquals, []string{"DeleteSubnet"})
	c.Assert(req.Form["SubnetId"], gocheck.DeepEquals, []string{"subnet-id"})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
}

func (s *S) TestSubnetsExample(c *gocheck.C) {
	testServer.Response(200, nil, DescribeSubnetsExample)

	ids := []string{"subnet-9d4a7b6c", "subnet-6e7f829e"}
	resp, err := s.ec2.Subnets(ids, nil)
	req := testServer.WaitRequest()

	c.Assert(req.Form["Action"], gocheck.DeepEquals, []string{"DescribeSubnets"})
	c.Assert(req.Form["SubnetId.1"], gocheck.DeepEquals, []string{ids[0]})
	c.Assert(req.Form["SubnetId.2"], gocheck.DeepEquals, []string{ids[1]})

	c.Assert(err, gocheck.IsNil)
	c.Assert(resp.RequestId, gocheck.Equals, "7a62c49f-347e-4fc4-9331-6e8eEXAMPLE")
	c.Assert(resp.Subnets, gocheck.HasLen, 2)
	subnet := resp.Subnets[0]
	c.Assert(subnet.Id, gocheck.Equals, "subnet-9d4a7b6c")
	c.Assert(subnet.State, gocheck.Equals, "available")
	c.Assert(subnet.VPCId, gocheck.Equals, "vpc-1a2b3c4d")
	c.Assert(subnet.CIDRBlock, gocheck.Equals, "10.0.1.0/24")
	c.Assert(subnet.AvailableIPCount, gocheck.Equals, 251)
	c.Assert(subnet.AvailZone, gocheck.Equals, "us-east-1a")
	c.Assert(subnet.DefaultForAZ, gocheck.Equals, false)
	c.Assert(subnet.MapPublicIPOnLaunch, gocheck.Equals, false)
	c.Assert(subnet.Tags, gocheck.HasLen, 0)
	subnet = resp.Subnets[1]
	c.Assert(subnet.Id, gocheck.Equals, "subnet-6e7f829e")
	c.Assert(subnet.State, gocheck.Equals, "available")
	c.Assert(subnet.VPCId, gocheck.Equals, "vpc-1a2b3c4d")
	c.Assert(subnet.CIDRBlock, gocheck.Equals, "10.0.0.0/24")
	c.Assert(subnet.AvailableIPCount, gocheck.Equals, 251)
	c.Assert(subnet.AvailZone, gocheck.Equals, "us-east-1a")
	c.Assert(subnet.DefaultForAZ, gocheck.Equals, false)
	c.Assert(subnet.MapPublicIPOnLaunch, gocheck.Equals, false)
	c.Assert(subnet.Tags, gocheck.HasLen, 0)
}

// Subnet tests run against either a local test server or live on EC2.

func (s *ServerTests) TestSubnets(c *gocheck.C) {
	resp, err := s.ec2.CreateVPC("10.2.0.0/16", "")
	c.Assert(err, gocheck.IsNil)
	vpcId := resp.VPC.Id
	defer s.deleteVPCs(c, []string{vpcId})

	resp1 := s.createSubnet(c, vpcId, "10.2.1.0/24", "")
	assertSubnet(c, resp1.Subnet, "", vpcId, "10.2.1.0/24")
	id1 := resp1.Subnet.Id

	resp2, err := s.ec2.CreateSubnet(vpcId, "10.2.2.0/24", "")
	c.Assert(err, gocheck.IsNil)
	assertSubnet(c, resp2.Subnet, "", vpcId, "10.2.2.0/24")
	id2 := resp2.Subnet.Id

	// We only check for the subnets we just created, because the user
	// might have others in his account (when testing against the EC2
	// servers). In some cases it takes a short while until both
	// subnets are created, so we need to retry a few times to make
	// sure.
	testAttempt := aws.AttemptStrategy{
		Total: 2 * time.Minute,
		Delay: 5 * time.Second,
	}
	var list *ec2.SubnetsResp
	done := false
	for a := testAttempt.Start(); a.Next(); {
		c.Logf("waiting for %v to be created", []string{id1, id2})
		list, err = s.ec2.Subnets(nil, nil)
		if err != nil {
			c.Logf("retrying; Subnets returned: %v", err)
			continue
		}
		found := 0
		for _, subnet := range list.Subnets {
			c.Logf("found subnet %v", subnet)
			switch subnet.Id {
			case id1:
				assertSubnet(c, subnet, id1, vpcId, resp1.Subnet.CIDRBlock)
				found++
			case id2:
				assertSubnet(c, subnet, id2, vpcId, resp2.Subnet.CIDRBlock)
				found++
			}
			if found == 2 {
				done = true
				break
			}
		}
		if done {
			c.Logf("all subnets were created")
			break
		}
	}
	if !done {
		c.Fatalf("timeout while waiting for subnets %v", []string{id1, id2})
	}

	list, err = s.ec2.Subnets([]string{id1}, nil)
	c.Assert(err, gocheck.IsNil)
	c.Assert(list.Subnets, gocheck.HasLen, 1)
	assertSubnet(c, list.Subnets[0], id1, vpcId, resp1.Subnet.CIDRBlock)

	f := ec2.NewFilter()
	f.Add("cidr", resp2.Subnet.CIDRBlock)
	list, err = s.ec2.Subnets(nil, f)
	c.Assert(err, gocheck.IsNil)
	c.Assert(list.Subnets, gocheck.HasLen, 1)
	assertSubnet(c, list.Subnets[0], id2, vpcId, resp2.Subnet.CIDRBlock)

	_, err = s.ec2.DeleteSubnet(id1)
	c.Assert(err, gocheck.IsNil)
	_, err = s.ec2.DeleteSubnet(id2)
	c.Assert(err, gocheck.IsNil)
}

// createSubnet ensures a subnet with the given vpcId and cidrBlock
// gets created, retrying a few times with a timeout. This needs to be
// done when testing against EC2 servers, because if the VPC was just
// created it might take some time for it to show up, so the subnet
// can be created.
func (s *ServerTests) createSubnet(c *gocheck.C, vpcId, cidrBlock, availZone string) *ec2.CreateSubnetResp {
	testAttempt := aws.AttemptStrategy{
		Total: 2 * time.Minute,
		Delay: 5 * time.Second,
	}
	for a := testAttempt.Start(); a.Next(); {
		resp, err := s.ec2.CreateSubnet(vpcId, cidrBlock, availZone)
		if err != nil && err.Error() == "InvalidVpcID.NotFound" {
			c.Logf("VPC %v not created yet; retrying", vpcId)
			continue
		}
		if err != nil {
			c.Logf("retrying; CreateSubnet returned: %v", err)
			continue
		}
		return resp
	}
	c.Fatalf("timeout while waiting for VPC and subnet")
	return nil
}

// deleteSubnets ensures the given subnets are deleted, by retrying
// until a timeout or all subnets cannot be found anymore.  This
// should be used to make sure tests leave no subnets around.
func (s *ServerTests) deleteSubnets(c *gocheck.C, ids []string) {
	testAttempt := aws.AttemptStrategy{
		Total: 2 * time.Minute,
		Delay: 5 * time.Second,
	}
	for a := testAttempt.Start(); a.Next(); {
		deleted := 0
		c.Logf("deleting subnets %v", ids)
		for _, id := range ids {
			_, err := s.ec2.DeleteSubnet(id)
			if err == nil || err.Error() == "InvalidSubnetID.NotFound" {
				c.Logf("subnet %s deleted", id)
				deleted++
				continue
			}
			if err != nil {
				c.Logf("retrying; DeleteSubnet returned: %v", err)
			}
		}
		if deleted == len(ids) {
			c.Logf("all subnets deleted")
			return
		}
	}
	c.Fatalf("timeout while waiting %v subnets to get deleted!", ids)
}

func assertSubnet(c *gocheck.C, obtained ec2.Subnet, expectId, expectVpcId, expectCidr string) {
	if expectId != "" {
		c.Assert(obtained.Id, gocheck.Equals, expectId)
	} else {
		c.Assert(obtained.Id, gocheck.Matches, `^subnet-[0-9a-f]+$`)
	}
	c.Assert(obtained.State, gocheck.Matches, "(available|pending)")
	if expectVpcId != "" {
		c.Assert(obtained.VPCId, gocheck.Equals, expectVpcId)
	} else {
		c.Assert(obtained.VPCId, gocheck.Matches, `^vpc-[0-9a-f]+$`)
	}
	if expectCidr != "" {
		c.Assert(obtained.CIDRBlock, gocheck.Equals, expectCidr)
	} else {
		c.Assert(obtained.CIDRBlock, gocheck.Matches, `^\d+\.\d+\.\d+\.\d+/\d+$`)
	}
	c.Assert(obtained.AvailZone, gocheck.Not(gocheck.Equals), "")
	c.Assert(obtained.AvailableIPCount, gocheck.Not(gocheck.Equals), 0)
	c.Assert(obtained.DefaultForAZ, gocheck.Equals, false)
	c.Assert(obtained.MapPublicIPOnLaunch, gocheck.Equals, false)
}
