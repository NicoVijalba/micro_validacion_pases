# ER Model

```mermaid
erDiagram
  RECORDS {
    BIGINT id PK
    VARCHAR external_id UK
    ENUM type
    DECIMAL amount
    VARCHAR description
    VARCHAR created_by
    TIMESTAMP created_at
  }
```
