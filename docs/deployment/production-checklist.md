# Production Checklist Cho Phuong An Sieu Re

## 1. Muc tieu

Checklist nay dung cho phuong an:

- website ca nhan `quorix.io.vn` giu nguyen tren Vercel
- frontend `one-time-link` deploy tren Vercel
- backend Go + Redis chay tren 1 VPS
- domain van do PA Viet Nam quan ly

## 2. Cau truc domain de xuat

- `quorix.io.vn`: website ca nhan Hugo/PaperMod
- `secret.quorix.io.vn`: frontend one-time-link
- `api.secret.quorix.io.vn`: backend API

## 3. Viec can chuan bi truoc khi deploy

### Tai khoan

- tai khoan Vercel
- tai khoan VPS provider
- quyen quan ly DNS cua domain tai PA Viet Nam

### May chu

- 1 VPS Linux nho
- SSH key
- IP public

### Du an

- frontend build duoc o local
- backend Go chay duoc o local
- Redis config tuong thich local va production
- env vars duoc tach rieng cho frontend va backend

## 4. DNS checklist

### Tai PA Viet Nam

- giu record cho `quorix.io.vn` nhu hien tai
- tao `CNAME` cho `secret.quorix.io.vn` theo huong dan Vercel
- tao `A` record cho `api.secret.quorix.io.vn` tro toi IP VPS

### Kiem tra

- `secret.quorix.io.vn` resolve dung ve Vercel
- `api.secret.quorix.io.vn` resolve dung ve VPS
- SSL tren ca hai subdomain hoat dong binh thuong

## 5. Vercel checklist cho frontend

### Project setup

- tao project rieng cho frontend one-time-link
- root directory tro dung vao folder frontend
- build command va output directory dung voi framework dang dung

### Environment variables

- `VITE_API_BASE_URL=https://api.secret.quorix.io.vn`

### Domain

- add custom domain `secret.quorix.io.vn`
- verify domain theo huong dan cua Vercel

### Kiem tra

- frontend load duoc
- frontend goi duoc API production
- khong co mixed content

## 6. VPS checklist

### He dieu hanh

- Ubuntu LTS hoac Debian stable

### Bao mat co ban

- tao user rieng, khong lam viec bang root
- tat dang nhap password qua SSH
- chi dung SSH key
- doi SSH port neu ban muon giam bot scan, nhung day la toi uu nho, khong phai lop bao mat chinh
- bat firewall

### Goi can cai

- Go runtime neu can build tren server, hoac copy binary build san
- Redis
- Caddy hoac Nginx
- systemd service files

### Kiem tra

- chi mo cong `80` va `443`, va `22` neu can
- Redis chi bind localhost
- backend chay o cong noi bo, vi du `127.0.0.1:8080`

## 7. Reverse proxy checklist

### De xuat

- dung `Caddy` cho nhanh va gon

### Cac viec can co

- HTTPS tu dong
- reverse proxy tu `api.secret.quorix.io.vn` vao app Go
- bat gzip hoac compression co ban
- log request o muc toi thieu

### Kiem tra

- `https://api.secret.quorix.io.vn/healthz` tra ve thanh cong
- chung chi TLS hop le

## 8. Backend checklist

### Configuration

- `PORT`
- `REDIS_ADDR`
- `REDIS_PASSWORD` neu co
- `ALLOWED_ORIGIN=https://secret.quorix.io.vn`
- `LOG_LEVEL`
- `MAX_SECRET_SIZE`
- `TTL_ALLOWED_VALUES`

### API can co toi thieu

- `POST /api/secrets`
- `GET /api/secrets/{id}/status`
- `POST /api/reveal-sessions`
- `POST /api/secrets/{id}/consume`
- `GET /healthz`

### Bao mat va on dinh

- validate request body
- gioi han kich thuoc payload
- CORS chi mo cho frontend cua ban
- rate limiting co ban
- khong log plaintext hay fragment key
- atomic consume voi Redis

## 9. Redis checklist

### Cau hinh

- bind `127.0.0.1`
- dat password neu can, du la chi local
- bat append-only neu ban muon giam rui ro mat du lieu do restart

### Nghiep vu

- luu secret voi TTL
- luu reveal session voi TTL ngan
- consume secret bang lenh atomic

### Kiem tra

- secret het han thi tu xoa
- hai request dong thoi chi co 1 request lay duoc secret

## 10. Frontend checklist ve UX va security

### UX

- trang tao secret ro rang
- trang reveal co nut bam ro rang
- state `already used`, `expired`, `invalid` hien thi than thien

### Security

- ma hoa o client
- key giai ma nam trong fragment `#`
- khong gui fragment len server
- khong luu plaintext trong localStorage neu khong can

## 11. Monitoring toi thieu

### Logs

- log co timestamp
- log co request id neu co
- log ket qua `created`, `revealed`, `expired`, `blocked`

### Health

- `/healthz` cho app
- kiem tra Redis connectivity trong health check hoac readiness

### Canh bao toi thieu bang tay

- neu frontend khong goi duoc API
- neu Redis khong ket noi duoc
- neu ty le `consume` loi tang bat thuong

## 12. Backup toi thieu

### Can backup gi

- file cau hinh server
- Caddyfile hoac Nginx config
- systemd unit files
- script deploy
- source code dang luu tren Git remote

### Can backup du lieu secret khong?

Khong nhat thiet cho MVP.

Ly do:

- day la du lieu ngan han
- backup secret da ma hoa van tao them do phuc tap van hanh
- voi san pham one-time-link, mat secret sau su co co the chap nhan hon so voi lo secret

## 13. Thu tu deploy de xuat

1. Chuan bi frontend build on o local.
2. Chuan bi backend Go + Redis chay on o local.
3. Thue VPS va harden co ban.
4. Cai Redis, Caddy, binary Go.
5. Cau hinh systemd cho backend.
6. Add DNS `api.secret.quorix.io.vn`.
7. Test API production bang `healthz`.
8. Tao project frontend tren Vercel.
9. Add DNS `secret.quorix.io.vn`.
10. Cau hinh env var frontend tro toi API production.
11. Test full flow create -> reveal -> already used.
12. Test expired flow.
13. Test preview-bot-safe flow bang cach mo trang ma khong bam reveal.

## 14. Tieu chi xem nhu deploy thanh cong

- `secret.quorix.io.vn` mo duoc
- tao duoc link moi
- mo link chi hien trang gate, chua reveal ngay
- bam reveal thi doc duoc secret
- bam lan hai nhan `already used`
- secret het han thi tra ve `expired`
- backend khong log plaintext

## 15. Khuyen nghi cuoi

Khong can lam day du tat ca ky thuat production nang ngay tu dau.

Thu quan trong nhat cho portfolio la:

- dung bai toan
- dung flow bao mat co ban
- deploy that
- domain that
- giai thich duoc trade-off

Neu 4 dieu nay tot, recruiter da co the nhin thay chat luong cua du an.
