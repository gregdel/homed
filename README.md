# Idea

homed - Home daemon suite

homed-temperatured ?

homed/
  devices
    device_name.json
  sensors
    sensor_name.json
  actions
    action_name.json
  rooms
    room_name.json
  apps/
    thermos/
      schedules/
        schedule_name_X.json
        schedule_name_Y.json
      rooms/
        room_1 -> schedule_name_X.json
        room_2 -> schedule_name_X.json
        room_2 -> schedule_name_Y.json

homed is just a lib to handle mqtt messages and configurations
thermos uses the homed lib to read / write configurations

homed (lib):
* notification channel to alert apps from config change
* a mechanism to republish events

homectl (bin):
* add/remove/list rooms
* add/remove/list sensors
* add/remove/list actions
* subscribe to events and execute commands

thermos (bin):
*


* Take sensor values from mqtt
* Take order from mqtt (e.g. room temperatures)

All orders / configs are in memory and flagged dirty until written to disk ?
