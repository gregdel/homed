import React, { useEffect } from "react";
import { useDispatch } from "react-redux";
import { fetchStuff } from "../actions/stuff";

import { HomedComponents } from "./HomedComponents/Components";

export const Dashboard = () => {
  const dispatch = useDispatch();

  useEffect(() => {
    dispatch(fetchStuff());
  }, [dispatch]);

  return <HomedComponents />;
};
