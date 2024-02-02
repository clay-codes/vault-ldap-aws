package cloud

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
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
	region 		    string
)



// SetRegion looks for the AWS region in the [default] profile of the AWS config file.
// If not found, it prompts the user to enter a region.
func SetRegion() {
    // Check environment variables
    region = os.Getenv("AWS_REGION")
    if region == "" {
        region = os.Getenv("AWS_DEFAULT_REGION")
    }

    // If not found in environment variables, try to read from AWS config file
    if region == "" {
        homeDir, err := os.UserHomeDir()
        if err == nil {
            configFile := fmt.Sprintf("%s/.aws/ponfig", homeDir)
            file, err := os.Open(configFile)
            if err == nil {
                defer file.Close()
                scanner := bufio.NewScanner(file)
                inDefaultProfile := false // Flag to track if we are in the [default] profile section
                for scanner.Scan() {
                    line := scanner.Text()
                    trimmedLine := strings.TrimSpace(line)

                    // Check if we've entered the [default] profile section
                    if trimmedLine == "[default]" {
                        inDefaultProfile = true
                        continue
                    }

                    // If another profile section starts, stop looking for the region
                    if inDefaultProfile && strings.HasPrefix(trimmedLine, "[") {
                        break
                    }

                    // If we're in the [default] profile section, look for the region setting
                    if inDefaultProfile && strings.HasPrefix(trimmedLine, "region") {
                        parts := strings.SplitN(trimmedLine, "=", 2)
                        if len(parts) == 2 {
                            region = strings.TrimSpace(parts[1])
							region = strings.ToLower(region)
                            break
                        }
                    }
                }
            }
        }
    }

    // If still not found, tell user to set env var
    if region == "" {
        fmt.Println("\nAWS default region not found in environment variables or config file ~/.aws/config")
        fmt.Println("Please set AWS region via below command")
        log.Fatal("\n\nexport AWS_REGION=<your-aws-region>\n\n ")
    }
}

func GetSession() *AWSSession { return instance }

func (s *AWSSession) GetAWSSession() *session.Session { return instance.session }

func GetServices() *Service { return svc }

func CreateSession() error {
	
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
// allows for specific service, if desired, or none to initialize all
func (s *AWSSession) CreateServices(serviceType ...string) error {
	sess := s.GetAWSSession()

	createServsOnce.Do(func() {
		svc = &Service{}

		if len(serviceType) == 0 {
			// Initialize all services if no specific service is provided
			svc.ec2 = ec2.New(sess)
			svc.iam = iam.New(sess)
			svc.ssm = ssm.New(sess)
		} else {
			// Initialize only the first specified service
			switch serviceType[0] {
			case "ec2":
				svc.ec2 = ec2.New(sess)
			case "iam":
				svc.iam = iam.New(sess)
			case "ssm":
				svc.ssm = ssm.New(sess)
			}
		}
	})

	return nil
}
