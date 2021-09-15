# Okta Hackauth0n 2021

## Use the AWS CLI

1. Set the necessary environment variables
```shell
EC2_KEYPAIR_NAME={ec2_key_pair_name}
ASA_ENROLL_TOKEN={token_value}
```

2. Deploy with the AWS CLI
```shell
aws cloudformation create-stack --template-body file://template.yml --capabilities CAPABILITY_IAM --parameters ParameterKey=KeyName,ParameterValue=$EC2_KEYPAIR_NAME ParameterKey=EnrollmentToken,ParameterValue=$ASA_ENROLL_TOKEN --stack-name {unique-stack-name}
```

## Update the Lambda Function that Deploys the Candidate Stacks:

1. Compile the executable for the Lambda environment: 
```shell
GOOS=linux go build main.go
```

2. re-zip the deployment package:
```shell
zip ../createStack.zip main
```
