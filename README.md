# Okta Hackauth0n 2021

This is a collection of cloudformation templates and lambda functions. The template in the root folder launches an EC2 instance that auto-enrolls in Okta's Advanced Server Access (ASA) and serves as a bastion server. The template also deploys lambda functions to create and delete additional EC2 instances that are pre-configured as Microsoft Domain Controller server demo environments with Microsoft Active Directory (AD) available. 

These Microsoft servers can only be reached via port `3389` (RDP) from the bastion server IP. The creation and destruction of these servers is managed by triggering the deployed Lambda functions. Lambda function integration is configurable in Okta Workflows.

The idea is to provide a quick way to launch any number of pre-configured Microsoft domain controllers with AD services enabled, but to restrict access to those services from the internet except through the Okta-protected bastion server. In other words, a user must be granted ASA privileges to log into the bastion server, then they can create a tcp proxy to send all traffic from their local RDP client (any port) to port `3389` on their target Microsoft server.

The practical use of this is to integrate it with Okta workflows. Workflows can trigger the deployed Lambda functions. With an Okta workflow, something as simple as a user creation or group assignment can kick off the spin-up of an entire demo environment.

## TODO

I ran out of time before I could write a `deleteStack` function and create another workflow for it, which triggers stack tear-down when a user is removed from the target group... that should come next.

I also need to clean up the roles - namely, reduce role privileges and principals down to the bare requirement.

## Use the AWS CLI to deploy

1. Set the necessary environment variables. You must bring your own ec2 key pair. Note that the template currently leaves port `22` open, but since the server auto-enrolls in Okta ASA, this is unnecessary and you may not ever need your key pair.
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
