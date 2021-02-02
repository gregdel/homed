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
import { ESPLight } from "./EspLight";

import { Row, Col, Card } from "antd";

export const HomedComponents = () => {
  const components = useSelector((state) =>
    [...state.components.components].map(([, value]) => ({
      id: value.values.id,
      roomName: value.values.room_name,
      deviceName: value.values.device_name,
      updatedAt: value.values.updated_at,
      type: value.type,
    }))
  );

  return (
    <Row gutter={[10, 10]}>
      {components.map(({ id, roomName, deviceName, updatedAt, type }) => (
        <Col key={id} xs={24} sm={12} lg={8}>
          <HomedComponent
            id={id}
            type={type}
            deviceName={deviceName}
            roomName={roomName}
            updatedAt={updatedAt}
          />
        </Col>
      ))}
    </Row>
  );
};

export const HomedComponent = ({
  id,
  type,
  deviceName,
  roomName,
  updatedAt,
}) => {
  var typedComponent;
  switch (type) {
    case "temperature":
      typedComponent = <Temperature id={id} />;
      break;
    case "humidity":
      typedComponent = <Humidity id={id} />;
      break;
    case "device_status":
      typedComponent = <DeviceStatus id={id} />;
      break;
    case "tuya_trv":
      typedComponent = <TuyaTRV id={id} />;
      break;
    case "wifi_signal":
      typedComponent = <WifiSignal id={id} />;
      break;
    case "rtl_433":
      typedComponent = <RTL433 id={id} />;
      break;
    case "boiler":
      typedComponent = <Boiler id={id} />;
      break;
    case "tasmota_switch":
      typedComponent = <TasmotaSwitch id={id} />;
      break;
    case "esphome_light":
      typedComponent = <ESPLight id={id} />;
      break;
    case "homed_temperature":
      // Don't display the homed temperature here
      return null;
    default:
      typedComponent = <>Unhandled {type}</>;
      break;
  }

  var prettyDate = "";
  const m = moment(updatedAt, "YYYY-MM-DD HH:mm:ss Z");
  if (m.isValid()) {
    prettyDate = m.fromNow();
  }

  const title = `${roomName} - ${type} - ${deviceName}`;
  return (
    <Card title={title} extra={prettyDate}>
      {typedComponent}
    </Card>
  );
};
HomedComponent.propTypes = {
  id: PropTypes.string.isRequired,
  type: PropTypes.string.isRequired,
  deviceName: PropTypes.string.isRequired,
  roomName: PropTypes.string.isRequired,
  updatedAt: PropTypes.string,
};
