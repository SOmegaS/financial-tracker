package org.ivangalochkin.expensewriter.kafka.config;

import org.apache.kafka.clients.consumer.ConsumerConfig;
import org.apache.kafka.common.serialization.StringDeserializer;
import org.ivangalochkin.expensewriter.proto.CreateBillProto;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.kafka.config.ConcurrentKafkaListenerContainerFactory;
import org.springframework.kafka.core.ConsumerFactory;
import org.springframework.kafka.core.DefaultKafkaConsumerFactory;

import java.util.HashMap;
import java.util.Map;

@Configuration
public class KafkaCreateBillConsumerConfig {
    @Value("${app.kafka.server:kafka-moscow:9092}")
    private String kafkaServer;
    @Bean
    public ConsumerFactory<String, CreateBillProto.CreateBill> createBillConsumerFactory() {
        Map<String, Object> props = new HashMap<>();
        props.put(ConsumerConfig.BOOTSTRAP_SERVERS_CONFIG, kafkaServer);
        props.put(ConsumerConfig.GROUP_ID_CONFIG, "write-bill-consumer");
        props.put(ConsumerConfig.KEY_DESERIALIZER_CLASS_CONFIG, StringDeserializer.class);
        props.put(ConsumerConfig.VALUE_DESERIALIZER_CLASS_CONFIG, KafkaCreateBillProtoDeserializer.class);
        props.put(ConsumerConfig.MAX_POLL_RECORDS_CONFIG, "30");
        props.put(ConsumerConfig.AUTO_OFFSET_RESET_CONFIG, "earliest");

        return new DefaultKafkaConsumerFactory<>(props);
    }

    @Bean
    public ConcurrentKafkaListenerContainerFactory<String, CreateBillProto.CreateBill> createBillListenerContainerFactory() {
        ConcurrentKafkaListenerContainerFactory<String, CreateBillProto.CreateBill> factory =
                new ConcurrentKafkaListenerContainerFactory<>();
        factory.setConsumerFactory(createBillConsumerFactory());
        factory.setConcurrency(1);
        factory.setBatchListener(true);
        return factory;
    }

}
