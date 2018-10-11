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
Therefor, it is important to ensure there are no contents that could expire
(e.g. url links in the template might change, one-time token might expire). If
the emails are permitting access (such as password recovery), ensure that the
email is embedded with a token that could only be used once

## TODO

- add different delivery mechanism (mobile notification, browser notification, email, etc)
- create interface for the delivery and proper decoupling through design patterns
- track the type of notification sent and the open rate for delivery intelligence
- determine best sending time of the day
- multi-armed bandit to determine the best content
- add monitoring to display statistics
