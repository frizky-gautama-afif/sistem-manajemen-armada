# sistem-manajemen-armada

CARA MENJALANKAN APLIKASI:

1. Pull repository dari github.
2. Siapkan/install docker pada device anda.
3. Buka terminal/gitbash dan arahkan pada directory aplikasi
4. Jalankan command "docker-compose up -d --build" pada terminal anda.
5. Jalankan command "docker-compose logs -f mqtt-publisher" untuk cek apakah mock MQTT berjalan dan berhasil dikirim ke MQTT.
6. Jalankan command "docker-compose logs -f backend" untuk cek apakah backend berhasil menerima pesan dari MQTT dan juga berhasil simpan ke Database Postgres SQL.
7. Hit API lokasi terakhir yang ada pada collection untuk mendapatkan data lokasi terahir armada.
8. Hit API history lokasi yang ada pada collection untuk mendapatkan data history lokasi armada.
9. Jalankan command "docker-compose logs -f rabbitmq-worker" untuk cek kendaraan yang masuk pada geofence, jika ada "ALERT" pada log, maka sudah berfungsi dengan benar.
   
