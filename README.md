# Erajaya BE Tech Test

## Project Overview

This backend service is built with **Go** and **Gin** web framework. The system is designed to be clean, easy to understand, scale, and maintain. It serves a simple API to manage `products` complete with user authentication.


## How it Works

The project follows a **layered structure**, which means each part of the system has a clear job:

```bash
root
├── configs/              # Loads configs from env
├── databases/            # DB migration logic
├── library/              # Shared helpers (transactions, responses, etc.)
├── middleware/           # Auth middleware
├── models/               # Data models
├── src/
│   └── app/
│       └── admin/
│           ├── product/  # Product-related code
│           └── user/     # User-related code
├── services/             # Business logic
├── routes/               # All route definitions
└── main.go/              # Starts everything
```

Here’s the general flow:

1. **`main.go`** starts the server, loads config, sets up DB, and runs migrations.
2. **Routes** are registered under `/admin/v1`, split into different modules, such as `product` and `user`.
3. **Handlers** are responsible for incoming API calls.
4. **Usecases** handle the business logic.
5. **Repositories** talk to the database.
6. **`dataManager`** is used to wrap DB operations in transactions, safely.


## Design Choices and Why

The architecture it self is a combination of **Clean Architechture** and some **Domain-Driven Design (DDD)** concepts. The **Clean Architechture** is used as a base for the whole framework design. You can see that by looking at the structure in the `src` folder. Located inside are 3 main folders, which are `app`, `routes`, and `service`. Here are their respective explanation:

- `routes`: Routes incoming requests into the correct API.
- `app`: Contains the controllers/handlers for each module. Handles API calls
- `services`: Contains the usecase & the repository, in which the core logic and DB queries are stored.

This architechture is chosen and used due to its simplicity of use, scale, and maintain. Everything is where it's supposed to be. Everything has its own place and shall not be anywhere else. This code structure helps lower learning curve and thus saves development time.


## API Highlights

### User APIs

- `GET /users` — List all users
- `POST /users` — Add a new user
- `PUT /users/:id` — Update a user
- `PUT /users/:id/password` — Change password
- `POST /users/auth/login` — Login endpoint
- And a few more…

### Product APIs

- `GET /products` — List all products
- `POST /products` — Create a product
- `PUT /products/:id` — Update a product
- `PUT /products/:id/status` — Change product status (works as a soft delete)


## Running It Locally

1. Set your `.env` or config:

   ```
   {
    "SERVER_NAME": "ERAJAYA BE TECH TEST - FI",

    "DB_CONNECTION_STRING":"<db_username>:<db_pass>@(<db_host>)/erajaya_be_tech_test?parseTime=true",

    "PORT_APPS": ":9034",
    "APP_URL": "http://localhost:9034",

    "ANDROID_APP_MINIMUM_VERSION": "1.0.0",
    "IOS_APP_MINIMUM_VERSION": "1.0.0",

    "EXTERNAL_URL": "",
    "EXTERNAL_TOKEN": "",
    "EXTERNAL_ACCESS_TOKEN": "",

    "REDIS_ADDR": "localhost:6379",
    "REDIS_TIME_OUT": "259200",
    "REDIS_DB": "0",
    "REDIS_PASSWORD": "",

    "SEND_WHATSAPP_API": "",
    "SEND_WHATSAPP_TOKEN": "",

    "TELE_BOT_TOKEN": "",
    "TELE_GROUP_ID": "",

    "FIREBASE_SERVER_KEY":"",
    "FIREBASE_SENDER_ID":"",
    "FIREBASE_BUCKET_URL":"",
    "FIREBASE_AUTH_FILE_PATH":"",

    "WHITELISTED_IPS": "0.0.0.0"
   }
   ```

2. Run `seeder.sql` in you db of choice

3. Run the app:
   ```bash
   go run main.go
   ```
4. Run `permission_seeder.sql` in the db you are using.

5. You're all set and ready to go


## Notes

- Passwords are hashed using MD5 right now. Works for the test, but it’s better to use bcrypt in real projects.
- Auth is done via middleware and applied to protected routes.
- You can add more modules easily by copying the pattern from `user` or `product`.
