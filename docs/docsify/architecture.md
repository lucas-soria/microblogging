# Architecture

This architecture is a microservices architecture that allows users to post tweets, follow users and view their own timelines.
It is based on the assumptions defined in the [intro.md](intro.md) file.

The system is has 4 main services:

- Feed Service: Handles the timeline of users.
- Tweets CRUD: Handles the creation, reading, updating and deletion of tweets.
- Users CRUD: Handles the creation, reading, updating and deletion of users.
- Analytics Service: Handles events, hosts analytics logic and updates caches based on user activity.

The system has 2 main databases:

- Tweets DB: Handles the tweets of users.
- Users DB: Handles the users of the system.

The system has 3 caches:

- Redis Timeline Cache: Handles the preloaded timelines of users.
- Redis Popular Cache: Handles the popular tweets of the system.
- Redis User Cache: Handles the users of the system.

The architecture looks like this:

```mermaid
flowchart TD
  %% Services
  subgraph Services
    FeedService[Feed Service]
    TweetsCRUD[Tweets CRUD]
    UsersCRUD[Users CRUD]
    Analytics[Analytics Service]
  end

  %% Databases
  subgraph Databases
    TweetsWrite[(Tweets DB - Write)]
    TweetsRead[(Tweets DB - Read Replica)]
    UsersDB[(Users DB)]
    EventsDB[(Events DB)]
  end

  %% Kafka Topics
  subgraph KafkaTopics
    K_TweetPosted([Kafka: TweetPosted])
    K_TimelineViewed([Kafka: TimelineViewed])
  end

  %% Redis Caches
  subgraph RedisCaches
    RedisTimelineCache[(Redis: Preloaded Timelines)]
    RedisPopularCache[(Redis: Popular Tweets)]
    RedisUserCache[(Redis: User Cache)]
  end

  %% Feed Service interactions
  FeedService -->|fetch tweets| TweetsCRUD
  FeedService -->|fetch timeline| RedisTimelineCache
  FeedService -->|fetch popular| RedisPopularCache
  FeedService -.->|publish TimelineViewed| K_TimelineViewed

  %% Kafka → Analytics
  K_TimelineViewed -.->|consume| Analytics
  K_TweetPosted -.->|consume| Analytics

  %% Tweets CRUD interactions
  TweetsCRUD -->|write tweet| TweetsWrite
  TweetsWrite -->|replicate| TweetsRead
  TweetsCRUD -->|read tweet| TweetsRead
  TweetsWrite -->|publish TweetPosted| K_TweetPosted

  %% Analytics interactions
  Analytics -->|update cache| RedisTimelineCache
  Analytics -->|update cache| RedisPopularCache
  Analytics -->|write event| EventsDB

  %% Users CRUD
  UsersCRUD -->|write user| UsersDB
  UsersCRUD -->|cache user| RedisUserCache

  %% Feed Service accessing user data
  FeedService -->|fetch users followed| UsersCRUD

```
