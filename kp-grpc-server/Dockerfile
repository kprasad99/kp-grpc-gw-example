FROM gradle:8.4-jdk17-jammy as builder

COPY . /app/

WORKDIR /app

RUN gradle build \
    && java -Djarmode=layertools \
    -jar build/libs/kp-grpc-server-0.0.1-SNAPSHOT.jar extract \
    --destination /output

FROM amazoncorretto/openjdk:21-alpine

ARG DEPENDENCY=/output
WORKDIR /app

ENV JAVA_TOOLS_OPTS=$JAVA_TOOLS_OPTS

COPY --from=builder ${DEPENDENCY}/dependencies/ ./
COPY --from=builder ${DEPENDENCY}/spring-boot-loader/ ./
COPY --from=builder ${DEPENDENCY}/snapshot-dependencies/ ./
COPY --from=builder ${DEPENDENCY}/application/ ./

VOLUME /tmp

ENTRYPOINT ["java", "org.springframework.boot.loader.JarLauncher"]