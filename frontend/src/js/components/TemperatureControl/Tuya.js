import React from "react";
import { useSelector } from "react-redux";
import { prettyName } from "../../utils";
import moment from "moment";

import { TuyaTRV } from "../HomedComponents/TuyaTRV";

import { Row, Col, Card } from "antd";

export const Tuya = () => {
  const components = useSelector((state) =>
    [...state.stuff.components]
      .filter(([, component]) => component && component.type === "tuya_trv")
      .map(([, value]) => ({
        room: value.room,
        uuid: value.values.uuid,
        updatedAt: moment(
          value.values.updated_at,
          "YYYY-MM-DD HH:mm:ss Z"
        ).fromNow(),
      }))
  );

  return (
    <Row gutter={[10, 10]}>
      {components.map(({ room, uuid, updatedAt }) => (
        <Col key={uuid} xs={24} sm={12} lg={8}>
          <Card title={prettyName(room)} extra={updatedAt}>
            <TuyaTRV uuid={uuid} />
          </Card>
        </Col>
      ))}
    </Row>
  );
};
