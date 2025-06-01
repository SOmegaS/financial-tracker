package org.ivangalochkin.expensewriter.kafka.config;

import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.Parser;
import lombok.NoArgsConstructor;
import org.apache.kafka.common.errors.SerializationException;
import org.apache.kafka.common.serialization.Deserializer;
import org.ivangalochkin.expensewriter.proto.CreateBillProto;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;

@NoArgsConstructor(force = true)
public class KafkaCreateBillProtoDeserializer implements Deserializer<CreateBillProto.CreateBill> {
    private static final Logger log = LoggerFactory.getLogger(KafkaCreateBillProtoDeserializer.class);
    private final Parser<CreateBillProto.CreateBill> parser = CreateBillProto.CreateBill.parser();

    @Override
    public void configure(Map<String, ?> configs, boolean isKey) {
        Deserializer.super.configure(configs, isKey);
    }

    @Override
    public CreateBillProto.CreateBill deserialize(String topic, byte[] data) {
        try {
            return parser.parseFrom(data);
        } catch (InvalidProtocolBufferException e) {
            log.error("Failed to deserialize message from topic {}", topic, e);
            throw new SerializationException("Error deserializing CreateBillProto.CreateBill", e);
        }
    }

    @Override
    public void close() {
        Deserializer.super.close();
    }
}
