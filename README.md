# MY DOJO - Automate and build CI/CD on AWS

Hello and welcome. This is the repository for the MY Dojo session 5 - Automate and Build CI/CD on AWS.

![pipeline workflow](https://i.imgur.com/vS5wOt9.png)

## Buildspecs

There are a couple of buildspecs here:

1. buildspec.yml - This is to create the docker image based on the Dockerfile, tag it with the commit hash and push the image to ECR.
1. buildspec-image-scan.yml - This is to start the container image scanning and get back the result. Alerts will stop the pipeline.
1. buildspec-sonar.yml - This is to start the SAST analysis using sonarcloud.io. If the analysis does not pass the quality gate, the pipeline will stop.
1. buildspec-owasp.yml - This is to start the DAST analysis using OWASP ZAP. If the analysis contains alerts, the pipeline will stop.

## Slack integration

There are 2 Lambda functions to achieve this:

1. request-apporval.py - During the manual approval stage, an SNS will trigger this Lambda function, and will construct a message to be sent to Slack, with a Approve and Reject button.
1. process-approval.py - When a Slack user clicks on the Approve or Reject button, the response will be sent to an API Gateway and into this function to process the pipeline stage.
