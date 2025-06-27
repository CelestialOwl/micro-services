FROM alpine:latest

RUN mkdir /app

# COPY --from=builder /app/brokerApp /app
COPY loggerServiceApp /app

CMD [ "/app/loggerServiceApp" ]