-- Backend Roadmap 的最小示例库表，供 5 门语言的 advanced/B01、C01、C02、F01 共用。
-- 执行：psql "$BRM_PG_DSN" -f schema.sql

CREATE TABLE IF NOT EXISTS accounts (
    id          BIGSERIAL PRIMARY KEY,
    owner       TEXT        NOT NULL UNIQUE,
    balance     BIGINT      NOT NULL DEFAULT 0 CHECK (balance >= 0),
    version     BIGINT      NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS orders (
    id            BIGSERIAL PRIMARY KEY,
    order_no      TEXT        NOT NULL UNIQUE,
    account_id    BIGINT      NOT NULL REFERENCES accounts(id),
    amount        BIGINT      NOT NULL CHECK (amount > 0),
    status        TEXT        NOT NULL DEFAULT 'PENDING',
    idempotency_key TEXT      UNIQUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_orders_account_created
    ON orders (account_id, created_at DESC);

CREATE TABLE IF NOT EXISTS outbox (
    id          BIGSERIAL PRIMARY KEY,
    topic       TEXT        NOT NULL,
    payload     JSONB       NOT NULL,
    published   BOOLEAN     NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO accounts (owner, balance) VALUES ('alice', 1000), ('bob', 1000)
ON CONFLICT (owner) DO NOTHING;

-- 事务与隔离级别实验（开两个 psql 会话分别执行）
-- 会话 A: BEGIN ISOLATION LEVEL REPEATABLE READ; SELECT balance FROM accounts WHERE owner='alice';
-- 会话 B: UPDATE accounts SET balance = balance - 100 WHERE owner='alice';
-- 会话 A: SELECT balance FROM accounts WHERE owner='alice';  -- 观察是否看到新值
-- 会话 A: COMMIT;

-- 乐观锁扣款（对照 advanced/B01 与 B04）
-- UPDATE accounts SET balance = balance - 100, version = version + 1
--  WHERE owner = 'alice' AND version = :expected_version AND balance >= 100;
-- 检查 affected rows：为 0 说明发生并发冲突，需要重试或返回冲突。

-- 幂等写入（对照 advanced/B02）
-- INSERT INTO orders (order_no, account_id, amount, idempotency_key)
-- VALUES ('ORD-1', 1, 100, 'key-from-client')
-- ON CONFLICT (idempotency_key) DO UPDATE SET order_no = EXCLUDED.order_no
-- RETURNING id, status;

-- N+1 与批量查询对照（对照 advanced/C02）
-- 慢：SELECT * FROM orders WHERE account_id = 1; 然后循环查每个 account
-- 快：SELECT o.*, a.owner FROM orders o JOIN accounts a ON a.id = o.account_id WHERE o.account_id = 1;

-- 游标分页（对照 advanced/C03）
-- SELECT * FROM orders WHERE (created_at, id) < (:last_created_at, :last_id)
--  ORDER BY created_at DESC, id DESC LIMIT 20;
