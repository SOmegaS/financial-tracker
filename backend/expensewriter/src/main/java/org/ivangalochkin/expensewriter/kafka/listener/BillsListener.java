package org.ivangalochkin.expensewriter.kafka.listener;

import lombok.RequiredArgsConstructor;
import lombok.extern.slf4j.Slf4j;
import org.ivangalochkin.expensewriter.proto.CreateBillProto;
import org.ivangalochkin.expensewriter.service.BillService;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

import java.util.List;

@Slf4j
@Component
@RequiredArgsConstructor
public class BillsListener {

    private final BillService billService;

    @KafkaListener(
            topics = "write-bills",
            containerFactory = "createBillListenerContainerFactory"
    )
    public void consume(List<CreateBillProto.CreateBill> messages) {
        log.info("Received {} messages {}", messages.size(), Thread.currentThread().getName());
        billService.batchCreateBills(messages);
    }
}
