PORT="1323"
APP_SECURE=false
APP_VERSION="beta"
AUTH_TOKEN_ACCESS_SECRET="jwtsecretkey"
CORS_ALLOW_ORIGIN="*"
CORS_MAX_AGE="3600"
HTTP_MAXCONNS=100
HTTP_INSECURESKIPVERIFY=true
HTTP_TIMEOUT="5s"
OTP_URL="https://ohgikdu5ed.execute-api.ap-southeast-1.amazonaws.com/smsgw-api"
OTP_USER="morchana2"
OTP_PASSWORD="Y8adfQzJfwUKGwUY"
OTP_FROM="Morchana"
OTP_MESSAGE="ทดสอบ ข้อความ ภาษาไทย"
DB_CONN_STRING="User ID=root;Password=myPassword;Host=localhost;Port=5432;Database=myDataBase;Pooling=true;Min Pool Size=0;Max Pool Size=100;Connection Lifetime=0;"

run:
	go run main.go

secure:
	APP_SUCURE=true go run main.go

test:
	go test ./...