# Project Notifications

## **Run the frontend locally, follow these steps:**

```
Clone the repository: git clone https://github.com/ragFurlan/notification.git
Install dependencies: go mod download
Run: go run ./cmd/main.go
This service run on : http://localhost:8080
```

## **Create mock**

### **Platform**
```
~/go/bin/mockgen -source=internal/platform/repositories/log.go -destination=test/platform/log.go -package=log
```

### **Usecase**
```
~/go/bin/mockgen -source=internal/usecase/notification/notificationUsecase.go -destination=test/usecase/notificationUsecase.go -package=log
```

## **Run tests**
```
go test ./... -gcflags=-l
```

## **Design Patterns Used**

### ***Factory Method Pattern***

The `NewNotificationUseCase()` method in the `NotificationUseCase` struct is an example of the Factory Method pattern. It encapsulates the creation of related objects, such as repositories, services, and other dependencies, by providing a consistent way to create instances of `NotificationUseCase` with their dependencies properly initialized.

### ***Repository Pattern***

 The `LogRepository` struct represents the Repository pattern. It encapsulates the logic for accessing data related to notification logs, providing methods to save and get logs. This promotes separation between business logic and the data persistence layer.


### ***Strategy Patterns***

The use of the `Notifier` interface and separate implementation for each type of notification (such as `SMSUseCase`, `EmailUseCase`, and `PushUseCase`) represents the Strategy pattern. This allows different notification algorithms to be encapsulated in separate classes, making it easy to add new types of notifications in the future.










