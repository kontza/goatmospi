echo "localhost/goatmospi/data/devices/temperature"
curl -H "Accept-Encoding: application/json" localhost/goatmospi/data/devices/temperature
sleep 0.5
echo "
localhost/goatmospi/data/devices/humidity"
curl localhost/goatmospi/data/devices/humidity
sleep 0.5
echo "
localhost/goatmospi/data/latest/temperature"
curl localhost/goatmospi/data/latest/temperature
sleep 0.5
echo "
localhost/goatmospi/data/latest/humidity"
curl localhost/goatmospi/data/latest/humidity
