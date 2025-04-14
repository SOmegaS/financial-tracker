package org.ivangalochkin.expensewriter.service;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.ivangalochkin.expensewriter.mapper.ProtoToJooqMapper;
import org.ivangalochkin.expensewriter.proto.CreateBillProto;
import org.jooq.DSLContext;
import org.jooq.generated.tables.records.BillsRecord;
import org.springframework.stereotype.Service;

import java.util.Arrays;
import java.util.List;

@Slf4j
@Service
@RequiredArgsConstructor
public class BillService {
    private final DSLContext dslContext;

    public void batchCreateBills(List<CreateBillProto.CreateBill> bills) {
        log.info("Creating bills...");
        List<BillsRecord> records = bills.parallelStream().map(ProtoToJooqMapper::map).toList();
        log.info(records.getFirst().toString());
        int[] res = dslContext.batchInsert(records).execute();
        log.info(Arrays.toString(res));
        log.info("Bills are created");
    }
}
