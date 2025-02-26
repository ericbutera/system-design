package com.example;

import org.apache.flink.api.common.eventtime.WatermarkStrategy;
import org.apache.flink.api.common.functions.MapFunction;
import org.apache.flink.api.java.tuple.Tuple2;
import org.apache.flink.streaming.api.datastream.DataStream;
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment;
import org.apache.flink.connector.kafka.source.KafkaSource;
import org.apache.flink.connector.kafka.source.enumerator.initializer.OffsetsInitializer;
import org.apache.flink.api.common.serialization.SimpleStringSchema;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.util.Properties;

public class AssetFlinkJob {
    public static void main(String[] args) throws Exception {
        final StreamExecutionEnvironment env = StreamExecutionEnvironment.getExecutionEnvironment();

        String broker = System.getenv("BROKER");
        String group = System.getenv("GROUP");
        String topic = System.getenv("TOPIC");
        if (broker == null || group == null || topic == null) {
            throw new IllegalArgumentException(
                    "Environment variables BROKER, GROUP, and TOPIC must be set");
        }

        Properties properties = new Properties();
        properties.setProperty("bootstrap.servers", broker);
        properties.setProperty("group.id", group);

        KafkaSource<String> source = KafkaSource.<String>builder()
                .setBootstrapServers(broker)
                .setTopics(topic)
                .setGroupId(group)
                .setStartingOffsets(OffsetsInitializer.earliest())
                .setValueOnlyDeserializer(new SimpleStringSchema())
                .build();
        DataStream<String> assetStream = env.fromSource(source, WatermarkStrategy.noWatermarks(), "Kafka Source");

        ObjectMapper objectMapper = new ObjectMapper();

        DataStream<Tuple2<String, String>> processedAssets = assetStream
                .map(new MapFunction<String, Tuple2<String, String>>() {
                    @Override
                    public Tuple2<String, String> map(String value) throws Exception {
                        JsonNode jsonNode = objectMapper.readTree(value);
                        String assetId = jsonNode.has("asset_id") ? jsonNode.get("asset_id").asText() : "unknown";
                        return new Tuple2<>(assetId, value);
                    }
                });

        processedAssets.print();

        env.execute("Flink Asset Job");
    }
}
