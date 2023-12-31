package cloud

import (
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws" // AWS-specific configurations
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// buildEC2.go contains functions for building an EC2 instance
// the necessary parameters for building an EC2 instance are: image ID, instance type, key name, and security group
// which are retrieved in the below

func GetImgID() (string, error) {
	input := &ssm.GetParameterInput{
		Name: aws.String("/aws/service/ami-amazon-linux-latest/amzn2-ami-hvm-x86_64-gp2"),
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

func GetVPC() (string, error) {
	vpcs, err := svc.ec2.DescribeVpcs(nil)
	if err != nil {
		return "", fmt.Errorf("error when calling ec2.DescribeVpcs: %w", err)
	}

	// Select the first VPC
	vpcID := vpcs.Vpcs[0].VpcId

	return *vpcID, nil
}

func CreateSG(ports []int64) (string, error) {
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

	_, err = svc.ec2.AuthorizeSecurityGroupIngress(authorizeIngressInput)
	if err != nil {
		return "", fmt.Errorf("error authorizing security group ingress: %v", err)
	}
	// Authorize all outbound traffic
	authorizeEgressInput := &ec2.AuthorizeSecurityGroupEgressInput{
		GroupId: createSGOutput.GroupId,
		IpPermissions: []*ec2.IpPermission{
			{
				IpProtocol: aws.String("tcp"),
				FromPort:   aws.Int64(0),
				ToPort:     aws.Int64(65535),
				IpRanges: []*ec2.IpRange{
					{
						CidrIp:      aws.String("0.0.0.0/0"),
						Description: aws.String("for ec2-vault-sg-egress"),
					},
				},
			},
		},
	}

	_, err = svc.ec2.AuthorizeSecurityGroupEgress(authorizeEgressInput)
	if err != nil {
		return "", fmt.Errorf("error authorizing security group egress: %v", err)
	}
	// return the security group ID
	return *createSGOutput.GroupId, nil
}

func GetSGID() ([]string, error) {
	// ec2Svc := ec2.New(sess)
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
	_, err := svc.iam.CreateRole(createRoleInput)
	if err != nil {
		return fmt.Errorf("error: %w", err)
	}

	// Create the instance profile
	createInstanceProfileInput := &iam.CreateInstanceProfileInput{
		InstanceProfileName: aws.String("vault-ec2-InstProf"),
	}
	_, err = svc.iam.CreateInstanceProfile(createInstanceProfileInput)
	if err != nil {
		return fmt.Errorf("error creating instance profile: %w", err)
	}

	// Attach the role to the instance profile
	addRoleToInstanceProfileInput := &iam.AddRoleToInstanceProfileInput{
		InstanceProfileName: aws.String("vault-ec2-InstProf"),
		RoleName:            aws.String("vault-ec2-metadata-role"),
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

// func waitInstancePending(instanceID string) error {
// 	for {
// 		input := &ec2.DescribeInstancesInput{
// 			InstanceIds: []*string{aws.String(instanceID)},
// 		}

// 		result, err := svc.ec2.DescribeInstances(input)
// 		if err != nil {
// 			return err
// 		}

// 		state := result.Reservations[0].Instances[0].State.Name
// 		if *state == ec2.InstanceStateNamePending {
// 			break
// 		}

// 		time.Sleep(15 * time.Second) // Wait for 15 seconds before checking again
// 	}

// 	return nil
// }

func BuildEC2() (string, error) {

	// Read user data from file
	// _, err := os.ReadFile("user-data.txt")
	// if err != nil {
	// 	log.Printf("error reading user data file: %v", err)
	// 	//userData = []byte("")
	// }
	// encodedUserData := base64.StdEncoding.EncodeToString(userData)
	imageID, err := GetImgID()
	if err != nil {
		log.Printf("error getting image ID: %v", err)
	}
	sgID, err := GetSGID()
	if err != nil {
		log.Printf("error getting security group ID: %v", err)
	}
	subnetID, err := GetSubnetID()
	if err != nil {
		log.Printf("error getting subnet ID: %v", err)
	}

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
		//UserData: aws.String(encodedUserData),
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
		// waitInstancePending(*result.Instances[0].InstanceId)
		return *result.Instances[0].InstanceId, nil
	}

	return "", fmt.Errorf("no instance was created")
}

//unused function but may be useful in the future
func createAlternateSubnet() (string, error) {
	// Describe subnets to get the AZ of the first subnet
	vpcID, err := GetVPC()
	if err != nil {
		return "", fmt.Errorf("failed to get VPC for creating subnet: %s", err)
	}

	subnetsOutput, err := svc.ec2.DescribeSubnets(&ec2.DescribeSubnetsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []*string{aws.String(vpcID)},
			},
		},
	})
	if err != nil {
		return "", err
	}

	az := *subnetsOutput.Subnets[0].AvailabilityZone
	region := az[:len(az)-1]
	azChar := az[len(az)-1:]
	var newAzChar string
	if azChar == "a" {
		newAzChar = "b"
	} else {
		newAzChar = "a"
	}
	newAz := region + newAzChar

	// Describe VPCs to get the CIDR block
	vpcsOutput, err := svc.ec2.DescribeVpcs(&ec2.DescribeVpcsInput{
		VpcIds: []*string{aws.String(vpcID)},
	})
	if err != nil {
		return "", err
	}

	cidrBlock := *vpcsOutput.Vpcs[0].CidrBlock
	cidrPrefix := cidrBlock[:6]
	newSubnetCidr := fmt.Sprintf("%s.255.208/28", cidrPrefix)

	// Create the new subnet
	newSubnetOutput, err := svc.ec2.CreateSubnet(&ec2.CreateSubnetInput{
		VpcId:            aws.String(vpcID),
		AvailabilityZone: aws.String(newAz),
		CidrBlock:        aws.String(newSubnetCidr),
	})
	if err != nil {
		return "", err
	}

	return *newSubnetOutput.Subnet.SubnetId, nil

}