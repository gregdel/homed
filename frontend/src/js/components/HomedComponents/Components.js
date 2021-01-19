import React from "react";
import PropTypes from "prop-types";
import { useSelector } from "react-redux";
import moment from "moment";

import { Temperature } from "./Temperature";
import { Humidity } from "./Humidity";
import { DeviceStatus } from "./DeviceStatus";
import { TuyaTRV } from "./TuyaTRV";
import { WifiSignal } from "./WifiSignal";
import { RTL433 } from "./RTL433";
import { Boiler } from "./Boiler";
import { TasmotaSwitch } from "./TasmotaSwitch";

import { Row, Col, Card } from "antd";

export const HomedComponents = () => {
  const HomedComponents = useSelector((state) => state.stuff.components);

  var items = [];
  HomedComponents.forEach((value, key) => {
    items.push(
      <Col key={key} xs={24} sm={12} lg={8}>
        <HomedComponent key={key} uuid={key} {...value} />
      </Col>
    );
  });

  return <Row gutter={[10, 10]}>{items}</Row>;
};

export const HomedComponent = ({ uuid, room, device, type, values }) => {
  var typedComponent;
  switch (type) {
    case "temperature":
      typedComponent = <Temperature {...values} />;
      break;
    case "humidity":
      typedComponent = <Humidity {...values} />;
      break;
    case "device_status":
      typedComponent = <DeviceStatus {...values} />;
      break;
    case "tuya_trv":
      typedComponent = <TuyaTRV room={room} {...values} />;
      break;
    case "wifi_signal":
      typedComponent = <WifiSignal {...values} />;
      break;
    case "rtl_433":
      typedComponent = <RTL433 uuid={uuid} {...values} />;
      break;
    case "boiler":
      typedComponent = <Boiler {...values} />;
      break;
    case "tasmota_switch":
      typedComponent = <TasmotaSwitch {...values} />;
      break;
    case "homed_temperature":
      // Don't display the homed temperature here
      return null;
    default:
      typedComponent = <>Unhandled {type}</>;
      break;
  }

  var prettyDate = "";
  const updatedAt = moment(values.updated_at, "YYYY-MM-DD HH:mm:ss Z");
  if (updatedAt.isValid()) {
    prettyDate = updatedAt.fromNow();
  }

  const title = `${room} - ${type} - ${device}`;
  return (
    <Card title={title} extra={prettyDate}>
      {typedComponent}
    </Card>
  );
};
HomedComponent.propTypes = {
  uuid: PropTypes.string.isRequired,
  room: PropTypes.string.isRequired,
  device: PropTypes.string.isRequired,
  type: PropTypes.string.isRequired,
  values: PropTypes.object.isRequired,
};
