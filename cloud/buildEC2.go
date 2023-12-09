package cloud

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws" // AWS-specific configurations
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// buildEC2.go contains functions for building an EC2 instance
// the necessary parameters for building an EC2 instance are: image ID, instance type, key name, and security group
// which are retrieved in the below
func createSession(region string) (*session.Session, error) {
	// first get region from describe subnet

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			return nil, awsErr
		}
	}
	return sess, nil
}

func GetImgID(sess *session.Session) (string, error) {
	ssmSvc := ssm.New(sess)

	input := &ssm.GetParameterInput{
		Name: aws.String("/aws/service/ami-amazon-linux-latest/amzn2-ami-hvm-x86_64-gp2"),
	}

	result, err := ssmSvc.GetParameter(input)
	//aws-specific error library https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/handling-errors.html
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			return "Error  getting image-id", awsErr
		}
	}

	return *result.Parameter.Value, nil
}

func CreateKP(sess *session.Session) (string, error) {
	// Initialize a session in us-west-2 that the SDK will use to load credentials
	svc := ec2.New(sess)

	// Create the key pair
	input := &ec2.CreateKeyPairInput{
		KeyName: aws.String("vault-kp"),
		KeyType: aws.String("rsa"),
	}

	result, err := svc.CreateKeyPair(input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			return "Error creating key pair:", awsErr
		}
	}

	// Write the key material to a file
	file, err := os.Create("key.pem")
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			return "Error creating file:", awsErr
		}
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

func getVPC(sess *session.Session) (*string, error) {

	svc := ec2.New(sess)

	vpcs, err := svc.DescribeVpcs(nil)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			fmt.Println("Error describing VPCs: ")
			return nil, awsErr
		}
	}

	// Select the first VPC
	vpcID := vpcs.Vpcs[0].VpcId

	return vpcID, nil
}

func createSecurityGroup(sess *session.Session, name string, vpcID string, ports []int64) error {
	// Create EC2 client
	svc := ec2.New(sess)

	// Define the security group parameters
	createSGInput := &ec2.CreateSecurityGroupInput{
        GroupName:   aws.String(name),
        Description: aws.String("allowings all traffic in/out"),
        VpcId:       aws.String(vpcID), // Replace with your VPC ID
    }

    createSGOutput, err := svc.CreateSecurityGroup(createSGInput)
    if err != nil {
        return fmt.Println("Error creating security group:", err)
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
        fmt.Println("Error authorizing security group ingress:", err)
        return
    }
    fmt.Println("Inbound traffic allowed for all protocols and ports")

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
        fmt.Println("Error authorizing security group egress:", err)
        return
    }
    fmt.Println("Outbound traffic allowed for all protocols and ports")
}

