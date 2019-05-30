# go-delivery
Delivery engine created in go

## Design

The delivery engine should be able to generate different templates based on the
given data. The data can be personalized user data that is fetched from the
Content Service, and also basic user information (name, honorifics etc). 

There should be a policy to decide which templates can be used. It can be
specific:

- email/browser/sms/mobile notification only
- a combination of more than 1 template, at a desired time
- a combination of more than 1 template, randomly selected
- a combination of more than 1 template, personalized (based on user activity,
preferences etc)


It should be easy to perform the following:
- adding users
- adding templates

## Dealing with Failures

In any case, failures means not sending any notification at all to the users,
which is okay.

The notification delivery is not idempotent, so it is probably not suggested to
attempt a retry on failure. This is because it may lead to double-send, in
which a notification is send multiple times. Unless there is a guaranteed way
to ensure acknowledgement of delivery, it is hard to know if the notification
is sent successfully (user might not open the notification immediately, hence
delaying the acknowledgement). 

When there is failure, just deliver it the next running time.

## Pipeline

We can design a streaming pipeline that will chain different services together in order to "save" individual steps.

```
# v1
Initiator -> ContentBuilder -> Dispatcher -> End User -> Action -> Stats


# v2, after sufficient data
Model -> ContentBuilder -> Dispatcher -> End User -> Action -> Stats -> Model
```

- initiator: Can be a simple cron job that is responsible for querying the subscribers and the type of media (simple version) and publishing it to the next stage
- content builder: responsible for building the content through different templates for different media type. May call several external apis (recommendations, content) to build the desired content
- dispatcher: responsible for delivering the messages through the different mediums and marking them as read.
- end user: will receive the message. Their action (opening the email, reading the message) should create an event that will be sent back to the stats service. The stats service will create a personalized delivery based on what the user performed, browsing habits/open email, open mobile time/last login. 

The `v2` version will have the generated model as the first step, which will consistently update the model and creating the delivery trigger. The `v2` will be running alongside the `v1` version, except that the model will decide on which medium the message will be posted at, when it will be posted, what will be posted etc. If the `Subscriber` already exist in the `v2`, it should not be present in the `v1` to avoid duplicate delivery. To simplify, we can add a logic before content builder to filter sent messages. Since the sending time of the `Model` is dynamic, it will keep running in a loop to check if there are messages to be sent and execute it when there is. Whenever new data is added by the `Stats`, the `Model` will only send it in the next cycle. In this case, it is better to skip sending one to many then to have duplicates - user might not know if they did not receive the message too.

Each steps in between should introduce a reliable message queue with caching and retry mechanism. Also, successful delivery should be acknowledge and not repeated.

## Cron vs Real-Time

It is normal for a delivery notification engine to be send notification at
specific-time. But this could lead to spikes (sending 1M notifications at 9:00
A.M. daily would crash the server due to the high rate of delivery). It is
wrong to assume that the emails could be delivered at exactly 9:00 A.M. If we
are sending the messages through a queue, and it could only deliver messages at
a rate of 1,000 notifications/s, it would take 17 minutes before completing the
delivery, assuming there are no failures.

Also, it there are additional steps before sending the notification, such as
building the template, it would take more time for the messages to propagate
through the pipeline, since there might be additional external calls to build
it (calling Elasticsearch recommendation engine etc). Remember that the role of
the delivery engine is only to send the notification, and should not have
knowledge on how to build the template.


One way for generating the content is to do it offline - the day before the
delivery, assuming we do not need real-time data to be reflected in the
content. If we are sending notifications with immediately expiring content (an
offer that lasts only an hour), it might be hard to do it the day before the
delivery.


## Content Generation

The other issue with content generation is that it is persistent on the client
side (emails that were sent 1 month ago can still be accessed in the future).
Therefore, it is important to ensure there are no contents that could expire
(e.g. url links in the template might change, one-time token might expire). If
the emails are permitting access (such as password recovery), ensure that the
email is embedded with a token that could only be used once.

The `ContentBuilder` will be responsible for generating the content, which may call external APIs in the request/response style (event streaming is overkill?). The APIs can be static data from another service, or a recommendation service/elasticsearch search results etc. Once the content has been generated, we can first choose to persist them first before sending, so that the next retry doesn't require rebuilding the content if it is an expensive operation (matching engine). If the server restarts, we will also have the content stored somewhere as backup. Another way is to persist them in the message queue. The content should have an expiry date, to avoid poison pill. 

Also, we need to check for contacts (email, phone number) that are not valid and blacklist them automatically.


## Initiator

A more detailed section on the initiator. We can keep the Subscriber model simple:

- allow user to select what kind of notifications they want to receive - send only to the ones they subscribed too. If the user only has one (email), just send to that particular subscription. 
- if the user has two subscription preferences, toggle between this two, and see which has a better CTR. Or just use bandit algorithm to deliver to the optimum one.
- if the user did not set the subscription preferences/has multiple preferences, use bandit algorithm or create a profile for different devices/open time

How frequent do we need to send the notifications?

- we can just run it once daily. But as the number of users grow, processing the huge bulk of user is going to have a thundering herd effect on every other services this initiator is going to call.
- we can run a cron every 5 minutes, and take the subscriptions with the sending time in the time window (time +  time % 5 min interval) to deliver. This can reduce the load of the system, since the users will be spread out. If we did not have the subscription time, we can just modulo the user id by 24 hr (hash load balancing)

## Rules 

What are the delivery rules that we can consider?
- send to channels (mobile, browser notification, email). Different channels can have different impact (see below)
- frequency rule (1 email weekly, 1 push notification daily)
- elapsed rule (cannot send email 1 hr after sending push notification, or send the push notification 1 hr after the email)
- exclusion rule (either email or push notification)
- inclusion rule (send both email and push notification)
- exception rule (send to the channel user least used)
- last activity. If the user last login 3 months ago, try to re-engage the user by sending them an notification

## Channels

- email: Useful to re-engage users that are not active on our platform.
- push notification: Active users that are always on the go with the app would want to receive push notification more.
- sms: Normally for OTP etc
- chatbot: Useful if the user is always on the social media.

## Fraud

Detect bad emails and add a blacklist service. Send only to emails that are verified. Rate-limit to prevent unwanted mass registration of subscription. 

## Creating the user profile

Create two profile, one is the average global user activity (read time/open time), another is the specific user activity time). This can avoid the cold start problem - take the average global user activity if the user does not have a personalized profile. Create a time map (counter per hours/minutes of when the user open the email).

## TODO

- add different delivery mechanism (mobile notification, browser notification, email, etc)
- create interface for the delivery and proper decoupling through design patterns
- track the type of notification sent and the open rate for delivery intelligence
- determine best sending time of the day
- multi-armed bandit to determine the best content
- add monitoring to display statistics
- ensure one-time only delivery for the specific subscriber (subscriber id + time window) using bloom filter. Message can be retried if they fail, but they should not send the same message twice
- design the rule to ensure that the subscriber only receive the message with the given content on just one media (sms, email etc) at a time
- make the sending time dynamic (based on user preference, not a cron job running at a fixed time). We can first deliver the notifications at a fixed time, then gather metrics on the opening/read time of the message for average/specific user to determine the best sending time)
- personalized message for specific users. Design a rule engine to handle that.
