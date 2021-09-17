# Okta Hackauth0n 2021

This is a collection of cloudformation templates and lambda functions. The template in the root folder launches an EC2 instance that auto-enrolls in Okta's Advanced Server Access (ASA) and serves as a bastion server. The template also deploys lambda functions to create and delete additional EC2 instances that are pre-configured as Microsoft Domain Controller server demo environments with Microsoft Active Directory (AD) available. 

These Microsoft servers can only be reached via port `3389` (RDP) from the bastion server IP. The creation and destruction of these servers is managed by triggering the deployed Lambda functions. Lambda function integration is configurable in Okta Workflows.

The idea is to provide a quick way to launch any number of pre-configured Microsoft domain controllers with AD services enabled, but to restrict access to those services from the internet except through the Okta-protected bastion server. In other words, a user must be granted ASA privileges to log into the bastion server, then they can create a tcp proxy to send all traffic from their local RDP client (any port) to port `3389` on their target Microsoft server.

The practical use of this is to integrate it with Okta workflows. Workflows can trigger the deployed Lambda functions. With an Okta workflow, something as simple as a user creation or group assignment can kick off the spin-up of an entire demo environment.

## ASA Server Enrollment

ASA Enrollment of the bastion server occurs automatically, so long as a valid `EnrollmentToken` is provided when launching the Cloudformation template.

The `template.yml` in the root folder creates resource `AsaBastion` and enrolls an AWS Linux ami in Okta's Advanced Server Access with the EnrollmentToken provided as an input param. The script below shows how this occurs.

```shell
mkdir -p /var/lib/sftd
echo '${EnrollmentToken}' > /var/lib/sftd/enrollment.token
curl -C - https://pkg.scaleft.com/scaleft_yum.repo | sudo tee /etc/yum.repos.d/scaleft.repo
sudo rpm --import https://dist.scaleft.com/pki/scaleft_rpm_key.asc
sudo yum install scaleft-server-tools -y
sudo yum groupinstall "Development Tools" -y
git clone https://github.com/vzaliva/simpleproxy.git
cd simpleproxy && ./configure && sudo make install
```

## Using Workflows to invoke the Lambda

The Lambda function expects an object in the payload with a single param:
```json
{
    "stackName": "demo-microsoft-server-stack-12345"
}
```

Each stack name must be unique, so it is recommended to include a timestamp or uuid in the stackName. **As a best practice, include a descriptor of the principal or event which caused the invocation of the Lambda function in the `stackName` - this way you can recognize the purpose of the stack in the future.**

A workflow can be created to form a `stackName`, then invoke the Lambda function to launch the Microsoft server

## Using the bastion as a TCP proxy

The script above also installs the `simpleproxy` package from source. Use the `simpleproxy` package to open a TCP proxy to the Microsoft servers that are launched by invoking the `createStack` Lambda function with the following command:

```shell
# Forward bastion server port 3000 to the Microsoft server RDP port 3389
simpleproxy -L 3000 -R {msft-server-ip}:3389
```

Once you've authenticated via ASA and started the TCP proxy, you can use an RDP client on your local machine to connect to {bastion-ip}:3000 - your RDP client will connect to the target Microsoft server.

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
