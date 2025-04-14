FROM openjdk:21-jdk
ADD . /src
WORKDIR /src
RUN ./mvnw package -DskipTests
EXPOSE 8080
ENTRYPOINT java ${MEMORY_LIMIT} -jar target/expensewriter-0.0.1-SNAPSHOT.jar