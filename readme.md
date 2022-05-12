# Domain Management
Blueprint for end-to-end domain management.\
Providers - Godaddy, Google\
Payments - Stripe\
Frontend - Web & iOS app, Sign in via Google/Apple. No Android yet.\

Automates querying, purchasing, sire verification, dns provision, container building, service creation with custom domain name.\

# GCP prerequisites:
Replace keys:
* SERVICE_ACCOUNT ... abc@xyz.iam.gserviceaccount.com
* PROJECT_NAME ... xyz
* SERVICE_ACCOUNT_NAME ... abc
* SERVICE_ACCOUNT_EMAIL ... abc@gmail.com
* GOOGLE_KEY ... json
The default record set quota for managed zones is 10,000.\

Set up & enable APIs
```
gcloud auth login
gcloud config set project PROJECT_NAME
gcloud config set region us-central1
gcloud services enable artifactregistry.googleapis.com
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable iam.googleapis.com
gcloud services enable domains.googleapis.com
gcloud services enable firestore.googleapis.com
gcloud services enable dns.googleapis.com
gcloud services enable compute.googleapis.com
gcloud services enable siteverification.googleapis.com
```

Double check container image repo is ready for builds. For production triggering builds by pushing to master is useful.
```
gcloud artifacts locations list
gcloud artifacts repositories create cloud-run-source-deploy --repository-format docker --location us --description "PROJECT_NAME-repos" --async
gcloud artifacts repositories list
```

Policy bind service account
```
gcloud iam service-accounts create SERVICE_ACCOUNT_NAME --description="DESCRIPTION" --display-name="DISPLAY_NAME"
gcloud projects add-iam-policy-binding PROJECT_NAME --member="serviceAccount:SERVICE_ACCOUNT" --role="roles/run.admin"
gcloud projects add-iam-policy-binding PROJECT_NAME --member="serviceAccount:SERVICE_ACCOUNT" --role="roles/domains.admin"
gcloud projects add-iam-policy-binding PROJECT_NAME --member="serviceAccount:SERVICE_ACCOUNT" --role="roles/datastore.user"
gcloud projects add-iam-policy-binding PROJECT_NAME --member="serviceAccount:SERVICE_ACCOUNT" --role="roles/dns.admin"
gcloud projects add-iam-policy-binding PROJECT_NAME --member="serviceAccount:SERVICE_ACCOUNT" --role="roles/compute.admin"
gcloud iam service-accounts add-iam-policy-binding SERVICE_ACCOUNT --member="user:SERVICE_ACCOUNT_EMAIL" --role="roles/iam.serviceAccountUser"
gcloud iam service-accounts enable SERVICE_ACCOUNT
gcloud projects add-iam-policy-binding PROJECT_NAME --member="SERVICE_ACCOUNT" --role="roles/iam.serviceAccountKeyAdmin"
gcloud iam service-accounts keys create GOOGLE_KEY --iam-account SERVICE_ACCOUNT
```
Retrieve App Credentials

# Prepare local environment
```
export GOOGLE_APPLICATION_CREDENTIALS
```

import GoogleSignOn. Firebase in Xcode.\
From firebase obtain URL scheme and place in Xcode.\
Retrieve Keys from Stripe.\


# Manual gcloud commands/flow

```
gcloud run deploy SERVICE_NAME --source . --service-account SERVICE_ACCOUNT --allow-unauthenticated --region us-central1
gcloud domains verify DOMAIN_NAME
gcloud beta run domain-mappings create --service SERVICE_NAME --domain DOMAIN_NAME
```

Create Global IP
```
gcloud compute addresses create SERVICE_IP_NAME --network-tier=PREMIUM --ip-version=IPV4 --global
gcloud compute addresses describe SERVICE_IP_NAME --format="get(address)" --global
```
-> SERVICE_IP

NEG
```
gcloud beta compute network-endpoint-groups create SERVICE_NAME-neg --region=us-central1 --network-endpoint-type=serverless --cloud-run-service=SERVICE_NAME
gcloud beta compute backend-services create SERVICE_NAME-backend --load-balancing-scheme=EXTERNAL_MANAGED --global
gcloud beta compute backend-services add-backend SERVICE_NAME-backend --global --network-endpoint-group=SERVICE_NAME-neg --network-endpoint-group-region=us-central1
gcloud beta compute url-maps create SERVICE_NAME-url --default-service SERVICE_NAME-backend
gcloud beta compute ssl-certificates create SERVICE_NAME-ssl --domains DOMAIN_NAME
gcloud beta compute target-https-proxies create SERVICE_NAME-proxy --ssl-certificates=SERVICE_NAME-ssl --url-map=SERVICE_NAME-url
gcloud beta compute forwarding-rules create SERVICE_NAME-forward --load-balancing-scheme=EXTERNAL_MANAGED --network-tier=PREMIUM --address=SERVICE_IP --target-https-proxy=SERVICE_NAME-proxy --global --ports=443
gcloud compute target-https-proxies update SERVICE_NAME-proxy --ssl-certificates SERVICE_NAME-ssl --global-ssl-certificates --global
gcloud compute target-https-proxies describe SERVICE_NAME-proxy --global --format="get(sslCertificates)"
```

DNS
```
gcloud dns managed-zones create DNS_ZONE --description="A description of service" --dns-name=DNS_NAME(can be domain name)  --visibility=public
gcloud dns record-sets transaction start --zone=DNS_ZONE
gcloud dns record-sets transaction add SERVICE_IP --name=DNS_NAME --ttl=300 --type=A --zone=DNS_NAME
gcloud dns record-sets transaction execute --zone=DNS_ZONE
gcloud beta compute network-endpoint-groups create DNS_MASK --region=us-central1 --network-endpoint-type=serverless --cloud-run-url-mask="DNS_NAME/<service>"
```

These are useful for testing but the idea is to automate this flow so a customers click can purchase domain and deploy cloud run service.\
Shopify requires manual entry of their IP addresses on your domain name provider to begin routing traffic to their service for your shop, in that case, you manage the domain. In this case domain would be managed by us, the limitation of this approach is that per project domain quotas on GCP. Multiple projects can be set up to manage customer segments, however correspondence with the GCP team has advised me this quota can be extended up to 100, which is better, but not ideal. I'm sure if it gained traction they could accomodate further quota increases.

For development I used json files to populate domain cache. Will switch this to SQL. The reason we need this is because you must supply price and currency when registering domain and this fluctuates. When registering a domain fetched fomr the cache therefore we must fetch again its pricing info. The cache should refresh its entry if the pricing is likely to be out of date... 1 month? 2 month?