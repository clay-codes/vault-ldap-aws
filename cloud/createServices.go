package cloud

import (
	"sync"

	"github.com/aws/aws-sdk-go/aws" // AWS-specific configurations
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type AWSSession struct {
	session *session.Session
}

type Service struct {
	ec2 *ec2.EC2
	iam *iam.IAM
	ssm *ssm.SSM
}

var (
	instance        *AWSSession
	svc             *Service
	createSessOnce  sync.Once
	createServsOnce sync.Once
)

func GetSession() *AWSSession { return instance }

func (s *AWSSession) GetAWSSession() *session.Session { return instance.session }

func GetServices() *Service { return svc }

func CreateSession(region string) error {
	var err error
	createSessOnce.Do(func() {
		sess, sessErr := session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
		if sessErr != nil {
			err = sessErr
			return
		}
		instance = &AWSSession{session: sess}
	})
	return err
}

func CreateServices(sess *session.Session) error {
	createServsOnce.Do(func() {
		svc = &Service{
			ec2: ec2.New(sess),
			iam: iam.New(sess),
			ssm: ssm.New(sess),
		}
	})
	return nil
}
