<div align="center">
  <h1>Zero One News API</h1>
</div>

## About

Zero One Group News API made using Go, Echo, PostgreSQL, and documented by Swagger.

## Documentation

For generated Swagger JSON and additional documentations, see `/docs` directory.

<details>
  <summary>Entity Relationship Diagram</summary>

  ![ERD](docs/erd.png?raw=true)

</details>

<details>
  <summary>Swagger</summary>

  ![Screenshot](docs/swagger.png?raw=true)

</details>

<details>
  <summary>Endpoints</summary>

  ```mermaid
flowchart TD
    Root["/api"] --> Version["/v1"]
    Version --> Articles["/articles"] & Topics["/topics"]
    Topics --> TopicsGet["GET"] & TopicsPost["POST"] & TopicId["/:topic_id"]
    TopicId --> TopicIdGet["GET"] & TopicIdPatch["PUT"] & TopicIdDelete["DELETE"] & TopicArticles["/articles"]
    TopicArticles --> TopicArticlesGet["GET"]
    Articles --> ArticlesGet["GET"] & ArticlesPost["POST"] & ArticleId["/:article_id"]
    ArticleId --> ArticleIdGet["GET"] & ArticleIdPatch["PUT"] & ArticleIdDelete["DELETE"] & ArticleTopics["/topics"]
    ArticleTopics --> ArticleTopicsGet["GET"] & ArticleTopicsPut["PUT"] & ArticleTopicsId["/:topic_id"]
    ArticleTopicsId --> ArticleTopicsIdPost["POST"] & ArticleTopicsIdDelete["DELETE"]

     Root:::Sky
     Root:::Aqua
     Version:::Sky
     Version:::Aqua
     Articles:::Aqua
     Articles:::Sky
     Topics:::Aqua
     Topics:::Sky
     TopicId:::Rose
     TopicArticles:::Peach
     ArticleId:::Rose
     ArticleTopics:::Peach
     ArticleTopicsId:::Ash
    classDef Aqua stroke-width:1px, stroke-dasharray:none, stroke:#46EDC8, fill:#DEFFF8, color:#378E7A
    classDef Sky stroke-width:1px, stroke-dasharray:none, stroke:#374D7C, fill:#E2EBFF, color:#374D7C
    classDef Rose stroke-width:1px, stroke-dasharray:none, stroke:#FF5978, fill:#FFDFE5, color:#8E2236
    classDef Peach stroke-width:1px, stroke-dasharray:none, stroke:#FBB35A, fill:#FFEFDB, color:#8F632D
    classDef Ash stroke-width:1px, stroke-dasharray:none, stroke:#999999, fill:#EEEEEE, color:#000000
  ```

</details>

## Prerequisites

- PostgreSQL
- Docker
- Moon
- Go
- pnpm

## Development

1. Clone the repository

```sh
git clone https://github.com/sglkc/01-monorepo.git
cd 01-monorepo
```

2. Install dependencies with pnpm

```sh
pnpm install
```

5. Copy and set environment variables

```sh
cp .env.example .env
```

3. Run local development server

```sh
pnpm compose:up       # Start local development server
pnpm compose:down     # Stop local development server
pnpm compose:cleanup  # Remove all local development server data
```

4. Install Go dependencies

```sh
moon run tidy
```

5. Migrate database

```sh
moon zog-news:migration-up
```

6. Run on development mode

```sh
moon zog-news:dev
```

7. Open documentation endpoint at `/swagger/index.html`

## Testing

```sh
moon zog-news:test
```

## Building

```sh
moon zog-news:build
```
