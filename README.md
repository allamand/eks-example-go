# eks-example-go

A Very simple Go Web server that print some vars from request and the CLUSTER_NAME env var.
The image is store in a Public ECR repository

## Manually Build and push the image

First you need to create your public ECR repository:

```
aws ecr-public create-repository --repository-name eks-example-go --region us-east-1
```

```
REPO=<your reponame> make
```

```
go work init
#go work use tools tools/gopls
```

## Github Automation

For integration with AWS using actions we recommend using the [aws-actions/configure-aws-credentials](https://github.com/aws-actions/configure-aws-credentials) to configure the GitHub Actions environment with a role using Github's OIDC provider and your desired region.

You can configure the authentication step in your workflow

```
    - name: Configure AWS Credentials
      uses: aws-actions/configure-aws-credentials@v1
      with:
        role-to-assume: arn:aws:iam::123456789100:role/my-github-actions-role
        aws-region: us-east-2
```

Some security considerations:

- Do not store credentials in your repository's code.
- Grant least privilege to the credentials used in GitHub Actions workflows. Grant only the permissions required to perform the actions in your GitHub Actions workflows.
- Monitor the activity of the credentials used in GitHub Actions workflows.

We recommend using GitHub's OIDC provider to get short-lived credentials needed for your actions. Specifying role-to-assume without providing an aws-access-key-id or a web-identity-token-file will signal to the action that you wish to use the OIDC provider

### Create IAM Role

We need to create an IAM role with ECR public access to our repository only

We will create the folloginw policy

```
# You need to have ACCOUNT_ID in your environment
ECR_REPO=eks-example-go
cat << EOF> github-policy.json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "ManageRepositoryContents",
            "Effect": "Allow",
            "Action": [
                "ecr-public:DescribeImageTags",
                "ecr-public:DescribeImages",
                "ecr-public:InitiateLayerUpload",
                "ecr-public:DescribeRepositories",
                "ecr-public:UploadLayerPart",
                "ecr-public:PutImage",
                "ecr-public:CompleteLayerUpload",
                "ecr-public:GetRepositoryPolicy",
                "ecr-public:BatchCheckLayerAvailability",
                "ecr-public:CreateRepository"
            ],
            "Resource": "arn:aws:ecr-public::${ACCOUNT_ID}:repository/$ECR_REPO"
        },
        {
            "Sid": "GetAuthorizationToken",
            "Effect": "Allow",
            "Action": [
                "ecr-public:GetAuthorizationToken",
                "sts:GetServiceBearerToken"
            ],
            "Resource": "*"
        }
    ]
}
EOF
```

> Change the Resource field to include your account and targeted repository

### Configure the role for GitHub OIDC identity provider

If you use GitHub as an OIDC IdP, best practice is to limit the entities that can assume the role associated with the IAM IdP. When you include a condition statement in the trust policy, you can limit the role to a specific GitHub organization, repository, or branch.

```
# You need to have ACCOUNT_ID in your environment
GITHUB_ORG=allamand
GITHUB_REPO=eks-example-go
GITHUB_BRANCH=master
cat << EOF > github-trust-policy.json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::${ACCOUNT_ID}:oidc-provider/token.actions.githubusercontent.com"
      },
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringEquals": {
          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com",
          "token.actions.githubusercontent.com:sub": "repo:${GITHUB_ORG}/${GITHUB_REPO}:ref:refs/heads/${GITHUB_BRANCH}"
        }
      }
    }
  ]
}
EOF
```

Create the Role with GitHub IDP trust relationship

```
aws iam create-role --role-name GitHubAction-$GITHUB_ORG-$GITHUB_REPO --assume-role-policy-document file://github-trust-policy.json --description "Role to create public Repo from Github action for $GITHUB_ORG/$GITHUB_REPO for branch $GITHUB_BRANCH"
```

Create inline policy on the role

```
aws iam put-role-policy --role-name GitHubAction-$GITHUB_ORG-$GITHUB_REPO  --policy-name public-ecr --policy-document file://github-policy.json
```
