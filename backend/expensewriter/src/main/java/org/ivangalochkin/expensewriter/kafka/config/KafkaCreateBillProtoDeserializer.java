package org.ivangalochkin.expensewriter.kafka.config;

import com.google.protobuf.InvalidProtocolBufferException;
import com.google.protobuf.Parser;
import lombok.NoArgsConstructor;
import org.apache.kafka.common.errors.SerializationException;
import org.apache.kafka.common.serialization.Deserializer;
import org.ivangalochkin.expensewriter.proto.CreateBillProto;

import java.util.Map;

@NoArgsConstructor(force = true)
public class KafkaCreateBillProtoDeserializer implements Deserializer<CreateBillProto.CreateBill> {

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
            throw new SerializationException();
        }
    }
}