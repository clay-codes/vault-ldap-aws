package cloud

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
)

func getEC2ID() (string, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String("vault-ad-server")},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String(ec2.InstanceStateNameRunning), aws.String(ec2.InstanceStateNamePending)},
			},
		},
	}

	result, err := svc.ec2.DescribeInstances(input)
	if err != nil {
		return "", err
	}

	for _, reservation := range result.Reservations {
		for _, instance := range reservation.Instances {
			return *instance.InstanceId, nil
		}
	}

	return "", fmt.Errorf("no instances with the specified name in running or pending state")
}

func waitForInstanceTermination(instanceID string) error {
	for {
		input := &ec2.DescribeInstancesInput{
			InstanceIds: []*string{aws.String(instanceID)},
		}

		result, err := svc.ec2.DescribeInstances(input)
		if err != nil {
			return err
		}

		state := result.Reservations[0].Instances[0].State.Name
		if *state == ec2.InstanceStateNameTerminated {
			break
		}

		time.Sleep(15 * time.Second) // Wait for 15 seconds before checking again
	}
	return nil
}

// terminateEC2Instance terminates the specified EC2 instance

func TerminateEC2Instance() error {
	instanceID, err := getEC2ID()
	if err != nil {
		return fmt.Errorf("error in getEC2ID(): %v", err)
	}
	_, err = svc.ec2.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	})
	if err != nil {
		return fmt.Errorf("error terminating EC2 instance: %v", err)
	}
	waitForInstanceTermination(instanceID)
	fmt.Println("\nEC2 instance terminated                       ", instanceID)
	return nil
}

// deleteKeyPair deletes the specified key pair
func DeleteKeyPair() error {
	_, err := svc.ec2.DeleteKeyPair(&ec2.DeleteKeyPairInput{
		KeyName: aws.String("vault-EC2-kp"),
	})
	if err != nil {
		return fmt.Errorf("error deleting key pair: %v", err)
	}
	fmt.Println("Key pair deleted                               vault-EC2-kp")
	return nil
}

func DetachPolicyFromRole() error {
	_, err := svc.iam.DetachRolePolicy(&iam.DetachRolePolicyInput{
		PolicyArn: aws.String("arn:aws:iam::aws:policy/service-role/AmazonSSMAutomationRole"),
		RoleName:  aws.String("ec2-admin-role-custom"),
	})
	if err != nil {
		return fmt.Errorf("error detaching policy from role: %v", err)
	}
	fmt.Println("AWS role detatched from custom role            AmazonSSMAutomationRole")
	return nil
}

// detachRoleFromInstanceProfile detaches the specified role from the instance profile
func DetachRoleFromInstanceProfile() error {
	_, err := svc.iam.RemoveRoleFromInstanceProfile(&iam.RemoveRoleFromInstanceProfileInput{
		InstanceProfileName: aws.String("ec2-InstProf-custom"),
		RoleName:            aws.String("ec2-admin-role-custom"),
	})
	if err != nil {
		return fmt.Errorf("error detaching role from instance profile: %v", err)
	}
	fmt.Println("Custom role detached from instance profile     ec2-admin-role-custom")
	return nil
}

// deleteInstanceProfile deletes the specified instance profile
func DeleteInstanceProfile() error {
	_, err := svc.iam.DeleteInstanceProfile(&iam.DeleteInstanceProfileInput{
		InstanceProfileName: aws.String("ec2-InstProf-custom"),
	})
	if err != nil {
		return fmt.Errorf("error deleting instance profile: %v", err)
	}
	fmt.Println("Instance profile deleted                       ec2-InstProf-custom")
	return nil
}

// deleteRole deletes the specified IAM role
func DeleteRole() error {
	_, err := svc.iam.DeleteRole(&iam.DeleteRoleInput{
		RoleName: aws.String("ec2-admin-role-custom"),
	})
	if err != nil {
		return fmt.Errorf("error deleting role: %v", err)
	}
	fmt.Println("Custom role deleted                            ec2-admin-role-custom")
	return nil
}

// deleteSecurityGroup deletes the specified security group
func DeleteSecurityGroup() error {
	sgID, err := GetSGID()
	if err != nil {
		return fmt.Errorf("error in getSGID(): %v", err)
	}
	_, err = svc.ec2.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(sgID[0]),
	})
	if err != nil {
		return fmt.Errorf("error deleting security group: %v", err)
	}
	fmt.Println("Security group deleted                        ", sgID)
	return nil
}
