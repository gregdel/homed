#!/bin/sh
set -e

_log() {
	echo "$@"
}

_err() {
	_log "$@"
	exit 1
}

usage="$(basename "$0") LOCATION DEVICE_NAME MQTT_BROKER"

location=$1
device_name=$2
mqtt_broker=$3
[ "$location" ]    || _err "$usage"
[ "$device_name" ] || _err "$usage"
[ "$mqtt_broker" ] || _err "$usage"

url="https://wttr.in/$location?format=%t;%h"

_publish() {
	data_type=$1
	value=$2
	mosquitto_pub \
		-h "$mqtt_broker" \
		-t "home/devices/$device_name/sensor/$data_type/state" \
		-m "$value"
}

data=$(curl -sS "$url")
[ "$data" ] || _err "Missing data.."

# data looks like this: '+3°C;87%'
temperature=${data%;*}				# +3°C
temperature=${temperature%°*}       # +3
temperature=$(( temperature * 1 ))  # 3
humidity=${data#*;}                 # 87%
humidity=${humidity%\%*}            # 87

_log "Temperature: $temperature°C"
_log "Humidity: $humidity%"

_publish "humidity" "$humidity"
_publish "temperature" "$temperature"
