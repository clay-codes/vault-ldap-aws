package cloud

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
)

// ... [Existing functions] ...

// CleanupResources performs the cleanup of all created resources

// terminateEC2Instance terminates the specified EC2 instance
func TerminateEC2Instance(sess *session.Session, instanceID string) error {
	ec2Svc := ec2.New(sess)
	_, err := ec2Svc.TerminateInstances(&ec2.TerminateInstancesInput{
		InstanceIds: []*string{aws.String(instanceID)},
	})
	if err != nil {
		return fmt.Errorf("error terminating EC2 instance: %v", err)
	}
	fmt.Println("EC2 instance terminated:", instanceID)
	return nil
}

// deleteKeyPair deletes the specified key pair
func DeleteKeyPair(sess *session.Session) error {
	ec2Svc := ec2.New(sess)
	_, err := ec2Svc.DeleteKeyPair(&ec2.DeleteKeyPairInput{
		KeyName: aws.String("vault-EC2-kp"),
	})
	if err != nil {
		return fmt.Errorf("error deleting key pair: %v", err)
	}
	fmt.Println("Key pair deleted:", "vault-EC2-kp")
	return nil
}

func getSGID(sess *session.Session) (string, error) {
	ec2Svc := ec2.New(sess)
	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("group-name"),
				Values: []*string{aws.String("EC2-Vault-SG")},
			},
		},
	}

	result, err := ec2Svc.DescribeSecurityGroups(input)
	if err != nil {
		return "", err
	}

	var groupIds []string
	for _, group := range result.SecurityGroups {
		groupIds = append(groupIds, *group.GroupId)
	}

	return groupIds[0], nil
}

// deleteSecurityGroup deletes the specified security group
func DeleteSecurityGroup(sess *session.Session) error {
	ec2Svc := ec2.New(sess)
	sgID, err := getSGID(sess)
	if err != nil {
		return fmt.Errorf("error in getSGID(): %v", err)
	}
	_, err = ec2Svc.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(sgID),
	})
	if err != nil {
		return fmt.Errorf("error deleting security group: %v", err)
	}
	fmt.Println("Security group deleted:", sgID)
	return nil
}

// detachRoleFromInstanceProfile detaches the specified role from the instance profile
func DetachRoleFromInstanceProfile(sess *session.Session) error {
	iamSvc := iam.New(sess)
	_, err := iamSvc.RemoveRoleFromInstanceProfile(&iam.RemoveRoleFromInstanceProfileInput{
		InstanceProfileName: aws.String("vault-ec2-InstProf"),
		RoleName:            aws.String("vault-ec2-metadata-role"),
	})
	if err != nil {
		return fmt.Errorf("error detaching role from instance profile: %v", err)
	}
	fmt.Println("Role detached from instance profile:", "vault-ec2-metadata-role", "vault-ec2-InstProf")
	return nil
}

// deleteInstanceProfile deletes the specified instance profile
func DeleteInstanceProfile(sess *session.Session) error {
	iamSvc := iam.New(sess)
	_, err := iamSvc.DeleteInstanceProfile(&iam.DeleteInstanceProfileInput{
		InstanceProfileName: aws.String("vault-ec2-InstProf"),
	})
	if err != nil {
		return fmt.Errorf("error deleting instance profile: %v", err)
	}
	fmt.Println("Instance profile deleted:", "vault-ec2-InstProf")
	return nil
}

// deleteRole deletes the specified IAM role
func DeleteRole(sess *session.Session) error {
	iamSvc := iam.New(sess)
	_, err := iamSvc.DeleteRole(&iam.DeleteRoleInput{
		RoleName: aws.String("vault-ec2-metadata-role"),
	})
	if err != nil {
		return fmt.Errorf("error deleting role: %v", err)
	}
	fmt.Println("IAM role deleted:", "vault-ec2-metadata-role")
	return nil
}
