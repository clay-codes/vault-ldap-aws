package cloud

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/ssm"
)
// user-data.yaml file b64 encoded, necessary for binary, as the actual user-data.yaml file will not be compiled to binary upon `go build`
const EncodedUserData = "dmVyc2lvbjogMS4xCnRhc2tzOgotIHRhc2s6IGVuYWJsZU9wZW5Tc2gKLSB0YXNrOiBleGVjdXRlU2NyaXB0CiAgaW5wdXRzOgogIC0gZnJlcXVlbmN5OiBvbmNlCiAgICB0eXBlOiBwb3dlcnNoZWxsCiAgICBydW5BczogYWRtaW4KICAgIGNvbnRlbnQ6IHwtCiAgICAgICR0ZW1wRmlsZVBhdGggPSBbU3lzdGVtLklPLlBhdGhdOjpHZXRUZW1wRmlsZU5hbWUoKQogICAgICBzZWNlZGl0IC9leHBvcnQgL2NmZyAkdGVtcEZpbGVQYXRoCiAgICAgICRjb25maWcgPSBHZXQtQ29udGVudCAtUGF0aCAkdGVtcEZpbGVQYXRoCiAgICAgICRjb25maWcgPSAkY29uZmlnIC1yZXBsYWNlICJeTWluaW11bVBhc3N3b3JkQWdlXHMqPVxzKlxkKyIsICJNaW5pbXVtUGFzc3dvcmRBZ2UgPSAwIgogICAgICAkY29uZmlnID0gJGNvbmZpZyAtcmVwbGFjZSAiXk1heGltdW1QYXNzd29yZEFnZVxzKj1ccypcZCsiLCAiTWF4aW11bVBhc3N3b3JkQWdlID0gLTEiCiAgICAgICRjb25maWcgPSAkY29uZmlnIC1yZXBsYWNlICJeTWluaW11bVBhc3N3b3JkTGVuZ3RoXHMqPVxzKlxkKyIsICJNaW5pbXVtUGFzc3dvcmRMZW5ndGggPSAwIgogICAgICAkY29uZmlnID0gJGNvbmZpZyAtcmVwbGFjZSAiXlBhc3N3b3JkQ29tcGxleGl0eVxzKj1ccypcZCsiLCAiUGFzc3dvcmRDb21wbGV4aXR5ID0gMCIKICAgICAgJGNvbmZpZyB8IFNldC1Db250ZW50IC1QYXRoICR0ZW1wRmlsZVBhdGgKICAgICAgc2VjZWRpdCAvY29uZmlndXJlIC9kYiAkZW52OndpbmRpclxzZWN1cml0eVxsb2NhbC5zZGIgL2NmZyAkdGVtcEZpbGVQYXRoIC9hcmVhcyBTRUNVUklUWVBPTElDWQogICAgICBncHVwZGF0ZSAvZm9yY2UKICAgICAgUmVtb3ZlLUl0ZW0gJHRlbXBGaWxlUGF0aAogICAgICBuZXQgdXNlciBBZG1pbmlzdHJhdG9yIGFkbWluCiAgICAgIEluc3RhbGwtV2luZG93c0ZlYXR1cmUgQUQtRG9tYWluLVNlcnZpY2VzIC1JbmNsdWRlTWFuYWdlbWVudFRvb2xzCiAgICAgIEltcG9ydC1Nb2R1bGUgQUREU0RlcGxveW1lbnQKICAgICAgSW5zdGFsbC1BRERTRm9yZXN0IGAKICAgICAgLUNyZWF0ZURuc0RlbGVnYXRpb246JGZhbHNlIGAKICAgICAgLURhdGFiYXNlUGF0aCAiQzpcV2luZG93c1xOVERTIiBgCiAgICAgIC1Eb21haW5Nb2RlICJXaW5UaHJlc2hvbGQiIGAKICAgICAgLURvbWFpbk5hbWUgInZhdWx0ZXN0LmNvbSIgYAogICAgICAtRG9tYWluTmV0Ymlvc05hbWUgIlZBVUxURVNUIiBgCiAgICAgIC1Gb3Jlc3RNb2RlICJXaW5UaHJlc2hvbGQiIGAKICAgICAgLUluc3RhbGxEbnM6JHRydWUgYAogICAgICAtTG9nUGF0aCAiQzpcV2luZG93c1xOVERTIiBgCiAgICAgIC1Ob1JlYm9vdE9uQ29tcGxldGlvbjokZmFsc2UgYAogICAgICAtU3lzdm9sUGF0aCAiQzpcV2luZG93c1xTWVNWT0wiIGAKICAgICAgLVNhZmVNb2RlQWRtaW5pc3RyYXRvclBhc3N3b3JkIChDb252ZXJ0VG8tU2VjdXJlU3RyaW5nIC1Bc1BsYWluVGV4dCAiVmF1bHREU1JNUGFzc3cwcmQhIiAtRm9yY2UpIGAKICAgICAgLUZvcmNlOiR0cnVlCgo="

func GetImgID() (string, error) {
	input := &ssm.GetParameterInput{
		Name: aws.String("/aws/service/ami-windows-latest/Windows_Server-2022-English-Full-Base"),
	}

	result, err := svc.ssm.GetParameter(input)
	//aws-specific error library https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/handling-errors.html
	if err != nil {
		return "", err
	}

	return *result.Parameter.Value, nil
}

func CreateKP() (string, error) {
	// Create the key pair
	input := &ec2.CreateKeyPairInput{
		KeyName: aws.String("vault-EC2-kp"),
		KeyType: aws.String("rsa"),
	}

	result, err := svc.ec2.CreateKeyPair(input)
	if err != nil {
		return "", fmt.Errorf("error creating key pair: %w", err)
	}

	// Write the key material to a file
	file, err := os.Create("key.pem")
	if err != nil {
		return "", fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	// write key material to file
	_, err = file.WriteString(*result.KeyMaterial)
	if err != nil {
		return "", fmt.Errorf("error writing to file: %w", err)
	}

	// modify key.pem permissions to be read-only
	if err = os.Chmod("key.pem", 0400); err != nil {
		return "", fmt.Errorf("error changing file permissions: %w", err)
	}
	return *result.KeyName, nil
}

func GetVPC() (string, error) {
	vpcs, err := svc.ec2.DescribeVpcs(nil)
	if err != nil {
		return "", fmt.Errorf("error when calling ec2.DescribeVpcs: %w", err)
	}

	// Select the first VPC
	vpcID := vpcs.Vpcs[0].VpcId

	return *vpcID, nil
}

func CreateSG() (string, error) {
	// Create EC2 client
	vpcID, err := GetVPC()
	if err != nil {
		return "", err
	}
	// Define the security group parameters
	createSGInput := &ec2.CreateSecurityGroupInput{
		GroupName:   aws.String("EC2-Vault-SG"),
		Description: aws.String("sg for vault instance"),
		VpcId:       aws.String(vpcID), // Replace with your VPC ID
	}

	createSGOutput, err := svc.ec2.CreateSecurityGroup(createSGInput)
	if err != nil {
		return "", fmt.Errorf("error creating security group: %v", err)
	}

	// Authorize all inbound traffic
	authorizeIngressInput := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: createSGOutput.GroupId,
		IpPermissions: []*ec2.IpPermission{
			{
				IpProtocol: aws.String("-1"),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String("0.0.0.0/0"),
						Description: aws.String("for ec2-vault-sg-ingress"),
					},
				},
			},
		},
	}

	_, err = svc.ec2.AuthorizeSecurityGroupIngress(authorizeIngressInput)
	if err != nil {
		return "", fmt.Errorf("error authorizing security group ingress: %v", err)
	}
	// NOTE: AWS ALREADY HAS A DEFAULT EGRESS RULE ALLOWING ALL TRAFFIC, SO NO NEED TO AUTHORIZE EGRESS
	return *createSGOutput.GroupId, nil
}

func GetSGID() ([]string, error) {
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("group-name"),
				Values: []*string{aws.String("EC2-Vault-SG")},
			},
		},
	}

	result, err := svc.ec2.DescribeSecurityGroups(input)
	if err != nil {
		return nil, err
	}

	var groupIds []string
	for _, group := range result.SecurityGroups {
		groupIds = append(groupIds, *group.GroupId)
	}

	return groupIds, nil
}

func CreateInstProf() error {
	// Define the trust policy document
	policyDocument := `{
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
		RoleName:                 aws.String("ec2-admin-role-custom"),
		AssumeRolePolicyDocument: aws.String(policyDocument),
	}
	_, err := svc.iam.CreateRole(createRoleInput)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	input := &iam.AttachRolePolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/service-role/AmazonSSMAutomationRole"),
		RoleName:  aws.String("ec2-admin-role-custom"),
	}

	_, err = svc.iam.AttachRolePolicy(input)
	if err != nil {
		return err
	}
	// Create the instance profile
	createInstanceProfileInput := &iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String("ec2-InstProf-custom"),
	}
	_, err = svc.iam.CreateInstanceProfile(createInstanceProfileInput)
	if err != nil {
		return fmt.Errorf("error creating instance profile: %w", err)
	}

	// wait for instance profile to be created
	time.Sleep(5 * time.Second)

	// Attach the role to the instance profile
	addRoleToInstanceProfileInput := &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: aws.String("ec2-InstProf-custom"),
		RoleName:            aws.String("ec2-admin-role-custom"),
	}
	_, err = svc.iam.AddRoleToInstanceProfile(addRoleToInstanceProfileInput)
	if err != nil {
		return fmt.Errorf("error adding role to instance profile: %w", err)
	}

	return nil
}

func GetSubnetID() (string, error) {
	vpcID, err := GetVPC()
	if err != nil {
		return "", err
	}
	// Describe subnets with the specified VPC ID
	input := &ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	}

	result, err := svc.ec2.DescribeSubnets(input)
	if err != nil {
		return "", fmt.Errorf("error describing subnets: %w", err)
	}
	// Check if there is at least one subnet and get its ID
	if len(result.Subnets) == 0 {
		return "", fmt.Errorf("no subnets found for given VPC ID: %s", vpcID)
	}
	return *result.Subnets[0].SubnetId, nil
}
//  use this function if you would like to customize the user-data.yaml file
//	func encodeUserData() (string, error) {
//		// Read user data from file
//		userData, err := os.ReadFile("user-data.yaml")
//		if err != nil {
//			return "", err
//		}
//		return base64.StdEncoding.EncodeToString(userData), nil
//	}
func GetEC2ID() (string, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String("vault-ad-server")},
			},
		},
	}

	result, err := svc.ec2.DescribeInstances(input)
	if err != nil {
		return "", fmt.Errorf("error describing instances: %v", err)
	}

	return *result.Reservations[0].Instances[0].InstanceId, nil
}

func GetPublicDNS(instanceID *string) (string, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{
			instanceID,
		},
	}

	result, err := svc.ec2.DescribeInstances(input)
	if err != nil {
		return "", fmt.Errorf("error describing instances: %v", err)
	}

	return *result.Reservations[0].Instances[0].PublicDnsName, nil
}

func BuildEC2() (string, error) {
	encodedUserData := EncodedUserData
	// can use instead of EncodedUserData if you would like to customize the user-data.yaml file
	// encodedUserData, err := encodeUserData()
	// if err != nil {
	// 	return "", fmt.Errorf("error encoding user data: %v", err)
	// }

	imageID, err := GetImgID()
	if err != nil {
		return "", fmt.Errorf("error getting image ID: %v", err)
	}

	sgID, err := GetSGID()
	if err != nil {
		return "", fmt.Errorf("error getting security group ID: %v", err)
	}

	subnetID, err := GetSubnetID()
	if err != nil {
		return "", fmt.Errorf("error getting subnet ID: %v", err)
	}

	input := &ec2.RunInstancesInput{
		ImageId:          aws.String(imageID),
		InstanceType:     aws.String("t2.micro"),
		KeyName:          aws.String("vault-EC2-kp"),
		SecurityGroupIds: aws.StringSlice(sgID),
		SubnetId:         aws.String(subnetID),
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			Name: aws.String("ec2-InstProf-custom"),
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

	result, err := svc.ec2.RunInstances(input)
	if err != nil {
		return "", fmt.Errorf("error running instances: %v", err)
	}

	// Assuming only one instance is created
	if len(result.Instances) > 0 {
		// Instance must be in the running state before we can get its public DNS
		err := svc.ec2.WaitUntilInstanceRunning(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{result.Instances[0].InstanceId},
		})
		if err != nil {
			return "", fmt.Errorf("error waiting for instance to run: %v", err)
		}
		pubDNS, err := GetPublicDNS(result.Instances[0].InstanceId)
		if err != nil {
			return "", fmt.Errorf("error getting public DNS: %v", err)
		}
		return pubDNS, nil
	}

	return "", fmt.Errorf("no instance was created")
}
