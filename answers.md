## Why is it important to measure percentile latencies in production systems (e.g. p99)?

It is important to track percentile latency in production systems to track response times of systems. Whether it is a frontend user loading a web page or a backend service responding to a API call. We as engineers need to ensure that users/services are receiving responses in a reasonable time (within a services SLA). If we are able to track percentiles latencies, this allows us as engineers to ensure that code changes are not creating a more slow (poor) experience. We can use percentile latencies as a way to monitor systems to ensure that during peak usage our systems are responding withing our SLAs defined for our services (P50, P95, P99).

## Which metrics are important to track for queues? Why?

The below metrics are important to track for queues especially during peak times to help understand how a queue (and systems interacting it) are performing.

- Queue Latency (how long a message stays on the queue)
  - A high latency for messages can indicate that the consumers need to be scaled to take on the load of the queue.
- Queue size (the number of messages in the queue)
  - This is important to track to ensure that a queue does not reach its maximum capacity.
- Produce Rate/Consumer Rate
  - How quickly messages are being added and removed from the queue
- Retries (the number of times a message has been attempted to be processed)
  - A high retry number can often indicate errors in the consumer processing messages.
- Dead Letter Queue
  - It is important to monitor your dead letter queue so that you can track more specific issues with process certain messages. I.E. if a system tries to process a message more the 10 times, it may send it to the dead letter queue once it hits a certain code path to allow engineers to look at the message.
