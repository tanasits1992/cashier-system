# Introduction

I develop cashier system (REST backend) with go language.
As the database is not ready, I will provide with a fake database.
When you start the application, it will initialize data of `voucher` from `voucher.json`
and `cash-store` from `cash-store.json`.

## Desgin
1. `./docs/Exchange-Logic.jpg` - Excahnge logic (Discount by voucher, E-wallet payment, 
Credit card payment, Cash payment, Exchange, Save sale transaction, Notification).
2. `./docs/Payment-Flow.jpg` - Payment flow (use case)
2. `./docs/Clear-Store-Flow.jpg` - Clear store flow (use case)


### API Specification
- `./docs/cashier.postman_collection.json`

