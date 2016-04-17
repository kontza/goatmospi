#!/bin/sh
__ATMOSPI_STATIC="$(pwd)/web/static"
rm 000-atmospi.conf >& /dev/null
cat << EOF > 000-atmospi.conf
<VirtualHost *:80>
    ProxyRequests Off
    Alias "/static" $__ATMOSPI_STATIC
    <Directory $__ATMOSPI_STATIC>
        Require all granted
    </Directory>
    <Location /atmospi>
        RequestHeader    set Atmospi-Prefix-Path "atmospi"
        ProxyPass        http://localhost:4002
        ProxyPassReverse http://localhost:4002
    </Location>
</VirtualHost>
EOF
echo Next: 'sudo ln -sf $PWD/000-atmospi.conf /etc/apache2/sites-enabled/'
