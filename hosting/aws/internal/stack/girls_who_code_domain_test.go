package stack

import (
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/assertions"
	"github.com/aws/jsii-runtime-go"
)

func TestGirlsWhoCodeProductionDomain(t *testing.T) {
	defer jsii.Close()
	app := awscdk.NewApp(&awscdk.AppProps{Outdir: jsii.String(t.TempDir())})
	domain := NewGirlsWhoCodeDomainStack(app, "GirlsWhoCodeDomainStack", &GirlsWhoCodeDomainStackProps{
		DomainName: "example.org",
	})
	frontend := NewGirlsWhoCodeHostingStack(app, "FrontendStack", &GirlsWhoCodeHostingStackProps{
		DomainName: "example.org", HostedZone: domain.HostedZone,
	})
	zoneTemplate := assertions.Template_FromStack(domain.Stack, nil)
	template := assertions.Template_FromStack(frontend, nil)
	zoneTemplate.ResourceCountIs(jsii.String("AWS::Route53::HostedZone"), jsii.Number(1))
	zoneTemplate.HasResource(jsii.String("AWS::Route53::HostedZone"), map[string]interface{}{
		"Properties":     map[string]interface{}{"Name": "example.org.", "VPCs": assertions.Match_Absent()},
		"DeletionPolicy": "Retain",
	})
	zoneTemplate.HasOutput(jsii.String("HostedZoneId"), map[string]interface{}{"Value": assertions.Match_AnyValue()})
	zoneTemplate.HasOutput(jsii.String("NameServers"), map[string]interface{}{
		"Value": map[string]interface{}{"Fn::Join": assertions.Match_AnyValue()},
	})
	if len(*domain.Dependencies()) != 0 {
		t.Fatal("domain stack must not depend on the frontend stack")
	}
	if dependencies := *frontend.Dependencies(); len(dependencies) != 1 || dependencies[0] != domain.Stack {
		t.Fatal("frontend must depend only on the domain stack")
	}

	zoneImport := map[string]interface{}{"Fn::ImportValue": assertions.Match_AnyValue()}
	template.ResourceCountIs(jsii.String("AWS::CertificateManager::Certificate"), jsii.Number(1))
	template.HasResourceProperties(jsii.String("AWS::CertificateManager::Certificate"), map[string]interface{}{
		"DomainName": "example.org", "SubjectAlternativeNames": []string{"www.example.org"},
		"ValidationMethod": "DNS",
		"DomainValidationOptions": []interface{}{
			map[string]interface{}{"DomainName": "example.org", "HostedZoneId": zoneImport},
			map[string]interface{}{"DomainName": "www.example.org", "HostedZoneId": zoneImport},
		},
	})
	template.ResourceCountIs(jsii.String("AWS::CloudFront::Distribution"), jsii.Number(2))
	productionID := template.GetResourceId(jsii.String("AWS::CloudFront::Distribution"), map[string]interface{}{
		"Properties": map[string]interface{}{"DistributionConfig": map[string]interface{}{
			"Aliases": []string{"example.org", "www.example.org"}, "IPV6Enabled": true,
			"ViewerCertificate": map[string]interface{}{"AcmCertificateArn": assertions.Match_AnyValue()},
		}},
	})
	if *productionID != "FrontendProduction57D7F36D" {
		t.Fatal("existing production distribution logical ID changed")
	}
	stagingID := template.GetResourceId(jsii.String("AWS::CloudFront::Distribution"), map[string]interface{}{
		"Properties": map[string]interface{}{"DistributionConfig": map[string]interface{}{
			"Aliases": assertions.Match_Absent(), "ViewerCertificate": assertions.Match_Absent(),
		}},
	})
	if *stagingID != "FrontendMain4FAF8302" {
		t.Fatal("existing staging distribution logical ID changed")
	}
	template.ResourceCountIs(jsii.String("AWS::Route53::RecordSet"), jsii.Number(4))
	for _, name := range []string{"example.org.", "www.example.org."} {
		for _, recordType := range []string{"A", "AAAA"} {
			template.HasResourceProperties(jsii.String("AWS::Route53::RecordSet"), map[string]interface{}{
				"Name": name, "Type": recordType, "HostedZoneId": zoneImport,
				"AliasTarget": map[string]interface{}{
					"DNSName": map[string]interface{}{"Fn::GetAtt": []string{*productionID, "DomainName"}},
				},
			})
		}
	}
}
