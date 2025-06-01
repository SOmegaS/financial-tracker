package org.ivangalochkin.expensewriter.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.ivangalochkin.expensewriter.mapper.ProtoToJooqMapper;
import org.ivangalochkin.expensewriter.proto.CreateBillProto;
import org.jooq.DSLContext;
import org.jooq.generated.tables.records.BillsRecord;
import org.springframework.stereotype.Service;
import org.ivangalochkin.expensewriter.metrics.MetricsConfig;

import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class BillService {
    private final DSLContext dslContext;
    private final MetricsConfig metrics;

    public void batchCreateBills(List<CreateBillProto.CreateBill> bills) {
        // измеряем латентность, успех и ошибку
        metrics.billProcessLatency.record(() -> {
            try {
                log.info("Creating bills batch of size {}", bills.size());
                List<BillsRecord> records = bills.parallelStream()
                        .map(ProtoToJooqMapper::map)
                        .toList();

                if (!records.isEmpty()) {
                    log.debug("First record preview: {}", records.get(0).toString());
                }
                int[] result = dslContext.batchInsert(records).execute();
                log.info("Inserted {} records", result.length);
                // метрика успешных батчей
                metrics.billsProcessedTotal.increment();
            } catch (Exception e) {
                // метрика неуспешных батчей
                metrics.billsFailedTotal.increment();
                log.error("Error in batchCreateBills", e);
                throw e;
            }
        });
    }
}
