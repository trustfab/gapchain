---
name: ginkgo
description: Viết, tổ chức và chạy test cho project Go bằng Ginkgo v2 và Gomega. Dùng skill này khi người dùng muốn: thêm test vào package Go, viết unit test / integration test / table-driven test, bootstrap Ginkgo suite cho package mới, refactor test Go sang style BDD, test async/HTTP handler/service layer, chạy test song song, hoặc debug flaky test. Trigger khi thấy: "viết test cho Go", "ginkgo", "gomega", "BDD test", "test suite golang", "bootstrap test", "It / Describe / BeforeEach", "test handler", "test service golang".
---

# Ginkgo — Go Testing Framework

Skill này hướng dẫn cách dùng **Ginkgo v2** + **Gomega** để viết test cho project Go. Áp dụng cho mọi loại test: unit, integration, HTTP handler, service layer, async code.

Tài liệu gốc: https://onsi.github.io/ginkgo/  
Gomega matchers: https://onsi.github.io/gomega/

---

## 1. Cài đặt

```bash
# Cài ginkgo CLI và gomega vào project
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go get github.com/onsi/gomega/...
```

Kiểm tra:
```bash
ginkgo version
```

> Đảm bảo version CLI khớp với version trong `go.mod`. Nếu không khớp, chạy lại `go install` từ thư mục project.

---

## 2. Bootstrap suite cho package mới

Mỗi package cần **một** file bootstrap duy nhất:

```bash
cd path/to/your-package
ginkgo bootstrap          # sinh ra yourpackage_suite_test.go
ginkgo generate <name>    # sinh ra name_test.go với Describe rỗng
```

File bootstrap sinh ra (`yourpackage_suite_test.go`):

```go
package yourpackage_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "testing"
)

func TestYourPackage(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "YourPackage Suite")
}
```

> Package kết thúc bằng `_test` là chuẩn — test được compile riêng, chỉ access public API của package gốc.

---

## 3. Cấu trúc node cốt lõi

Ginkgo dùng 3 loại node để xây dựng spec tree:

### Container nodes — tổ chức theo nhóm

```go
Describe("tên tính năng", func() {      // mô tả cái gì
    Context("khi điều kiện X", func() { // context nào
        When("trạng thái Y", func() {   // alias của Context
            // specs ở đây
        })
    })
})
```

### Setup nodes — chuẩn bị / dọn dẹp

```go
BeforeEach(func() { /* chạy trước mỗi It */ })
AfterEach(func()  { /* chạy sau mỗi It  */ })

BeforeSuite(func() { /* chạy 1 lần trước toàn suite */ })
AfterSuite(func()  { /* chạy 1 lần sau  toàn suite */ })
```

### Subject nodes — viết assertion

```go
It("mô tả hành vi", func() {
    Expect(result).To(Equal(expected))
})

Specify("alias của It", func() { ... })
```

**Quy tắc vàng:** Khai báo biến trong container node, khởi tạo trong `BeforeEach`.

```go
var _ = Describe("UserService", func() {
    var (                          // ← khai báo ở đây
        svc  *UserService
        repo *MockUserRepo
    )

    BeforeEach(func() {           // ← khởi tạo ở đây
        repo = NewMockUserRepo()
        svc  = NewUserService(repo)
    })

    It("trả về user khi tìm thấy", func() {
        repo.AddUser(&User{ID: 1, Name: "Alice"})
        user, err := svc.GetByID(1)
        Expect(err).NotTo(HaveOccurred())
        Expect(user.Name).To(Equal("Alice"))
    })
})
```

---

## 4. Ví dụ thực tế theo từng use case

### 4a. Unit test service đơn giản

```go
var _ = Describe("OrderService", func() {
    var (
        svc   *OrderService
        store *MockOrderStore
    )

    BeforeEach(func() {
        store = &MockOrderStore{}
        svc   = NewOrderService(store)
    })

    Describe("CreateOrder", func() {
        Context("khi input hợp lệ", func() {
            It("lưu order và trả về ID", func() {
                order, err := svc.CreateOrder(CreateOrderInput{
                    UserID:   42,
                    Amount:   150_000,
                    Currency: "VND",
                })
                Expect(err).NotTo(HaveOccurred())
                Expect(order.ID).NotTo(BeZero())
                Expect(store.SaveCalled).To(BeTrue())
            })
        })

        Context("khi amount âm", func() {
            It("trả về lỗi validation", func() {
                _, err := svc.CreateOrder(CreateOrderInput{Amount: -1})
                Expect(err).To(MatchError(ErrInvalidAmount))
            })
        })
    })
})
```

### 4b. Test HTTP handler

```go
var _ = Describe("GET /users/:id", func() {
    var (
        router *gin.Engine  // hoặc http.Handler tùy framework
        w      *httptest.ResponseRecorder
    )

    BeforeEach(func() {
        router = setupRouter()
        w      = httptest.NewRecorder()
    })

    Context("user tồn tại", func() {
        It("trả về 200 và user JSON", func() {
            req, _ := http.NewRequest("GET", "/users/1", nil)
            router.ServeHTTP(w, req)

            Expect(w.Code).To(Equal(http.StatusOK))
            Expect(w.Body.String()).To(ContainSubstring(`"id":1`))
        })
    })

    Context("user không tồn tại", func() {
        It("trả về 404", func() {
            req, _ := http.NewRequest("GET", "/users/9999", nil)
            router.ServeHTTP(w, req)

            Expect(w.Code).To(Equal(http.StatusNotFound))
        })
    })
})
```

### 4c. Table-driven test với DescribeTable

Dùng khi cùng logic nhưng nhiều bộ input/output:

```go
var _ = Describe("Validator", func() {
    DescribeTable("ValidateEmail",
        func(email string, expectValid bool) {
            err := ValidateEmail(email)
            if expectValid {
                Expect(err).NotTo(HaveOccurred())
            } else {
                Expect(err).To(HaveOccurred())
            }
        },
        Entry("email hợp lệ",        "user@example.com",  true),
        Entry("thiếu @",             "userexample.com",   false),
        Entry("thiếu domain",        "user@",             false),
        Entry("có khoảng trắng",     "user @example.com", false),
        Entry("subdomain hợp lệ",    "a@b.co.vn",         true),
    )
})
```

### 4d. Test async / goroutine

```go
It("xử lý message trong vòng 2 giây", func(ctx SpecContext) {
    ch := make(chan string, 1)
    go processAsync(ch)

    Eventually(ctx, ch).WithTimeout(2 * time.Second).Should(Receive(Equal("done")))
}, SpecTimeout(3*time.Second))
```

Hoặc dùng `Eventually` / `Consistently` trực tiếp:

```go
It("cache được populate sau khi gọi", func() {
    go warmupCache()
    Eventually(func() bool {
        return cache.Has("key")
    }).WithTimeout(5 * time.Second).WithPolling(100 * time.Millisecond).Should(BeTrue())
})
```

### 4e. Integration test với database (BeforeSuite / AfterSuite)

```go
// suite_test.go
var db *sql.DB

var _ = BeforeSuite(func() {
    var err error
    db, err = sql.Open("postgres", testDSN())
    Expect(err).NotTo(HaveOccurred())
    runMigrations(db)
})

var _ = AfterSuite(func() {
    db.Close()
})

// order_repo_test.go
var _ = Describe("OrderRepository", func() {
    BeforeEach(func() {
        cleanDB(db) // truncate tables
    })

    It("insert và query thành công", func() {
        repo := NewOrderRepository(db)
        err := repo.Insert(&Order{UserID: 1, Amount: 100})
        Expect(err).NotTo(HaveOccurred())

        orders, err := repo.FindByUserID(1)
        Expect(err).NotTo(HaveOccurred())
        Expect(orders).To(HaveLen(1))
    })
})
```

---

## 5. Gomega matchers thường dùng

```go
// Equality
Expect(x).To(Equal(y))
Expect(x).To(BeEquivalentTo(y))   // so sánh sau khi convert type

// Nil / Zero
Expect(err).NotTo(HaveOccurred())
Expect(ptr).NotTo(BeNil())
Expect(val).To(BeZero())          // zero value của kiểu

// Boolean
Expect(ok).To(BeTrue())
Expect(disabled).To(BeFalse())

// Số
Expect(n).To(BeNumerically(">", 0))
Expect(n).To(BeNumerically("~", 3.14, 0.01)) // xấp xỉ

// String
Expect(s).To(ContainSubstring("hello"))
Expect(s).To(HavePrefix("GET"))
Expect(s).To(HaveSuffix(".json"))
Expect(s).To(MatchRegexp(`^\d{4}-\d{2}-\d{2}$`))

// Collection
Expect(slice).To(HaveLen(3))
Expect(slice).To(ContainElement("foo"))
Expect(slice).To(ConsistOf("a", "b", "c"))   // không quan tâm thứ tự
Expect(m).To(HaveKey("id"))
Expect(m).To(HaveKeyWithValue("status", "ok"))

// Error
Expect(err).To(MatchError("expected error message"))
Expect(err).To(MatchError(ErrNotFound))

// Struct / JSON
Expect(obj).To(MatchFields(IgnoreExtras, Fields{
    "Name":  Equal("Alice"),
    "Email": ContainSubstring("@"),
}))
```

---

## 6. Focus và Skip

```go
// Chỉ chạy spec/block này (debug tạm thời)
FIt("...", func() { ... })
FDescribe("...", func() { ... })
FContext("...", func() { ... })

// Bỏ qua
XIt("...", func() { ... })
XDescribe("...", func() { ... })

// Pending (chưa implement)
It("todo: xử lý edge case", Pending)
```

> **Quan trọng:** Không commit code có `F` prefix — toàn bộ suite chỉ chạy spec được focus.

---

## 7. Chạy test

```bash
# Chạy toàn bộ suite ở package hiện tại
ginkgo

# Verbose: hiện tên từng spec
ginkgo -v

# Chạy toàn bộ project (đệ quy)
ginkgo -r

# Chạy song song (N workers)
ginkgo -p
ginkgo --procs=4

# Chỉ chạy spec khớp tên
ginkgo --focus="CreateOrder"

# Bỏ qua spec khớp tên
ginkgo --skip="integration"

# Chạy lại flaky test nhiều lần để phát hiện lỗi ngẫu nhiên
ginkgo --repeat=5

# Seed cố định để reproduce flaky test
ginkgo --seed=12345

# Fail ngay khi gặp lỗi đầu tiên
ginkgo --fail-fast

# Tích hợp CI
ginkgo -r --randomize-all --fail-on-pending -cover
```

---

## 8. Decorator hữu ích

```go
// Timeout cho 1 spec
It("gọi API chậm", func(ctx SpecContext) {
    ...
}, SpecTimeout(10*time.Second))

// Retry tự động (cho flaky test)
It("kết nối mạng có thể thất bại", func() {
    ...
}, FlakeAttempts(3))

// Label để filter khi chạy
It("test database", Label("integration", "slow"), func() {
    ...
})
// Chạy: ginkgo --label-filter="integration"

// Ordered — đảm bảo thứ tự trong block
Describe("setup trình tự", Ordered, func() {
    It("bước 1", func() { ... })
    It("bước 2", func() { ... }) // luôn chạy sau bước 1
})
```

---

## 9. Cấu trúc file thực tế cho project

```
your-project/
├── internal/
│   ├── user/
│   │   ├── user.go
│   │   ├── user_suite_test.go    ← bootstrap (ginkgo bootstrap)
│   │   ├── user_service_test.go  ← (ginkgo generate user_service)
│   │   └── user_repo_test.go
│   └── order/
│       ├── order.go
│       ├── order_suite_test.go
│       └── order_handler_test.go
└── test/
    └── integration/              ← integration tests riêng
        ├── suite_test.go
        └── e2e_test.go
```

---

## 10. Checklist khi viết spec mới

- [ ] Đã có `*_suite_test.go` trong package chưa? Nếu chưa, chạy `ginkgo bootstrap`.
- [ ] Biến dùng chung: khai báo ở container node, khởi tạo ở `BeforeEach`.
- [ ] Mỗi `It` độc lập — không phụ thuộc thứ tự chạy.
- [ ] Dùng `AfterEach` để cleanup (close file, rollback DB, reset mock).
- [ ] Không commit `FIt`, `FDescribe`, `FContext`.
- [ ] Spec có tên mô tả rõ: `It("trả về lỗi khi token hết hạn")` thay vì `It("test auth")`.
- [ ] Async code dùng `Eventually` / `SpecContext` thay vì `time.Sleep`.
