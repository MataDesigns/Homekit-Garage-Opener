# Garage Opener

This is a simple implementation of [brutella/hc](https://github.com/brutella/hc). Turning a **Liftmaster** into a homekit compatible garage door opener.

- [Parts](#parts)
- [Features](#features)
- [Development](#development)
    - [**Basics**](#basics) - [Install Location](#deb-install-path), [Logs](#logs), [Config](#config-file)
    - [Executable](#executable)
    - [Packaging](#packaging)
    - [Installing](#installing)

## Parts
| Name                | Link    | Price Estimate |
| --------------------|:-------:|:--------------:|
| Raspberry Pi        | [Zero W](https://www.amazon.com/ELEMENT-Element14-Raspberry-Pi-Motherboard/dp/B07BDR5PDW/ref=sr_1_5?s=pc&ie=UTF8&qid=1545888929&sr=1-5&keywords=raspberry+pi) /  [3B+](https://www.raspberrypi.org/products/raspberry-pi-3-model-b-plus/) | Zero W ~ $14 / 3b+ ~$35 |
| 5V 2.5A Charger (For Pi) | [Newark](https://www.newark.com/stontronics/t5875dv/psu-raspberry-pi-5v-2-5a-multi/dp/77Y6535?src=raspberrypi) | ~ $9
| 5v Relay Module     | [Amazon](https://www.amazon.com/Tolako-Arduino-Indicator-Channel-Official/dp/B00VRUAHLE/ref=sr_1_6?ie=UTF8&qid=1544893749&sr=8-6&keywords=relay+module) / [Frys](https://www.frys.com/product/9410451?site=sr:SEARCH:MAIN_RSLT_PG) | ~ $5
| Jumper Wires F to F | [Amazon](https://www.amazon.com/EDGELEC-Optional-Breadboard-Assorted-Multicolored/dp/B07GD2BWPY/ref=sr_1_1_sspa?ie=UTF8&qid=1545889448&sr=8-1-spons&keywords=jumper+wire+female+to+female&psc=1) | ~ $7
| Wire any guage | [Amazon](https://www.amazon.com/Electronix-Express-Hook-Wire-Solid/dp/B00B4ZRPEY/ref=sr_1_3?s=industrial&ie=UTF8&qid=1545889628&sr=1-3&keywords=solid+wire) | ~ $14
||**Total**|  $49 ~ $70|


## Features

- ✅ Opening and Closing
    - 🔜 True stateful opening and closing (will require distance sensor)
- 🔜 Obstacle Detection

## Development
### Basics
#### Deb install path
The deb package installs files to `/opt/garage/`
#### Logs
Logs can be found at `/opt/garage/logs`
#### Config File
The config file allows you make minor changes without touching code

Optional Values 

- **Name (Default=Garage Door)**: The text that will appear in Homekit as the name of the accessory
- **Manufacturer (Default=Liftmaster)**: The manufacturer of the accessory
- **Model (Default=Professional 1/3 HP)**: The model of the accessory
- **SerialNumber (Default=l1f7m4573r)**: The serial number of the accessory

**IMPORTANT Manufacturer, Model, SerialNumber has no effect besides showing it under accessory information in Homekit.**

Required Values
- **Pin**: The BCM gpio pin that will active the relay.

Simple Config
```json
{
    "pin": 2
}
```

### Executable
If you would rather build to a binary/executable you can just use `go build` like
```
go build -o garage main.go
```
### Packaging
This project uses nfpm to make a deb package for installing on linux with armhf architecture. It is fairly simple to enable other architectures, requested if would be useful.
```
make package-linux arch=armhf
```

### Installing
```
sudo dpkg -i INSERTDEBFILEHERE
```

