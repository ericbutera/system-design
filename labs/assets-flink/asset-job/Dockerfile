FROM maven:3.8.6-openjdk-11-slim AS dependencies

WORKDIR /opt/flink-app
COPY pom.xml .

RUN mvn dependency:go-offline

FROM maven:3.8.6-openjdk-11-slim AS build

WORKDIR /opt/flink-app
COPY --from=dependencies /opt/flink-app /opt/flink-app
COPY src src

RUN mvn clean install -DskipTests -T 2C

FROM openjdk:11-jre-slim

WORKDIR /opt/flink-app
COPY --from=build /opt/flink-app/target/flink-asset-job-1.0-SNAPSHOT.jar app.jar

CMD ["java", "-jar", "app.jar"]