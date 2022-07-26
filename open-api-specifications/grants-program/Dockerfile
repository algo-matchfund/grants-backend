FROM swaggerapi/swagger-ui

RUN mkdir /api
COPY grants-program/grants-program.yaml /api/
COPY grants-program/schemas /api/schemas
COPY grants-program/paths /api/paths
COPY grants-program/parameters /api/parameters
COPY grants-program/responses /api/responses

ENV SWAGGER_JSON=/api/grants-program.yaml

