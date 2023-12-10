ARG JAVA_RELEASE

FROM eclipse-temurin:${JAVA_RELEASE}-jre-alpine
RUN addgroup -S app && adduser -S app -G app
USER app:app
WORKDIR /home/app
COPY lib lib
COPY bin bin

ENTRYPOINT [ \
    "java", "-classpath", "bin/*:lib/*", \
    "-Dconfig.file=application.conf", \
    "smecalculus.bezmen.construction.App" \
]
