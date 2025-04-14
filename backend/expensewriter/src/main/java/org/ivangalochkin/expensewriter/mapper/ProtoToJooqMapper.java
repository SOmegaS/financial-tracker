package org.ivangalochkin.expensewriter.mapper;

import org.ivangalochkin.expensewriter.proto.CreateBillProto;
import org.jooq.generated.tables.records.BillsRecord;

import java.time.Instant;
import java.time.ZoneId;
import java.util.UUID;

public class ProtoToJooqMapper {
    public static BillsRecord map(CreateBillProto.CreateBill createBill) {
        BillsRecord billsRecord = new BillsRecord();
        billsRecord.setAmount(createBill.getAmount());
        billsRecord.setCategory(createBill.getCategory());
        billsRecord.setTmstmp(Instant
                .ofEpochSecond(
                        createBill.getTimestamp().getSeconds(),
                        createBill.getTimestamp().getNanos()
                )
                .atZone(ZoneId.systemDefault()).toLocalDateTime());
        billsRecord.setName(createBill.getName());
        billsRecord.setUserId(UUID.fromString(createBill.getUserId()));
        return billsRecord;
    }
}
