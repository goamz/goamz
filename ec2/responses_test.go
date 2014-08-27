package ec2_test

var ErrorDump = `
<?xml version="1.0" encoding="UTF-8"?>
<Response><Errors><Error><Code>UnsupportedOperation</Code>
<Message>AMIs with an instance-store root device are not supported for the instance type 't1.micro'.</Message>
</Error></Errors><RequestID>0503f4e9-bbd6-483c-b54f-c4ae9f3b30f4</RequestID></Response>
`

// http://goo.gl/Mcm3b
var RunInstancesExample = `
<RunInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2011-12-15/">
  <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
  <reservationId>r-47a5402e</reservationId>
  <ownerId>999988887777</ownerId>
  <groupSet>
      <item>
          <groupId>sg-67ad940e</groupId>
          <groupName>default</groupName>
      </item>
  </groupSet>
  <instancesSet>
    <item>
      <instanceId>i-2ba64342</instanceId>
      <imageId>ami-60a54009</imageId>
      <instanceState>
        <code>0</code>
        <name>pending</name>
      </instanceState>
      <privateDnsName></privateDnsName>
      <dnsName></dnsName>
      <keyName>example-key-name</keyName>
      <amiLaunchIndex>0</amiLaunchIndex>
      <instanceType>m1.small</instanceType>
      <launchTime>2007-08-07T11:51:50.000Z</launchTime>
      <placement>
        <availabilityZone>us-east-1b</availabilityZone>
      </placement>
      <monitoring>
        <state>enabled</state>
      </monitoring>
      <virtualizationType>paravirtual</virtualizationType>
      <clientToken/>
      <tagSet/>
      <hypervisor>xen</hypervisor>
    </item>
    <item>
      <instanceId>i-2bc64242</instanceId>
      <imageId>ami-60a54009</imageId>
      <instanceState>
        <code>0</code>
        <name>pending</name>
      </instanceState>
      <privateDnsName></privateDnsName>
      <dnsName></dnsName>
      <keyName>example-key-name</keyName>
      <amiLaunchIndex>1</amiLaunchIndex>
      <instanceType>m1.small</instanceType>
      <launchTime>2007-08-07T11:51:50.000Z</launchTime>
      <placement>
         <availabilityZone>us-east-1b</availabilityZone>
      </placement>
      <monitoring>
        <state>enabled</state>
      </monitoring>
      <virtualizationType>paravirtual</virtualizationType>
      <clientToken/>
      <tagSet/>
      <hypervisor>xen</hypervisor>
    </item>
    <item>
      <instanceId>i-2be64332</instanceId>
      <imageId>ami-60a54009</imageId>
      <instanceState>
        <code>0</code>
        <name>pending</name>
      </instanceState>
      <privateDnsName></privateDnsName>
      <dnsName></dnsName>
      <keyName>example-key-name</keyName>
      <amiLaunchIndex>2</amiLaunchIndex>
      <instanceType>m1.small</instanceType>
      <launchTime>2007-08-07T11:51:50.000Z</launchTime>
      <placement>
         <availabilityZone>us-east-1b</availabilityZone>
      </placement>
      <monitoring>
        <state>enabled</state>
      </monitoring>
      <virtualizationType>paravirtual</virtualizationType>
      <clientToken/>
      <tagSet/>
      <hypervisor>xen</hypervisor>
    </item>
  </instancesSet>
</RunInstancesResponse>
`

// http://goo.gl/3BKHj
var TerminateInstancesExample = `
<TerminateInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2011-12-15/">
  <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
  <instancesSet>
    <item>
      <instanceId>i-3ea74257</instanceId>
      <currentState>
        <code>32</code>
        <name>shutting-down</name>
      </currentState>
      <previousState>
        <code>16</code>
        <name>running</name>
      </previousState>
    </item>
  </instancesSet>
</TerminateInstancesResponse>
`

// http://goo.gl/mLbmw
var DescribeInstancesExample1 = `
<DescribeInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2011-12-15/">
  <requestId>98e3c9a4-848c-4d6d-8e8a-b1bdEXAMPLE</requestId>
  <reservationSet>
    <item>
      <reservationId>r-b27e30d9</reservationId>
      <ownerId>999988887777</ownerId>
      <groupSet>
        <item>
          <groupId>sg-67ad940e</groupId>
          <groupName>default</groupName>
        </item>
      </groupSet>
      <instancesSet>
        <item>
          <instanceId>i-c5cd56af</instanceId>
          <imageId>ami-1a2b3c4d</imageId>
          <instanceState>
            <code>16</code>
            <name>running</name>
          </instanceState>
          <privateDnsName>domU-12-31-39-10-56-34.compute-1.internal</privateDnsName>
          <dnsName>ec2-174-129-165-232.compute-1.amazonaws.com</dnsName>
          <reason/>
          <keyName>GSG_Keypair</keyName>
          <amiLaunchIndex>0</amiLaunchIndex>
          <productCodes/>
          <instanceType>m1.small</instanceType>
          <launchTime>2010-08-17T01:15:18.000Z</launchTime>
          <placement>
            <availabilityZone>us-east-1b</availabilityZone>
            <groupName/>
          </placement>
          <kernelId>aki-94c527fd</kernelId>
          <ramdiskId>ari-96c527ff</ramdiskId>
          <monitoring>
            <state>disabled</state>
          </monitoring>
          <privateIpAddress>10.198.85.190</privateIpAddress>
          <ipAddress>174.129.165.232</ipAddress>
          <architecture>i386</architecture>
          <rootDeviceType>ebs</rootDeviceType>
          <rootDeviceName>/dev/sda1</rootDeviceName>
          <blockDeviceMapping>
            <item>
              <deviceName>/dev/sda1</deviceName>
              <ebs>
                <volumeId>vol-a082c1c9</volumeId>
                <status>attached</status>
                <attachTime>2010-08-17T01:15:21.000Z</attachTime>
                <deleteOnTermination>false</deleteOnTermination>
              </ebs>
            </item>
          </blockDeviceMapping>
          <instanceLifecycle>spot</instanceLifecycle>
          <spotInstanceRequestId>sir-7a688402</spotInstanceRequestId>
          <virtualizationType>paravirtual</virtualizationType>
          <clientToken/>
          <tagSet/>
          <hypervisor>xen</hypervisor>
       </item>
      </instancesSet>
      <requesterId>854251627541</requesterId>
    </item>
    <item>
      <reservationId>r-b67e30dd</reservationId>
      <ownerId>999988887777</ownerId>
      <groupSet>
        <item>
          <groupId>sg-67ad940e</groupId>
          <groupName>default</groupName>
        </item>
      </groupSet>
      <instancesSet>
        <item>
          <instanceId>i-d9cd56b3</instanceId>
          <imageId>ami-1a2b3c4d</imageId>
          <instanceState>
            <code>16</code>
            <name>running</name>
          </instanceState>
          <privateDnsName>domU-12-31-39-10-54-E5.compute-1.internal</privateDnsName>
          <dnsName>ec2-184-73-58-78.compute-1.amazonaws.com</dnsName>
          <reason/>
          <keyName>GSG_Keypair</keyName>
          <amiLaunchIndex>0</amiLaunchIndex>
          <productCodes/>
          <instanceType>m1.large</instanceType>
          <launchTime>2010-08-17T01:15:19.000Z</launchTime>
          <placement>
            <availabilityZone>us-east-1b</availabilityZone>
            <groupName/>
          </placement>
          <kernelId>aki-94c527fd</kernelId>
          <ramdiskId>ari-96c527ff</ramdiskId>
          <monitoring>
            <state>disabled</state>
          </monitoring>
          <privateIpAddress>10.198.87.19</privateIpAddress>
          <ipAddress>184.73.58.78</ipAddress>
          <architecture>i386</architecture>
          <rootDeviceType>ebs</rootDeviceType>
          <rootDeviceName>/dev/sda1</rootDeviceName>
          <blockDeviceMapping>
            <item>
              <deviceName>/dev/sda1</deviceName>
              <ebs>
                <volumeId>vol-a282c1cb</volumeId>
                <status>attached</status>
                <attachTime>2010-08-17T01:15:23.000Z</attachTime>
                <deleteOnTermination>false</deleteOnTermination>
              </ebs>
            </item>
          </blockDeviceMapping>
          <instanceLifecycle>spot</instanceLifecycle>
          <spotInstanceRequestId>sir-55a3aa02</spotInstanceRequestId>
          <virtualizationType>paravirtual</virtualizationType>
          <clientToken/>
          <tagSet/>
          <hypervisor>xen</hypervisor>
       </item>
      </instancesSet>
      <requesterId>854251627541</requesterId>
    </item>
  </reservationSet>
</DescribeInstancesResponse>
`

// http://goo.gl/mLbmw
var DescribeInstancesExample2 = `
<DescribeInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2011-12-15/">
  <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
  <reservationSet>
    <item>
      <reservationId>r-bc7e30d7</reservationId>
      <ownerId>999988887777</ownerId>
      <groupSet>
        <item>
          <groupId>sg-67ad940e</groupId>
          <groupName>default</groupName>
        </item>
      </groupSet>
      <instancesSet>
        <item>
          <instanceId>i-c7cd56ad</instanceId>
          <imageId>ami-b232d0db</imageId>
          <instanceState>
            <code>16</code>
            <name>running</name>
          </instanceState>
          <privateDnsName>domU-12-31-39-01-76-06.compute-1.internal</privateDnsName>
          <dnsName>ec2-72-44-52-124.compute-1.amazonaws.com</dnsName>
          <keyName>GSG_Keypair</keyName>
          <amiLaunchIndex>0</amiLaunchIndex>
          <productCodes/>
          <instanceType>m1.small</instanceType>
          <launchTime>2010-08-17T01:15:16.000Z</launchTime>
          <placement>
              <availabilityZone>us-east-1b</availabilityZone>
          </placement>
          <kernelId>aki-94c527fd</kernelId>
          <ramdiskId>ari-96c527ff</ramdiskId>
          <monitoring>
              <state>disabled</state>
          </monitoring>
          <privateIpAddress>10.255.121.240</privateIpAddress>
          <ipAddress>72.44.52.124</ipAddress>
          <architecture>i386</architecture>
          <rootDeviceType>ebs</rootDeviceType>
          <rootDeviceName>/dev/sda1</rootDeviceName>
          <blockDeviceMapping>
              <item>
                 <deviceName>/dev/sda1</deviceName>
                 <ebs>
                    <volumeId>vol-a482c1cd</volumeId>
                    <status>attached</status>
                    <attachTime>2010-08-17T01:15:26.000Z</attachTime>
                    <deleteOnTermination>true</deleteOnTermination>
                </ebs>
             </item>
          </blockDeviceMapping>
          <virtualizationType>paravirtual</virtualizationType>
          <clientToken/>
          <tagSet>
              <item>
                    <key>webserver</key>
                    <value></value>
             </item>
              <item>
                    <key>stack</key>
                    <value>Production</value>
             </item>
          </tagSet>
          <hypervisor>xen</hypervisor>
        </item>
      </instancesSet>
    </item>
  </reservationSet>
</DescribeInstancesResponse>
`

// http://goo.gl/icuXh5
var ModifyInstanceExample = `
<ModifyImageAttributeResponse xmlns="http://ec2.amazonaws.com/doc/2013-06-15/">
  <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
  <return>true</return>
</ModifyImageAttributeResponse>
`

// http://goo.gl/9rprDN
var AllocateAddressExample = `
<AllocateAddressResponse xmlns="http://ec2.amazonaws.com/doc/2013-10-15/">
   <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
   <publicIp>198.51.100.1</publicIp>
   <domain>vpc</domain>
   <allocationId>eipalloc-5723d13e</allocationId>
</AllocateAddressResponse>
`

// http://goo.gl/3Q0oCc
var ReleaseAddressExample = `
<ReleaseAddressResponse xmlns="http://ec2.amazonaws.com/doc/2013-10-15/">
   <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
   <return>true</return>
</ReleaseAddressResponse>
`

// http://goo.gl/uOSQE
var AssociateAddressExample = `
<AssociateAddressResponse xmlns="http://ec2.amazonaws.com/doc/2013-10-15/">
   <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
   <return>true</return>
   <associationId>eipassoc-fc5ca095</associationId>
</AssociateAddressResponse>
`

// http://goo.gl/LrOa0
var DisassociateAddressExample = `
<DisassociateAddressResponse xmlns="http://ec2.amazonaws.com/doc/2013-10-15/">
   <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
   <return>true</return>
</DisassociateAddressResponse>
`

//http://goo.gl/zW7J4p
var DescribeAddressesExample = `
<DescribeAddressesResponse xmlns="http://ec2.amazonaws.com/doc/2013-10-01/">
   <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
   <addressesSet>
      <item>
         <publicIp>192.0.2.1</publicIp>
         <domain>standard</domain>
         <instanceId>i-f15ebb98</instanceId>
      </item>
      <item>
         <publicIp>198.51.100.2</publicIp>
         <domain>standard</domain>
         <instanceId/>
      </item>
      <item>
         <publicIp>203.0.113.41</publicIp>
         <allocationId>eipalloc-08229861</allocationId>
         <domain>vpc</domain>
         <instanceId>i-64600030</instanceId>
         <associationId>eipassoc-f0229899</associationId>
         <networkInterfaceId>eni-ef229886</networkInterfaceId>
         <networkInterfaceOwnerId>053230519467</networkInterfaceOwnerId>
         <privateIpAddress>10.0.0.228</privateIpAddress>
     </item>
   </addressesSet>
</DescribeAddressesResponse>
`

var DescribeAddressesAllocationIdExample = `
<DescribeAddressesResponse xmlns="http://ec2.amazonaws.com/doc/2013-10-01/">
   <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
   <addressesSet>
      <item>
         <publicIp>203.0.113.41</publicIp>
         <allocationId>eipalloc-08229861</allocationId>
         <domain>vpc</domain>
         <instanceId>i-64600030</instanceId>
         <associationId>eipassoc-f0229899</associationId>
         <networkInterfaceId>eni-ef229886</networkInterfaceId>
         <networkInterfaceOwnerId>053230519467</networkInterfaceOwnerId>
         <privateIpAddress>10.0.0.228</privateIpAddress>
     </item>
     <item>
         <publicIp>146.54.2.230</publicIp>
         <allocationId>eipalloc-08364752</allocationId>
         <domain>vpc</domain>
         <instanceId>i-64693456</instanceId>
         <associationId>eipassoc-f0348693</associationId>
         <networkInterfaceId>eni-da764039</networkInterfaceId>
         <networkInterfaceOwnerId>053230519467</networkInterfaceOwnerId>
         <privateIpAddress>10.0.0.102</privateIpAddress>
     </item>
   </addressesSet>
</DescribeAddressesResponse>
`

// http://goo.gl/cxU41
var CreateImageExample = `
<CreateImageResponse xmlns="http://ec2.amazonaws.com/doc/2013-02-01/">
   <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
   <imageId>ami-4fa54026</imageId>
</CreateImageResponse>
`

// http://goo.gl/V0U25
var DescribeImagesExample = `
<DescribeImagesResponse xmlns="http://ec2.amazonaws.com/doc/2012-08-15/">
         <requestId>4a4a27a2-2e7c-475d-b35b-ca822EXAMPLE</requestId>
    <imagesSet>
        <item>
            <imageId>ami-a2469acf</imageId>
            <imageLocation>aws-marketplace/example-marketplace-amzn-ami.1</imageLocation>
            <imageState>available</imageState>
            <imageOwnerId>123456789999</imageOwnerId>
            <isPublic>true</isPublic>
            <productCodes>
                <item>
                    <productCode>a1b2c3d4e5f6g7h8i9j10k11</productCode>
                    <type>marketplace</type>
                </item>
            </productCodes>
            <architecture>i386</architecture>
            <imageType>machine</imageType>
            <kernelId>aki-805ea7e9</kernelId>
            <imageOwnerAlias>aws-marketplace</imageOwnerAlias>
            <name>example-marketplace-amzn-ami.1</name>
            <description>Amazon Linux AMI i386 EBS</description>
            <rootDeviceType>ebs</rootDeviceType>
            <rootDeviceName>/dev/sda1</rootDeviceName>
            <blockDeviceMapping>
                <item>
                    <deviceName>/dev/sda1</deviceName>
                    <ebs>
                        <snapshotId>snap-787e9403</snapshotId>
                        <volumeSize>8</volumeSize>
                        <deleteOnTermination>true</deleteOnTermination>
                    </ebs>
                </item>
            </blockDeviceMapping>
            <virtualizationType>paravirtual</virtualizationType>
            <tagSet>
                <item>
                    <key>Purpose</key>
                    <value>EXAMPLE</value>
                </item>
            </tagSet>
            <hypervisor>xen</hypervisor>
        </item>
    </imagesSet>
</DescribeImagesResponse>
`

// http://goo.gl/bHO3z
var ImageAttributeExample = `
<DescribeImageAttributeResponse xmlns="http://ec2.amazonaws.com/doc/2013-07-15/">
   <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
   <imageId>ami-61a54008</imageId>
   <launchPermission>
      <item>
         <group>all</group>
      </item>
      <item>
         <userId>495219933132</userId>
      </item>
   </launchPermission>
</DescribeImageAttributeResponse>
`

// http://goo.gl/ttcda
var CreateSnapshotExample = `
<CreateSnapshotResponse xmlns="http://ec2.amazonaws.com/doc/2012-10-01/">
  <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
  <snapshotId>snap-78a54011</snapshotId>
  <volumeId>vol-4d826724</volumeId>
  <status>pending</status>
  <startTime>2008-05-07T12:51:50.000Z</startTime>
  <progress>60%</progress>
  <ownerId>111122223333</ownerId>
  <volumeSize>10</volumeSize>
  <description>Daily Backup</description>
</CreateSnapshotResponse>
`

// http://goo.gl/vwU1y
var DeleteSnapshotExample = `
<DeleteSnapshotResponse xmlns="http://ec2.amazonaws.com/doc/2012-10-01/">
  <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
  <return>true</return>
</DeleteSnapshotResponse>
`

// http://goo.gl/nkovs
var DescribeSnapshotsExample = `
<DescribeSnapshotsResponse xmlns="http://ec2.amazonaws.com/doc/2012-10-01/">
   <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
   <snapshotSet>
      <item>
         <snapshotId>snap-1a2b3c4d</snapshotId>
         <volumeId>vol-8875daef</volumeId>
         <status>pending</status>
         <startTime>2010-07-29T04:12:01.000Z</startTime>
         <progress>30%</progress>
         <ownerId>111122223333</ownerId>
         <volumeSize>15</volumeSize>
         <description>Daily Backup</description>
         <tagSet>
            <item>
               <key>Purpose</key>
               <value>demo_db_14_backup</value>
            </item>
         </tagSet>
      </item>
   </snapshotSet>
</DescribeSnapshotsResponse>
`

// http://goo.gl/YUjO4G
var ModifyImageAttributeExample = `
<ModifyImageAttributeResponse xmlns="http://ec2.amazonaws.com/doc/2013-06-15/">
  <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
  <return>true</return>
</ModifyImageAttributeResponse>
`

// http://goo.gl/hQwPCK
var CopyImageExample = `
<CopyImageResponse xmlns="http://ec2.amazonaws.com/doc/2013-06-15/">
   <requestId>60bc441d-fa2c-494d-b155-5d6a3EXAMPLE</requestId>
   <imageId>ami-4d3c2b1a</imageId>
</CopyImageResponse>
`

var CreateKeyPairExample = `
<CreateKeyPairResponse xmlns="http://ec2.amazonaws.com/doc/2013-02-01/">
  <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
  <keyName>foo</keyName>
  <keyFingerprint>
     00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00:00
  </keyFingerprint>
  <keyMaterial>---- BEGIN RSA PRIVATE KEY ----
MIICiTCCAfICCQD6m7oRw0uXOjANBgkqhkiG9w0BAQUFADCBiDELMAkGA1UEBhMC
VVMxCzAJBgNVBAgTAldBMRAwDgYDVQQHEwdTZWF0dGxlMQ8wDQYDVQQKEwZBbWF6
b24xFDASBgNVBAsTC0lBTSBDb25zb2xlMRIwEAYDVQQDEwlUZXN0Q2lsYWMxHzAd
BgkqhkiG9w0BCQEWEG5vb25lQGFtYXpvbi5jb20wHhcNMTEwNDI1MjA0NTIxWhcN
MTIwNDI0MjA0NTIxWjCBiDELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAldBMRAwDgYD
VQQHEwdTZWF0dGxlMQ8wDQYDVQQKEwZBbWF6b24xFDASBgNVBAsTC0lBTSBDb25z
b2xlMRIwEAYDVQQDEwlUZXN0Q2lsYWMxHzAdBgkqhkiG9w0BCQEWEG5vb25lQGFt
YXpvbi5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMaK0dn+a4GmWIWJ
21uUSfwfEvySWtC2XADZ4nB+BLYgVIk60CpiwsZ3G93vUEIO3IyNoH/f0wYK8m9T
rDHudUZg3qX4waLG5M43q7Wgc/MbQITxOUSQv7c7ugFFDzQGBzZswY6786m86gpE
Ibb3OhjZnzcvQAaRHhdlQWIMm2nrAgMBAAEwDQYJKoZIhvcNAQEFBQADgYEAtCu4
nUhVVxYUntneD9+h8Mg9q6q+auNKyExzyLwaxlAoo7TJHidbtS4J5iNmZgXL0Fkb
FFBjvSfpJIlJ00zbhNYS5f6GuoEDmFJl0ZxBHjJnyp378OD8uTs7fLvjx79LjSTb
NYiytVbZPQUQ5Yaxu2jXnimvw3rrszlaEXAMPLE=
-----END RSA PRIVATE KEY-----
</keyMaterial>
</CreateKeyPairResponse>
`

var DeleteKeyPairExample = `
<DeleteKeyPairResponse xmlns="http://ec2.amazonaws.com/doc/2013-02-01/">
  <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
  <return>true</return>
</DeleteKeyPairResponse>
`

// http://goo.gl/Eo7Yl
var CreateSecurityGroupExample = `
<CreateSecurityGroupResponse xmlns="http://ec2.amazonaws.com/doc/2011-12-15/">
   <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
   <return>true</return>
   <groupId>sg-67ad940e</groupId>
</CreateSecurityGroupResponse>
`

// http://goo.gl/k12Uy
var DescribeSecurityGroupsExample = `
<DescribeSecurityGroupsResponse xmlns="http://ec2.amazonaws.com/doc/2011-12-15/">
  <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
  <securityGroupInfo>
    <item>
      <ownerId>999988887777</ownerId>
      <groupName>WebServers</groupName>
      <groupId>sg-67ad940e</groupId>
      <groupDescription>Web Servers</groupDescription>
      <ipPermissions>
        <item>
           <ipProtocol>tcp</ipProtocol>
           <fromPort>80</fromPort>
           <toPort>80</toPort>
           <groups/>
           <ipRanges>
             <item>
               <cidrIp>0.0.0.0/0</cidrIp>
             </item>
           </ipRanges>
        </item>
      </ipPermissions>
    </item>
    <item>
      <ownerId>999988887777</ownerId>
      <groupName>RangedPortsBySource</groupName>
      <groupId>sg-76abc467</groupId>
      <groupDescription>Group A</groupDescription>
      <ipPermissions>
        <item>
           <ipProtocol>tcp</ipProtocol>
           <fromPort>6000</fromPort>
           <toPort>7000</toPort>
           <groups/>
           <ipRanges/>
        </item>
      </ipPermissions>
    </item>
  </securityGroupInfo>
</DescribeSecurityGroupsResponse>
`

// A dump which includes groups within ip permissions.
var DescribeSecurityGroupsDump = `
<?xml version="1.0" encoding="UTF-8"?>
<DescribeSecurityGroupsResponse xmlns="http://ec2.amazonaws.com/doc/2011-12-15/">
    <requestId>87b92b57-cc6e-48b2-943f-f6f0e5c9f46c</requestId>
    <securityGroupInfo>
        <item>
            <ownerId>12345</ownerId>
            <groupName>default</groupName>
            <groupDescription>default group</groupDescription>
            <ipPermissions>
                <item>
                    <ipProtocol>icmp</ipProtocol>
                    <fromPort>-1</fromPort>
                    <toPort>-1</toPort>
                    <groups>
                        <item>
                            <userId>12345</userId>
                            <groupName>default</groupName>
                            <groupId>sg-67ad940e</groupId>
                        </item>
                    </groups>
                    <ipRanges/>
                </item>
                <item>
                    <ipProtocol>tcp</ipProtocol>
                    <fromPort>0</fromPort>
                    <toPort>65535</toPort>
                    <groups>
                        <item>
                            <userId>12345</userId>
                            <groupName>other</groupName>
                            <groupId>sg-76abc467</groupId>
                        </item>
                    </groups>
                    <ipRanges/>
                </item>
            </ipPermissions>
        </item>
    </securityGroupInfo>
</DescribeSecurityGroupsResponse>
`

// http://goo.gl/QJJDO
var DeleteSecurityGroupExample = `
<DeleteSecurityGroupResponse xmlns="http://ec2.amazonaws.com/doc/2011-12-15/">
   <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
   <return>true</return>
</DeleteSecurityGroupResponse>
`

// http://goo.gl/u2sDJ
var AuthorizeSecurityGroupIngressExample = `
<AuthorizeSecurityGroupIngressResponse xmlns="http://ec2.amazonaws.com/doc/2011-12-15/">
  <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
  <return>true</return>
</AuthorizeSecurityGroupIngressResponse>
`

// http://goo.gl/Mz7xr
var RevokeSecurityGroupIngressExample = `
<RevokeSecurityGroupIngressResponse xmlns="http://ec2.amazonaws.com/doc/2011-12-15/">
  <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
  <return>true</return>
</RevokeSecurityGroupIngressResponse>
`

// http://goo.gl/Vmkqc
var CreateTagsExample = `
<CreateTagsResponse xmlns="http://ec2.amazonaws.com/doc/2011-12-15/">
   <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
   <return>true</return>
</CreateTagsResponse>
`

// http://goo.gl/awKeF
var StartInstancesExample = `
<StartInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2011-12-15/">
  <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
  <instancesSet>
    <item>
      <instanceId>i-10a64379</instanceId>
      <currentState>
          <code>0</code>
          <name>pending</name>
      </currentState>
      <previousState>
          <code>80</code>
          <name>stopped</name>
      </previousState>
    </item>
  </instancesSet>
</StartInstancesResponse>
`

// http://goo.gl/436dJ
var StopInstancesExample = `
<StopInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2011-12-15/">
  <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
  <instancesSet>
    <item>
      <instanceId>i-10a64379</instanceId>
      <currentState>
          <code>64</code>
          <name>stopping</name>
      </currentState>
      <previousState>
          <code>16</code>
          <name>running</name>
      </previousState>
    </item>
  </instancesSet>
</StopInstancesResponse>
`

// http://goo.gl/baoUf
var RebootInstancesExample = `
<RebootInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2011-12-15/">
  <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
  <return>true</return>
</RebootInstancesResponse>
`

var DescribeRouteTablesExample = `
<DescribeRouteTablesResponse xmlns="http://ec2.amazonaws.com/doc/2014-02-01/">
   <requestId>6f570b0b-9c18-4b07-bdec-73740dcf861aEXAMPLE</requestId>
   <routeTableSet>
      <item>
         <routeTableId>rtb-13ad487a</routeTableId>
         <vpcId>vpc-11ad4878</vpcId>
         <routeSet>
            <item>
               <destinationCidrBlock>10.0.0.0/22</destinationCidrBlock>
               <gatewayId>local</gatewayId>
               <state>active</state>
               <origin>CreateRouteTable</origin>
            </item>
         </routeSet>
         <associationSet>
             <item>
                <routeTableAssociationId>rtbassoc-12ad487b</routeTableAssociationId>
                <routeTableId>rtb-13ad487a</routeTableId>
                <main>true</main>
             </item>
         </associationSet>
        <tagSet/>
      </item>
      <item>
         <routeTableId>rtb-f9ad4890</routeTableId>
         <vpcId>vpc-11ad4878</vpcId>
         <routeSet>
            <item>
               <destinationCidrBlock>10.0.0.0/22</destinationCidrBlock>
               <gatewayId>local</gatewayId>
               <state>active</state>
               <origin>CreateRouteTable</origin>
            </item>
            <item>
               <destinationCidrBlock>0.0.0.0/0</destinationCidrBlock>
               <gatewayId>igw-eaad4883</gatewayId>
               <state>active</state>
            </item>
         </routeSet>
         <associationSet>
            <item>
                <routeTableAssociationId>rtbassoc-faad4893</routeTableAssociationId>
                <routeTableId>rtb-f9ad4890</routeTableId>
                <subnetId>subnet-15ad487c</subnetId>
            </item>
         </associationSet>
         <tagSet/>
      </item>
   </routeTableSet>
</DescribeRouteTablesResponse>
`

var CreateRouteTableExample = `
<CreateRouteTableResponse xmlns="http://ec2.amazonaws.com/doc/2014-02-01/">
   <requestId>59abcd43-35bd-4eac-99ed-be587EXAMPLE</requestId>
   <routeTable>
      <routeTableId>rtb-f9ad4890</routeTableId>
      <vpcId>vpc-11ad4878</vpcId>
      <routeSet>
         <item>
            <destinationCidrBlock>10.0.0.0/22</destinationCidrBlock>
            <gatewayId>local</gatewayId>
            <state>active</state>
         </item>
      </routeSet>
      <associationSet/>
      <tagSet/>
   </routeTable>
</CreateRouteTableResponse>
`

var DeleteRouteTableExample = `
<DeleteRouteTableResponse xmlns="http://ec2.amazonaws.com/doc/2014-02-01/">
   <requestId>49dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
   <return>true</return>
</DeleteRouteTableResponse>
`

var AssociateRouteTableExample = `
<AssociateRouteTableResponse xmlns="http://ec2.amazonaws.com/doc/2014-02-01/">
   <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
   <associationId>rtbassoc-f8ad4891</associationId>
</AssociateRouteTableResponse>
`

var DisassociateRouteTableExample = `
<DisassociateRouteTableResponse xmlns="http://ec2.amazonaws.com/doc/2014-02-01/">
   <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId>
   <return>true</return>
</DisassociateRouteTableResponse>
`

var ReplaceRouteTableAssociationExample = `
<ReplaceRouteTableAssociationResponse xmlns="http://ec2.amazonaws.com/doc/2014-02-01/">
   <requestId>59dbff89-35bd-4eac-88ed-be587EXAMPLE</requestId>
   <newAssociationId>rtbassoc-faad2958</newAssociationId>
</ReplaceRouteTableAssociationResponse>
`
var DescribeReservedInstancesExample = `
<DescribeReservedInstancesResponse xmlns="http://ec2.amazonaws.com/doc/2014-06-15/">
   <requestId>59dbff89-35bd-4eac-99ed-be587EXAMPLE</requestId> 
   <reservedInstancesSet>
      <item>
         <reservedInstancesId>e5a2ff3b-7d14-494f-90af-0b5d0EXAMPLE</reservedInstancesId>
         <instanceType>m1.xlarge</instanceType>
         <availabilityZone>us-east-1b</availabilityZone>
         <duration>31536000</duration>
         <fixedPrice>61.0</fixedPrice>
         <usagePrice>0.034</usagePrice>
         <instanceCount>3</instanceCount>
         <productDescription>Linux/UNIX</productDescription>
         <state>active</state> 
         <instanceTenancy>default</instanceTenancy>
         <currencyCode>USD</currencyCode>
         <offeringType>Light Utilization</offeringType>
         <recurringCharges/>
      </item>
   </reservedInstancesSet> 
</DescribeReservedInstancesResponse>
`
