#!/bin/bash
set -l __ATMOSPI_STATIC "$PWD/web/static"
set -l __ATMOSPI_PREFIX (basename $PWD)
set -l __CONF_NAME "000-$__ATMOSPI_PREFIX.conf"
rm $__CONF_NAME ^&1
echo -n "<VirtualHost *:80>
    ProxyRequests Off
    <Location /$__ATMOSPI_PREFIX>
        RequestHeader    set Atmospi-Prefix-Path \"$__ATMOSPI_PREFIX\"
        ProxyPass        http://localhost:4002
        ProxyPassReverse http://localhost:4002
    </Location>
</VirtualHost>" > $__CONF_NAME
echo Next:
echo "sudo ln -sf \$PWD/$__CONF_NAME /etc/apache2/sites-enabled/"
echo "sudo apache2ctl restart"
