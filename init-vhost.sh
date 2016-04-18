#!/bin/bash
__ATMOSPI_STATIC="$(pwd)/web/static"
__ATMOSPI_PREFIX="goatmospi"
__CONF_NAME="000-$__ATMOSPI_PREFIX.conf"
rm 000-atmospi.conf >& /dev/null
cat << EOF > $__CONF_NAME
<VirtualHost *:80>
    ProxyRequests Off
    <Location /$__ATMOSPI_PREFIX>
        RequestHeader    set Atmospi-Prefix-Path "$__ATMOSPI_PREFIX"
        ProxyPass        http://localhost:4002
        ProxyPassReverse http://localhost:4002
    </Location>
</VirtualHost>
EOF
echo Next:
echo "sudo ln -sf \$PWD/$__CONF_NAME /etc/apache2/sites-enabled/"
echo "sudo apache2ctl restart"
