# 빠른시작

## 프리스마 세팅하기

1) Go 모듈을 이용해서 새 프로젝트 생성

    기존의 프로젝트를 사용하고 싶다면 건너뛰어도 됨
    ```shell script
    go mod init github.com/your/repo
    ```

2) Prisma client Go 패키지를 다운로드
    
    Prisma client Go 는 Prisma 에서 분리되었기 때문에 Prisma CLI를 별도로 설치하지 않아도 사용할 수 있음. 대신에 Go module 을 통해서 받을 수 있음.
    ```shell script
    go get github.com/prisma/prisma-client-go
    ```

3) `schema.prisma` 파일에서 [Prisma schema 준비하기](https://www.prisma.io/docs/reference/tools-and-interfaces/prisma-schema/prisma-schema-file). 예를 들어, sqlite 데이터베이스와 Prisma Client Go를 통해 두 가지 모델을 생성하는 경우에는 다음과 같이 보일 수 있음

    ```prisma
    datasource db {
        provider = "sqlite"
        url      = "file:dev.db"
    }

    generator db {
        provider = "go run github.com/prisma/prisma-client-go"
        output = "./db/db_gen.go"
        package = "db"
    }

    model User {
        id        String   @default(cuid()) @id
        createdAt DateTime @default(now())
        email     String   @unique
        name      String?
        age       Int?

        posts     Post[]
    }

    model Post {
        id        String   @default(cuid()) @id
        createdAt DateTime @default(now())
        updatedAt DateTime @updatedAt
        published Boolean
        title     String
        content   String?

        author   User @relation(fields: [authorID], references: [id])
        authorID String
    }
    ```

    이 schema 정의를 너의 데이터베이스에 적용시킬때, Prisma Client Go 는 Prisma migration tool[`migrate`](https://github.com/prisma/migrate) (Note: this tool is experimental) 을 사용해야함.
    ```shell script
    # 첫 번째 migration 초기화
    go run github.com/prisma/prisma-client-go migrate save --experimental --create-db --name "init"
    # migration 적용
    go run github.com/prisma/prisma-client-go migrate up --experimental
    ```

4) Prisma Client Go 클라이언트를 프로젝트 안에서 생성하기

    ```shell script
    go run github.com/prisma/prisma-client-go generate
    ```

    Prisma Client Go 클라이언트는 "output" 옵션 (여기서는 `"./db/db_gen.go"`) 을 통해 지정한 파일경로에 생성됨.  
    만약 prisma schema 에 변경사항이 있으면, 이 명령어를 한 번 더 실행해야함


## 사용

Prisma Client Go 클라이언트를 일단 한 번 생성하고, Prisma를 통해 데이터소스를 세팅하면, 시작할 준비가 됨!

Prisma Client Go 에서는 클라이언트를 `./db/db_gen.go`에 `db`라는 이름의 패키지로 생성하는걸 추천함 (위의 step3 를 확인).  물론 이 세팅을 원하는 아무곳에나 적용할 수도 있음.



### 클라이언트를 생성하고 prisma 엔진에 연결하기

```go
client := db.NewClient()
err := client.Connect()
if err != nil {
    handle(err)
}

defer func() {
    err := client.Disconnect()
    if err != nil {
        panic(fmt.Errorf("연결을 끊을 수 없습니다 %w", err))
    }
}()
```


### 예시 코드 전체

```go
package main

import (
    "context"
    "log"
    "github.com/your/repo/db"
)

func main() {
    client := db.NewClient()
    err := client.Connect()
    if err != nil {
        panic(err)
    }

    defer func() {
        err := client.Disconnect()
        if err != nil {
            panic(err)
        }
    }()

    ctx := context.Background()

    // 유저 생성
    createdUser, err := client.User.CreateOne(
        db.User.Email.Set("gildong@example.com"),
        db.User.Name.Set("홍길동"),

        // ID 는 옵션, 입력을 안해도 나중에 자동으로 입력되기 때문
        db.User.ID.Set("123"),
    ).Exec(ctx)

    log.Printf("생성된 유저: %+v", createdUser)

    // 한명의 유저 찾기
    user, err := client.User.FindOne(
        db.User.Email.Equals("gildong@example.com"),
    ).Exec(ctx)
    if err != nil {
        panic(err)
    }

    log.Printf("유저: %+v", user)

    // optional/nullable 한 값에 대해서는, 두 개의 리턴 값을 받아야함
    // `name` 은 string 값이고, `ok` 는 저장된 값이 null 인지 아닌지를 체크하는 bool 값임
    // 만약에 가져온 값이 null이면 `ok` 는 false 이고, `name` 은 Go의 default 값이 됨 (이 경우에는 빈 문자열 "")
    // 가져온 값이 null 이 아닌경우에 `ok` 는 true 이고 `name` 은 "홍길동" 이 됨.
    name, ok := user.Name()

    if !ok {
        log.Printf("user의 name이 null 입니다.")
        return
    }

    log.Printf("user의 이름은: %s", name)
}
```

이 문서는 prisma-client-go 공식 quickstart 의 번역입니다.
[원본 문서보기](https://github.com/prisma/prisma-client-go/blob/master/docs/quickstart.md)
