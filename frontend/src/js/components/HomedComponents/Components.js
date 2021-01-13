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

import { Row, Card } from "antd";

export const HomedComponents = () => {
  const HomedComponents = useSelector((state) => state.stuff.components);

  var items = [];
  HomedComponents.forEach((value, key) => {
    items.push(
      // <Col key={key} flex={1}>
      <HomedComponent key={key} uuid={key} {...value} />
      // </Col>
    );
  });

  return <Row gutter={16}>{items}</Row>;
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
    <Card
      style={{ width: 300, margin: "0.5em" }}
      title={title}
      extra={prettyDate}
    >
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
