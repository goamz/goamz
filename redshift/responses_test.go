package redshift_test

var ErrorDump = `
<?xml version="1.0" encoding="UTF-8"?>
<Response><Error><Code>UnsupportedOperation</Code>
<Message>something failed</Message>
</Error><RequestID>0503f4e9-bbd6-483c-b54f-c4ae9f3b30f4</RequestID></Response>
`

var DescribeClustersExample1 = `
<DescribeClustersResponse xmlns="http://redshift.amazonaws.com/doc/2012-12-01/">
  <DescribeClustersResult>
    <Clusters>
      <Cluster>
        <PendingModifiedValues>
          <MasterUserPassword>****</MasterUserPassword>
        </PendingModifiedValues>
        <ClusterVersion>1.0</ClusterVersion>
        <VpcSecurityGroups/>
        <ClusterStatus>creating</ClusterStatus>
        <NumberOfNodes>2</NumberOfNodes>
        <AutomatedSnapshotRetentionPeriod>1</AutomatedSnapshotRetentionPeriod>
        <PubliclyAccessible>true</PubliclyAccessible>
        <Encrypted>false</Encrypted>
        <DBName>dev</DBName>
        <PreferredMaintenanceWindow>sun:10:30-sun:11:00</PreferredMaintenanceWindow>
        <ClusterParameterGroups>
          <ClusterParameterGroup>
            <ParameterApplyStatus>in-sync</ParameterApplyStatus>
            <ParameterGroupName>default.redshift-1.0</ParameterGroupName>
          </ClusterParameterGroup>
        </ClusterParameterGroups>
        <ClusterSecurityGroups>
          <ClusterSecurityGroup>
            <Status>active</Status>
            <ClusterSecurityGroupName>default</ClusterSecurityGroupName>
          </ClusterSecurityGroup>
        </ClusterSecurityGroups>
        <AvailabilityZone>us-east-1a</AvailabilityZone>
        <NodeType>dw1.xlarge</NodeType>
        <ClusterIdentifier>examplecluster</ClusterIdentifier>
        <AllowVersionUpgrade>true</AllowVersionUpgrade>
        <MasterUsername>masteruser</MasterUsername>
      </Cluster>
    </Clusters>
  </DescribeClustersResult>
  <ResponseMetadata>
    <RequestId>837d45d6-64f0-11e2-b07c-f7fbdd006c67</RequestId>
  </ResponseMetadata>
</DescribeClustersResponse>
`

var CreateClusterExample1 = `
<CreateClusterResponse xmlns="http://redshift.amazonaws.com/doc/2012-12-01/">
  <CreateClusterResult>
    <Cluster>
      <PendingModifiedValues>
        <MasterUserPassword>****</MasterUserPassword>
      </PendingModifiedValues>
      <ClusterVersion>1.0</ClusterVersion>
      <VpcSecurityGroups/>
      <ClusterStatus>creating</ClusterStatus>
      <NumberOfNodes>2</NumberOfNodes>
      <AutomatedSnapshotRetentionPeriod>1</AutomatedSnapshotRetentionPeriod>
      <PubliclyAccessible>true</PubliclyAccessible>
      <Encrypted>false</Encrypted>
      <DBName>dev</DBName>
      <PreferredMaintenanceWindow>sun:10:30-sun:11:00</PreferredMaintenanceWindow>
      <ClusterParameterGroups>
        <ClusterParameterGroup>
          <ParameterApplyStatus>in-sync</ParameterApplyStatus>
          <ParameterGroupName>default.redshift-1.0</ParameterGroupName>
        </ClusterParameterGroup>
      </ClusterParameterGroups>
      <ClusterSecurityGroups>
        <ClusterSecurityGroup>
          <Status>active</Status>
          <ClusterSecurityGroupName>default</ClusterSecurityGroupName>
        </ClusterSecurityGroup>
      </ClusterSecurityGroups>
      <NodeType>dw1.xlarge</NodeType>
      <ClusterIdentifier>examplecluster</ClusterIdentifier>
      <AllowVersionUpgrade>true</AllowVersionUpgrade>
      <MasterUsername>masteruser</MasterUsername>
    </Cluster>
  </CreateClusterResult>
  <ResponseMetadata>
    <RequestId>e69b1294-64ef-11e2-b07c-f7fbdd006c67</RequestId>
  </ResponseMetadata>
</CreateClusterResponse>
`

var CreateClusterExample2 = `
<CreateClusterResponse xmlns="http://redshift.amazonaws.com/doc/2012-12-01/">
  <CreateClusterResult>
    <Cluster>
      <PendingModifiedValues>
        <MasterUserPassword>****</MasterUserPassword>
      </PendingModifiedValues>
      <ClusterSubnetGroupName>mysubnetgroup1</ClusterSubnetGroupName>
      <ClusterVersion>1.0</ClusterVersion>
      <VpcSecurityGroups/>
      <ClusterStatus>creating</ClusterStatus>
      <NumberOfNodes>2</NumberOfNodes>
      <AutomatedSnapshotRetentionPeriod>1</AutomatedSnapshotRetentionPeriod>
      <PubliclyAccessible>false</PubliclyAccessible>
      <Encrpyted>false</Encrypted>
      <DBName>dev</DBName>
      <PreferredMaintenanceWindow>sat:08:30-sat:09:00</PreferredMaintenanceWindow>
      <ClusterParameterGroups>
        <ClusterParameterGroup>
          <ParameterApplyStatus>in-sync</ParameterApplyStatus>
          <ParameterGroupName>default.redshift-1.0</ParameterGroupName>
        </ClusterParameterGroup>
      </ClusterParameterGroups>
      <VpcId>vpc-796a5913</VpcId>
      <ClusterSecurityGroups/>
      <NodeType>dw1.xlarge</NodeType>
      <ClusterIdentifier>exampleclusterinvpc</ClusterIdentifier>
      <AllowVersionUpgrade>true</AllowVersionUpgrade>
      <MasterUsername>master</MasterUsername>
    </Cluster>
  </CreateClusterResult>
  <ResponseMetadata>
    <RequestId>fa337bb4-6a4d-11e2-a12a-cb8076a904bd</RequestId>
  </ResponseMetadata>
</CreateClusterResponse>
`

var DeleteClusterExample1 = `
<DeleteClusterResponse xmlns="http://redshift.amazonaws.com/doc/2012-12-01/">
  <DeleteClusterResult>
    <Cluster>
      <PendingModifiedValues/>
      <ClusterVersion>1.0</ClusterVersion>
      <VpcSecurityGroups/>
      <Endpoint>
        <Port>5439</Port>
        <Address>examplecluster2.cobbanlpscsn.us-east-1.redshift.amazonaws.com</Address>
      </Endpoint>
      <ClusterStatus>deleting</ClusterStatus>
      <NumberOfNodes>2</NumberOfNodes>
      <AutomatedSnapshotRetentionPeriod>1</AutomatedSnapshotRetentionPeriod>
      <PubliclyAccessible>true</PubliclyAccessible>
      <Encrypted>true</Encrypted>
      <DBName>dev</DBName>
      <PreferredMaintenanceWindow>sun:10:30-sun:11:00</PreferredMaintenanceWindow>
      <ClusterParameterGroups>
        <ClusterParameterGroup>
          <ParameterApplyStatus>in-sync</ParameterApplyStatus>
          <ParameterGroupName>default.redshift-1.0</ParameterGroupName>
        </ClusterParameterGroup>
      </ClusterParameterGroups>
      <ClusterCreateTime>2013-01-23T00:11:32.804Z</ClusterCreateTime>
      <ClusterSecurityGroups>
        <ClusterSecurityGroup>
          <Status>active</Status>
          <ClusterSecurityGroupName>default</ClusterSecurityGroupName>
        </ClusterSecurityGroup>
      </ClusterSecurityGroups>
      <AvailabilityZone>us-east-1a</AvailabilityZone>
      <NodeType>dw1.xlarge</NodeType>
      <ClusterIdentifier>examplecluster2</ClusterIdentifier>
      <AllowVersionUpgrade>true</AllowVersionUpgrade>
      <MasterUsername>masteruser</MasterUsername>
    </Cluster>
  </DeleteClusterResult>
  <ResponseMetadata>
    <RequestId>f2e6b87e-6503-11e2-b343-393adc3f0a21</RequestId>
  </ResponseMetadata>
</DeleteClusterResponse>
`

var CreateClusterSubnetGroupExample1 = `
<CreateClusterSubnetGroupResponse xmlns="http://redshift.amazonaws.com/doc/2012-12-01/">
  <CreateClusterSubnetGroupResult>
    <ClusterSubnetGroup>
      <VpcId>vpc-796a5913</VpcId>
      <Description>My subnet group 1</Description>
      <ClusterSubnetGroupName>mysubnetgroup1</ClusterSubnetGroupName>
      <SubnetGroupStatus>Complete</SubnetGroupStatus>
      <Subnets>
        <Subnet>
          <SubnetStatus>Active</SubnetStatus>
          <SubnetIdentifier>subnet-756a591f</SubnetIdentifier>
          <SubnetAvailabilityZone>
            <Name>us-east-1c</Name>
          </SubnetAvailabilityZone>
        </Subnet>
      </Subnets>
    </ClusterSubnetGroup>
  </CreateClusterSubnetGroupResult>
  <ResponseMetadata>
    <RequestId>0a60660f-6a4a-11e2-aad2-71d00c36728e</RequestId>
  </ResponseMetadata>
</CreateClusterSubnetGroupResponse>
`

var DeleteClusterSubnetGroupExample1 = `
<DeleteClusterSubnetGroupResponse xmlns="http://redshift.amazonaws.com/doc/2012-12-01/">
  <ResponseMetadata>
    <RequestId>3a63806b-6af4-11e2-b27b-4d850b1c672d</RequestId>
  </ResponseMetadata>
</DeleteClusterSubnetGroupResponse>
`
