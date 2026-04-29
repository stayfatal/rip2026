ALTER TABLE sharding_strategies
ADD COLUMN IF NOT EXISTS short_description_en VARCHAR(100) NOT NULL DEFAULT '';

UPDATE sharding_strategies
SET short_description_en = CASE title
    WHEN 'Range Sharding' THEN 'Range shards group neighboring keys and speed up ordered scans for time-series data.'
    WHEN 'Hash Sharding' THEN 'Hash shards spread keys uniformly, balancing write load across partitions in OLTP systems.'
    WHEN 'Geo Sharding' THEN 'Geo shards place data near users, reducing latency and supporting regional compliance needs.'
    WHEN 'Directory-Based Sharding' THEN 'Directory shards map each tenant key to a shard for flexible routing and controlled movement.'
    WHEN 'Composite Sharding' THEN 'Composite shards combine rules like geo plus hash to scale globally with balanced traffic.'
    WHEN 'Dynamic Sharding' THEN 'Dynamic shards split and merge automatically as traffic changes to keep cluster usage stable.'
    ELSE 'Database sharding strategy profile.'
END
WHERE COALESCE(TRIM(short_description_en), '') = '';
