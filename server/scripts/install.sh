#! /bin/bash

installSoftware() {
    apt -qq -y install nginx
}

installWeather() {
    mkdir -p /etc/weather
    curl -Lo- https://github.com/sunshineplan/weather/releases/latest/download/release.tar.gz | tar zxC /etc/weather
    cd /etc/weather
    chmod +x weather
}

configWeather() {
    read -p 'Please enter metadata server: ' server
    read -p 'Please enter VerifyHeader header: ' header
    read -p 'Please enter VerifyHeader value: ' value
    read -p 'Please enter unix socket(default: /run/weather.sock): ' unix
    [ -z $unix ] && unix=/run/weather.sock
    read -p 'Please enter host(default: 127.0.0.1): ' host
    [ -z $host ] && host=127.0.0.1
    read -p 'Please enter port(default: 12345): ' port
    [ -z $port ] && port=12345
    read -p 'Please enter log path(default: /var/log/app/weather.log): ' log
    [ -z $log ] && log=/var/log/app/weather.log
    read -p 'Please enter update URL: ' update
    sed "s,\$server,$server," /etc/weather/config.ini.default > /etc/weather/config.ini
    sed -i "s/\$header/$header/" /etc/weather/config.ini
    sed -i "s/\$value/$value/" /etc/weather/config.ini
    sed -i "s,\$unix,$unix," /etc/weather/config.ini
    sed -i "s,\$log,$log," /etc/weather/config.ini
    sed -i "s/\$host/$host/" /etc/weather/config.ini
    sed -i "s/\$port/$port/" /etc/weather/config.ini
    sed -i "s,\$update,$update," /etc/weather/config.ini
    ./weather install || exit 1
    service weather start
}

writeLogrotateScrip() {
    if [ ! -f '/etc/logrotate.d/app' ]; then
	cat >/etc/logrotate.d/app <<-EOF
		/var/log/app/*.log {
		    copytruncate
		    rotate 12
		    compress
		    delaycompress
		    missingok
		    notifempty
		}
		EOF
    fi
}

setupNGINX() {
    cp -s /etc/weather/scripts/weather.conf /etc/nginx/conf.d
    sed -i "s/\$domain/$domain/" /etc/weather/scripts/weather.conf
    sed -i "s,\$unix,$unix," /etc/weather/scripts/weather.conf
    service nginx reload
}

main() {
    read -p 'Please enter domain:' domain
    installSoftware
    installWeather
    configWeather
    writeLogrotateScrip
    setupNGINX
}

main
