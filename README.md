# OneToNine

Go로 쓰여진 카카오톡 숫자야구 봇입니다. PostgreSQL이 필요합니다. 


## 테스트 환경

OS: Ubuntu 16.04.5 LTS armv7l

Host: Raspberry Pi 2 Model B Rev 1.1


OS: macOS High Sierra 10.13.6 17G65 x86_64

Host: MacBookPro11,1

그 외의 환경에선 실행하지 않았습니다.

## 설치
```shell
git clone https://github.com/JedBeom/OneToNine
```

먼저 PostgreSQL에 'onetonine'이라는 이름의 DB를 생성해주셔야합니다.
DB만 생성하시고 테이블은 하지 말아주세요. 프로그램 첫 실행 시 자동 생성됩니다.

main.go:15 해당 DB 수정 권한이 있는 username과 password를 수정해주세요.

main.go:78 기본 포트는 80이지만, 다른 개방된 포트로 수정하여도 됩니다.

```shell
go get github.com/jinzhu/gorm
go get github.com/lib/pq
```

```shell
cd OneToNine
go build
```

## 사용
80포트를 쓸 경우:
```shell
sudo ./OneToNine
```

그 외의 포트를 사용할 경우:
```shell
./OneToNine
```

<https://center-pf.kakao.com>에 접속하여 플러스친구를 생성한 뒤, 스마트채팅 > Api형 설정에 들어가 챗봇을 실행할 서버의 IP주소와 포트, 또는 도메인과 포트를 입력하고 'Api형 저장하기'를 클릭합니다.

카카오톡에 들어가 해당 플러스친구와 '1:1채팅하기'로 들어갑니다.

