# To-Do-Golang

install Go, ini linknya "https://go.dev/doc/install"
masuk ke folder back-end 
    -cd back-end
di terminal

terus di terminal ketik go run ./main.go\
nanti dia bakal ngejalanin kode main.go yang dimana otomatis buka port 8080 buat localhost buat ngakses apinya

database harus ada, kalo nggak nanti dia nggak bisa connect
di folder database di directory back-end ada database.sql, copy code itu terus paste di query aplikasi sql manager atau apalah itu namanya, saya sendiri pake HeidiSQL, terus run, dia nanti akan membuat database + tables yang diperlukan di project ini

buat testing api bisa pake postman, tinggal masukin aja http://localhost:8080/api/<api_yang_mau_dipake>

list API di aplikasi ini
api/token = login
api/work = crud work
api/task = crud task

buat liat fullnya isi crud bisa diliat di file main.go, di bagian router (r.HandleFunc...)

contoh yang paling gampang itu login
tinggal buka postman, bikin reqest baru, set request ke POST terus masukin url
"http://localhost:8080/api/token", terus pergi ke Body, pilih raw, pastiin formatnya JSON, tinggal ketik

{
    "username": "admin",
    "password": "admin"
}

nanti tinggak send, harusnya dia return token

buat api yang lain, itu karena hanya bisa di akses pake token, kamu harus pake Headers di postmannya, ada di jajaran yang sama kayak Body, kamu masuk ke Headers, kalo udah tinggal tambahin key -- Authorization --, terus valuenya pake format -- Bearer kodetokenmu --. 


