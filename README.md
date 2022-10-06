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