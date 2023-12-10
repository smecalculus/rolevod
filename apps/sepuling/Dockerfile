ARG JAVA_RELEASE
FROM eclipse-temurin:${JAVA_RELEASE}-jre-alpine
RUN addgroup -S bezmen && adduser -S bezmen -G bezmen
USER bezmen:bezmen
WORKDIR /home/bezmen
COPY libs libs
COPY sepuling*.jar app.jar
ENTRYPOINT [ \
    "java", "-cp", "app.jar:libs/*", \
    "-Dconfig.file=application.conf", \
    "smecalculus.bezmen.construction.App" \
]
