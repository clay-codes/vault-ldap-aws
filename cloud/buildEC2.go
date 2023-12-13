package cloud

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws" // AWS-specific configurations
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// buildEC2.go contains functions for building an EC2 instance
// the necessary parameters for building an EC2 instance are: image ID, instance type, key name, and security group
// which are retrieved in the below

type AWSSession struct {
	session *session.Session
	mu      sync.Mutex
}

var instance *AWSSession

// InitializeAWSSession creates and stores a new AWS session if it doesn't exist.
func CreateSession(region string) error {
	instance = &AWSSession{}
	instance.mu.Lock()
	defer instance.mu.Unlock()
	if instance.session == nil {
		sess, err := session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
		if err != nil {
			return err
		}
		instance.session = sess
	}
	return nil
}

// GetSession returns the entire AWS session struct.
func GetSession() *AWSSession {
	if instance == nil {
		CreateSession("us-west-2")
		log.Println("no current session previously, so CreateSession called with default region")
	}
	return instance
}

// returns session.session needed by the aws sdk
func (s *AWSSession) GetAWSSession() *session.Session {
	return instance.session
}

func GetImgID(sess *session.Session) (string, error) {
	ssmSvc := ssm.New(sess)

	input := &ssm.GetParameterInput{
		Name: aws.String("/aws/service/ami-amazon-linux-latest/amzn2-ami-hvm-x86_64-gp2"),
	}

	result, err := ssmSvc.GetParameter(input)
	//aws-specific error library https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/handling-errors.html
	if err != nil {
		return "Error  getting image-id", err
	}

	return *result.Parameter.Value, nil
}

func CreateKP(sess *session.Session) (string, error) {
	// Initialize a session in us-west-2 that the SDK will use to load credentials
	svc := ec2.New(sess)

	// Create the key pair
	input := &ec2.CreateKeyPairInput{
		KeyName: aws.String("vault-EC2-kp"),
		KeyType: aws.String("rsa"),
	}

	result, err := svc.CreateKeyPair(input)
	if err != nil {
		return "Error creating key pair:", err
	}

	// Write the key material to a file
	file, err := os.Create("key.pem")
	if err != nil {
		return "Error creating file:", err
	}
	defer file.Close()

	// write key material to file
	_, err = file.WriteString(*result.KeyMaterial)
	if err != nil {
		return "Error writing to file: ", err
	}

	// modify key.pem permissions to be read-only
	if err = os.Chmod("key.pem", 0400); err != nil {
		return "Error changing permissions: ", err
	}

	return "Created key pair", nil
}

func GetVPC(sess *session.Session) (*string, error) {

	svc := ec2.New(sess)

	vpcs, err := svc.DescribeVpcs(nil)
	if err != nil {
		fmt.Println("Error describing VPCs: ")
		return nil, err
	}

	// Select the first VPC
	vpcID := vpcs.Vpcs[0].VpcId

	return vpcID, nil
}

func CreateSG(sess *session.Session, vpcID *string, ports []int64) (string, error) {
	// Create EC2 client
	svc := ec2.New(sess)

	// Define the security group parameters
	createSGInput := &ec2.CreateSecurityGroupInput{
		GroupName:   aws.String("EC2-Vault-SG"),
		Description: aws.String("allowing all traffic in/out"),
		VpcId:       aws.String(*vpcID), // Replace with your VPC ID
	}

	createSGOutput, err := svc.CreateSecurityGroup(createSGInput)
	if err != nil {
		return "", fmt.Errorf("error creating security group: %v", err)
	}
	fmt.Println("Security Group Created with ID:", *createSGOutput.GroupId)

	// Authorize all inbound traffic
	authorizeIngressInput := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: createSGOutput.GroupId,
		IpPermissions: []*ec2.IpPermission{
			{
				IpProtocol: aws.String("-1"),
				FromPort:   aws.Int64(0),
				ToPort:     aws.Int64(65535),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp: aws.String("0.0.0.0/0"),
					},
				},
			},
		},
	}

	_, err = svc.AuthorizeSecurityGroupIngress(authorizeIngressInput)
	if err != nil {
		return "", fmt.Errorf("error authorizing security group ingress: %v", err)
	}
	// Authorize all outbound traffic
	authorizeEgressInput := &ec2.AuthorizeSecurityGroupEgressInput{
		GroupId: createSGOutput.GroupId,
		IpPermissions: []*ec2.IpPermission{
			{
				IpProtocol: aws.String("-1"),
				FromPort:   aws.Int64(0),
				ToPort:     aws.Int64(65535),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp: aws.String("0.0.0.0/0"),
					},
				},
			},
		},
	}

	_, err = svc.AuthorizeSecurityGroupEgress(authorizeEgressInput)
	if err != nil {
		return "", fmt.Errorf("error authorizing security group egress: %v", err)
	}
	// return the security group ID
	return *createSGOutput.GroupId, nil
}

func CreateInstProf(sess *session.Session) error {
	// Create an IAM client
	svc := iam.New(sess)

	// Define the trust policy document
	trustPolicyDocument := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {
					"Service": "ec2.amazonaws.com"
				},
				"Action": "sts:AssumeRole"
			}
		]
	}`

	// Create the role
	createRoleInput := &iam.CreateRoleInput{
		RoleName:                 aws.String("vault-ec2-metadata-role"),
		AssumeRolePolicyDocument: aws.String(trustPolicyDocument),
	}
	_, err := svc.CreateRole(createRoleInput)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	// Create the instance profile
	createInstanceProfileInput := &iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String("vault-ec2-InstProf"),
	}
	_, err = svc.CreateInstanceProfile(createInstanceProfileInput)
	if err != nil {
		return fmt.Errorf("error creating instance profile: %w", err)
	}

	// Attach the role to the instance profile
	addRoleToInstanceProfileInput := &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: aws.String("vault-ec2-InstProf"),
		RoleName:            aws.String("vault-ec2-metadata-role"),
	}
	_, err = svc.AddRoleToInstanceProfile(addRoleToInstanceProfileInput)
	if err != nil {
		return fmt.Errorf("error adding role to instance profile: %w", err)
	}

	return nil
}

func GetSubnetID(sess *session.Session, vpcID string) (string, error) {
	ec2Svc := ec2.New(sess)

	// Describe subnets with the specified VPC ID
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	}

	result, err := ec2Svc.DescribeSubnets(input)
	if err != nil {
		return "", fmt.Errorf("error describing subnets: %w", err)
	}
	if err != nil {
		return "", fmt.Errorf("error describing subnets: %w", err)
	}

	// Check if there is at least one subnet and get its ID
	if len(result.Subnets) == 0 {
		return "", fmt.Errorf("no subnets found for given VPC ID: %s", vpcID)
	}
	return *result.Subnets[0].SubnetId, nil
}

func BuildEC2(sess *session.Session, sgID []string, imageID, subnetID string) (string, error) {
	ec2Svc := ec2.New(sess)

	// Read user data from file
	userData, err := os.ReadFile("user-data.txt")
	if err != nil {
		return "", fmt.Errorf("error reading user data file: %v", err)
	}
	encodedUserData := base64.StdEncoding.EncodeToString(userData)

	// Run Instances
	input := &ec2.RunInstancesInput{
		ImageId:          aws.String(imageID),
		InstanceType:     aws.String("t2.micro"),
		KeyName:          aws.String("vault-EC2-kp"),
		SecurityGroupIds: aws.StringSlice(sgID),
		SubnetId:         aws.String(subnetID),
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String("vault-ec2-InstProf"),
		},
		UserData: aws.String(encodedUserData),
		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("vault-ad-server"),
					},
				},
			},
		},
		MinCount: aws.Int64(1),
		MaxCount: aws.Int64(1),
	}

	result, err := ec2Svc.RunInstances(input)
	if err != nil {
		return "", fmt.Errorf("error running instances: %v", err)
	}

	// Assuming only one instance is created
	if len(result.Instances) > 0 {
		return *result.Instances[0].InstanceId, nil
	}

	return "", fmt.Errorf("no instance was created")
}
