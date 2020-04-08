# Konnek GCP
Transform GcP events into CloudEvents – and send them somewhere.

> **This is a proof of concept, so everything – from code to instructions – are not production ready. If the idea is valid, I'll bring it to an alpha version soon.**

## Idea
Konnek is a PoC trying to encapsulate cloud provider events into the [CloudEvents](https://cloudevents.io/) specification and forward them through the CloudEvents HTTP protocol-binding. The idea is to receive the events inside the cloud provider FaaS platform, parse, format and them send off.

The original idea was to feed these events into the Knative Eventing system, see [here](https://github.com/jonatasbaldin/konnek-event-receiver) for more info.

This repository contains the Google Cloud Functions implementation.

## Things to be aware
- Google Cloud Functions only allows _one trigger_ from _one resource_ per function. If you want to receive an event from a Storage Bucket and Pub Sub, you'll need two konnek deployments. Even if you want to receive events from two different buckets, you'll need two konnek deployments. _Even_ if you want to receive an event when a file is created _or_ deleted in the _same_ Storage Bucket, you'll need multiple konnek deployments.

- You can't change the _function trigger type_ of an already deployed function.

## Installing

### Setting up a local receiver 
Before we deploy the Lambda function, let's setup a place for these events to arrive. You'll need [ngrok](https://ngrok.com/).
```bash
docker run --rm -p 8080:8080 jonatasbaldin/konnek-knative-consumer

# Open another terminal and:
ngrok http 8080
```

Let ngrok run in this terminal. Note down your Ngrok address `https://xxxxxxx.ngrok.io`.

### Deploying the Google Cloud Function
```bash
# Clone this repository
git clone git@github.com:jonatasbaldin/konnek-gcp

# Deploy a function triggered by a Bucket
gcloud functions deploy konnek --runtime go111 --entry-point Handler --trigger-bucket <bucket-name> --set-env-vars KONNEK_CE_CONSUMER=<your-ngrok-address>


# OR, receiving an event from Pub/Sub
gcloud functions deploy konnek --runtime go111 --entry-point Handler --trigger-topic <topic-name> --set-env-vars KONNEK_CE_CONSUMER=<your-ngrok-address>
```

Once deployed, test it by uploading a file to the Bucket or sending a message into the Pub/Sub Topic.

## Events Tested
Events that were tested in this PoC (list from `gcloud functions event-types list`):

```
   EVENT_PROVIDER                   EVENT_TYPE                                                
✓  cloud.pubsub                     google.pubsub.topic.publish                               
✓  cloud.pubsub                     providers/cloud.pubsub/eventTypes/topic.publish           
✓  cloud.storage                    google.storage.object.archive                             
✓  cloud.storage                    google.storage.object.delete                              
✓  cloud.storage                    google.storage.object.finalize                            
✓  cloud.storage                    google.storage.object.metadataUpdate                      
✓  cloud.storage                    providers/cloud.storage/eventTypes/object.change          
   google.firebase.analytics.event  providers/google.firebase.analytics/eventTypes/event.log  
   google.firebase.database.ref     providers/google.firebase.database/eventTypes/ref.create  
   google.firebase.database.ref     providers/google.firebase.database/eventTypes/ref.delete  
   google.firebase.database.ref     providers/google.firebase.database/eventTypes/ref.update  
   google.firebase.database.ref     providers/google.firebase.database/eventTypes/ref.write   
   google.firestore.document        providers/cloud.firestore/eventTypes/document.create      
   google.firestore.document        providers/cloud.firestore/eventTypes/document.delete      
   google.firestore.document        providers/cloud.firestore/eventTypes/document.update      
   google.firestore.document        providers/cloud.firestore/eventTypes/document.write       
```