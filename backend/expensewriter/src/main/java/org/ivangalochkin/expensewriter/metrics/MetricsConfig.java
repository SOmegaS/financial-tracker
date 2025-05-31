package org.ivangalochkin.expensewriter.metrics;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.Timer;
import io.micrometer.core.instrument.MeterRegistry;
import org.springframework.stereotype.Component;

@Component
public class MetricsConfig {
    public final Counter billsProcessedTotal;
    public final Counter billsFailedTotal;
    public final Timer billProcessLatency;

    public MetricsConfig(MeterRegistry registry) {
        billsProcessedTotal = Counter.builder("expense_writer_bills_processed_total")
            .description("Total number of bills successfully processed from Kafka")
            .register(registry);

        billsFailedTotal = Counter.builder("expense_writer_bills_failed_total")
            .description("Total number of bills processing failures")
            .register(registry);

        billProcessLatency = Timer.builder("expense_writer_bill_process_latency_seconds")
            .description("Latency of batchCreateBills method in seconds")
            .register(registry);
    }
}
