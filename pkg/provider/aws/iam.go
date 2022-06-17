package aws

import "github.com/aws/aws-sdk-go-v2/service/iam/types"

func listPoliciesScope() types.PolicyScopeType {
	return types.PolicyScopeTypeLocal
}
