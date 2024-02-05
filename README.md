# Summary

policy-gen is a utility that allows you to generate cloud policies from file markers.  It is 
inspired by the [controller-gen](https://book.kubebuilder.io/reference/controller-gen) project 
which allows creation of Kubernetes CRDs by using inline markers.


## Background

Developing against public cloud APIs is relatively easy nowadays.  There are several intuitive
SDKs and many examples out there to get a developer started, as well as tooling such as 
Terraform and Anisble which allow more operations-fucused individuals to develop 
against public cloud APIs.  However, we oftentimes forget that actions taken against a 
public cloud API require specific access and get caught up in the tactical "make it work" 
mindset.  What this leads to is a working program, and developers scrambling to figure out 
permissions needed and the reasons behind those permissions.  This can leave consumers of 
developed projects curious and questioning as to why the specific project needs all these permissions!

This utility is designed to allow developers to write markers inline with their code, while 
developing, so that they may generate their policies before they forget the reasoning behind 
the "what is the permissions and why do I need it".

In addition to policies, documentation may also be produced.  This makes it handy to generate 
policies and documentation while developing your project.  This could even be ran as part of a 
CI/CD process to generate policies and documentation for each release.


## Installation

Installation of this project can be done via direct download, or via `brew`:

*Direct Download*:

```
VERSION=v0.0.1
OS=Linux
ARCH=x86_64
wget https://github.com/scottd018/policy-gen/releases/download/${VERSION}/policy-gen_v${VERSION}_${OS}_${ARCH}.tar.gz -O - |\
    tar -xz && sudo mv policy-gen /usr/local/bin/policy-gen
```

*Brew*:

```
brew install scottd018/tap/policy-gen
```


## Usage

All usage is stored as part of the command.  You can see how to use the command by running:

```
policy-gen <cloud> help

# Example:
policy-gen aws help
```


## Examples

Currently, this only works for AWS policies.  Others may be able to be developed as needed.

### AWS

As a developer, I am writing a custom go application that creates an AWS network for
my application.  The code, including markers (see [markers](#markers)), looks like this:

```go
// Create creates an AWS VPC.
func (v *vpc) Create(ctx context.Context, client client.Cloud) error {
	// create the vpc
	cidrString := v.CIDR.String()
	input := &ec2.CreateVpcInput{
		CidrBlock:         &cidrString,
		TagSpecifications: tags.Specifications(v, tags.WithAdditional(v, v.Config.Tags)...),
	}

    // +policy-gen:aws:iam:policy:name=createvpc,action=`ec2:CreateVpc*`,effect=Allow,reason=`User needs to init a 
    // VPC where this application resides in.`
	output, err := client.AWS().EC2().CreateVpc(ctx, input)
	if err != nil {
		return resources.CreateError(v, err)
	}

	v.Object = output.Vpc

	return nil
}
```

You can see above that I have defined a `policy-gen:aws:iam:policy` marker inline with my code.  I was able to write 
this as I developed the application so that the justification for the policy, as well as the needed policy definition 
is not lost during the development cycle.

Once defined, I can run a command for my code (assuming the above is at `internal/pkg/aws/testinput`) to generate my 
policies to the `internal/pkg/aws/testoutput` directory with an accompanying documentation file:

```
policy-gen aws \
    --input-path=internal/pkg/aws/testinput \
    --output-path=internal/pkg/aws/testoutput \
    --documentation=internal/pkg/aws/testoutput/README.md \
    --force \
    --debug
```

And I can view my generated policy document:

```
cat internal/pkg/aws/testoutput/createvpc.json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "Default",
            "Effect": "Allow",
            "Action": [
                "ec2:CreateVpc*"
            ],
            "Resource": [
                "*"
            ]
        }
    ]
}
```


## Markers

Markers are parsed from the `--input-path` flag of the `policy-gen` command.  This utility 
will scan all files for the existence of markers that are prefixed with 
`+policy-gen`.

### AWS

*Sample*:

```
+policy-gen:aws:iam:policy:name=test,id=ME,action=`ec2:Describe*`,effect=Allow,resource=`*`,reason=`because i said so`
```

Please be familiar with [AWS IAM Policies](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/iam-policy-structure.html)
to better understand what these fields mean and what effect they have on your policies.

Specific to AWS, markers that have `+policy-gen:aws:iam:policy` are parsed.  The following 
arguments to this command are accepted:

| Field    | Type                           | Default   | Required |
| ---------| ------------------------------ | --------- | -------- |
| name     | string                         | ""        | true     |
| id       | string                         | "Default" | false    |
| action   | string                         | ""        | true     |
| resource | string                         | "*"       | false    |
| effect   | string ("Allow" or "Deny")     | "Allow"   | false    |
| reason   | string                         | ""        | false    |

* **name**: name of the specific policy.  This will be used as the generated file name.  Markers
which shared the same name value will become separate statements within the same file.

* **id**: the statement ID of the specific policy.  Actions and resources are merged if a 
statement ID matches for the marker.

* **action**: the action or permission that this policy allows or denies, as specified 
by the `effect` field.

* **resource**: the resource that the action applies to.

* **effect**: whether this policy should `Allow` or `Deny` the `action` field.  If an existing 
statement ID has a mismatched effect, a new statement ID is created with an appended effect.  This 
is because we cannot have Allow/Deny effects in the same statement.

* **reason**: the reason for the necessary permission.  Needed when generation of 
markdown documentation with the `--documentation` flag.
