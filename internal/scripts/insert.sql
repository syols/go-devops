INSERT INTO metrics
    (
        id, metric_type, counter_value, gauge_value, hash
    )
VALUES
    (
        :id, :metric_type, :counter_value, :gauge_value, :hash
    )
ON CONFLICT (id) DO
    UPDATE
    SET
        metric_type = excluded.metric_type,
        counter_value = excluded.counter_value,
        gauge_value = excluded.gauge_value,
        hash = excluded.hash;