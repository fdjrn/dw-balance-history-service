# Balance History Service
---
### Deskripsi
Service ini berfungsi untuk meng-handle riwayat transaksi saldo pengguna, seperti:

    - Topup Saldo
    - Deduct Saldo
    - Transfer / Distribusi Saldo

### Service Type
    - Message Consumer
    - RestAPI endpoint

### Subscribed/Consumed Topic
    - mdw.transaction.topup.result                                  ✅
    - mdw.transaction.deduct.result                                 ✅
    - mdw.transaction.distribute.result                             ✅
    - mdw.transaction.distribute.result.members                     ✅

### RESTFul API Endpoint
    - POST /api/v1/account/transaction/history/all                  ✅
    - POST /api/v1/account/transaction/history/last-transaction     ✅
    - POST /api/v1/account/transaction/history/periods              ✅

### Build Docker Image
    docker build -t dw-history:1.0.0 -f Dockerfile .

### Available Environment Value:
    - DATABASE_MONGODB_URI : conncetion uri to mongodb cluster
        
        example: mongodb+srv://<user>:<password>@<cluster-host>/?retryWrites=true&w=majority

    - DATABASE_MONGODB_DB_NAME : Database Name used for parameter service

        example: mdw-balance-history

    - KAFKA_BROKERS : kafka cluster address

        example: touching-ghoul-8389-us1-kafka.upstash.io:9092

    - KAFKA_SASL_USER : kafka cluster username

    - KAFKA_SASL_PASSWORD : kafka cluster password

### Docker Run Command
    docker run -d -p 8010:8010 --name dw-history-service --env "DATABASE_MONGODB_DB_NAME=mdw-balance-history" --restart unless-stopped dw-history:1.0.0
    