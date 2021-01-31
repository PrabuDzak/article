FROM alpine:latest  
COPY ./_output/restapi .
CMD ["./restapi"]  
