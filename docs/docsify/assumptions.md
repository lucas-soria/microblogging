# Assumptions

1. All users are valid, no need to create a signin module or handle sessions.
2. Think of a solution that can scale to millions of users (e.g., 100M users).
3. The application must be optimized for reads.

## Additional assumptions

### Functionalities

- No processing of mentions inside tweets.
- No processing of hashtags inside tweets.
- Users can't interact with tweets (like, retweet, reply).
- Users can't follow themselves.

### Data

#### Tweets

- text (280 characters) -> 280 bytes.
- user_id -> 20 bytes.
- timestamp -> 25 bytes.
- id -> 32 bytes.

> 280 + 20 + 25 + 32 = 357 bytes.

#### Users

- name (20 characters) -> 20 bytes.
- lastname (20 characters) -> 20 bytes.
- id -> 20 bytes.

> 20 + 20 + 20 = 60 bytes.

### Traffic

The system has arround 100M users and they tweet daily 10 tweets per user. That means that the system will receive 1B tweets per day, each tweet has a size of 357 bytes. 

- 100M users.
- 10 tweets per user per day.
- 357 bytes per tweet.

> 100M * 10 tweets/day * 357 bytes/tweet = 357GB/day.

Each user follows an average of 1000 users. If the user views their timeline, the system will need to fetch the tweets of the users they follow.

- 10 tweets per user per day.
- 1000 followers per user.

> 10 tweets/day * 1000 followers = 10k tweets.

If the timeline first loads the tweets from the last 24hs, that means a user will need to fetch 10k tweets.

> 10k tweets * 357 bytes/tweet = 3.57MB.

Assuming active users reprensent 60% of the total users and view their timeline 20 times per day.

> 100M * 0.6 * 3.57MB * 20 = 4.28PB/day.

This is why the system should be read optimized.
