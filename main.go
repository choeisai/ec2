package main

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"

	"fmt"
)

// EnableInstanceDisableAPITermination to enable termination protection
func EnableInstanceDisableAPITermination(instanceID string, ec2Svc *ec2.EC2) error {
	// Defind as type *bool
	valueToModify := true
	// & == pointer
	input := &ec2.ModifyInstanceAttributeInput{
		InstanceId: aws.String(instanceID),
		DisableApiTermination: &ec2.AttributeBooleanValue{
			Value: &valueToModify,
		},
	}

	_, err := ec2Svc.ModifyInstanceAttribute(input)

	if err != nil {
		return err
	}
	return nil
}

func CreateTag(instanceID, tag, value string, ec2Svc *ec2.EC2) error {
	input := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(instanceID),
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String(tag),
				Value: aws.String(value),
			},
		},
	}

	_, err := ec2Svc.CreateTags(input)
	if err != nil {
		return err
	}
	return nil
}

type EC2Info struct {
	Customer   string
	InstanceID string
	CoreDNS    string
}

func main() {
	// Load session from shared config
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create new EC2 client
	ec2Svc := ec2.New(sess)

	result, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		panic(err)
	}

	customers := map[string][]EC2Info{}
	for _, v := range result.Reservations {
		instance := v.Instances[0]
		tags := instance.Tags
		// name := ""
		// customer := ""
		serverType := ""
		// coreDNS := ""

		info := new(EC2Info)

		for _, v := range tags {
			// if *v.Key == "Name" {
			//	name = *v.Value
			// }
			if *v.Key == "Customer" {
				// customer = *v.Value
				info.Customer = *v.Value
			}
			if *v.Key == "Type" {
				serverType = *v.Value
			}
			if *v.Key == "CoreDNS" {
				// coreDNS = *v.Value
				info.CoreDNS = *v.Value
			}
		}
		info.InstanceID = *instance.InstanceId
		// fmt.Println(name)
		// fmt.Println(customer)
		// fmt.Println(serverType)
		// fmt.Println(*instance.InstanceId)
		tagValue := fmt.Sprintf(`km-%s`, serverType)
		tagValue = strings.Replace(tagValue, "_server", "", -1)
		tagValue = strings.Replace(tagValue, "elasticsearch", "es", -1)
		fmt.Printf("%s\t%s\t%s\n", "km", *instance.InstanceId, tagValue)
		tag := "CoreDNS"
		if err := CreateTag(*instance.InstanceId, tag, tagValue, ec2Svc); err != nil {
			fmt.Println("Err", err)
		}

		// if customer != "" && serverType != "" && serverType != "http_asg_server" {
		//	customers[info.Customer] = append(customers[info.Customer], *info)
		//	tagValue := fmt.Sprintf(`%s-%s`, customer, serverType)
		//	tagValue = strings.Replace(tagValue, "_server", "", -1)
		//	tagValue = strings.Replace(tagValue, "elasticsearch", "es", -1)
		//	fmt.Printf("%s\t%s\t%s\n", customer, *instance.InstanceId, coreDNS)
		//	fmt.Printf("%s\t%s\t%s\n", info.Customer, info.InstanceID, info.CoreDNS)
		//	tag := "CoreDNS"
		//	if err := CreateTag(*instance.InstanceId, tag, tagValue, ec2Svc); err != nil {
		//		fmt.Println("Err", err)
		//	}
		// }

		//	// if err := EnableInstanceDisableAPITermination(*instance.InstanceId, ec2Svc); err != nil {
		//	//	fmt.Println("Err", err)
		//	// }
		// fmt.Println("")

	}
	for _, v1 := range customers {
		fmt.Printf("%s", v1[0].Customer)
		for _, v2 := range v1 {
			fmt.Printf("\t%s\t%s.sea\n", v2.InstanceID, v2.CoreDNS)
		}
	}

}
